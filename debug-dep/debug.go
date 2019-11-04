package debug

import (
	"fmt"
	"strings"
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

var debugString = []string{"Silent", "Warning", "Update", "All"}

// String returns a textual representation of a debug constant
func String(debug int) string {
	if debug >= DebugSilent && debug <= DebugAll {
		return debugString[debug]
	}
	return "Invalid"
}

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

func (u ui) String() string {
	return String(u.Debug)
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
