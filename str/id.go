package str

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func tgoIdentifier() {
	fmt.Println(ToGoIdentifier("89He  llo / there?"))
}

func xmain() {
	original := "cat"
	fmt.Println(original)

	// Get rune slice.
	// ... Modify the slice's data.
	r := []rune(original)
	r[0] = 'm'
	r = append(r, 'e')

	// Create new string from modified slice.
	result := string(r)
	fmt.Println(result)
}

// ToGoIdentifier converts a string to a valid go identifier: only unicode 
// letters, digits and _ are kept. case is preserved; spaces are replaced 
// with single _; string is prefixed with _ if it starts with digits
func ToGoIdentifier(s string) string {
	var prevC rune
	s = strings.Map(func(c rune) rune {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			prevC = c
			return c
		}
		if unicode.IsSpace(c) {
			if prevC == '_' {
				return -1
			}
			prevC = '_'
			return prevC
		}
		return -1
	}, s)
	if len(s) > 0 {
		r, _ := utf8.DecodeRuneInString(s)
		if unicode.IsDigit(r) {
			return string(append([]rune{'_'}, []rune(s)...))
		}
	}
	return s
}

func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != ""
}

func Clean(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}
