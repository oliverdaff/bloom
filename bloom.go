package bloom

import (
	"math"
)

//
func findBitCoords(index uint) (uint, uint) {
	byteIndex := uint(math.Floor(float64(index) / float64(8)))
	bitOffset := index % 8
	return byteIndex, bitOffset
}

func readBit(bitsArray []byte, index uint) uint8 {
	element, bit := findBitCoords(index)
	return uint8((bitsArray[element] & (1 << bit)) >> bit)
}

func writeBit(bitsArray []byte, index uint) []byte {
	element, bit := findBitCoords(index)
	bitsArray[element] = bitsArray[element] | (1 << bit)
	return bitsArray
}
