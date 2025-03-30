package main

import (
	"context"
	"fmt"
	"sync"

	udev "github.com/jochenvg/go-udev"
)

func main() {
	u := udev.Udev{}
	m := u.NewMonitorFromNetlink("udev")

	// Create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	// Start monitor goroutine and get receive channel
	ch, _, err := m.DeviceChan(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		fmt.Println("Started listening on channel")
		for d := range ch {
			fmt.Println("Event:", d.Syspath(), d.Action())
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("Channel closed")
}
