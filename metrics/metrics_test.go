package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSanity(t *testing.T) {
	key := "test"
	val := time.Duration(10)

	m := NewMetrics()
	for i := 0; i < 1000; i++ {
		m.Put(key, val)
	}

	assert.Equal(t, val, m.Get(key))
}

func TestAvg(t *testing.T) {
	key := "test"

	m := NewMetrics()
	m.Put(key, 1)
	m.Put(key, 2)
	m.Put(key, 3)
	m.Put(key, 4)
	m.Put(key, 5)

	assert.Equal(t, time.Duration(3), m.Get(key))
}
