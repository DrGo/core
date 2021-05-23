package md

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"testing"

	"github.com/matryer/is"
)

 const testDir="test/content"


func Test_ParseContentFile(t *testing.T) {
	p, err := ParseContentFile(os.DirFS(testDir),"team/investigators/salah-mahmud.md")
	if err != nil {
		t.Fatalf("error %v", err)
	}
  for k, v := range p.FrontMatter {
	  fmt.Printf("%v=%v\n", k, v)
  }
	fmt.Printf("The body is:\n%q\n", p.Body)
  const tpl=`{{.Prop "title"}}`
  tm := template.Must(template.New("test").Parse(tpl))
  tm.Execute(os.Stdout, p)
  var buf bytes.Buffer
  is := is.New(t)
  is.NoErr(tm.Execute(&buf, p))
  is.Equal(buf.String(),"Salah Mahmud") 
}

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
	 {{.Prop "xtitle"}}
  </head>
	<body>
	   {{.Body}}
  </body>
</html>`
