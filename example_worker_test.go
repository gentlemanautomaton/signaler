package signaler_test

import (
	"fmt"
	"os"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func Example_worker() {
	var (
		// Start manually
		start = signaler.New()
		// Stop via ctrl+c, after 20ms or manually
		stop = signaler.New().Capture(os.Interrupt).Timeout(time.Millisecond * 20)
		// Post-messaging signals
		starting = start.Then(func() { fmt.Println("Starting") })
		stopping = stop.Then(func() { fmt.Println("Stopping") })
	)

	// Simulated worker bee
	doWork := func(starting, stopping signaler.Signal) signaler.Signal {
		stopped := signaler.New()
		go func() {
			defer stopped.Trigger()

			starting.Wait()
			for {
				if stopping.Signaled() {
					return
				}
				time.Sleep(time.Millisecond)
			}
		}()
		return stopped.Signal
	}

	// Make sure we always run our shutdown procedure
	defer stop.Trigger()

	// Indicate that we're ready to start working
	start.Trigger()

	// Do the work then print a completion message (and wait until finished)
	doWork(starting, stopping).Then(func() {
		fmt.Println("Stopped")
	}).Wait()
	// Output:
	// Starting
	// Stopping
	// Stopped
}
