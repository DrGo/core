package watcher

import (
	"context"
	"fmt"
	"os"
	"time"
)

func ExampleNew() {
	err := os.Mkdir("test", 0755)
	if err != nil && !os.IsExist(err){
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
	  os.Remove("test/test")
		os.WriteFile("test/test.txt", []byte("test"), 0666)
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = New(ctx, "test", events)
	if err != nil {
		fmt.Println("Error", err)
	}
	// time.Sleep(3 * time.Second)
	
	// Output:
	// Event: "test/test.txt": CHMOD
	// Event: "test/test.txt": WRITE
}
