package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// CleanUpString cleans a utf8 string by removing redundant punctuation
// TODO: do two passes, 1) remove all empty spaces 2) clean-up all redundant punctations.
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

func Lower(ch rune) rune { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
func ByteSlice2String(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

func IsASCIIAlphaNumeric(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || '0' <= ch && ch <= '9'
}
func OnlyASCIAlphaNumeric(s string) string {
	b := make([]byte, len(s))
	i := 0
	for _, ch := range s {
		ch := Lower(ch)
		if IsASCIIAlphaNumeric(ch) {
			b[i] = byte(ch)
			i++
		}
	}
	return ByteSlice2String(b[:i])
}


func TrimAffixes(b []byte, spacesOnly bool) string {
	start, end := 0, len(b)
	if end == 0 {
		return ""
	}
	found := false
loop:
	for ; start < end; start++ {
		switch b[start] {
		case ' ', '\t', '\n', '\r': // consume
		case '{', '"': //consume the outermost delimiter only
			if spacesOnly || found {
				break loop
			}
			found = true
		default:
			break loop
		}
	}
	end--
	found = false
	comma := false
loop1:
	for ; end >= 0; end-- {
		switch b[end] {
		case ' ', '\t', '\n', '\r': // consume
		case ',': //consume the outermost comma
			if spacesOnly || comma {
				break loop1
			}
			comma = true
		case '}', '"': //consume the outermost delimiter only
			if spacesOnly || found {
				break loop1
			}
			found = true
		default:
			break loop1
		}
	}
	// fmt.Printf("%s: %d,%d\n", string(ob), start, end)
	return (string(b[start : end+1]))
}
