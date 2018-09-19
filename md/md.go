package md

import (
	"fmt"

	"github.com/drgo/core/stacks"
)

//InlinedMdToHTMLOptions options for Markdown to HMTL conversion
type InlinedMdToHTMLOptions struct {
}

const (
	fmtNone int = iota
	fmtItalic
	fmtBold
	fmtBoldItalic
	fmtSubscript
	fmtSuperscript
)

// InlinedMdToHTML converts certain md element into html
// rules:
//^Superscript^
//~Subscript~
//* italic *
//**  bold **
//_ italic _
//__  bold __
// These characters can be escaped using \. \ itself can be escaped using \\
func InlinedMdToHTML(md string, opts *InlinedMdToHTMLOptions) ([]byte, error) {
	const null = '\x00'
	formats := make(stacks.IntStack, 0, 100)
	mdLen := len(md)
	if mdLen == 0 {
		return []byte{}, 0
	}
	html := make([]byte, mdLen*3)
	i := 0                          // is the source (md) byte index
	j := 0                          //is the dest (html) byte index
	cp := func(s string, inc int) { //copy string to the html slice
		k := 0
		for ; k < len(s); k++ {
			//fmt.Printf("s[k]=%c-->", s[k])  //DEBUG
			html[j+k] = s[k]
			//fmt.Printf("html=%s\n", html)  //DEBUG
		}
		j += k
		i += inc //increase i by inc; could be zero
	}
	next := func() byte {
		if i+1 < mdLen {
			return md[i+1]
		}
		return null
	}
	cpCurrent := func() {
		html[j] = md[i]
		j++
	}
	//TODO: check for errors
	emit := func(fmtType, inc int, tag string) {
		//fmt.Println("format=", fmtType, "top is", formats[len(formats)-1].format)
		if formats.TopIs(fmtType) { //we are in a stretch of type fmtType; so close it
			formats.Pop()
			cp("</"+tag+">", inc) //close the stretch and skip inc bytes in the source string
			return
		}
		//we are in a new strtech
		formats.Push(fmtType) //keep track of it
		cp("<"+tag+">", inc)  //open the stretch and skip inc bytes in the source string
	}
	//outer:
	for ; i < mdLen; i++ { //parsing as bytes because all the characters we are looking for are bytes in utf8
		switch md[i] {
		case '&':
			cp("&amp;", 0)
		case '"':
			cp("&#34;", 0)
		case '\'': //apostrophe
			cp("&#39;", 0)
		case '=':
			switch next() {
			case '>':
				cp("&ge;", 1) //increase i by 1 because we consumed the next byte
			case '<':
				cp("&le;", 1) //ditto
			default:
				cpCurrent() // copy "=" only
			}
		case '<':
			switch next() {
			case '=':
				cp("&le;", 1)
			default:
				cp("&lt;", 0)
			}
		case '>':
			switch next() {
			case '=':
				cp("&ge;", 1)
			default:
				cp("&gt;", 0)
			}
			// Markdown provides backslash escapes for the following characters:
			// \   backslash
			// `   backtick
			// *   asterisk
			// _   underscore
			// {}  curly braces
			// []  square brackets
			// ()  parentheses
			// #   hash mark
			// +   plus sign
			// -   minus sign (hyphen)
			// .   dot
			// !   exclamation mark
		case '\\': //backslash; TODO: handle escaping
			switch c := next(); c {
			case '\\', '^', '~', '*', '_':
				i++ //skip "/"
				cpCurrent()
			default: //following Dingus example
				cpCurrent()
			}
		case '^':
			emit(fmtSuperscript, 0, "sup")
		case '~':
			emit(fmtSubscript, 0, "sub")
		case '*': //italic or bold or both
			switch next() {
			case '*': //bold
				emit(fmtBold, 1, "strong")
			default: //italic
				emit(fmtItalic, 0, "em")
			}
		case '_': //italic or bold or both
			switch next() {
			case '_': //bold
				emit(fmtBold, 1, "strong")
			default: //italic
				emit(fmtItalic, 0, "em")
			}
		default:
			cpCurrent()
		}
	}
	var err error
	if !formats.Empty() {
		err = fmt.Errorf("improper nesting of markdown tags, %d tags remain unclosed", len(formats))
	}
	return html[:j], err
}
