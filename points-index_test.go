package geoindex

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	all = func(_ Point) bool { return true }
)

func TestClonePointsIndex(t *testing.T) {
	index := NewPointsIndex(Km(1.0))

	for _, point := range tubeStations() {
		index.Add(point)
	}

	clone := index.Clone()

	// C.R.U.F.T.
	if fmt.Sprintf("%p", clone.index) == fmt.Sprintf("%p", index.index) {
		t.Errorf("Clone point index should be pointing to separate geoindex [index-geoindex=%p, clone-geoindex=%p]", index.index, clone.index)
	}
	if fmt.Sprintf("%p", clone.currentPosition) == fmt.Sprintf("%p", index.currentPosition) {
		t.Errorf("Clone currentPosition should be pointing to separate map [index-currentpos=%p, clone-currentpos=%p]", index.currentPosition, clone.currentPosition)
	}
	if !reflect.DeepEqual(index.currentPosition, clone.currentPosition) {
		t.Errorf("Clone currentPosition should have same data as original [index-currentpos=%v, clone-currentpos=%v]", index.currentPosition, clone.currentPosition)
	}

	if !reflect.DeepEqual(index.index.resolution, clone.index.resolution) {
		t.Errorf("Original points index and clone points index do not have the same resolution [original=%+v, clone=%+v]", index.index.resolution, clone.index.resolution)
	}
	if !reflect.DeepEqual(index.index.index, clone.index.index) {
		t.Errorf("Original points index and clone points index are not the same [original=%+v, clone=%+v]", index.index.index, clone.index.index)
	}
	if !reflect.DeepEqual(index.index.newEntry(), clone.index.newEntry()) {
		t.Errorf("Original points index and clone points index new entry functions produce different results [original=%+v, clone=%+v]", index.index.newEntry, clone.index.newEntry)
	}

}

func BenchmarkClone(b *testing.B) {
	b.StopTimer()

	index := NewPointsIndex(Km(1.0))

	for c, point := range worldCapitals() {
		if c > 20 {
			break
		}
		for i := 0; i < 5000; i++ {
			index.Add(&GeoPoint{
				Pid:  point.Id() + fmt.Sprintf("%d", i),
				Plat: point.Lat() + rand.Float64()/3.0,
				Plon: point.Lon() + rand.Float64()/3.0,
			})
		}
	}

	b.StartTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		index.Clone()
	}
}

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
	b.ReportAllocs()
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
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		//index.Add(RandomPoint())
		index.KNearest(randomPoint(), 5, Km(5), all)
	}
}

func BenchmarkExpiringPointIndexRange(b *testing.B) {
	expiration := Minutes(15)
	bench(b).CentralLondonExpiringRange(NewExpiringPointsIndex(Km(0.5), expiration), expiration)
}
