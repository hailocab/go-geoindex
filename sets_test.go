package geoindex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	set := newSet()
	assert.Equal(t, set.Size(), 0)

	set.Add(charring.Id(), charring)
	assert.Equal(t, set.Size(), 1)

	set.Add(embankment.Id(), embankment)
	assert.Equal(t, set.Size(), 2)

	set.Remove(charring.Id())
	assert.Equal(t, set.Size(), 1)

	value, ok := set.Get(charring.Id())
	assert.False(t, ok)
	assert.Nil(t, value)

	value, ok = set.Get(embankment.Id())
	assert.True(t, ok)
	assert.NotNil(t, value)
	assert.Equal(t, value.(Point).Id(), "Embankment")

	set.Add(picadilly.Id(), picadilly)
	set.Add(oxford.Id(), oxford)

	assert.True(t, pointsEqualIgnoreOrder(toPoints(set.Values()), []Point{picadilly, embankment, oxford}))
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
	set.Add(picadilly.Id(), picadilly)

	now = currentTime.Add(5 * time.Minute)
	set.Add(oxford.Id(), oxford)
	assert.Equal(t, set.Size(), 2)
	assert.Equal(t, len(set.Values()), 2)

	set.Remove(picadilly.Id())
	assert.Equal(t, set.Size(), 1)

	now = currentTime.Add(11 * time.Minute)
	assert.Equal(t, set.Size(), 1)

	set.Add(oxford.Id(), oxford)
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.Values()), 1)

	now = currentTime.Add(16 * time.Minute)
	assert.Equal(t, set.Size(), 1)
	assert.Equal(t, len(set.Values()), 1)

	now = currentTime.Add(22 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	assert.Equal(t, len(set.Values()), 0)

	now = currentTime.Add(24 * time.Minute)
	assert.Equal(t, set.Size(), 0)
	set.Add(oxford.Id(), oxford)
	now = currentTime.Add(25 * time.Minute)
	set.Add(oxford.Id(), oxford)
	now = currentTime.Add(26 * time.Minute)
	set.Add(oxford.Id(), oxford)
	assert.Equal(t, set.Size(), 1)

	set.Remove(oxford.Id())
	assert.Equal(t, set.Size(), 0)
}
