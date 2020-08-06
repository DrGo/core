package md

import (
	"fmt"
	"strings"

	"github.com/drgo/core/stacks"
)

//InlinedMdToHTMLOptions options for Markdown to HTML conversion
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
func InlinedMdToHTML(md string, opts *InlinedMdToHTMLOptions) (string, error) {
	const null = '\x00'
	formats := make(stacks.IntStack, 0, 100)
	mdLen := len(md)
	if mdLen == 0 {
		return "", nil
	}
	var html strings.Builder
	html.Grow(mdLen * 3)
	i := 0                          // is the source (md) byte index
	cp := func(s string, inc int) { //copy string to the html slice
		html.WriteString(s)
		i += inc //increase i by inc; could be zero
	}
	next := func() byte {
		if i+1 < mdLen {
			return md[i+1]
		}
		return null
	}
	cpCurrent := func() {
		html.WriteByte(md[i])
	}

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
	for ; i < mdLen; i++ { //parsing as bytes
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
		case '\\': //backslash;
			switch c := next(); c {
			case '\\', '^', '~', '*', '_':
				i++ //skip "/"
				cpCurrent()
			default:
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
	return html.String(), err
}

// EscapeAsTex returns a string representation of s that is suitable for Latex output
func EscapeAsTex(s string) string {
	var sb strings.Builder
	var es string
	written := 0
	// The byte loop below assumes that all URLs use UTF-8 as the
	// content-encoding.
	for i, n := 0, len(s); i < n; i++ {
		c := s[i]
		switch c {
		case '%', '#', '$', '&', '_', '{', '}':
			es = "\\" + string(c)
		case '\\':
			es = `\textbackslash{}`
		case '^':
			es = `\^{}` // `\textasciicircum{}`
		case '~':
			es = `\~{}` // `\textasciitilde{}`
		default: // do not process any other char
			continue
		}
		if written == 0 {
			sb.Grow(len(s) + 16)
		}
		sb.WriteString(s[written:i] + es)
		written = i + 1
	}
	sb.WriteString(s[written:])
	if written != 0 {
		return sb.String()
	}
	return s
}
