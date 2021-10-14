package errors

import "fmt"

var Debug = true

//TODO:  add support for logging errors.D

func Logln(a ...interface{}) bool {
	fmt.Println(a...)
	return true
}

func Logf(format string, a ...interface{}) bool {
	fmt.Printf(format, a...)
	return true
}
