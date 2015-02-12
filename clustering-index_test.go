package geoindex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testPoints = []Point{picadilly, oxford, londonBridge, regentsPark, charring}
)

func TestClusteringIndex(t *testing.T) {
	index := NewClusteringIndex()

	for _, point := range testPoints {
		index.Add(point)
	}

	assert.True(t, distance(regentsPark, londonBridge) < streetLevel)
	assert.Equal(t, index.Within(regentsPark, londonBridge), []Point{charring, londonBridge, picadilly, oxford, regentsPark})

	assert.True(t, distance(aylesbury, aylesford) < cityLevel)
	assert.True(t, distance(aylesbury, aylesford) > streetLevel)

	expected := []Point{&CountPoint{&GeoPoint{"", 51.514200, -0.136751}, 4}, &CountPoint{&GeoPoint{"", 51.504674, -0.086006}, 1}}
	actual := index.Within(aylesbury, aylesford)
	assert.True(t, pointsEqual(expected, actual))

	assert.True(t, distance(reykjavik, ankara) > cityLevel)

	expected = []Point{&CountPoint{&GeoPoint{"", 51.512295, -0.126602}, 5}}
	actual = index.Within(reykjavik, ankara)
	assert.True(t, pointsEqual(expected, actual))

	// test remove
	index.Remove(oxford)
	expected = []Point{charring, londonBridge, picadilly, regentsPark}
	actual = index.Within(regentsPark, londonBridge)
	assert.True(t, pointsEqual(expected, actual))

	expected = []Point{&CountPoint{&GeoPoint{"", 51.513896, -0.135101}, 3}, &CountPoint{&GeoPoint{"", 51.504674, -0.086006}, 1}}
	actual = index.Within(aylesbury, aylesford)
	assert.True(t, pointsEqual(actual, expected))

	expected = []Point{&CountPoint{&GeoPoint{"", 51.511591, -0.122827}, 4}}
	actual = index.Within(reykjavik, ankara)
	assert.True(t, pointsEqual(actual, expected))
}

// Benchmark adding points to the clustering index
func BenchmarkClusterIndexAdd(b *testing.B) {
	bench(b).AddWorldWide(NewClusteringIndex())
}

// Benchmark doing within query on the street level
func BenchmarkClusterIndexWithinStreet(b *testing.B) {
	bench(b).WithinCentralLondon(NewClusteringIndex())
}

// Benchmark doing within query on the city level
func BenchmarkClusterIndexWithinCity(b *testing.B) {
	bench(b).WithinLondon(NewClusteringIndex())
}

// Benchmark doing within query on the world level
func BenchmarkClusterIndexWithinWorld(b *testing.B) {
	bench(b).WithinEurope(NewClusteringIndex())
}

// Benchmark adding points to the clustering index
func BenchmarkExpiringClusterIndexAdd(b *testing.B) {
	expiration := Minutes(15)
	bench(b).AddLondonExpiring(NewExpiringClusteringIndex(expiration), expiration)
}

// Benchmark doing within query on the street level
func BenchmarkExpiringClusterIndexWithinStreet(b *testing.B) {
	expiration := Minutes(15)
	bench(b).WithinCentralLondon(NewExpiringClusteringIndex(expiration))
}

// Benchmark doing within query on the city level
func BenchmarkExpiringClusterIndexWithinCity(b *testing.B) {
	expiration := Minutes(15)
	bench(b).WithinLondon(NewExpiringClusteringIndex(expiration))
}

// Benchmark doing within query on the world level
func BenchmarkExpiringClusterIndexWithinWorld(b *testing.B) {
	expiration := Minutes(15)
	bench(b).WithinEurope(NewExpiringClusteringIndex(expiration))
}
