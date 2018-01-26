package signaler_test

import (
	"fmt"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func Example_sequence() {
	var (
		count    = func() { fmt.Println("1 2 3 4") }
		sing     = func() { fmt.Println("la la la la la") }
		delay    = func() { time.Sleep(time.Millisecond * 50) }
		start    = signaler.New()
		finished = start.Then(count).Then(delay).Then(sing)
	)
	go start.Trigger()
	finished.Wait()
	// Output:
	// 1 2 3 4
	// la la la la la
}
