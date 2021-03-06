// Package bloom is an implementation of Bloom filter.
package bloom

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/big"

	"github.com/spaolacci/murmur3"
)

// BloomFilter is a datastructure to tell if a element is present in a set.
// The BloomFilter trades accuracy for space.  The contains method returns
// false if the element is definitely not present but true if the element
// might be in the set.
type BloomFilter struct {
	maxSize, size, seed, numBits, numHashFunctions, numElements uint32
	bitsArray                                                   []byte
	hashFunctions                                               []func(uint32, uint32) uint32
}

// NewBloomFilter allocates a new BloomFilter parameterised by the arguments.
//
// 	maxSize - the maximum number of elements the filter is expected to hold (must be > 0)
// 	maxTolerance - the expected accuracy (a sensible default is 0.01)
// 	seed - the seed to use for hashFunctions.
func NewBloomFilter(maxSize uint32, maxTolerance float64, seed uint32) (*BloomFilter, error) {
	if maxSize == 0 {
		return nil, fmt.Errorf("Max Size is 0")
	}
	bigLog2 := big.NewFloat(math.Log(2))

	bigMax := new(big.Float)
	bigMax.SetUint64(uint64(maxSize))
	bigMaxTolerance := big.NewFloat(math.Log(maxTolerance))
	mutResult := new(big.Float)
	numBits64, acc := mutResult.Mul(bigMax, bigMaxTolerance).
		Quo(mutResult, bigLog2).Quo(mutResult, bigLog2).
		Neg(mutResult).Int64()
	if acc != big.Exact {
		numBits64++
	}
	if numBits64 > int64(^uint32(0)) {
		return nil, fmt.Errorf("Number of bits too large than %d", ^uint32(0))
	}

	numBits := uint32(numBits64)

	numElements := uint32(math.Ceil(float64(numBits) / 8))
	numHashFunctions := uint32(-math.Ceil(math.Log2(maxTolerance)))
	return &BloomFilter{
		size:             0,
		maxSize:          maxSize,
		seed:             seed,
		numBits:          numBits,
		numElements:      numElements,
		numHashFunctions: numHashFunctions,
		bitsArray:        make([]byte, numElements),
		hashFunctions:    initHashFunctions(numHashFunctions, numBits),
	}, nil
}

// Contains returns false if the key is definitely not contained in the set
// else returns true if the key might be in the set.
func (bf *BloomFilter) Contains(key []byte) bool {
	positions := key2Position(bf.hashFunctions, bf.seed, key)
	return bf.positionContains(key, positions)
}

// Insert adds the key to the set in the BloomFilter.
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
		if readBit(bf.bitsArray, pos) == 0 {
			return false
		}
	}
	return true
}

// FalsePositiveProbability returns the probability of a false
// positive being returned by BloomFilter.
func (bf *BloomFilter) FalsePositiveProbability() float64 {
	return math.Pow(
		1-math.Exp(float64(bf.numHashFunctions)*float64(bf.size)/
			float64(bf.numBits)),
		float64(bf.numHashFunctions))
}

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
	hasherM := murmur3.New32WithSeed(seed)
	hasherM.Write(key)
	hM := hasherM.Sum32()
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
