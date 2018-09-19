package errors

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// ErrorList is a list of errors; a more generic and concurrent-safe version of Go's parser ErrorList
// The zero value for an ErrorList is an empty ErrorList ready to use.
type ErrorList struct {
	list []error
	sync.RWMutex
}

//NewErrorList returns an empty ErrorList
func NewErrorList() *ErrorList {
	return &ErrorList{}
}

// Add adds an error to an ErrorList. NOT Safe to use concurrently
func (el *ErrorList) Add(e error) {
	el.Lock()
	defer el.Unlock()
	el.list = append(el.list, e)
}

// Get returns the error at index or nil if index is out of bounds of list
func (el *ErrorList) Get(index int) error {
	el.RLock()
	defer el.RUnlock()
	if index < 0 || index >= len(el.list) {
		return nil
	}
	return el.list[index]
}

// // Count returns the number of errors in the List; preferred to len(errorList)
// func (el *ErrorList) Count() int {
// 	el.RLock()
// 	defer el.RUnlock()
// 	return len(el.list)
// }

// Reset resets an ErrorList to no errors.
func (el *ErrorList) Reset() {
	el.Lock()
	defer el.Unlock()
	el.list = (el.list)[0:0]
}

// ErrorList implements the sort Interface; override Less for any special sorting needs
func (el *ErrorList) Len() int {
	el.RLock()
	defer el.RUnlock()
	return len(el.list)
}
func (el *ErrorList) Swap(i, j int) {
	el.Lock()
	defer el.Unlock()
	el.list[i], el.list[j] = el.list[j], el.list[i]
}

func (el *ErrorList) Less(i, j int) bool {
	el.RLock()
	defer el.RUnlock()
	return el.list[i].Error() < el.list[j].Error()
}

// Sort sorts an ErrorList. The default sorting is by error message (asc)
func (el *ErrorList) Sort() {
	el.Lock()
	defer el.Unlock()
	sort.Sort(el)
}

// RemoveMultiples sorts an ErrorList and removes all but the first error.
func (el *ErrorList) RemoveMultiples() {
	el.Lock()
	defer el.Unlock()
	el.Sort()
	var last string
	i := 0
	for _, e := range el.list {
		if e.Error() != last {
			last = e.Error()
			(el.list)[i] = e
			i++
		}
	}
	(el.list) = (el.list)[0:i]
}

// An ErrorList implements the error interface.
func (el *ErrorList) Error() string {
	el.RLock()
	defer el.RUnlock()
	switch len(el.list) {
	case 0:
		return "no errors"
	case 1:
		return el.list[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", el.list[0], len(el.list)-1)
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (el *ErrorList) Err() error {
	if len(el.list) == 0 {
		return nil
	}
	return el
}

// Print prints a list of errors to w,
func (el *ErrorList) Print(w io.Writer) {
	el.RLock()
	defer el.RUnlock()
	for _, e := range el.list {
		fmt.Fprintf(w, "%s\n", e)
	}
}

// PrintError is a utility function that prints a list of errors to w,
// one error per line, if the err parameter is an ErrorList. Otherwise
// it prints the err string.
func PrintError(w io.Writer, err error) {
	if list, ok := err.(*ErrorList); ok {
		list.Print(w)
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}

//ErrorsToError converts list of error into one error string
func ErrorsToError(err error) error {
	switch e := err.(type) {
	case *ErrorList:
		var w strings.Builder
		PrintError(&w, e)
		return fmt.Errorf("%s", w.String())
	default:
		return e
	}
}

type errFunc func() error

func checkErrors(warpIn string, efs ...errFunc) error {
	for _, ef := range efs {
		if err := ef(); err != nil {
			return fmt.Errorf(warpIn, err)
		}
	}
	return nil
}

// // //usage
// // return checkError(
// 	os.open(path),
// 	func(){return os.Stat},
// 	os.close()
// // )
