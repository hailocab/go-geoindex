package geoindex

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
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
