package files

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type FileDescriptor struct {
	Name    string
	Version string `json:",omitempty"`
}

func NewFileDescriptor(fileName string) *FileDescriptor {
	return &FileDescriptor{
		Name:    fileName,
		Version: "",
	}
}

func FileDescriptorsToStrings(fds []*FileDescriptor) (ss []string) {
	for _, f := range fds {
		ss = append(ss, f.Name)
	}
	return ss
}

//CloseAndRename closes an os.File and save it to newFileName overwriting if it exists if overWrite is true.
//Useful for closing and renaming a temp file to a permanent path
//FIXME: replace with atomic package because this one may not work under Windows
func CloseAndRename(f *os.File, newFileName string, overWrite bool) error {
	if err := f.Close(); err != nil {
		return err
	}
	//if the newFileName already exists and overWrite is false, return an error
	if !overWrite {
		if _, err := os.Stat(newFileName); err == nil {
			return fmt.Errorf("file already exists: %s", newFileName)
		}
	}
	//otherwise, return the result of attempting to rename it
	//	fmt.Printf("in closeandrename(): saving %s as\n %s\n", f.Name(), newFileName)
	return os.Rename(f.Name(), newFileName)
}

func DetectContentType(r io.ReadSeeker) (string, error) {
	// Only the first 512 bytes are needed
	buffer := make([]byte, 512)
	n, err := r.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	// Reset the read pointer if necessary.
	r.Seek(0, 0)
	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	return http.DetectContentType(buffer[:n]), nil
}

func PrintFileStat(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	fmt.Printf("file name=%s\n", path)
	fmt.Printf("file size=%d\n", fi.Size())
	return nil
}

func PrintFileContent(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, os.Stdout)
	return err
}

// GetTempWriter returns a pointer to os.File writing to a temp file or stdout if filename==""
func GetTempWriter(fileName string) (*os.File, error) {
	if fileName == "<stdout>" {
		return os.Stdout, nil
	}
	out, err := ioutil.TempFile("", "rwtmp") //create in the system default temp folder, a file prefixed with rwtmp
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetInputReader returns a pointer to os.File reading from filename or stdout if filename==""
func GetInputReader(fileName string) (*os.File, error) {
	if fileName == "" || fileName == "<stdin>" {
		return os.Stdin, nil
	}
	return os.Open(fileName)
}

// WriteBufferToFile atomically writes a byte array into fileName. fileName can be "" in which case
// a temp file is created and its name is returned
func WriteBufferToFile(fileName string, buf []byte, overWrite bool) (usedFileName string, err error) {
	file, err := ioutil.TempFile("", "rwtmp") //created in the system default temp folder
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			_ = file.Close()  // attempt to close the file ignoring any errors &
			usedFileName = "" // return original error and empty usedFileName
		} else {
			usedFileName = fileName
			if usedFileName == "" {
				usedFileName = file.Name()
			}
			err = CloseAndRename(file, usedFileName, overWrite)
		}
	}()

	w := bufio.NewWriter(file)
	_, err = w.Write(buf)
	if err != nil {
		return "", err
	}
	return "", w.Flush() // return vars will be overwritten by the deferred func
}

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package renameio writes files atomically by renaming temporary files.

const patternSuffix = "*.tmp"

// Pattern returns a glob pattern that matches the unrenamed temporary files
// created when writing to filename.
func Pattern(filename string) string {
	return filepath.Join(filepath.Dir(filename), filepath.Base(filename)+patternSuffix)
}

// WriteFile is like ioutil.WriteFile, but first writes data to an arbitrary
// file in the same directory as filename, then renames it atomically to the
// final name.
//
// That ensures that the final location, if it exists, is always a complete file.
func WriteFile(filename string, data []byte) (err error) {
	return WriteToFile(filename, bytes.NewReader(data))
}

// WriteToFile is a variant of WriteFile that accepts the data as an io.Reader
// instead of a slice.
func WriteToFile(filename string, data io.Reader) (err error) {
	f, err := ioutil.TempFile(filepath.Dir(filename), filepath.Base(filename)+patternSuffix)
	if err != nil {
		return err
	}
	defer func() {
		// Only call os.Remove on f.Name() if we failed to rename it: otherwise,
		// some other process may have created a new file with the same name after
		// that.
		if err != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}()

	if _, err := io.Copy(f, data); err != nil {
		return err
	}
	// Sync the file before renaming it: otherwise, after a crash the reader may
	// observe a 0-length file instead of the actual contents.
	// See https://golang.org/issue/22397#issuecomment-380831736.
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), filename)
}
