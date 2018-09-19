package files

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func WriteBufferToFile(fileName string, buf []byte, overWrite bool) (err error) {
	file, err := ioutil.TempFile("", "rwtmp") //create in the system default temp folder, a file prefixed with
	defer func() {
		err = CloseAndRename(file, fileName, overWrite)
	}()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return w.Flush()
}
