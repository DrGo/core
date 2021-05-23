package template

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/drgo/core/errors"
	"github.com/drgo/core/files"
)

type Template= template.Template 

// {{/* using staff-card.gohtml */}}
var usingComment = "{{/* using "
var closeComment = "*/}}"

// returns nil, nil if no using comments were found
// FIXME: add tests for zero or multiple using lines
func extractUsingPaths(input io.Reader) (files []string, err error) {
	s := bufio.NewScanner(input)
	// parse for custom using comment
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		fmt.Println("extractUsingPath line=", line)
// strings.TrimLeft(text, " \t\n")
		if !strings.HasPrefix(line, usingComment) {
			return files, s.Err()
		}
		fileName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, usingComment), closeComment))
		// fileName, _ = strconv.Unquote(fileName)
		fmt.Printf("extractUsingPath found filename:[%s]\n", fileName)
		if fileName == "" {
		  return nil, errors.Errorf("invalid using statement %s", line)
		}
		files = append(files, fileName)
	}
	return files, s.Err()
}

// ParseGlob creates a new Template and parses the template definitions from
// the files identified by the pattern. The files are matched according to the
// semantics of filepath.Match, and the pattern must match at least one file.
// The returned template will have the (base) name and (parsed) contents of the
// first file matched by the pattern. ParseGlob is equivalent to calling
// ParseFiles with the list of files matched by the pattern.
//
// When parsing multiple files with the same name in different directories,
// the last one mentioned will be the one that results.
// ParseGlob is like ParseFiles or ParseGlob but reads from the file system fs
// instead of the host operating system's file system.
// It accepts a list of glob patterns.
// (Note that most file names serve as glob patterns matching only themselves.)
func ParseGlob(fsys fs.FS, patterns ...string) (*template.Template, error) {
	filenames, err := files.GlobFS(fsys, patterns...)
	if err != nil {
		return nil, err
	}
	return parseFiles(fsys, filenames...)
}

// ParseFiles parses the named files and associates the resulting templates with
// t. If an error occurs, parsing stops and the returned template is nil;
// otherwise it is t. There must be at least one file.
//
// When parsing multiple files with the same name in different directories,
// the last one mentioned will be the one that results.
//
// ParseFiles returns an error if t or any associated template has already been executed.
func ParseFiles(fsys fs.FS, filenames ...string) (*template.Template, error) {
	return parseFiles(fsys, filenames...)
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(fsys fs.FS, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, errors.Errorf("no files named in call to ParseFiles")
	}

	// collect dependencies from each file
	deps, err := collectDeps(fsys, filenames...)
	if err != nil {
		return nil, errors.Wrapf(err, "error identifying template dependencies")
	}
	// duplicating code from
	var t *template.Template
	for _, filename := range deps {
		fmt.Printf("parsing:[%s]\n", filename)
		b, err := fs.ReadFile(fsys, filename)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing file %s", filename)
		}
		name := path.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *template.Template

		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(string(b))
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing file %s", filename)
		}
	}
	return t, nil
}

func collectDeps(fsys fs.FS, filenames ...string) ([]string, error) {
	dgraph := make(map[string][]string, len(filenames)*2)
	for _, filename := range filenames {
		f, err := fsys.Open(filename)
		if err != nil {
			return nil, err
		}
		deps, err := extractUsingPaths(f)
		if err != nil {
			return nil, err
		}
		dgraph[filename] = deps
	}
	deps, cyclics, err := KahnSort(dgraph)
	if err != nil {
		var cnames string
		if len(cyclics) > 0 {
			cnames = strings.Join(cyclics, ",")
		}
		return nil, errors.Wrapf(err, "cyclics if any:%s", cnames)
	}
	return deps, err
}

// Must is a helper that wraps a call to a function returning (*template.Template, error)
// and panics if the error is non-nil. It is intended for use in variable initializations
// such as
//	var t = template.Must(template.New("name").Parse("html"))
func Must(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return t
}
