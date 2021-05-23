package template

import (
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

const testDir="tests"

func TestParseGlob(t *testing.T){
  
	// pattern is the glob pattern used to find all the template files.
	pattern := "*.tmpl" // filepath.Join(testDir, "*.tmpl")

	// Here starts the example proper.
	// T0.tmpl is the first name matched, so it becomes the starting template,
	// the value returned by ParseGlob.
	tmpl := Must(ParseGlob(os.DirFS(testDir),pattern))
  // tList := listTemplateNodes(tmpl.Tree.Root, nil)
  // fmt.Println("Templates found:", strings.Join(tList, ","))
	err := tmpl.ExecuteTemplate(os.Stdout,"parent.tmpl", nil)
	
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	// Output:
	// T0 invokes T1: (T1 invokes T2: (This is T2))
}


var layout = `
<html>
  <head>
    <title>{{.Title}}</title>
  </head>
  <body>
    <div>{{ template "Content" . }}</div>
  </body>
</html>`

var page1 = `
{{/* using "folder/file.ext" */}}
"{{ .Message }}" from Page 1`

var page2 = `
"{{ .Message }}" from Page 2`


func xTestExtractUsingPaths(t *testing.T){
  f, err := os.Open("examples/parent.tmpl")
  if err!= nil {
     t.Fatalf("%v", err)
  }
  defer f.Close()
  files, err := extractUsingPaths(f)
  is := is.New(t)
  is.NoErr(err)
  if len(files) == 0 {
    t.Fatal("empty files")
  }
  is.Equal(files[0],"folder/file.ext")
  t.Logf("%s", files[0])
}

