package wg_test

import (
	"errors"
	"testing"

	"github.com/tech10/wg"
)

func TestGo(t *testing.T) {
	var wgGroup wg.WaitGroup
	counter := 0
	wgGroup.Go(func() {
		counter++
	})
	wgGroup.Wait()
	if counter != 1 {
		t.Errorf("expected counter = 1, got %d", counter)
	}
	t.Log("No deadlock.")
}

func TestRunRecover_NoPanic(t *testing.T) {
	var wgGroup wg.WaitGroup
	err := wgGroup.RunRecover(func() {
		// No panic here
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	wgGroup.Wait()
	t.Log("No deadlock.")
}

func TestRunRecover_WithPanic(t *testing.T) {
	var wgGroup wg.WaitGroup
	panicErr := errors.New("test panic")
	err := wgGroup.RunRecover(func() {
		panic(panicErr)
	})
	if err == nil || !errors.Is(err, panicErr) {
		t.Errorf("expected wrapped error, got %v", err)
	}
	wgGroup.Wait()
	t.Log("No deadlock.")
}

func TestGoRecover_withNoCustomPanicHandler(t *testing.T) {
	var wgGroup wg.WaitGroup
	wgGroup.GoRecover(func() {
		panic("catch this panic in the default manner")
	})
	wgGroup.Wait()
	t.Log("No deadlock. Panic caught.")
}

func TestGoRecover_withCustomPanicHandler(t *testing.T) {
	var wgGroup wg.WaitGroup
	wgGroup.PanicHandler = func(r interface{}) {
		if r != nil {
			t.Log("Panic caught successfully.")
		} else {
			t.Errorf("Expected panic to be caught, received nil value.")
		}
	}
	wgGroup.GoRecover(func() {
		panic("catch this panic")
	})
	wgGroup.Wait()
	t.Log("No deadlock. Panic caught.")
}

func TestGoRecover_withCustomPanicHandlerAndNestedPanic(t *testing.T) {
	var wgGroup wg.WaitGroup
	wgGroup.PanicHandler = func(r interface{}) {
		if r != nil {
			t.Log("Original panic caught successfully.")
		} else {
			t.Errorf("Expected panic to be caught, received nil value.")
		}
		panic("this is a nested panic")
	}
	wgGroup.GoRecover(func() {
		panic("catch this panic")
	})
	wgGroup.Wait()
	t.Log("No deadlock. Nested panic caught.")
}
