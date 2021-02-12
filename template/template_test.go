package template

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)




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


func TestExtractUsingPaths(t *testing.T){
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

func ignoreExampleTemplate() {
	data := struct {
		Title   string
		Message string
	}{
		Title:   "Some Title",
		Message: "Hello, World",
	}

	parser := &Parser{layout: layout}
	output1, err := parser.Parse(page1, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Page 1 output: ", output1)

	output2, err := parser.Parse(page2, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Page 2 output: ", output2)

	// Output:
	// Page 1 output:
	// <html>
	//   <head>
	//     <title>Some Title</title>
	//   </head>
	//   <body>
	//     <div>"Hello, World" from Page 1</div>
	//   </body>
	// </html>
	// Page 2 output:
	// <html>
	//   <head>
	//     <title>Some Title</title>
	//   </head>
	//   <body>
	//     <div>"Hello, World" from Page 2</div>
	//   </body>
	// </html>
}
