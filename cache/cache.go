// A simple non-concurrent file cache
package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileCache struct {
	cache map[string]string
}

func NewFileCache(dirName string) *FileCache {
	fc := &FileCache{
		cache: make(map[string]string),
	}
	if dirName != "" {
		if err := fc.Populate(dirName); err != nil {
			panic("NewFileCache failed: " + err.Error())
		}
	}
	return fc
}

func (fc *FileCache) Populate(dirName string) error {
	visit := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fc.cache[path] = string(contents)
		fmt.Println(path)
		return nil
	}
	err := filepath.Walk(dirName, visit)
	return err
}

func (fc *FileCache) Get(fileName string) (contents string, ok bool) {
	contents, ok = fc.cache[fileName]
	return
}

func (fc *FileCache) Len() int {
	return len(fc.cache)
}
