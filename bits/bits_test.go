package bits // import "github.com/solarsea/misc/bits"

import (
	"strconv"
	"testing"
)

func TestHamming(t *testing.T) {
	for i := 0; i < 256; i++ {
		b, result := i, 0
		for b > 0 {
			if (b & 1) > 0 {
				result++
			}
			b = b >> 1
		}
		if result != Hamming(byte(i)) {
			t.Error("Invalid hamming weight reported for ", i)
		}
	}
}

func TestReverse(t *testing.T) {
	for i := 0; i < 256; i++ {
		input, result := i, byte(0)
		for j := 0; j < 7; j++ {
			if (input & 1) > 0 {
				result |= 1
			}
			result <<= 1
			input >>= 1
		}
		if (input & 1) > 0 {
			result |= 1
		}
		if result != Reverse(byte(i)) {
			t.Error("Invalid reverse byte reported for ", i)
		}
	}
}

var sizes []uint = []uint{0, 31, 32, 33, 61, 63, 64, 127, 128, 129}

func TestBitSize(t *testing.T) {
	for _, size := range sizes {
		v := NewBit(size)
		if v.Len() < size || v.Len() > size+strconv.IntSize {
			t.Error("Invalid length", v.Len(), "expected", size)
		}
	}
}

func TestBitEmpty(t *testing.T) {
	var size uint = 128
	v := NewBit(size)

	// Check if it is empty by default
	for i := uint(0); i < size; i++ {
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}

func TestBitBasic(t *testing.T) {
	var size uint = 128
	v := NewBit(size)

	// Raise and lower each position explicitly
	for i := uint(0); i < size; i++ {
		v.Poke(i, true)
		if !v.Peek(i) {
			t.Error("Invalid lowered bit at", i)
		}

		v.Poke(i, false)
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}

func TestBitFlip(t *testing.T) {
	var size uint = 128
	v := NewBit(size)

	// Raise and lower each position by flipping
	for i := uint(0); i < size; i++ {
		v.Flip(i)
		if !v.Peek(i) {
			t.Error("Invalid lowered bit at", i)
		}

		v.Flip(i)
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}

func TestBoolSize(t *testing.T) {
	for _, size := range sizes {
		v := NewBool(size)
		if v.Len() != size {
			t.Error("Invalid length", v.Len(), "expected", size)
		}
	}
}

func TestBoolEmpty(t *testing.T) {
	var size uint = 128
	v := NewBool(size)

	// Check if it is empty by default
	for i := uint(0); i < size; i++ {
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}

func TestBoolBasic(t *testing.T) {
	var size uint = 128
	v := NewBool(size)

	// Raise and lower each position explicitly
	for i := uint(0); i < size; i++ {
		v.Poke(i, true)
		if !v.Peek(i) {
			t.Error("Invalid lowered bit at", i)
		}

		v.Poke(i, false)
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}

func TestBoolFlip(t *testing.T) {
	var size uint = 128
	v := NewBool(size)

	// Raise and lower each position by flipping
	for i := uint(0); i < size; i++ {
		v.Flip(i)
		if !v.Peek(i) {
			t.Error("Invalid lowered bit at", i)
		}

		v.Flip(i)
		if v.Peek(i) {
			t.Error("Invalid raised bit at", i)
		}
	}
}
