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
		{"", "a", Delta{Added: []Mark{Mark{0, 1}}}},
		{"a", "", Delta{Removed: []Mark{Mark{0, 1}}}},
		//
		{"a", "a", Delta{}},
		{"a", "aa", Delta{Added: []Mark{Mark{1, 2}}}},
		{"aa", "a", Delta{Removed: []Mark{Mark{1, 2}}}},
		//
		{"abcdefgh", "abbcedfh", Delta{
			Added:   []Mark{Mark{2, 3}, Mark{5, 6}},
			Removed: []Mark{Mark{3, 4}, Mark{6, 7}},
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
