package stacks

import "strings"

type StringStack struct {
	data []string
}

func (ss *StringStack) Push(s string) {
	ss.data = append(ss.data, s)
}
func (ss *StringStack) Pop() string { //no error handling
	s := ss.data[len(ss.data)-1]
	ss.data = ss.data[:len(ss.data)-1]
	return s
}

func (ss *StringStack) String() string {
	return strings.Join(ss.data, "|")
}

func (ss *StringStack) GetAllStrings() []string {
	return ss.data
}

// containsAll reports whether StringStack contains the elements of list, in order.
func (ss *StringStack) ContainsAll(list []string) bool {
	x := ss.data //copy to avoid changing
	for len(list) <= len(x) {
		if len(list) == 0 {
			return true
		}
		if x[0] == list[0] {
			list = list[1:]
		}
		x = x[1:]
	}
	return false
}
