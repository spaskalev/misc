package diff // import "github.com/solarsea/misc/diff"

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	data := []struct {
		seq1  string
		seq2  string
		delta Delta
	}{
		{"", "", Delta{}},
		{"", "a", Delta{Added: [][]int{[]int{0}}}},
		{"a", "", Delta{Removed: [][]int{[]int{0}}}},
		//
		{"a", "a", Delta{}},
		{"a", "aa", Delta{Added: [][]int{[]int{1}}}},
		{"aa", "a", Delta{Removed: [][]int{[]int{1}}}},
		//
		{"abcdefgh", "abbcedfh", Delta{
			Added:   [][]int{[]int{2}, []int{5}},
			Removed: [][]int{[]int{3}, []int{6}},
		}},
	}

	for _, testCase := range data {
		delta := Diff(D{len(testCase.seq1), len(testCase.seq2), func(i, j int) bool {
			return testCase.seq1[i] == testCase.seq2[j]
		}})

		if fmt.Sprintf("%v", delta) != fmt.Sprintf("%v", testCase.delta) {
			t.Errorf("Unexpected delta for data\n[%s]\n[%s]\nGot %v\nExpected %v",
				testCase.seq1, testCase.seq2, delta, testCase.delta)
		}
	}
}
