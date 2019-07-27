package core

import (
	"fmt"
	"io"

	"github.com/drgo/mdson"
)

const defaultDebugLevel = 1

//RunOptions holds info related to a specific app run
type RunOptions struct {
	Debug                 int
	ConfigFileName        string //if not empty, points to source config file that was used to load the config
	Command               string
	InputFileNames        []string
	OverWriteOutputFile   bool
	OutputFileName        string
	OutputFormat          string
	WorkDirName           string
	TempDirName           string
	PreserveWorkFiles     bool
	DefaultConfigFileName string
	ExecutableVersion     string `mdson:"-"`
	LibVersion            string `mdson:"-"`
	ExecutableName        string `mdson:"-"`
	DefaultScriptFileName string `mdson:"-"`
}

//DefaultOptions returns a default option setting
func DefaultOptions() *RunOptions {
	return &RunOptions{
		Debug: defaultDebugLevel,
	}
}

//NewOptions returns a new option setting
func NewOptions(ConfigFileName string, debug int) *RunOptions {
	return &RunOptions{
		ConfigFileName: ConfigFileName,
		Debug:          debug,
	}
}

//SetDebug sets the debug level
func (opt *RunOptions) SetDebug(level int) *RunOptions {
	opt.Debug = level
	return opt
}

//SaveToMDSon saves job configuration to specified MDSon writer
func (opt *RunOptions) SaveToMDSon(w io.Writer) error {
	buf, err := mdson.Marshal(opt)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	_, err = w.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to save job configuration: %s", err)
	}
	return nil
}
