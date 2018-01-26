package signaler

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

// Signaler is responsible for signaling any number of copies of its signal.
// It should be created with New() or derived from an existing signal via
// signal.Derive().
//
// Pass signalers by reference (via pointers). Pass signals by value.
//
// Each signaler will consume resources until signaled in one of three ways:
// by calling Trigger(), through an operating system signal via Capture() or
// implicitly when a parent is signaled.
type Signaler struct {
	Signal
	once sync.Once
	ch   chan struct{}
}

// New returns a new signaler.
func New() *Signaler {
	s := &Signaler{ch: make(chan struct{})}
	s.Signal = s.ch
	return s
}

// Trigger will trigger the signal. It can safely be called more than once.
func (s *Signaler) Trigger() {
	s.once.Do(func() { close(s.ch) })
}

// Deadline will trigger the signal when the current time reaches deadline.
func (s *Signaler) Deadline(deadline time.Time) *Signaler {
	return s.Timeout(time.Until(deadline))
}

// Timeout will trigger the signal after timeout elapses.
func (s *Signaler) Timeout(timeout time.Duration) *Signaler {
	t := time.NewTimer(timeout)

	go func() {
		select {
		case <-s.Signal:
			if !t.Stop() {
				<-t.C
			}
		case <-t.C:
			s.Trigger()
		}
	}()

	return s
}

// Capture listens for the requested operating system signals and will cause
// s to be signaled when one occurs.
//
// The receiver s will be returned for convenient function chaining.
func (s *Signaler) Capture(sig ...os.Signal) *Signaler {
	if len(sig) > 0 {
		go func() {
			osChan := make(chan os.Signal, 1)
			signal.Notify(osChan, sig...)
			select {
			case <-s.Signal:
				signal.Stop(osChan)
			case <-osChan:
				signal.Stop(osChan)
				s.Trigger()
			}
		}()
	}
	return s
}
