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
	value, hasCache, alive, err := m.Load("key1")
	existed := hasCache && alive
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
	_, hasCache, alive, _ = m.Load("key2")
	existed = hasCache && alive
	if existed {
		t.Fatal("Expected key2 to not exist")
	}

	// Store a value with an error
	m.Store("key3", nil, errors.New("some error"))
	m.SetErrTTL(-1)

	// Load the value with an error
	_, _, _, err = m.Load("key3")
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if err.Error() != "some error" {
		t.Fatalf("Expected 'some error', got: %v", err)
	}
}

func TestCacheLru_SetTTL(t *testing.T) {
	m := NewCacheLru(100)
	m.SetTTL(time.Second)

	// Store a value
	m.Store("key1", "value1", nil)

	// Load the value before ttl
	_, hasCache, alive, _ := m.Load("key1")
	existed := hasCache && alive
	if !existed {
		t.Errorf("Expected key1 to exist")
	}

	// Set a shorter ttl
	m.SetTTL(time.Millisecond)

	// Wait for the ttl to expire
	time.Sleep(time.Millisecond * 10)

	// Load the value after ttl
	_, hasCache, alive, _ = m.Load("key1")
	if hasCache && alive {
		t.Errorf("Expected key1 to not exist after ttl")
	}
}

func TestCacheLru_SetReuseTTL(t *testing.T) {
	m := NewCacheLru(100)
	m.SetTTL(time.Millisecond * 10)
	m.SetReuseTTL(time.Millisecond * 10)

	// Store a value
	m.Store("key1", "value1", nil)

	// Load the value before ttl
	_, hasCache, alive, _ := m.Load("key1")
	existed := hasCache && alive
	if !existed {
		t.Errorf("Expected key1 to exist")
	}

	// Wait for the ttl to expire
	time.Sleep(time.Millisecond * 10)
	_, hasCache, alive, _ = m.Load("key1")
	if !(hasCache && !alive) {
		t.Errorf("Unexpected cache: hasCache=%v, alive=%v", hasCache, alive)
	}

	// Wait for the reuse ttl to expire
	time.Sleep(time.Millisecond * 10)
	_, hasCache, alive, _ = m.Load("key1")
	if hasCache {
		t.Errorf("there should be no cache")
	}

}
