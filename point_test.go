package geoindex

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistance(t *testing.T) {
	assert.InDelta(t, float64(distance(waterloo, kingsCross)), 3080.074, 0.001)

	d := math.Sqrt(float64(approximateSquareDistance(waterloo, kingsCross)))
	assert.InDelta(t, d, 3074.987, 0.001)

	d = math.Sqrt(float64(approximateSquareDistance(leicester, coventGarden)))
	assert.InDelta(t, d, 305.662, 0.001)

	d = math.Sqrt(float64(approximateSquareDistance(oxford, embankment)))
	assert.InDelta(t, d, 1593.763, 0.001)
}

func TestDirection(t *testing.T) {
	assert.Equal(t, North, DirectionTo(waterloo, kingsCross))
	assert.Equal(t, NorthEast, DirectionTo(leicester, coventGarden))
	assert.Equal(t, SouthEast, DirectionTo(oxford, embankment))
}

func TestBearing(t *testing.T) {
	b := BearingTo(waterloo, kingsCross)
	assert.InDelta(t, b, -12.659, 0.001)

	b = BearingTo(leicester, coventGarden)
	assert.InDelta(t, b, 57.706, 0.001)

	b = BearingTo(oxford, embankment)
	assert.InDelta(t, b, 122.939, 0.001)
}
