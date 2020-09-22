package bloom

import (
	"fmt"
	"testing"
)

func TestFindBitCoords(t *testing.T) {
	var tests = []struct {
		index                uint
		byteIndex, bitOffset uint
	}{
		{4, 0, 4},
		{8, 0, 8},
		{9, 0, 9},
		{63, 0, 63},
		{64, 1, 0},
		{65, 1, 1},
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
