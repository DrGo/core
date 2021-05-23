package template

import (
	"html/template"
	"testing"
)


func TestAST(t *testing.T) {
	txt := `text
        {{range $i, $a := .a}}
         ( {{$a.Bio}} ) {{comma $i (len $a) }}
        {{end}}
    `
	txt = `text
        {{range $i, $a := .a}}
         ( {{$a.Bio}} ) {{len $a | comma $i}}
        {{end}}
    `
	funcMap := map[string]interface{}{
		"comma": func(index, max int) string {
			if index < max-1 {
				return ","
			}
			return ""
		},
	}
	tmpl := template.Must(template.New("").Funcs(funcMap).Parse(txt))
	visit(tmpl.Tree.Root, printer)
}

