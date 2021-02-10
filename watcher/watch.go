// Package watcher a thin wrapper over fsnotify
package watcher

import (
	"context"
	"fmt"
	"time"

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

// WatchDir watches for changes in a dir and send notification through the result channel.
// runs in a loop until ctx is cancelled.
func WatchDir(ctx context.Context, dir string, result chan Event) error {
	towatch := make(chan string, 1)
	towatch <- dir
	return Watch(ctx, towatch, result)
}
//TODO: make delay configurable
const delay time.Duration = 250 * time.Millisecond

func Watch(ctx context.Context, towatch chan string, result chan Event) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watching files failed:%v", err)
	}
	defer watcher.Close()
	go func() {
		sender := time.NewTicker(delay)
		received := make(map[fsnotify.Event]bool)
		send := func(event fsnotify.Event) {
			// fmt.Println("sending event:", event.String())
			result <- Event{event, nil}
			delete(received, event)
		}
		for {
			select {
			//FIXME: unwatch deleted files
			case event, ok := <-watcher.Events:
				// fmt.Println("inside watch():" + event.String())
				if !ok {
					return
				}
				if !received[event]{ //we have not receieved this event
					received[event] = true 
				}
			case err, ok := <-watcher.Errors:
				// fmt.Println("inside watch(): failure" + err.Error())
				if !ok {
					return
				}
				result <- Event{fsnotify.Event{}, err}
			case <-sender.C:  //send any unsent event
        for event := range received {
           send(event)
        } 
			case <-ctx.Done():
				close(towatch)
				return
			}
		}
	}()
	for f := range towatch {
		// fmt.Println("watch(): adding "+ f)
		err = watcher.Add(f)
		if err != nil {
			return fmt.Errorf("watching [%s] failed:%v", f, err)
		}
	}
	return nil
}
