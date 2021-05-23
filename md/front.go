package md

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/drgo/core"
	"github.com/drgo/core/errors"
)

type FrontMatter map[string]string

// ContentFile represents a HUGO-compatible md file
type ContentFile struct {
	Path        string
	FrontMatter FrontMatter
	Summary     []byte
	Body        []byte
}

// Prop returns the value of the specified front matter entry
func (p *ContentFile) Prop(name string) string {
	if val, ok := p.FrontMatter[name]; ok {
		return val
	}
	return "" 
}

// ParseContentFilesGlob parses markdown-with-front-matter files identified by the pattern
// The files are matched according to the semantics of filepath.Match, and the pattern
// must match at least one file. It is equivalent to calling t.ParseContentFiles with the
// list of files matched by the pattern.
func ParseContentFilesGlob(fsys fs.FS, pattern string) ([]*ContentFile, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, errors.Errorf("pattern matches no files: %#q", pattern)
	}
	return ParseContentFiles(fsys, filenames...)
}

// ParseContentFiles pareses markdown-with-front-matter files into a slice of ContentFile pointers
func ParseContentFiles(fsys fs.FS, filenames ...string) ([]*ContentFile, error) {
	if len(filenames) == 0 {
		return nil, errors.Errorf("no files named in call to ParseContentFiles")
	}
	var lst []*ContentFile
	for _, filename := range filenames {
		c, err := ParseContentFile(fsys, filename)
		if err != nil {
			return nil, err
		}
		lst = append(lst, c)
	}
	return lst, nil
}

// ParseContentFile pareses markdown-with-front-matter files into a ContentFile pointer
func ParseContentFile(fsys fs.FS, path string) (*ContentFile, error) {
	r, err := fsys.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse md file %s", path)
	}
	defer r.Close()
	cf, err := ParseFrontMatter(r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse md file %s", path)
	}
	cf.Path= path
  // fmt.Println("summary:", string(cf.Summary))
	return cf, nil
}

var Separator = []byte("+++")
var ColSeparator = []byte("=")
var MoreMarker = []byte("<!--more-->")

func ParseFrontMatter(input io.Reader) (*ContentFile, error) {
	const errmsg = "failed to parse front matter"
	s := bufio.NewScanner(input)
	// Necessary so we can handle larger than default 4096b buffer
	// bufsize := 1024 * 1024
	// buf := make([]byte, bufsize)
	// s.Buffer(buf, bufsize)
	//the first line should start by a separator
	var (
		fMap    map[string]string
		body    []byte
		summary []byte
	)
	s.Scan()
	if err := s.Err(); err != nil {
		return nil, errors.Wrap(err, errmsg)
	}
	fMap = make(map[string]string)
	if line := s.Bytes(); !bytes.HasPrefix(line, Separator) {
		return nil, errors.Errorf("%s:file does not start with valid separator [%s]", errmsg, string(line))
	}
	// parse the front matter entries in this form: key=value
	for s.Scan() {
		line := s.Bytes()
		if bytes.HasPrefix(line, Separator) { //found second separator, the body starts
			break
		}
		entry := bytes.SplitN(line, ColSeparator,2)  //split on the first colseparator
		if len(entry) != 2 {
			return nil, errors.Errorf("%s:%s is not front matter entry", errmsg, string(line))
		}
		cleaned := strings.TrimSpace(string(entry[1]))
		val, err := strconv.Unquote(cleaned)
		if err != nil { //not a string
			// errors.Logln(string(line), cleaned, err.Error())
			val = cleaned
		}
		fMap[strings.TrimSpace(string(entry[0]))] = val
	}
	// read the rest of the file into body
	for s.Scan() {
		line := s.Bytes()
		if bytes.HasPrefix(line, MoreMarker) { //found more...
			summary = append(summary, body...) // move parsed lines to summary
			body = body[:0]                    //clear the body slice
			continue                           //ignore the more... line
		}
		body = append(body, s.Bytes()...)
		body = append(body, []byte(core.LineBreak)...)
		
	}
	return &ContentFile{
		FrontMatter: fMap,
		Body:        body,
		Summary:     summary,
	}, s.Err()
}

func UnmarshalStringTo(data string, typ interface{}) (err error) {
	data = strings.TrimSpace(data)
	// We only check for the possible types in YAML, JSON and TOML.
	switch typ.(type) {
	case string:
		typ = data
	case bool:
		typ, err = strconv.ParseBool(data)
	case int, int32, int64:
		typ, err = strconv.ParseInt(data, 0, 0)
	case uint, uint32, uint64:
		typ, err = strconv.ParseUint(data, 0, 0)
	case float32:
		typ, err = strconv.ParseFloat(data, 32)
	case float64:
		typ, err = strconv.ParseFloat(data, 64)
	default:
		err = errors.Errorf("unmarshal: %T not supported", typ)
	}
	return err
}

// This is a subset of the formats allowed by the regular expression
// defined at http://yaml.org/type/timestamp.html.
var allowedTimestampFormats = []string{
	"2006-1-2T15:4:5.999999999Z07:00", // RCF3339Nano with short date fields.
	"2006-1-2t15:4:5.999999999Z07:00", // RFC3339Nano with short date fields and lower-case "t".
	"2006-1-2 15:4:5.999999999",       // space separated with no time zone
	"2006-1-2",                        // date only
	// Notable exception: time.Parse cannot handle: "2001-12-14 21:59:43.10 -5"
	// from the set of examples.
}

// parseTimestamp parses s as a timestamp string and
// returns the timestamp and reports whether it succeeded.
// Timestamp formats are defined at http://yaml.org/type/timestamp.html
func parseTimestamp(s string) (time.Time, bool) {
	// TODO write code to check all the formats supported by
	// http://yaml.org/type/timestamp.html instead of using time.Parse.

	// Quick check: all date formats start with YYYY-.
	i := 0
	for ; i < len(s); i++ {
		if c := s[i]; c < '0' || c > '9' {
			break
		}
	}
	if i != 4 || i == len(s) || s[i] != '-' {
		return time.Time{}, false
	}
	for _, format := range allowedTimestampFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func mapToStruct(m map[string]interface{}, res interface{}) error {
	jsonbody, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonbody, &res)
}
