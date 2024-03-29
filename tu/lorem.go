package tu

import (
	"bytes"
	"strings"
)

var r = 12345

func rand() int {
	r = r*1103515245 + 12345
	return r & 0x7fffffff
}

func main() {
}

func Lorem(count int) []byte {
	buf := &bytes.Buffer{}
	buf.Write([]byte(strings.Join(initial[:min(count, lenInitial)], " ")))
	n := count - lenInitial
	for i := 0; i <= n; i++ {
		buf.Write([]byte(loremIpsumWords[rand()%lenWordsList]))
		if rand()%20 == 1 {
			buf.Write([]byte(". "))
		} else {
			buf.WriteByte(' ')
		}
	}
	return buf.Bytes()
}

const lenInitial = 8
var initial = [...]string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit"}

var lenWordsList = len(loremIpsumWords)
var loremIpsumWords = [...]string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"a", "ac", "accumsan", "ad", "aenean", "aliquam", "aliquet", "ante",
	"aptent", "arcu", "at", "auctor", "augue", "bibendum", "blandit",
	"class", "commodo", "condimentum", "congue", "consequat", "conubia",
	"convallis", "cras", "cubilia", "curabitur", "curae", "cursus",
	"dapibus", "diam", "dictum", "dictumst", "dignissim", "dis", "donec",
	"dui", "duis", "efficitur", "egestas", "eget", "eleifend", "elementum",
	"enim", "erat", "eros", "est", "et", "etiam", "eu", "euismod", "ex",
	"facilisi", "facilisis", "fames", "faucibus", "felis", "fermentum",
	"feugiat", "finibus", "fringilla", "fusce", "gravida", "habitant",
	"habitasse", "hac", "hendrerit", "himenaeos", "iaculis", "id",
	"imperdiet", "in", "inceptos", "integer", "interdum", "justo",
	"lacinia", "lacus", "laoreet", "lectus", "leo", "libero", "ligula",
	"litora", "lobortis", "luctus", "maecenas", "magna", "magnis",
	"malesuada", "massa", "mattis", "mauris", "maximus", "metus", "mi",
	"molestie", "mollis", "montes", "morbi", "mus", "nam", "nascetur",
	"natoque", "nec", "neque", "netus", "nibh", "nisi", "nisl", "non",
	"nostra", "nulla", "nullam", "nunc", "odio", "orci", "ornare",
	"parturient", "pellentesque", "penatibus", "per", "pharetra",
	"phasellus", "placerat", "platea", "porta", "porttitor", "posuere",
	"potenti", "praesent", "pretium", "primis", "proin", "pulvinar",
	"purus", "quam", "quis", "quisque", "rhoncus", "ridiculus", "risus",
	"rutrum", "sagittis", "sapien", "scelerisque", "sed", "sem", "semper",
	"senectus", "sociosqu", "sodales", "sollicitudin", "suscipit",
	"suspendisse", "taciti", "tellus", "tempor", "tempus", "tincidunt",
	"torquent", "tortor", "tristique", "turpis", "ullamcorper", "ultrices",
	"ultricies", "urna", "ut", "varius", "vehicula", "vel", "velit",
	"venenatis", "vestibulum", "vitae", "vivamus", "viverra", "volutpat",
	"vulputate"}
