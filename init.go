package core

import (
	"reflect"
	"runtime"
)

//LineBreak OS-specific line break string
var LineBreak = "\n"

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}

//IsNil checks that an interface holds a non-nil concrete value
func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
