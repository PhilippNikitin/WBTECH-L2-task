package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	// Helper function to create a channel that closes after a specific duration
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("No channels", func(t *testing.T) {
		done := or()
		if done != nil {
			t.Error("expected nil, but got a non-nil channel")
		}
	})

	t.Run("Single channel", func(t *testing.T) {
		done := or(sig(1 * time.Millisecond))
		select {
		case <-done:
			// Expected behavior
		case <-time.After(10 * time.Millisecond):
			t.Error("expected the channel to close within 1 millisecond")
		}
	})

	t.Run("Multiple channels", func(t *testing.T) {
		done := or(
			sig(5*time.Second),
			sig(500*time.Millisecond),
			sig(1*time.Second),
		)
		select {
		case <-done:
			// Expected behavior
		case <-time.After(1 * time.Second):
			t.Error("expected the channel to close within 1 second")
		}
	})

	t.Run("All channels already closed", func(t *testing.T) {
		c1 := make(chan interface{})
		c2 := make(chan interface{})
		close(c1)
		close(c2)

		done := or(c1, c2)
		select {
		case <-done:
			// Expected behavior
		case <-time.After(10 * time.Millisecond):
			t.Error("expected the channel to be already closed")
		}
	})

	t.Run("Some channels already closed", func(t *testing.T) {
		c1 := make(chan interface{})
		close(c1)

		done := or(c1, sig(1*time.Second))
		select {
		case <-done:
			// Expected behavior
		case <-time.After(10 * time.Millisecond):
			t.Error("expected the channel to be already closed")
		}
	})
}
