package files

import (
	"io"
	"io/ioutil"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GlobFS takes a file system and one or more globing pattern and
// return matched files [could be 0] or nil and error 
func GlobFS(fsys fs.FS, patterns ...string) ([]string, error) {
	var filenames []string
	for _, pattern := range patterns {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return nil, err
		}
		filenames = append(filenames, list...)
	}
	return Unique(filenames), nil
}

// FIXME: use binary search?! See https://github.com/mpvl/unique/blob/cbe035fff7de56b8185768b119ee94a9e42dd938/unique.go#L61
func Unique(e []string) []string {
    r := []string{}

    for _, s := range e {
        if !contains(r[:], s) {
            r = append(r, s)
        }
    }
    return r
}

func contains(e []string, c string) bool {
    for _, s := range e {
        if s == c {
            return true
        }
    }
    return false
}

func ParseWindowsFileArguments(args string) ([]string, error) {
	if !strings.ContainsAny(args, "*?") {
		return []string{args}, nil
	}
	return filepath.Glob(args)
}

//GetFileList get list of files from a pattern
func GetFileList(pattern string) (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	return
}

/*
Linux rules
/ 		=> means start at root of filesystem (absolute reference).
./ 		=> current directory
../ 	=> means go up one directory from the current directory then proceed.
../../ 	=> means go up two directories then proceed.
~/ 		=> start from home directory
None of the above: path relative to current directory
*/

//returns filePath if absolute, otherwise it constructs one from base and filepath.
//no error checking here.
func GetUsableFilePath(filePath string, base string) string {
	if strings.TrimSpace(filePath) == "" {
		return filePath
	}
	// if this is a glob, return it with base if one is specified
	if IsFileGlobPattern(filePath) {
		return filepath.Join(base, filePath)
	}
	// if this is a valid filename on its own, use it regardless of the base
	if IsValidFileName(filePath) {
		return filePath
	}
	// if not, perhaps adding the base (if any) would help
	return filepath.Join(base, filePath)
}

func IsValidFileGlob(pattern string) bool {
	// The only possible returned error from Match is ErrBadPattern, when pattern is malformed.
	if _, err := filepath.Match(pattern, "dummyfilename"); err == nil {
		return true
	}
	return false
}

// IsFileGlobPattern determines if "pattern" has 1 or more characters [', ']', '*' or '//' used for globbing in Go.
// hasMeta reports whether path contains any of the magic characters
// recognized by path.Match.
func IsFileGlobPattern(pattern string) bool {
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		if c == '*' || c == '?' || c == '[' || runtime.GOOS == "windows" && c == '\\' {
			return true
		}
	}
	return false
}

//may not always work!!
func IsValidFileName(fileName string) bool {
	// Check if file already exists
	if _, err := os.Stat(fileName); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fileName, d, 0644); err == nil {
		os.Remove(fileName) // And delete it
		return true
	}

	return false
}

func NewOutWriter(outFileName string) (wc io.WriteCloser, err error) {
	if outFileName == "" {
		return os.Stdout, nil
	}
	wc, err = os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		wc.Close() // don't care about the error here
		return nil, err
	}
	return wc, nil
}


