package template

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"strings"
)

var usingComment = "{{/* using "
var closeComment = "*/}}"

// returns nil, nil if no using comments were found
// FIXME: fragile does not tolerate spaces
// FIXME: add tests for zero or multiple using lines
func extractUsingPaths(input io.Reader) (files []string, err error) {
	s := bufio.NewScanner(input)
	// parse for custom using comment
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		fmt.Println("line:", line)
		if !strings.HasPrefix(line, usingComment) {
			return files, s.Err()
		}
		fileName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, usingComment), closeComment))
		fileName, _ = strconv.Unquote(fileName)
		fmt.Println("filename:", fileName)
		// FIXME: check for errors
		files = append(files, fileName)
	}
	return files, s.Err()
}

type Parser struct {
	layout string
}

func (p *Parser) Parse(tmplText string, data interface{}) (string, error) {
	layout, err := template.New("layout").Parse(p.layout)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("Content").Parse(tmplText)
	if err != nil {
		return "", err
	}

	if _, err := layout.AddParseTree("Content", tmpl.Tree); err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := layout.ExecuteTemplate(&buffer, "layout", data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
