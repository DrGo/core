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
	"strconv"
	"strings"
)
// NameFromParts using an a fileName it returns its full path overriding its components
// as appropriate by dir, baseName and ext if provided; if dir, basename and ext
// are all empty, it returns f.Name(). if dir="." it returns os.GetWd() + f.Name() 
func NameFromParts(fileName, dir, baseName, ext string) (string, error) {
	var err error
	if baseName == "" {
		baseName = BaseNameNoExt(fileName)
	}
	if ext == "" {
		ext = filepath.Ext(fileName) 
	}
	// if dir is empty, assume current dir
	switch dir {
	case ".":
		if dir, err = os.Getwd(); err != nil {
			return "", err
		}
	case "":
		dir = filepath.Dir(fileName) 
	}
	//try and come up with a valid name
	return filepath.Join(dir, baseName+ext), nil
}
	
//CloseAndRename closes an os.File and save it to newFileName overwriting if it exists if overWrite is true.
//Useful for closing and renaming a temp file to a permanent path
//FIXME: replace with atomic package because this one may not work under Windows
//FIXME: add err option and rename only if err==nil
func CloseAndRename(f *os.File, newFileName string, overWrite bool) error {
	if err := f.Close(); err != nil {
		return err
	}
	if f.Name() == newFileName {
		return nil
	}
	//if the newFileName already exists and overWrite is false, return an error
	if !overWrite {
		if _, err := os.Stat(newFileName); err == nil {
			return fmt.Errorf("file already exists: %s", newFileName)
		}
	}
	//otherwise, return the result of attempting to rename it
	// fmt.Printf("in closeandrename(): saving %s as\n %s\n", f.Name(), newFileName) /*DEBUG*/
	return os.Rename(f.Name(), newFileName)
}

// CloseAndName close an os.File and rename it to fileName and returns its full path
// If the file exists, it attempts to find unused filename by appending a number to basename
func CloseAndName(f *os.File, fileName string) (string, error) {
	var err error
	if err = f.Close(); err != nil {
		return "", err
	}
	dir, baseName := filepath.Split(fileName)
  ext := filepath.Ext(baseName)
  baseName = BaseNameNoExt(baseName)
	for i := 1; i < 10001; i++ {
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return fileName, os.Rename(f.Name(), fileName)
		}
		fileName = filepath.Join(dir, baseName+strconv.Itoa(i)+ext)
	}
	return "", fmt.Errorf("failed to rename using a unique file name")
}

// GetTempWriter returns a pointer to os.File writing to a temp file
// func GetTempWriter(fileName string) (*os.File, error) {
// 	// if fileName == "<stdout>" {
// 	// 	return os.Stdout, nil
// 	// }
// 	out, err := ioutil.TempFile("", "rwtmp") //create in the system default temp folder, a file prefixed with rwtmp
// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }

// GetInputReader returns a pointer to os.File reading from filename or stdout if filename==""
func GetInputReader(fileName string) (*os.File, error) {
	if fileName == "" || fileName == "<stdin>" {
		return os.Stdin, nil
	}
	return os.Open(fileName)
}

// WriteBufferToFile atomically writes a byte array into fileName. fileName can be "" in which case
// a temp file is created and its name is returned
// FIXME: optimize
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

// CreateFile creates the named file with mode 0666 (before umask), truncating
// it if it already exists and overWrite is true. If successful, methods on the returned
// File can be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
func CreateFile(fileName string, overWrite bool) (out *os.File, err error) {
	mode := os.O_RDWR | os.O_CREATE | os.O_TRUNC //read-write, create if none exists or truncate existing one
	if !overWrite {
		mode |= os.O_EXCL //file must not exist
	}
	return os.OpenFile(fileName, mode, 0666)
}

//CheckTextStream returns an error if r does contain text otherwise return nil
func CheckTextStream(r io.Reader, streamMinSize int) error {
	first512Bytes := make([]byte, 512)
	n, err := r.Read(first512Bytes)
	switch {
	case err != nil && err != io.EOF:
		return err
	case n < streamMinSize:
		return fmt.Errorf("stream is empty or does not contain sufficient data, size=%d", n)
	case !strings.Contains(http.DetectContentType(first512Bytes[0:n]), "text"):
		return fmt.Errorf("file does not contain text (possibly a binary file)")
	default:
		return nil
	}
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
