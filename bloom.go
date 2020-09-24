package bloom

import (
	"hash/fnv"
	"math"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	maxSize, size, seed, numBits, numHashFunctions, numElements uint32
	bitsArray                                                   []byte
	hashFunctions                                               []func(uint32, uint32) uint32
}

func NewBloomFilter(maxSize uint32, maxTolerance float64, seed uint32) *BloomFilter {
	numBits := uint32(math.Ceil(float64(maxSize) * math.Log(maxTolerance) / math.Log(2) / math.Log(2)))
	numElements := uint32(math.Ceil(float64(numBits) / 8))
	numHashFunctions := uint32(-math.Ceil(math.Log2(maxTolerance)))
	return &BloomFilter{
		size:             0,
		maxSize:          maxSize,
		seed:             seed,
		numBits:          numBits,
		numHashFunctions: numHashFunctions,
		numElements:      numElements,
		bitsArray:        make([]byte, numElements),
		hashFunctions:    initHashFunctions(numHashFunctions, numBits),
	}
}

func (bf *BloomFilter) Contains(key []byte) bool {
	positions := key2Position(bf.hashFunctions, bf.seed, key)
	return bf.positionContains(key, positions)
}

func (bf *BloomFilter) Insert(key []byte) {
	positions := key2Position(bf.hashFunctions, bf.seed, key)
	if !bf.positionContains(key, positions) {
		bf.size++
		for _, pos := range positions {
			writeBit(bf.bitsArray, pos)
		}
	}
}

func (bf *BloomFilter) positionContains(key []byte, positions []uint32) bool {
	for _, pos := range positions {
		if readBit(bf.bitsArray, pos) != 0 {
			return false
		}
	}
	return true
}

//
func findBitCoords(index uint32) (uint32, uint32) {
	byteIndex := uint32(math.Floor(float64(index) / float64(8)))
	bitOffset := index % 8
	return byteIndex, bitOffset
}

func readBit(bitsArray []byte, index uint32) uint8 {
	element, bit := findBitCoords(index)
	return uint8((bitsArray[element] & (1 << bit)) >> bit)
}

func writeBit(bitsArray []byte, index uint32) []byte {
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
