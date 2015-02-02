package diff // import "github.com/solarsea/misc/diff"

import (
	_ "0dev.org/bits"
	_ "fmt"
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

// func TestWTF(t *testing.T) {
// 	data := diffSlice{
// 		[]rune("abcdefgh"),
// 		[]rune("abbcedfh"),
// 	}

// 	var len1, len2 int = data.Len()
// 	var mat matrix = matrix{v: bits.NewBool(uint(len1 * len2)), lenX: len1, lenY: len2}

// 	for i := 0; i < len1; i++ {
// 		for j := 0; j < len2; j++ {
// 			mat.v.Poke(mat.at(i, j), data.Equal(i, j))
// 		}
// 	}

// 	debugPrint(box{5, 5, 6, 6}, mat) // visual debugging as its finest
// 	debugPrint(box{5, 5, 7, 7}, mat)
// 	debugPrint(box{5, 5, 8, 8}, mat) // ZO RELAXEN UND WATSCHEN DER BLINKENLICHTEN.
// }

// func debugPrint(bounds box, mat matrix) {
// 	// Debug print
// 	fmt.Printf("-%d-%d--%d-%d--\n", bounds.x, bounds.y, bounds.lenX, bounds.lenY)
// 	for i := bounds.x; i < bounds.lenX; i++ {
// 		fmt.Print("| ")
// 		for j := bounds.y; j < bounds.lenY; j++ {
// 			//if vector.Peek(uint(j + (i * bounds.lenY))) {
// 			if mat.v.Peek(mat.at(i, j)) {
// 				fmt.Print("\\")
// 			} else {
// 				fmt.Print(".")
// 			}
// 		}
// 		fmt.Println(" |")
// 	}
// 	fmt.Println("------------")
// }
