package template

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/matryer/is"
)

// prereqs maps computer science courses to their prerequisites.
// source: gopl book
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func xTestTopoSort(t *testing.T) {
	sorted,_, err := KahnSort(prereqs)
	// is := is.New(t)
	// is.True(err==nil)
	fmt.Printf("length of array: %d, err:%v\n", len(sorted), err)
	reverseSlice(sorted)
	for i, course := range sorted {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
	fmt.Println(results)
}

const results = `
1:	intro to programming
2:	discrete math
3:	data structures
4:	algorithms
5:	linear algebra
6:	calculus
7:	formal languages
8:	computer organization
9:	compilers
10:	databases
11:	operating systems
12:	networks
13:	programming languages
`

func reverseSlice(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

var (
	m = map[string][]string{
		"0": {"1", "4"},
		"1": {"3", "5"},
		"2": {"5"},
		"3": {"5", "7"},
		"4": {},
		"5": {"6"},
		"6": {"7"},
		"7": {}}
	r = [][]string{
		{"2", "0", "4", "1", "3", "5", "6", "7"},
		{"0", "4", "1", "3", "2", "5", "6", "7"}}
)

func xTestKhan(t *testing.T) {
	sorted, _, err := KahnSort(m)
	is := is.New(t)
	is.NoErr(err)
	// fmt.Printf("length of array: %d\n", len(sorted))
	is.True(reflect.DeepEqual(sorted, r[0]) || reflect.DeepEqual(sorted, r[1]))
	// reverseSlice(sorted)
	// for i, course := range sorted {
	// 	fmt.Printf("%d:\t%s\n", i+1, course)
	// }
	// fmt.Println(results)
}

var libs = map[string][]string{
	"des_system_lib": {"std", "synopsys", "std_cell_lib", "dw02", "dw01", "ramlib", "ieee"},
	"dw01":           {"ieee", "dware", "gtech"},
	"dw02":           {"ieee", "dware"},
	"dw03":           {"std", "synopsys", "dware", "dw02", "dw01", "ieee", "gtech"},
	"dw04":           {"ieee", "dw01", "dware", "gtech"},
	"dw05":           {"ieee", "dware"},
	"dw06":           {"ieee", "dware"},
	"dw07":           {"ieee", "dware"},
	"dware":          {"ieee"},
	"gtech":          {"ieee"},
	"ramlib":         {"std", "ieee"},
	"std_cell_lib":   {"ieee"},
	"synopsys":       {},
  "orphan": {},
}

var libres = []string{
	"ieee",
	"std_cell_lib",
	"gtech",
	"dware",
	"dw07",
	"dw06",
	"dw05",
	"dw02",
	"dw01",
	"dw04",
	"std",
	"ramlib",
	"synopsys",
	"dw03",
	"des_system_lib",
}

func xTestKhanLibs(t *testing.T) {
	sorted, _, err := KahnSort(libs)
	is := is.New(t)
	is.NoErr(err)
	// is.True(reflect.DeepEqual(sorted, libres)) 
	for i, course := range sorted {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
	fmt.Println(libres)
}
