// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package core

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// import(
// 	"testing"
// 	"io/ioutil"
// 	"path/filepath"
// )
// func helperLoadBytes(t *testing.T, name string) []byte {
// 	path := filepath.Join("testdata", name) // relative path
// 	bytes, err := ioutil.ReadFile(path)
// 	if err != nil {
// 	  t.Fatal(err)
// 	}
// 	return bytes
//   }

//   //var update = flag.Bool("update", false, "update .golden files")
//   func TestSomething(t *testing.T) {
// 	actual := doSomething()
// 	golden := filepath.Join(“testdata”, tc.Name+”.golden”)
// 	if *update {
// 	  ioutil.WriteFile(golden, actual, 0644)
// 	}
// 	expected, _ := ioutil.ReadFile(golden)

// 	if !bytes.Equal(actual, expected) {
// 	  // FAIL!
// 	}
//   }

//   var testHasGit bool
// func init() {
//   if _, err := exec.LookPath("git"); err == nil {
//     testHasGit = true
//   }
// }
// func TestGitGetter(t *testing.T) {
//   if !testHasGit {
//     t.Log("git not found, skipping")
//     t.Skip()
//   }
//   // ...
// }

// func TestFailingGit(t *testing.T) {
// 	if os.Getenv("BE_CRASHING_GIT") == "1" {
// 	  CrashingGit()
// 	  return
// 	}
// 	cmd := exec.Command(os.Args[0], "-test.run=TestFailingGit")
// 	cmd.Env = append(os.Environ(), "BE_CRASHING_GIT=1")
// 	err := cmd.Run()
// 	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
// 	  return
// 	}
// 	t.Fatalf("Process ran with err %v, want os.Exit(1)", err)
//   }

//   alias gtest="go test \$(go list ./… | grep -v /vendor/)
// -tags=integration"
