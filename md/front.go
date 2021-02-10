package md

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/drgo/core/errors"
)

type FrontMatter map[string]string

// ContentFile represents a HUGO-compatible md file
type ContentFile struct {
	Path        string
	FrontMatter FrontMatter
	Body        string
}

// Prop returns the value of the specified front matter entry
func (p *ContentFile) Prop(name string) string {
	if val, ok := p.FrontMatter[name]; ok {
		return val
	}
	return "!!! Error= no such property: " + name
}

// ParseContentFilesGlob parses markdown-with-front-matter files identified by the pattern
// The files are matched according to the semantics of filepath.Match, and the pattern
// must match at least one file. It is equivalent to calling t.ParseContentFiles with the
// list of files matched by the pattern.
func ParseContentFilesGlob(pattern string) ([]*ContentFile, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, errors.Errorf("pattern matches no files: %#q", pattern)
	}
	return ParseContentFiles(filenames...)
}

// ParseContentFiles pareses markdown-with-front-matter files into a slice of ContentFile pointers
func ParseContentFiles(filenames ...string) ([]*ContentFile, error) {
	if len(filenames) == 0 {
		return nil, errors.Errorf("no files named in call to ParseContentFiles")
	}
	var lst []*ContentFile
	for _, filename := range filenames {
		c, err := ParseContentFile(filename)
		if err != nil {
			return nil, err
		}
		lst = append(lst, c)
	}
	return lst, nil
}

// ParseContentFile pareses markdown-with-front-matter files into a ContentFile pointer
func ParseContentFile(path string) (*ContentFile, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse md file %s", path)
	}
	defer r.Close()
	fm, body, err := ParseFrontMatter(r)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse md file %s", path)
	}
	return &ContentFile{path, fm, string(body)}, err
}

var Separator = []byte("+++")
var ColSeparator = []byte("=")

func ParseFrontMatter(input io.Reader) (fMap map[string]string, body []byte, err error) {
	const errmsg = "failed to parse front matter"
	s := bufio.NewScanner(input)
	// Necessary so we can handle larger than default 4096b buffer
	// bufsize := 1024 * 1024
	// buf := make([]byte, bufsize)
	// s.Buffer(buf, bufsize)
	//the first line should start by a separator
	s.Scan()
	if err := s.Err(); err != nil {
		return nil, nil, errors.Wrap(err, errmsg)
	}
	fMap = make(map[string]string)
	if line := s.Bytes(); !bytes.HasPrefix(line, Separator) {
		return nil, nil, errors.Errorf("%s:%s is not valid separator", errmsg, string(line))
	}
	// parse the front matter entries in this form: key=value
	for s.Scan() {
		line := s.Bytes()
		if bytes.HasPrefix(line, Separator) { //found second separator
			break
		}
		entry := bytes.Split(line, ColSeparator)
		if len(entry) != 2 {
			return nil, nil, errors.Errorf("%s:%s is not front matter entry", errmsg, string(line))
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
		body = append(body, s.Bytes()...)
	}
	return fMap, body, s.Err()
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
