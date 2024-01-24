package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

type Data struct {
	Files       []string
	Texts       []string
	LineNumbers []int
	CurrentFile string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var files []string
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range entries {
		if f.Name()[0] == '.' || f.IsDir() {
			continue
		}
		files = append(files, f.Name())
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := Data{Files: files}
		filename := r.URL.Path[1:]
		if filename == "" {
			err := tmpl.Execute(w, data)
			if err != nil {
				log.Println(err)
			}
			return
		}

		if !slices.Contains(files, filename) {
			w.WriteHeader(404)
			return
		}

		data.CurrentFile = filename
		bs, err := os.ReadFile(filename)
		if err != nil {
			w.WriteHeader(500)
			log.Println(err)
			return
		}
		data.Texts = render(string(bs))
		for i := range data.Texts {
			data.LineNumbers = append(data.LineNumbers, i+1)
		}
		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
			return
		}
	})

	addr := "127.0.0.1:8000"
	log.Printf("listening %s", addr)
	log.Println(http.ListenAndServe(addr, nil))
}

func render(text string) (lines []string) {
	for _, rawLine := range strings.Split(text, "\n") {
		if rawLine == "" {
			lines = append(lines, "<br>")
			continue
		}
		var line string
		syntaxs := parseSyntax(rawLine)
		for _, syntax := range syntaxs {
			kind, token := syntax[0], syntax[1]
			token = template.HTMLEscapeString(token)
			span := fmt.Sprintf(`<span style="color: %s;">%s</span>`, syntaxToColor[kind], token)
			line += span
		}
		lines = append(lines, line)
	}
	return lines
}

// not using iota here, can make the result of syntax parsing readble
const (
	SyntaxKeyord   = "keyword"
	SyntaxString   = "string"
	SyntaxFunction = "function"
	SyntaxComment  = "comment"
	SyntaxPlain    = "plain"
	SyntaxType     = "type"
)

var (
	syntaxToColor = map[string]string{
		SyntaxKeyord:   "purple",
		SyntaxString:   "rgb(196,27,22)",
		SyntaxPlain:    "black",
		SyntaxComment:  "gray",
		SyntaxFunction: "rgb(62,117,122)",
		SyntaxType:     "#4419A6",
	}
	// "byte", "int", "string", "nil",
	keywords = []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}
	types      = []string{"int", "string", "true", "false"}
	delimiters = []byte{' ', '(', ')', '{', '}', '.', '\t', '[', ']'}
)

func parseSyntax(line string) [][2]string {
	if line == "" {
		return nil
	}

	var result [][2]string
	var delim byte
	var inString bool
	token := make([]byte, 0, len(line))
	for i := 0; i < len(line); i++ {
		// comment
		if !inString && (len(line[i:]) >= 2 && line[i:i+2] == "//") {
			result = append(result, [2]string{SyntaxComment, line[i:]})
			break
		}
		// inside string
		if line[i] == '"' {
			if !inString {
				inString = true
			} else {
				inString = false
				result = append(result, [2]string{SyntaxString, string(append(token, '"'))})
				token = token[:0]
				continue
			}
		}
		if inString || !slices.Contains(delimiters, line[i]) {
			token = append(token, line[i])
			continue
		}
		delim = line[i]
		if len(token) > 0 {
			kind := SyntaxPlain
			st := string(token)
			if slices.Contains(keywords, st) {
				kind = SyntaxKeyord
			} else if slices.Contains(types, st) {
				kind = SyntaxType
			} else if _, err := strconv.Atoi(st); err == nil {
				kind = SyntaxType
			} else if delim == '(' {
				kind = SyntaxFunction
			}
			result = append(result, [2]string{kind, st})
		}
		result = append(result, [2]string{SyntaxPlain, string(delim)})
		token = token[:0]
	}
	if len(token) > 0 {
		kind := SyntaxPlain
		st := string(token)
		if slices.Contains(keywords, st) {
			kind = SyntaxKeyord
		} else if slices.Contains(types, st) {
			kind = SyntaxType
		} else if _, err := strconv.Atoi(st); err == nil {
			kind = SyntaxType
		}
		result = append(result, [2]string{kind, st})
	}
	return result
}
