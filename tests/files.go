// Package tests has testing helpers
package tests

import (
	"io/ioutil"
	"os"
	"testing"
)

// MkTempDir makes a temporary directory
func MkTempDir(t *testing.T)(string, func()) {
	dir, err := ioutil.TempDir("", "sm-tests-")
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir, func(){
     os.RemoveAll(dir)
	}
}

// MkTempFile create a temporary file and returns its name
func MkTempFile(t *testing.T, dir string) (*os.File, func()) {
	f, err := ioutil.TempFile(dir, "sm-tests-")
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return f,func(){
	  f.Close()
	  os.Remove(f.Name())
	}
}


//go test sets the package dir as pwd
