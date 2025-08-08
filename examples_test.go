package wg_test

import (
	"fmt"

	"github.com/tech10/wg"
)

// ExampleWaitGroup_Go demonstrates using the Go method to run a goroutine
// with automatic Add/Done handling.
func ExampleWaitGroup_Go() {
	var wg wg.WaitGroup

	wg.Go(func() {
		fmt.Println("Hello from goroutine!")
	})

	// Wait for goroutines to finish
	wg.Wait()

	// Output:
	// Hello from goroutine!
}

// ExampleWaitGroup_Go_withRecover demonstrates using the Go method to run a goroutine
// with automatic Add/Done handling, and manually recovering from panics.
// You must defer recover() from within f, as f is running in its own goroutine.
// If you do not, the panic caused by f will not be caught.
func ExampleWaitGroup_Go_withRecover() {
	var wg wg.WaitGroup

	wg.Go(func() {
		defer func() {
			if r := recover(); r != nil {
				// Do something with r here.
				fmt.Println("Panic from goroutine!")
			}
		}()
		panic("testing")
	})

	// Wait for goroutines to finish
	wg.Wait()

	// Output:
	// Panic from goroutine!
}

// ExampleWaitGroup_GoRecover demonstrates using the GoRecover method
// which runs a goroutine and recovers from any panic, printing to stderr by default.
func ExampleWaitGroup_GoRecover() {
	var wg wg.WaitGroup

	wg.GoRecover(func() {
		panic("something bad happened")
	})

	// Wait for goroutine to finish
	wg.Wait()

	// Output:
	//
}

// ExampleWaitGroup_GoRecover_withCustomPanicHandler demonstrates GoRecover usage with
// a custom PanicHandler to handle panics differently than printing to stderr.
func ExampleWaitGroup_GoRecover_withCustomPanicHandler() {
	var wg wg.WaitGroup

	wg.PanicHandler = func(p interface{}) {
		fmt.Printf("Custom panic handler caught: %v\n", p)
	}

	wg.GoRecover(func() {
		panic("custom panic")
	})

	wg.Wait()

	// Output:
	// Custom panic handler caught: custom panic
}

// ExampleWaitGroup_RunRecover demonstrates running a function in the current goroutine,
// recovering from panics and returning them as an error.
func ExampleWaitGroup_RunRecover() {
	var wg wg.WaitGroup

	err := wg.RunRecover(func() {
		fmt.Println("Running without panic")
	})

	fmt.Println("Error:", err)

	err = wg.RunRecover(func() {
		panic("panic in RunRecover")
	})

	fmt.Println("Error:", err)

	// Output:
	// Running without panic
	// Error: <nil>
	// Error: PANIC CAUGHT:
	// panic in RunRecover
}
