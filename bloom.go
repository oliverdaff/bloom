package bloom

import "math"

func findBitCoords(index uint) (uint, uint) {
	byteIndex := uint(math.Ceil(float64(index) / float64(8)))
	bitOffset := index % 8
	return byteIndex, bitOffset
}
