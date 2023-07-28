package ui

import (
	"fmt"
	"strings"
)

type Debug int
const (
	//DebugSilent print errors only
	DebugSilent Debug = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)

func (d Debug) String() string {
	//[...] creates an array rather than a slice
	return [...]string{"Off", "Warn", "Info",
		"All"}[d]
}

//TODO: add error, errorf and move crash and crashf here

// UI handles interactions with the user
type UI interface {
	Log(a ...interface{})
	Logf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Warn(a ...interface{})
	Level() Debug
	SetLevel(level Debug)
}

//TODO: add error and log os.File to allow changing output destination (https://github.com/github/hub/tree/master/ui)
type ui struct {
	Depth int
	debug Debug
}

// NewUI initializes a new default UI
func NewUI(debug Debug) UI {
	return &ui{
		debug: debug,
	}
}

func (u ui) String() string {
	return fmt.Sprintf("ui.debug: %d, ui.Depth: %d", u.debug, u.Depth)
}

func (u ui) Log(a ...interface{}) {
	if u.debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", u.Depth))
		fmt.Println(a...)
	}
}

func (u ui) Logf(format string, a ...interface{}) {
	if u.debug >= DebugAll {
		fmt.Printf(strings.Repeat("  ", u.Depth))
		fmt.Printf(format, a...)
	}
}

func (u ui) Warn(a ...interface{}) {
	if u.debug >= DebugWarning {
		fmt.Println(a...)
	}
}

func (u ui) Info(a ...interface{}) {
	if u.debug >= DebugUpdates {
		fmt.Println(a...)
	}
}

func (u ui) Infof(format string, a ...interface{}) {
	if u.debug >= DebugUpdates {
		fmt.Printf(strings.Repeat("  ", u.Depth))
		fmt.Printf(format, a...)
	}
}

func (u ui) Level() Debug {
	return u.debug
}

func (u *ui) SetLevel(level Debug) {
	u.debug = level
}
