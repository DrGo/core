package tu

import (
	"fmt"
	"reflect"
	"testing"

	"slices"
)

type Options int

const (
	FailNow Options = iota
	Panic
)

var Debug bool

func Equal[T comparable](t *testing.T, actual, expected T, options ...Options) {
	t.Helper() //report error in the file that calls this func
	if expected != actual {
		const errmsg = "wanted:%v\ngot:%v\n"
		Assert(!slices.Contains(options, Panic), fmt.Sprintf(errmsg, expected, actual))
		t.Errorf(errmsg, expected, actual)
		if slices.Contains(options, FailNow) {
			t.FailNow()
		}
	}
}

// isNil gets whether the object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func NotNil(t *testing.T, obj any, options ...Options) {
	t.Helper()
	if isNil(obj) {
		t.Errorf("%v is nil", obj)
	}
}

func Assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func Assertf(cond bool, format string, v ...any) {
	if cond {
		return
	}
	panic(fmt.Sprintf(format, v...))
}

func P(format string, a ...any) {
	if Debug {
		fmt.Printf(format, a...)
	}
}
func PL(a ...any) {
	if Debug {
		fmt.Println(a...)
	}
}
