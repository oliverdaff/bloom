package bloom

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFindBitCoords(t *testing.T) {
	var tests = []struct {
		index                uint32
		byteIndex, bitOffset uint32
	}{
		{4, 0, 4},
		{8, 1, 0},
		{9, 1, 1},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.index)
		t.Run(testname, func(t *testing.T) {
			byteIndex, bitOffset := findBitCoords(tt.index)
			if byteIndex != tt.byteIndex {
				t.Errorf("index %d, got byte index %d, want %d", tt.index, byteIndex, tt.byteIndex)
			}
			if bitOffset != tt.bitOffset {
				t.Errorf("index %d, got bitOffset %d, want %d", tt.index, bitOffset, tt.bitOffset)
			}
		})
	}
}

func TestWriteBit(t *testing.T) {
	var tests = []struct {
		inbits, outbits []byte
		index           uint32
	}{
		{[]byte{0}, []byte{1}, 0},
		{[]byte{0}, []byte{2}, 1},
		{[]byte{0}, []byte{4}, 2},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.index)
		t.Run(testname, func(t *testing.T) {
			result := writeBit(tt.inbits, tt.index)
			if !bytes.Equal(result, tt.outbits) {
				t.Errorf("Index %d input %d got %d expected %d",
					tt.index, tt.inbits, result, tt.outbits)
			}
		})

	}
}

func TestReadBit(t *testing.T) {
	var tests = []struct {
		inbits   []byte
		index    uint32
		expected uint8
	}{
		{[]byte{1}, 0, 1},
		{[]byte{1, 1}, 8, 1},
		{[]byte{1, 1}, 9, 0},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.index)
		t.Run(testname, func(t *testing.T) {
			result := readBit(tt.inbits, tt.index)
			if result != tt.expected {
				t.Errorf("Index %d input %d got %d expected %d",
					tt.index, tt.inbits, result, tt.expected)
			}
		})

	}
}

func TestNewBloomFilter(t *testing.T) {
	var tests = []struct {
		maxSize, seed uint32
		maxTolerance  float64
		numBits       uint32
	}{
		{100, 1, 0.01, 958},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("maxSize %d seed %d, maxTolerance %f",
			tt.maxSize, tt.seed, tt.maxTolerance)
		t.Run(testname, func(t *testing.T) {
			bf := NewBloomFilter(tt.maxSize, tt.maxTolerance, tt.seed)
			if bf.size != 0 {
				t.Errorf("Size %d and expected 0", bf.size)
			}
			if bf.maxSize != tt.maxSize {
				t.Errorf("MaxSize %d and expected %d", bf.maxSize, tt.maxSize)
			}
			if bf.numBits != tt.numBits {
				t.Errorf("NumBits %d and expected %d", bf.numBits, tt.numBits)
			}
		})

	}
}
