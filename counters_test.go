package geoindex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExpiringCounter(t *testing.T) {
	cur := time.Now().Truncate(1 * time.Minute)
	now = cur
	counter := newExpiringCounter(Minutes(3))

	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(int), 1)

	now = cur.Add(50 * time.Second)
	counter.Add(picadilly)
	assert.Equal(t, counter.Point().Count.(int), 2)

	now = cur.Add(61 * time.Second)
	counter.Add(picadilly)
	assert.Equal(t, counter.Point().Count.(int), 3)

	now = cur.Add(70 * time.Second)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(int), 4)

	now = cur.Add(4 * time.Minute)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(int), 3)

	now = cur.Add((60*4 + 30) * time.Second)
	counter.Add(picadilly)
	assert.Equal(t, counter.Point().Count.(int), 4)

	now = cur.Add(5 * time.Minute)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(int), 5)

	now = cur.Add(6 * time.Minute)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(int), 4)

	now = cur.Add(7 * time.Minute)
	assert.Equal(t, counter.Point().Count.(int), 4)

	now = cur.Add(8 * time.Minute)
	assert.Equal(t, counter.Point().Count.(int), 2)

	now = time.Time{}
}

func TestMultiCounter(t *testing.T) {
	cur := time.Now().Truncate(1 * time.Minute)
	now = cur

	counter := newExpiringMultiCounter(Minutes(3))

	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(map[string]int)[oxford.Id()], 1)
	now = cur.Add(1 * time.Minute)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(map[string]int)[oxford.Id()], 2)
	now = cur.Add(2 * time.Minute)
	counter.Add(oxford)
	assert.Equal(t, counter.Point().Count.(map[string]int)[oxford.Id()], 3)

	now = cur.Add(4 * time.Minute)
	assert.Equal(t, counter.Point().Count.(map[string]int)[oxford.Id()], 2)

	now = cur.Add(5 * time.Minute)
	assert.Equal(t, counter.Point().Count.(map[string]int)[oxford.Id()], 1)
}

func assertCountPoint(t *testing.T, point *CountPoint, lat, lon, count float64) {
	assert.Equal(t, point.Lat(), lat)
	assert.Equal(t, point.Lon(), lon)
	assert.Equal(t, point.Count.(float64), count)
}

func TestAverageAccumulatingCounter(t *testing.T) {
	counter := newAverageAccumulatingCounter(&CountPoint{&GeoPoint{Plat: 1.0, Plon: 2.0, Pid: ""}, 3.0})

	counter.Add(&CountPoint{&GeoPoint{Plat: 2.0, Plon: 4.0, Pid: ""}, 6.0})
	counter.Add(&CountPoint{&GeoPoint{Plat: 3.0, Plon: 6.0, Pid: ""}, 9.0})
	assertCountPoint(t, counter.Point(), 2.0, 4.0, 6.0)

	counter.Remove(&CountPoint{&GeoPoint{Plat: 3.0, Plon: 6.0, Pid: ""}, 9.0})
	assertCountPoint(t, counter.Point(), 1.5, 3.0, 4.5)

	anotherCounter := newAverageAccumulatingCounter(&CountPoint{&GeoPoint{Plat: 3.0, Plon: 6.0, Pid: ""}, 9.0})
	counter.Plus(anotherCounter)
	assertCountPoint(t, counter.Point(), 2.0, 4.0, 6.0)

	counter.Minus(anotherCounter)
	assertCountPoint(t, counter.Point(), 1.5, 3.0, 4.5)
}
