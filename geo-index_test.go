package geoindex

import (
	"strconv"
	"testing"
)

type TestEntry struct {
	id    int
	count int
}

func (e *TestEntry) Add(p Point) {
	e.count++
}

func (e *TestEntry) Remove(_ Point) {
	e.count--
}

func (e *TestEntry) String() string {
	return strconv.Itoa(e.count)
}

var totalEntries = 0
var newTestEntry = func() interface{} {
	totalEntries++
	result := TestEntry{totalEntries, 0}
	return &result
}

func TestGeoIndexRange(t *testing.T) {
	index := newGeoIndex(Km(0.1), newTestEntry)

	for _, point := range tubeStations() {
		indexEntry := index.AddEntryAt(point)
		entry := (indexEntry).(*TestEntry)
		entry.Add(point)
	}

	entries := index.Range(oxford, embankment)

	count := 0
	for _, entry := range entries {
		testEntry := (entry).(*TestEntry)
		count += testEntry.count
	}

	if count != 6 {
		t.Error("Invalid number of stations")
	}
}
