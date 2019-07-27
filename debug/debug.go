package debug

import (
	"fmt"
	"strings"
)

const (
	ExitWithError = 1
	ExitSuccess   = 0
)

const (
	//DebugSilent print errors only
	DebugSilent int = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)

// UI handles interactions with the user
type UI interface {
	Log(a ...interface{})
	Warn(a ...interface{})
}

type ui struct {
	Depth int
	Debug int
}

// NewUI initializes a new default UI
func NewUI(debug int) UI {
	return &ui{
		Debug: debug,
	}
}

func (u ui) Log(a ...interface{}) {
	if u.Debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", u.Depth))
		fmt.Println(a...)
	}
}

// TODO: optimize by removing unnecessary calls to strings.Repeat
func (u ui) Warn(a ...interface{}) {
	if u.Debug >= DebugWarning {
		fmt.Printf(strings.Repeat("  ", u.Depth) + "warning: ")
		fmt.Println(a...)
	}
}


//Crash prints errors and exit with code 2. First line is printed in bold red
func Crash(err error) {
	lines := strings.Split(err.Error(), "\n")
	if len(lines) > 0 {
		//"\033[31;1;4m turn color red and bold. \033[0m reset colors"
		fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", lines[0])
		for i := 1; i < len(lines); i++ {
			fmt.Fprintf(os.Stderr, "%s\n", lines[i])
		}
	}
	os.Exit(ExitWithError)
}

//Crashf prints errors and exit with code 2. First line is printed in bold red
func Crashf(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", err)
	os.Exit(ExitWithError)
}