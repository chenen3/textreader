package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var data struct {
		Files []string
	}
	files, err := os.ReadDir(".")
	if err != nil {
		log.Println(err)
		return
	}
	for _, file := range files {
		if file.Name()[0] == '.' || file.IsDir() {
			continue
		}
		data.Files = append(data.Files, file.Name())
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
			return
		}
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/", func(w http.ResponseWriter, r *http.Request) {
		filename, _ := strings.CutPrefix(r.URL.Path, "/file/")
		// should not read static files on every request,
		// but this is a demo, it's fine
		bs, err := os.ReadFile(filename)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(bs)
	})
	addr := "127.0.0.1:8000"
	log.Printf("listening %s", addr)
	log.Println(http.ListenAndServe(addr, nil))
}
