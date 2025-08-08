# Introduction

This is a clean, minimal wrapper around sync.WaitGroup, with a few functions for convenience:

1. Go (runs func() inside its own goroutine, automatically adding wg.Add and defer wg.Done logic)
2. GoRecover (same as Go, but will recover from panics and print them to os.Stderr)
3. RunRecover (will run func() within the same goroutine and return any recovered panics as an error)

sync.WaitGroup is convenient for waiting on concurrent operations, but this will make it more so under certain circumstances. Go 1.25 is planning to have sync.WaitGroup.Go available, which is nearly identical to my version, with the exception of no panic recoveries. It's from Go's implementation that I got the idea to create this wrapper.

# Examples

Take a look at `wg_test.go` and `examples_test.go` for examples on how to use this package correctly. The tests should pass race testing, and they provide a minimal, but example worthy test framework to test this wrapper.

# Contributions

Follow Go formatting guidelines and submit an issue or PR.
