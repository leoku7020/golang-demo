package bootkit

import (
	"context"
	"fmt"
	"time"
)

// Example simulates a boot procedure.
func Example() {
	AddShutdownHandler(func() error {
		fmt.Println("shutdown handler with default")
		return nil
	})

	defer func() {
		fmt.Println("defer at top")
	}()

	AddShutdownHandler(func() error {
		fmt.Println("shutdown handler with index one")
		return nil
	}, WithShutdownLevel(1))

	Register(func(shutdown ShutdownFunc) error {
		pause := make(chan struct{})
		shutdown(func() error {
			fmt.Println("register func closing")
			close(pause)
			return nil
		})

		fmt.Println("register func starting")

		<-pause
		return nil
	})

	defer func() {
		fmt.Println("defer at bottom")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	if err := ListenAndServe(ctx); err != nil {
		panic("ListenAndServe failed")
	}

	// Wait waits the signal triggered from os.Signal or manually.
	Wait()

	// the order to shut down whole processes begins at those registerd in Register().
	// then those handled by AddShutdownHandler().
	// the final would be defer functions.

	// Output:
	// register func starting
	// register func closing
	// shutdown handler with default
	// shutdown handler with index one
	// defer at bottom
	// defer at top
}
