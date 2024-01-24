package main

import (
	"reflect"
	"testing"
)

func TestRender(t *testing.T) {
	cases := []struct {
		line string
		want [][2]string
	}{
		{line: "package main", want: [][2]string{{}}},
		{
			line: `var s = fmt.Sprintf("%s", "hi")`,
			want: [][2]string{
				{"keyword", "var"}, {"plain", " "}, {"plain", "s"},
				{"plain", " "}, {"plain", "="}, {"plain", " "},
				{"plain", "fmt"}, {"plain", "."}, {"function", "Sprintf"},
				{"plain", "("}, {"string", "\"%s\""}, {"plain", ","},
				{"plain", " "}, {"string", "\"hi\""}, {"plain", ")"},
			},
		},
	}

	for _, c := range cases {
		got := parseSyntax(c.line)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("got %#v, want %#v", got, c.want)
		}
	}
	var token string
	token += string([]byte{'a'})
}
