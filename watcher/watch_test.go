package watcher_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/drgo/core/tests"
	"github.com/drgo/core/watcher"
)

func TestNew(t *testing.T) {
  dir, clean:= tests.MkTempDir(t)
  defer clean()
	events := make(chan watcher.Event)
	// monitor events
	go func() {
		for e := range events {
			fmt.Println("Event:", e.String())
		}
	}()
	// start watcher
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
    err := watcher.New(ctx, dir, events)
		if err != nil {
			t.Fatalf("Error", err)
		}
	}()
	// induce watcher events
	// time.Sleep(1 * time.Second)
  f, del := tests.MkTempFile(t, dir) 
	f.WriteString("testing...")
	f.Sync()
  f.Close()
  time.Sleep(50 * time.Millisecond) // give system time to sync write change before delete
  del()
}
