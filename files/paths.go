package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//GetExeDir returns exe dir
// may fail if exe invoked through symlink see Executable() docs
func GetExeDir() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

func BaseNameNoExt(path string) string {
	s := filepath.Base(path)
	dot := strings.LastIndexByte(s, '.')
	if dot >= 0 {
		return s[:dot]
	}
	return s
}

//ConstructFileName if called with "dir/filename", ".newext", "prefix" , "postfix". it returns
//os-appropriate name "dir/prefixfilenamepostfix.newext"
func ConstructFileName(path, newExt, prefix, postfix string) string {
	//TODO: validate
	dir := filepath.Dir(path)
	base := BaseNameNoExt(path)
	base = prefix + base + postfix + newExt
	return filepath.Join(dir, base)
}

//GetFullPath returns fileName if it has a full path. If not, it returns a filename relative to the current working dir
func GetFullPath(fileName string) (string, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return fileName, nil
	}
	//FIXME: replace with Abs()
	//// Abs returns an absolute representation of path.
	// If the path is not absolute it will be joined with the current
	// working directory to turn it into an absolute path. The absolute
	// path name for a given file is not guaranteed to be unique.
	if filepath.IsAbs(fileName) || filepath.Dir(fileName) != "." { // absolute path or dir provided
		return fileName, nil
	}
	//add support to dir argument to be used if not empty before using os.Getwd()
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, fileName), nil
}

//GetOutputDir returns dirname or current working directory if preserveTempFiles is true
// otherwise creates a temp dir as per os defaults and return its name
func GetOutputDir(dirName string, preserveTempFiles bool) (dir string, temp bool, err error) {
	dir = strings.TrimSpace(dirName)
	if preserveTempFiles { //save in workDir or currentDir if workDir is null
		if dir != "" {
			return dir, false, nil
		}
		if dir, err = os.Getwd(); err != nil {
			return "", false, err
		}
		return dir, false, nil
	}
	// if preserveTempFiles is false, create temp folder and return
	if dir, err = ioutil.TempDir("", "rw-temp101"); err != nil { //create temp dir in the os default temp dir
		return "", false, err
	}
	return dir, true, nil
}

// ReplaceFileExt returns fileName with extension replace with newExt (both ext and .ext are ok)
func ReplaceFileExt(fileName, newExt string) string {
	ext := path.Ext(fileName)
	newExt = strings.Trim(newExt, ".")
	return fileName[0:len(fileName)-len(ext)] + "." + newExt
}

// RemoveIfExists removes a file, returning no error if it does not exist.
func RemoveIfExists(fileName string) error {
	err := os.Remove(fileName)
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return err
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

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// func FileCompare(file1, file2 string) (error, bool) {
// 	const chunckSize = 64 * 1024
// 	f1s, err := os.Stat(file1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	f2s, err := os.Stat(file2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if f1s.Size() != f2s.Size() {
// 		return nil, false
// 	}

// 	f1, err := os.Open(file1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	f2, err := os.Open(file2)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for {
// 		b1 := make([]byte, chunckSize)
// 		_, err1 := f1.Read(b1)

// 		b2 := make([]byte, chunckSize)
// 		_, err2 := f2.Read(b2)

// 		if err1 != nil || err2 != nil {
// 			if err1 == io.EOF && err2 == io.EOF {
// 				return nil, true
// 			} else if err1 == io.EOF && err2 == io.EOF {
// 				return nil, false
// 			} else {
// 				log.Fatal(err1, err2)
// 			}
// 		}

// 		if !bytes.Equal(b1, b2) {
// 			return nil, false
// 		}
// 	}
// }
