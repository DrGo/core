// Package watcher a thin wrapper over fsnotify
package watcher

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
)

// Event encapsulate fsnotify events and errors
type Event struct {
	fsnotify.Event
	err error
}

// IsWrite checks if the triggered event is fsnotify.Write|fsnotify.Create.
func (event Event) IsWrite() bool {
	return event.Op&fsnotify.Write == fsnotify.Write ||
		event.Op&fsnotify.Create == fsnotify.Create
}

// IsRemove checks if the triggered event is fsnotify.Remove.
func (event Event) IsRemove() bool {
	return event.Op&fsnotify.Remove == fsnotify.Remove
}

func (event Event) String() string {
	return fmt.Sprintf("%q: %s", event.Name, event.Op.String())
}

// New watches for changes in a dir and send notification through the result channel.
// runs in a loop until ctx is cancelled.
func New(ctx context.Context, dir string, result chan Event) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watching files failed:%v", err)
	}
	defer watcher.Close()
	done := make(chan bool) //used below to keep the go routine alive
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				result <- Event{event, nil}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				result <- Event{fsnotify.Event{}, err}
			case <-ctx.Done():
			  close(done)
				return
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("watching files failed:%v", err)
	}
	<-done
	return nil
}




