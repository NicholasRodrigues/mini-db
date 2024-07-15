package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageSetAndGet(t *testing.T) {
	s := NewStorage()

	s.Set("key1", "value1")
	s.Set("key2", "value2")

	value, exists := s.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	value, exists = s.Get("key2")
	assert.True(t, exists)
	assert.Equal(t, "value2", value)

	value, exists = s.Get("key3")
	assert.False(t, exists)
	assert.Equal(t, "", value)
}

func TestStorageStore(t *testing.T) {
	s := NewStorage()

	s.Set("key1", "value1")
	s.Set("key2", "value2")

	data := s.Store()
	assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, data)
}
