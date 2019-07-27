package files

import (
	"io"
	"os"
)

// Cat prints file content to stdout
func Cat(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}
	return nil
}
