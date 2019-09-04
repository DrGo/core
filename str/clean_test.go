package str

import "testing"

func Test_cleanUpString(t *testing.T) {
	tests := []struct {
		args string
		want string
	}{
		{`should be unchanged`, `should be unchanged`},
		{`should be unchanged.`, `should be unchanged.`},
		{`(good parens) should be unchanged`, `(good parens) should be unchanged`},
		{`((duplicate parens)) should be removed`, `(duplicate parens) should be removed`},
		{`duplicate,, commas,, should be removed`, `duplicate, commas, should be removed`},
		{`duplicate  space   should be removed`, `duplicate space should be removed`},
		{`duplicate  space   should             be removed`, `duplicate space should be removed`},
		{`only ??duplicate? question marks should be removed??`, `only ?duplicate? question marks should be removed?`},
		{`duplicate period should be removed..`, `duplicate period should be removed.`},
		{`\\duplicate slashes should be removed.`, `\duplicate slashes should be removed.`},
		{`empty () should be removed`, `empty should be removed`},
		{`empty (] should be removed`, `empty should be removed`},
		{`[10.1016/j.lungcan.2016.07.016](http://dx.doi.org/10.1016/j.lungcan.2016.07.016).`,
			`[10.1016/j.lungcan.2016.07.016](http://dx.doi.org/10.1016/j.lungcan.2016.07.016).`},
		//TODO: remove parens only with space within them
		//{`empty parens ( ) should be removed`, `empty parens should be removed`},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			if got := CleanUpString(tt.args); got != tt.want {
				t.Errorf("\ngot  %s \nwant %s", got, tt.want)
			}
		})
	}
}
