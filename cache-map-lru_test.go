package gofnext

import (
	"errors"
	"testing"
	"time"
)

func TestCacheLru_StoreAndLoad(t *testing.T) {
	m := NewCacheLru(100).SetTTL(time.Second)

	// Store a value
	m.Store("key1", "value1", nil)

	// Load the value
	value, existed, err := m.Load("key1")
	if !existed {
		t.Errorf("Expected key1 to exist")
	}
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value1, got: %v", value)
	}

	// Load a non-existent key
	_, existed, _ = m.Load("key2")
	if existed {
		t.Errorf("Expected key2 to not exist")
	}

	// Store a value with an error
	m.Store("key3", nil, errors.New("some error"))

	// Load the value with an error
	_, _, err = m.Load("key3")
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if err.Error() != "some error" {
		t.Errorf("Expected 'some error', got: %v", err)
	}
}

func TestCacheLru_SetTTL(t *testing.T) {
	m := NewCacheLru(100)
	m.SetTTL(time.Second)

	// Store a value
	m.Store("key1", "value1", nil)

	// Load the value before ttl
	_, existed, _ := m.Load("key1")
	if !existed {
		t.Errorf("Expected key1 to exist")
	}

	// Set a shorter ttl
	m.SetTTL(time.Millisecond)

	// Wait for the ttl to expire
	time.Sleep(time.Millisecond * 10)

	// Load the value after ttl
	_, existed, _ = m.Load("key1")
	if existed {
		t.Errorf("Expected key1 to not exist after ttl")
	}
}
