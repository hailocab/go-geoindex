package geoindex

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var (
	all = func(_ Point) bool { return true }
)

func TestRange(t *testing.T) {
	index := NewPointsIndex(Km(1.0))

	for _, point := range tubeStations() {
		index.Add(point)
	}

	within := index.Range(oxford, embankment)
	expected := []Point{picadilly, charring, coventGarden, embankment, leicester, oxford}

	assert.True(t, pointsEqualIgnoreOrder(expected, within))

	for _, point := range points {
		index.Remove(point.Id())
	}

	assert.Equal(t, len(index.Range(oxford, embankment)), 0)
}

func TestKNearest(t *testing.T) {
	index := NewPointsIndex(Km(0.5))

	for _, point := range tubeStations() {
		index.Add(point)
	}

	assert.Equal(t, index.KNearest(charring, 3, Km(1), all), []Point{charring, embankment, leicester}, true)
	assert.Equal(t, index.KNearest(charring, 5, Km(20), all), []Point{charring, embankment, leicester, coventGarden, picadilly}, true)

	noPicadilly := func(p Point) bool {
		return !strings.Contains(p.Id(), "Piccadilly")
	}
	assert.Equal(t, index.KNearest(charring, 5, Km(20), noPicadilly), []Point{charring, embankment, leicester, coventGarden, westminster}, true)

	assert.Equal(t, index.KNearest(charring, 5, Km(20), all), []Point{charring, embankment, leicester, coventGarden, picadilly}, true)
	assert.Equal(t, len(index.KNearest(charring, 100, Km(1), all)), 9)
}

func TestExpiringIndex(t *testing.T) {
	index := NewExpiringPointsIndex(Km(1.0), Minutes(5))

	currentTime := time.Now()

	now = currentTime
	index.Add(picadilly)

	now = currentTime.Add(1 * time.Minute)
	index.Add(charring)

	now = currentTime.Add(2 * time.Minute)
	index.Add(embankment)

	now = currentTime.Add(3 * time.Minute)
	index.Add(coventGarden)

	now = currentTime.Add(4 * time.Minute)
	index.Add(leicester)

	assert.True(t, pointsEqualIgnoreOrder(index.Range(oxford, embankment), []Point{picadilly, charring, embankment, coventGarden, leicester}))
	assert.Equal(t, index.KNearest(charring, 3, Km(5), all), []Point{charring, embankment, leicester})

	assert.NotNil(t, index.Get(picadilly.Id()))
	assert.NotNil(t, index.Get(charring.Id()))

	now = currentTime.Add(7 * time.Minute)

	assert.Nil(t, index.Get(picadilly.Id()))
	assert.Nil(t, index.Get(charring.Id()))

	assert.NotNil(t, index.Get(embankment.Id()))
	assert.NotNil(t, index.Get(coventGarden.Id()))
	assert.NotNil(t, index.Get(leicester.Id()))

	assert.True(t, pointsEqualIgnoreOrder(index.Range(oxford, embankment), []Point{embankment, coventGarden, leicester}))
	assert.Equal(t, index.KNearest(charring, 3, Km(5), all), []Point{embankment, leicester, coventGarden})
}

func BenchmarkPointIndexRange(b *testing.B) {
	bench(b).CentralLondonRange(NewPointsIndex(Km(1.0)))
}

func BenchmarkPointIndexAdd(b *testing.B) {
	bench(b).AddLondon(NewPointsIndex(Km(0.5)))
}

func BenchmarkPointIndexKNearest(b *testing.B) {
	b.StopTimer()

	index := NewPointsIndex(Km(0.5))

	for i := 0; i < 10000; i++ {
		index.Add(randomPoint())
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		index.KNearest(randomPoint(), 5, Km(5), all)
	}
}

func BenchmarkExpiringPointIndexAdd(b *testing.B) {
	expiration := Minutes(15)
	bench(b).AddLondonExpiring(NewExpiringPointsIndex(Km(0.5), expiration), expiration)
}

func BenchmarkExpiringPointIndexKNearest(b *testing.B) {
	index := NewExpiringPointsIndex(Km(0.5), Minutes(15))

	b.StopTimer()
	for i := 0; i < 10000; i++ {
		index.Add(randomPoint())
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		//index.Add(RandomPoint())
		index.KNearest(randomPoint(), 5, Km(5), all)
	}
}

func BenchmarkExpiringPointIndexRange(b *testing.B) {
	expiration := Minutes(15)
	bench(b).CentralLondonExpiringRange(NewExpiringPointsIndex(Km(0.5), expiration), expiration)
}
