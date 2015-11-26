package geoindex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	set := newSet()
	assert.Equal(t, set.Size(), 0)

	set.Add(charring.Id())
	assert.Equal(t, set.Size(), 1)

	set.Add(embankment.Id())
	assert.Equal(t, set.Size(), 2)

	set.Remove(charring.Id())
	assert.Equal(t, set.Size(), 1)

	ok := set.Has(charring.Id())
	assert.False(t, ok)

	ok = set.Has(embankment.Id())
	assert.True(t, ok)

	set.Add(picadilly.Id())
	set.Add(oxford.Id())

	assert.Equal(t, set.Size(), 3)
}

func toPoints(values []interface{}) []Point {
	result := make([]Point, 0)
	for _, value := range values {
		result = append(result, value.(Point))
	}
	return result
}

func TestExpiringSet(t *testing.T) {
	set := newExpiringSet(Minutes(10))

	currentTime := time.Now()

	now = currentTime
	set.Add(picadilly.Id())

	now = currentTime.Add(5 * time.Minute)
	set.Add(oxford.Id())
	assert.Equal(t, set.Size(), 2)
	assert.Equal(t, len(set.IDs()), 2)

	set.Remove(picadilly.Id())
	assert.Equal(t, set.Size(), 1)

	now = currentTime.Add(11 * time.Minute)
	assert.Equal(t, set.Size(), 1)

	set.Add(oxford.Id())
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.IDs()), 1)

	now = currentTime.Add(16 * time.Minute)
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.IDs()), 1)

	now = currentTime.Add(22 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	assert.Equal(t, len(set.IDs()), 0)

	now = currentTime.Add(24 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	set.Add(oxford.Id())
	now = currentTime.Add(25 * time.Minute)
	set.Add(oxford.Id())
	now = currentTime.Add(26 * time.Minute)
	set.Add(oxford.Id())
	assert.Equal(t, set.Size(), 1)

	set.Remove(oxford.Id())
	assert.Equal(t, set.Size(), 0)
}
