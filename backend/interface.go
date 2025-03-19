package backend

import (
	"fmt"
	"strings"
)

// Backend is a interface for s3-compatible backends. We only need read for now.
type Backend interface {
	// Get returns the object for the given key
	Get(key string) ([]byte, error)
	Exists(key string) (bool, error)
	List(prefix string) ([]string, error)
}

var ErrKeyNotFound = fmt.Errorf("key not found")

type MockBackend struct {
	Data map[string]string
}

func (b *MockBackend) Get(key string) ([]byte, error) {
	it, ok := b.Data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return []byte(it), nil
}

func (b *MockBackend) Exists(key string) (bool, error) {
	_, ok := b.Data[key]
	return ok, nil
}

func (b *MockBackend) List(prefix string) ([]string, error) {
	keys := []string{}
	for k := range b.Data {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func NewMockBackend() Backend {
	return &MockBackend{
		Data: map[string]string{},
	}
}

func NewMockBackendWithData(data map[string]string) Backend {
	return &MockBackend{
		Data: data,
	}
}

func (b *MockBackend) Put(key string, data string) error {
	b.Data[key] = data
	return nil
}
