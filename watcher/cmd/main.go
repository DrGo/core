package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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

// New watches for changes in a dir and send notification through the result channel.
// runs in a loop until ctx is cancelled.
func New(ctx context.Context, dir string, result chan Event) error {
	var mu sync.Mutex
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
				result <- Event{event, err}
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
  mu.Lock()
	err = watcher.Add(dir)
	mu.Unlock()
	if err != nil {
		return fmt.Errorf("watching files failed:%v", err)
	}
	<-done
	return nil
}

func main() {
	err := os.Mkdir("test", 0755)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	events := make(chan Event)
	go func() {
		for e := range events {
			fmt.Println("Event:", e.String())
		}
	}()
	go func() {
		time.Sleep(1 * time.Second)
		os.Remove("test/test.txt")
		ioutil.WriteFile("test/test.txt", []byte("test"), 0666)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	err = New(ctx, "test", events)
	if err != nil {
		fmt.Println("Error", err)
	}
	time.Sleep(3 * time.Second)
	cancel()
	// Output:
	// Event: "test/test.txt": CHMOD
	// Event: "test/test.txt": WRITE
}
