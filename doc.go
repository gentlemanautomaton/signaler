// Package signaler is built around the idea of closing channels to signal
// events. It is similar in purpose to the context package, but has a different
// style. Signals can be converted to contexts via signal.Context().
//
// It provides a lightweight API that hides the complexities of channel
// management.
package signaler
