package signaler

import "context"

// Signal is a channel that will be closed when signaled.
//
// Pass signals by value. Copy signals by value.
//
// Signals do not consume resources and may be discarded.
type Signal <-chan struct{}

// Signaled returns true if s has been signaled.
func (s Signal) Signaled() bool {
	select {
	case <-s:
		return true
	default:
		return false
	}
}

// Wait blocks until s is signaled.
//
// If s has already been signaled Wait will return immediately.
func (s Signal) Wait() {
	<-s
}

// Then causes f() to be executed in its own goroutine when s is signaled.
//
// If s has already been signaled f() will be executed immediately.
//
// Then returns a child that will be signaled after the call to f() has
// completed.
//
// If then is called more than once, the order in which the functions are
// executed is undefined.
func (s Signal) Then(f func()) (child Signal) {
	c := make(chan struct{})
	go func() {
		<-s
		f()
		close(c)
	}()
	return c
}

// Derive returns a signaler that will automatically be signaled when s is.
func (s Signal) Derive() (child *Signaler) {
	child = New()
	go func() {
		<-s
		child.Trigger()
	}()
	return
}

// Context returns a context.Context that will be cancelled when s is signaled.
func (s Signal) Context() (ctx context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-s
		cancel()
	}()
	return
}
