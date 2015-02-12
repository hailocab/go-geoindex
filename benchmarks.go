package geoindex

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var (
	lonCenterLat = 51.512161
	lonCenterLon = -0.123811
	pointIndex   = 0
)

func randSign() float64 {
	if rand.Float64() > 0.5 {
		return 1.0
	} else {
		return -1.0
	}
}

func RandomPoint() Point {
	lat := lonCenterLat + rand.Float64()/4.0*randSign()
	lon := lonCenterLon + rand.Float64()/4.0*randSign()
	pointIndex++
	return &GeoPoint{strconv.Itoa(pointIndex), lat, lon}
}

var (
	capitals []Point = nil
)

func RandomPointWorldWide() Point {
	if capitals == nil {
		capitals = worldCapitals()
	}

	index := rand.Int() % len(capitals)
	pointIndex++
	lat := capitals[index].Lat() + rand.Float64()/4.0*randSign()
	lon := capitals[index].Lon() + rand.Float64()/4.0*randSign()

	return &GeoPoint{strconv.Itoa(pointIndex), lat, lon}
}

type Index interface {
	Add(point Point)
	Within(topLeft Point, bottomRight Point) []Point
	KNearest(point Point, k int, maxDistance Meters, accept func(p Point) bool) []Point
}

type benchmarks struct {
	b *testing.B
}

func bench(b *testing.B) benchmarks {
	return benchmarks{b}
}

func add(index Index, n int, generatePoint func() Point) {
	for i := 0; i < n; i++ {
		index.Add(generatePoint())
	}
}

func addStopTimer(index Index, n int, generatePoint func() Point, b benchmarks) {
	b.b.StopTimer()
	add(index, n, generatePoint)
	b.b.StartTimer()
}

func (b benchmarks) AddWorldWide(index Index) {
	add(index, b.b.N, RandomPointWorldWide)
}

func (b benchmarks) AddLondon(index Index) {
	add(index, b.b.N, RandomPoint)
}

func toMinute(i int, n int, expiration Minutes) time.Duration {
	count := int(float64(i) / float64(n) * float64(expiration) * 2)
	return time.Duration(count) * time.Minute
}

func (b benchmarks) AddLondonExpiring(index Index, expiration Minutes) {
	currentTime := time.Now()
	now = currentTime

	for i := 0; i < b.b.N; i++ {
		minute := toMinute(i, b.b.N, expiration)
		now = currentTime.Add(minute)
		index.Add(RandomPoint())
	}
}

func (b benchmarks) WithinCentralLondon(index Index) {
	addStopTimer(index, 10000, RandomPoint, b)

	for i := 0; i < b.b.N; i++ {
		index.Within(regentsPark, londonBridge)
	}
}

func (b benchmarks) WithinCentralLondonExpiring(index Index, expiration Minutes) {
	currentTime := time.Now()
	now = currentTime

	addStopTimer(index, 10000, RandomPoint, b)

	for i := 0; i < b.b.N; i++ {
		minute := toMinute(i, b.b.N, expiration)
		now = currentTime.Add(minute)
		index.Within(regentsPark, londonBridge)
	}
}

func (b benchmarks) WithinLondon(index Index) {
	addStopTimer(index, 10000, RandomPoint, b)

	for i := 0; i < b.b.N; i++ {
		index.Within(watford, swanley)
	}
}

func (b benchmarks) WithinEurope(index Index) {
	addStopTimer(index, 200000, RandomPointWorldWide, b)

	for i := 0; i < b.b.N; i++ {
		index.Within(reykjavik, ankara)
	}
}
