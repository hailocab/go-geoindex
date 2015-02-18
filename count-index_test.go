package geoindex

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccumulatingCounter(t *testing.T) {
	counter := &singleValueAccumulatingCounter{}

	assert.Equal(t, counter.String(), "0.000000 0.000000 0")
	counter.Add(oxford)
	assert.Equal(t, counter.String(), "51.515110 -0.141700 1")
	counter.Add(embankment)
	assert.Equal(t, counter.String(), "103.022422 -0.264067 2")
	counter.Add(picadilly)
	assert.Equal(t, counter.String(), "154.532282 -0.397767 3")

	counter.Remove(embankment)
	assert.Equal(t, counter.String(), "103.024970 -0.275400 2")
	counter.Remove(oxford)
	assert.Equal(t, counter.String(), "51.509860 -0.133700 1")
	counter.Remove(picadilly)
	assert.Equal(t, counter.String(), "0.000000 0.000000 0")
}

func TestCountIndex(t *testing.T) {
	countIndex := NewCountIndex(Km(3.0))

	for _, station := range tubeStations() {
		countIndex.Add(station)
	}

	counters := countIndex.Range(oxford, embankment)
	expected := []Point{
		&CountPoint{&GeoPoint{"", 51.500776, -0.158290}, 7},
		&CountPoint{&GeoPoint{"", 51.504882, -0.124143}, 11},
		&CountPoint{&GeoPoint{"", 51.522122, -0.157587}, 12},
		&CountPoint{&GeoPoint{"", 51.523935, -0.129213}, 10},
	}

	assert.True(t, pointsEqual(counters, expected))
}

func TestExpiringCountIndex(t *testing.T) {
	countIndex := NewExpiringCountIndex(Km(0.5), Minutes(1))

	currentTime := time.Now().Truncate(1 * time.Minute)

	londonTopLeft := &GeoPoint{"", 51.747439, -0.704713}
	londonBottomRight := &GeoPoint{"", 51.249023, 0.484557}

	for i, station := range tubeStations() {
		now = currentTime.Add(time.Duration(i) * time.Minute)
		countIndex.Add(station)
		count := len(countIndex.Range(londonTopLeft, londonBottomRight))

		t.Log(count, " ", station)
	}

	now = time.Time{}
}

func BenchmarkCountIndexAdd(b *testing.B) {
	bench(b).AddLondon(NewCountIndex(Km(0.5)))
}

func BenchmarkCountIndexCityRange(b *testing.B) {
	bench(b).LondonRange(NewCountIndex(Km(10)))
}

func BenchmarkExpiringCountIndexAdd(b *testing.B) {
	expiration := Minutes(15)
	bench(b).AddLondonExpiring(NewExpiringCountIndex(Km(0.5), expiration), expiration)
}

func BenchmarkExpiringCountIndexRange(b *testing.B) {
	bench(b).CentralLondonRange(NewExpiringCountIndex(Km(0.5), Minutes(15)))
}
