//Package com implements a thread-safe inter-package communication object
package com

const (
	Debug int = iota
	Info
	Warning
	Error
	Fatal
)

type Message struct {
	Contents string
	Severity int
}

type Conduit struct {
	backCh   chan Message
	severity int
}

func NewConduit(severity int) *Conduit {
	c := Conduit{
		backCh:   make(chan Message, 10),
		severity: severity,
	}
	return &c
}

func (c Conduit) Warn(msg string) bool {
	c.backCh <- Message{Contents: msg, Severity: Warning}
	return true
}

func (c Conduit) WarnOn() bool {
	return c.severity >= Warning
}
func (c Conduit) InfoOn() bool {
	return c.severity >= Info
}
func (c Conduit) DebugOn() bool {
	return c.severity >= Debug
}

func (c Conduit) Channel() chan Message {
	return c.backCh
}

func (c Conduit) Severity() int {
	return c.severity
}

func (c *Conduit) SetSeverity(severity int) *Conduit {
	c.severity = severity
	return c
}

func (c *Conduit) Close() {
	close(c.backCh)
}
