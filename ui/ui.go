package ui

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

// UI handles interactions with the user
type UI interface {
	Log(a ...interface{})
	Logf(format string, a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Level() int
	SetLevel(level int)
}

type ui struct {
	Depth int
	debug int
}

// NewUI initializes a new default UI
func NewUI(debug int) UI {
	return &ui{
		debug: debug,
	}
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
	if u.debug >= DebugWarning {
		fmt.Println(a...)
	}
}

func (u ui) Level() int {
	return u.debug
}

func (u ui) SetLevel(level int) {
	u.debug = level
}
