package str

import (
	"strings"
	"unicode/utf8"
)

//CleanUpString cleans a utf8 string by removing redundant punctuation
//TODO: do two passes, 1) remove all empty spaces 2) clean-up all redundant punctations.
func CleanUpString(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	var nextC rune
	i := 0
	c, width := utf8.DecodeRuneInString(s[i:])
	//write the current rune and advance the string's byte position
	emit := func() {
		sb.WriteRune(c)
		c, width = utf8.DecodeRuneInString(s[i:])
		i += width
	}

	i += width //pos of next rune in s
	for {
		if i >= len(s) {
			sb.WriteRune(c) //end of string, write current rune and break
			break
		}
		switch c {
		case ' ':
			sb.WriteRune(c) //write one space and skip the next
			for c == ' ' {
				c, width = utf8.DecodeRuneInString(s[i:])
				i += width //skip it
			}
		case '*', '#', '_', '\\', '/', '(', ')', '`', '~', '[', ']', '!', '?', '{', '}', ';', ':', ',', '-', '.':
			nextC, width = utf8.DecodeRuneInString(s[i:])
			switch {
			case nextC == c: //duplicate
				i += width //skip it
				continue
			default:
				switch {
				case c == ']' && nextC == '(':
					emit() //markdown link marker [](), do not remove}
				case c == ':' && (nextC == '/' || nextC == '\\'):
					// following, followingWidth := utf8.DecodeRuneInString(s[i+width:])
					// if following == '/' { //if this :// do not remove
					// 	sb.WriteString("://")
					// 	i = i + width + followingWidth
					// 	continue
					// }
					//fallthrough //otherwise remove both :/
					emit() //http or folder marker
				case strings.ContainsAny(string(nextC), `*#_\\/()~[]!?{};:-`):
					//move to the next char
					i += width
					c, width = utf8.DecodeRuneInString(s[i:])
					i += width //skip it
				}
			}
			emit()
		default:
			emit()
		}
	}
	return strings.Join(strings.Fields(sb.String()), " ")
}
