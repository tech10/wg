// Package wg provides a minimal wrapper around sync.WaitGroup with
// convenient goroutine launching and panic recovery support.
package wg

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

// PanicHandler allows you to define what will happen with panics when GoRecover is used.
// Otherwise, it will fall back to printing to os.Stderr.
type PanicHandler func(interface{})

// WaitGroup wraps sync.WaitGroup, allowing the use of Add, Wait, and Done as usual.
// It adds a couple other features for convenience.
type WaitGroup struct {
	sync.WaitGroup
	PanicHandler PanicHandler
}

// Go runs a function in a new goroutine using the usual Add, defer Done pattern.
// f must not panic.
// If you think f will panic, place a defer function within f, containing recover() and handle panics as you like.
func (wg *WaitGroup) Go(f func()) {
	wg.goRunner(f, false)
}

// GoRecover will run f in a new goroutine just like Go, but it will recover from any panics.
// Any panics caught by this function are printed to stderr, or will use the custom PanicHandler function you define in PanicHandler.
func (wg *WaitGroup) GoRecover(f func()) {
	wg.goRunner(f, true)
}

// RunRecover runs f in the current goroutine to catch any panics
// and return them as an error. The wait group counter is adjusted accordingly.
// This function will not spawn a new goroutine, and will block until f has been executed.
func (wg *WaitGroup) RunRecover(f func()) (err error) {
	wg.Add(1)
	defer wg.Done()
	defer func() { err = errFormat(recover()) }()
	f()
	return
}

// goRunner runs f in a new goroutine, and optionally captures panics with recover.
// The captured panics are printed to os.Stderr or handled by PanicHandler if set.
func (wg *WaitGroup) goRunner(f func(), handlePanics bool) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if handlePanics {
			defer func() {
				if r := recover(); r != nil {
					if wg.PanicHandler == nil {
						fmt.Fprintf(os.Stderr, "Panic recovered in goroutine: %v\n%s\n", r, debug.Stack())
					} else {
						wg.PanicHandler(r)
					}
				}
			}()
		}
		f()
	}()
}

// errFormat formats an error based upon the return value of recover.
func errFormat(r interface{}) error {
	if r == nil {
		return nil
	}
	const prefix = "PANIC CAUGHT:"
	if err, ok := r.(error); ok {
		return fmt.Errorf("%s\n%w", prefix, err)
	}
	return fmt.Errorf("%s\n%v", prefix, r)
}
