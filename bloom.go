package bloom

import (
	"hash/fnv"
	"math"

	"github.com/spaolacci/murmur3"
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

func key2Position(hashFunctions []func(uint32, uint32) uint32,
	seed uint32, key []byte) []uint32 {
	hM := murmur3.Sum32WithSeed(key, seed)
	hasherV := fnv.New32()
	hasherV.Write(key)
	hV := hasherV.Sum32()
	positions := make([]uint32, len(hashFunctions))
	for i, h := range hashFunctions {
		positions[i] = h(hM, hV)
	}
	return positions
}

func initHashFunctions(numHashFunctions uint32, numBits uint32) []func(uint32, uint32) uint32 {
	hashFunctions := make([]func(uint32, uint32) uint32, numHashFunctions)
	for i := uint32(0); i < numHashFunctions; i++ {
		hashFunctions[i] = func(h1 uint32, h2 uint32) uint32 {
			return (h1 + i*h2 + i*i) % numBits
		}
	}
	return hashFunctions
}
