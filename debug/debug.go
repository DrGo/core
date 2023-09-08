package debug

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	//DebugSilent print errors only
	DebugSilent int = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)

var debugString = []string{"Silent", "Warning", "Update", "All"}

// String returns a textual representation of a debug constant
func String(debug int) string {
	if debug >= DebugSilent && debug <= DebugAll {
		return debugString[debug]
	}
	return "Invalid"
}

// UI handles interactions with the user
type UI interface {
	Log(a ...interface{})
	Warn(a ...interface{})
}

type ui struct {
	Depth int
	Debug int
}

// NewUI initializes a new default UI
func NewUI(debug int) UI {
	return &ui{
		Debug: debug,
	}
}

func (u ui) String() string {
	return String(u.Debug)
}

func (u ui) Log(a ...interface{}) {
	if u.Debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", u.Depth))
		fmt.Println(a...)
	}
}

// TODO: optimize by removing unnecessary calls to strings.Repeat
func (u ui) Warn(a ...interface{}) {
	if u.Debug >= DebugWarning {
		fmt.Printf(strings.Repeat("  ", u.Depth) + "warning: ")
		fmt.Println(a...)
	}
}

//Source: Mcvaden sh 
// DebugPrint prints the provided syntax tree, spanning multiple lines and with
// indentation. Can be useful to investigate the content of a syntax tree.
func DebugPrint(w io.Writer, node any) error {
	p := debugPrinter{out: w}
	p.print(reflect.ValueOf(node))
	return p.err
}

type debugPrinter struct {
	out   io.Writer
	level int
	err   error
}

func (p *debugPrinter) printf(format string, args ...any) {
	_, err := fmt.Fprintf(p.out, format, args...)
	if err != nil && p.err == nil {
		p.err = err
	}
}

func (p *debugPrinter) newline() {
	p.printf("\n")
	for i := 0; i < p.level; i++ {
		p.printf(".  ")
	}
}

func (p *debugPrinter) print(x reflect.Value) {
	switch x.Kind() {
	case reflect.Interface:
		if x.IsNil() {
			p.printf("nil")
			return
		}
		p.print(x.Elem())
	case reflect.Ptr:
		if x.IsNil() {
			p.printf("nil")
			return
		}
		p.printf("*")
		p.print(x.Elem())
	case reflect.Slice:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.level++
			p.newline()
			for i := 0; i < x.Len(); i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				if i == x.Len()-1 {
					p.level--
				}
				p.newline()
			}
		}
		p.printf("}")

	case reflect.Struct:
		// if v, ok := x.Interface().(Pos); ok {
		// 	p.printf("%v:%v", v.Line(), v.Col())
		// 	return
		// }
		t := x.Type()
		p.printf("%s {", t)
		p.level++
		p.newline()
		for i := 0; i < t.NumField(); i++ {
			p.printf("%s: ", t.Field(i).Name)
			p.print(x.Field(i))
			if i == x.NumField()-1 {
				p.level--
			}
			p.newline()
		}
		p.printf("}")
	default:
		p.printf("%#v", x.Interface())
	}
}
