package com_test

import (
	"fmt"
	"testing"

	"github.com/drgo/core/com"
)

func TestCom_Conduit(t *testing.T) {
	c := com.NewConduit(com.Warning)
	_ = c.WarnOn() && c.Warn("hello")
	fmt.Println(<-c.Channel())
}
