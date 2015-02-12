package geoindex

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

const (
	numIterations = 1000
)

func TestQueue(t *testing.T) {
	queue := newQueue(1)

	assert.True(t, queue.IsEmpty())
	assert.Equal(t, queue.Size(), 0)

	queue.Push(1)

	assert.False(t, queue.IsEmpty())
	assert.Equal(t, queue.Size(), 1)

	queue.Push(2)
	assert.Equal(t, queue.Size(), 2)
	assert.Equal(t, queue.Peek().(int), 1)
	assert.Equal(t, queue.PeekBack().(int), 2)

	assert.Equal(t, queue.Pop().(int), 1)
	assert.Equal(t, queue.Pop().(int), 2)

	for i := 1; i < 13; i++ {
		queue.Push(i)
	}

	assert.Equal(t, queue.Peek().(int), 1)
	assert.Equal(t, queue.PeekBack().(int), 12)

	for i := 1; i < 10; i++ {
		assert.Equal(t, queue.Pop().(int), i)
	}

	for i := 1; i < 10; i++ {
		queue.Push(i)
	}

	for i := 10; i < 13; i++ {
		assert.Equal(t, queue.Pop().(int), i)
	}

	for i := 1; i < 10; i++ {
		assert.Equal(t, queue.Pop().(int), i)
	}
}

func TestQueueRandomly(t *testing.T) {
	testRandomly(2, 1)
	testRandomly(3, 2)
	testRandomly(3, 1)
}

func testRandomly(mod int, pushThreshold int) {
	queue := newQueue(1)

	for i := 1; i < numIterations; i++ {
		op := rand.Int() % mod
		if op < pushThreshold {
			queue.Push(i)
		} else {
			queue.Pop()
		}
	}
}
