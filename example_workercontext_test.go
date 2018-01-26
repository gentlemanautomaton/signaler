package signaler_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func Example_workerContext() {
	var (
		// Prepare a cancellation that we'll trigger manually or through an OS signal
		cancel = signaler.New().Capture(os.Interrupt, os.Kill)

		// Prepare a derived cancellation with a deadline (which we won't meet)
		when     = time.Now().Add(time.Second * 5)
		deadline = cancel.Derive().Deadline(when).Then(func() {
			fmt.Println("Stopping")
		})
	)

	// Always clean up your mess
	// It's safe to call this more than once
	defer cancel.Trigger()

	// Simulated worker bee that needs a context
	doWork := func(ctx context.Context) <-chan error {
		ch := make(chan error)
		go func() {
			fmt.Println("Started")
			<-ctx.Done()
			fmt.Println("Stopped")
			ch <- ctx.Err()
		}()
		return ch
	}

	// Start the worker and have it stop when deadline is signaled
	result := doWork(deadline.Context())

	// Cancel immediately instead of waiting for the deadline to expire
	// This will signal deadline because deadline is derived from cancel
	cancel.Trigger()

	// Print the result
	fmt.Printf("Reason: %v", <-result)
	// Output:
	// Started
	// Stopping
	// Stopped
	// Reason: context canceled
}
