package bloom

import (
	"math"
)

func findBitCoords(index uint) (uint, uint) {
	byteIndex := uint(math.Floor(float64(index) / float64(64)))
	bitOffset := index % 64
	return byteIndex, bitOffset
}
