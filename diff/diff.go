// Package diff provides a diff algorithm implementation
// for finite, indexable sequences with comparable elements.
package diff // import "github.com/solarsea/misc/diff"

import (
	bits "github.com/solarsea/misc/bits"
)

// Interface abstracts the required knowledge to perform a diff
// on any two fixed-length sequences with comparable elements.
type Interface interface {
	// The sequences' lengths
	Len() (int, int)
	// True when the sequences' elements at those indices are equal
	Equal(int, int) bool
}

// A Delta struct is the result of a Diff operation
type Delta struct {
	// Indices of added elements in the second sequence
	Added   [][]int
	// Indices of removed elements from the first sequence
	Removed [][]int
}

// Diffs the provided data and returns e Delta struct
// with added entries' indices in the second sequence and removed from the first
func Diff(data Interface) Delta {
	var len1, len2 = data.Len()
	var mat matrix = matrix{v: bits.NewBit(uint(len1 * len2)), lenX: len1, lenY: len2}

	for i := 0; i < len1; i++ {
		for j := 0; j < len2; j++ {
			mat.v.Poke(mat.at(i, j), data.Equal(i, j))
		}
	}

	return recursiveDiff(box{0, 0, len1, len2}, mat)
}

type match struct {
	x, y   int
	length int
}

type box struct {
	x, y       int
	lenX, lenY int
}

// A helper structure that stores absolute dimension along a linear bit vector
// so that it can always properly translate (x, y) -> z on the vector
type matrix struct {
	v          bits.Vector
	lenX, lenY int
}

// Translates (x, y) to an absolute position on the bit vector
func (m *matrix) at(x, y int) uint {
	return uint(y + (x * m.lenY))
}

func recursiveDiff(bounds box, mat matrix) Delta {
	var m match = largest(bounds, mat)

	if m.length == 0 { // Recursion terminates
		var immediate Delta
		if bounds.lenY-bounds.y > 0 {
			adds := make([]int, bounds.lenY - bounds.y)
			for i := 0; i < len(adds); i++ {
				adds[i] = bounds.y+i
			}
			immediate.Added = [][]int{adds}
		}
		if bounds.lenX-bounds.x > 0 {
			removes := make([]int, bounds.lenX - bounds.x)
			for i := 0; i < len(removes); i++ {
				removes[i] = bounds.x+i
			}
			immediate.Removed = [][]int{removes}
		}
		return immediate
	}

	var left Delta = recursiveDiff(box{bounds.x, bounds.y, m.x, m.y}, mat)
	var right Delta = recursiveDiff(box{m.x + m.length, m.y + m.length, bounds.lenX, bounds.lenY}, mat)

	var result Delta

	result.Added = append(left.Added, right.Added...)
	result.Removed = append(left.Removed, right.Removed...)

	return result
}

// Finds the largest common substring by looking at the provided match matrix
// starting from (bounds.x, bounds.y) with lengths bounds.lenX, bounds.lenY
func largest(bounds box, mat matrix) match {
	var result match

	// Look for LCS in the too-right half, including the main diagonal
	for i := bounds.x; i < bounds.lenX && result.length < (bounds.lenX-i); i++ {
		var m match = search(i, bounds.y, bounds.lenX, bounds.lenY, mat)
		if m.length > result.length {
			result = m
		}
	}

	// Look for LCS in the bottom-left half, excluding the main diagonal
	for j := bounds.y + 1; j < bounds.lenY && result.length < (bounds.lenY-j); j++ {
		var m match = search(bounds.x, j, bounds.lenX, bounds.lenY, mat)
		if m.length > result.length {
			result = m
		}
	}
	return result
}

// Searches the main diagonal for the longest sequential match line
func search(x, y, lenX, lenY int, mat matrix) (result match) {
	var inMatch bool
	var m match
	for step := 0; step+x < lenX && step+y < lenY; step++ {
		if mat.v.Peek(mat.at(step+x, step+y)) {
			if !inMatch { // Create a new current record if there is none ...
				inMatch, m.x, m.y, m.length = true, step+x, step+y, 1
			} else { // ... otherwise just increment the existing
				m.length++
			}

			if m.length > result.length {
				result = m // Store it if it is longer ...
			}
		} else { // End of current of match
			inMatch = false // ... and reset the current one
		}
	}
	return
}

// A diff.Interface implementation with plugable Equal function
type D struct {
	Len1, Len2 int
	EqualFunc  func(i, j int) bool
}

// Required per diff.Interface
func (d D) Len() (int, int) {
	return d.Len1, d.Len2
}

// Required per diff.Interface
func (d D) Equal(i, j int) bool {
	return d.EqualFunc(i, j)
}
