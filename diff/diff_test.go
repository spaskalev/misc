package diff // import "github.com/solarsea/misc/diff"

import (
	"testing"
)

// A diff.Interface implementation for testing purposes
type diffSlice struct {
	first  []rune
	second []rune
}

// Required per diff.Interface
func (d diffSlice) Len() (int, int) {
	return len(d.first), len(d.second)
}

// Required per diff.Interface
func (d diffSlice) Equal(i, j int) bool {
	return d.first[i] == d.second[j]
}

func TestDiff(t *testing.T) {
	data := diffSlice{
		[]rune("abcdefgh"),
		[]rune("abbcedfh"),
	}

	result := Diff(data)
	if len(result.Added) != 2 ||
		result.Added[0].From != 2 || result.Added[0].Length != 3 ||
		result.Added[1].From != 5 || result.Added[1].Length != 6 ||
		result.Removed[0].From != 3 || result.Removed[0].Length != 4 ||
		result.Removed[1].From != 6 || result.Removed[1].Length != 7 {
		t.Error("Unexpected diff results", result)
	}
}
