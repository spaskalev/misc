// purely a prototype :)
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
)

var (
	pt = rand.Perm(256)

	streams = []string{
		"\x1b[1;37;40m",
		"\x1b[1;37;4;40m",

		"\x1b[1;32;41m",
		"\x1b[1;33;4;41m",

		"\x1b[1;32;42m",
		"\x1b[34;4;42m",

		"\x1b[1;33;43m",
		"\x1b[1;34;4;43m",

		"\x1b[1;33;44m",
		"\x1b[1;36;4;44m",

		"\x1b[1;37;45m",
		"\x1b[30;4;45m",

		"\x1b[1;32;46m",
		"\x1b[1;34;4;46m",

		"\x1b[1;37;47m",
		"\x1b[1;34;4;47m",
	}
)

func pearsonHash(content []byte) (hash byte) {
	for _, c := range content {
		hash = byte(pt[hash^c])
	}
	return hash
}

func hashStrings(data []string) []byte {
	var result = make([]byte, len(data))
	for i, s := range data {
		result[i] = pearsonHash([]byte(s))
	}
	return result
}

func hamming(a byte, b byte) (result byte) {
	for i := uint(0); i < 8; i++ {
		if (a & (1 << i)) == (b & (1 << i)) {
			result++
		}
	}
	return result
}

func main() {
	// TODO there could be better tokenizers
	tokenExp := regexp.MustCompile("[^\\s]+")

	reader := bufio.NewReader(os.Stdin)

	// keep the last XXX hashes in a ring
	var last [16]byte
	var pos int

	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		tokens := tokenExp.FindAllString(string(line), -1)

		if len(tokens) <= 1 {
			continue
		}

		// cut the tokens to the first X ...
		if len(tokens) > 16 {
			tokens = tokens[:16]
		}

		// Cut timestamp
		tokens = tokens[1:]

		hashes := hashStrings(tokens)

		// TODO use position-based weights
		var counter [8]int
		for _, hash := range hashes {
			for i := uint(0); i < 8; i++ {
				if (hash & (1 << i)) > 0 {
					counter[i]++
				} else {
					counter[i]--
				}
			}
		}

		var hash byte
		for i := 0; i < len(counter); i++ {
			if counter[i] >= 0 {
				hash |= 1
			}
			hash <<= 1
		}

		closest := pos
		var minham byte = 9 // lol
		a := last[:pos]
		for i := len(a) - 1; i >= 0; i-- {
			ham := hamming(hash, a[i])
			if ham < minham {
				minham = ham
				closest = i
			}
		}
		b := last[pos:]
		for i := len(b) - 1; i >= 0; i-- {
			ham := hamming(hash, b[i])
			if ham < minham {
				minham = ham
				closest = i
			}
		}

		fmt.Println(streams[closest], string(line))

		// the wheel spins
		last[pos] = hash
		pos++
		pos %= 16
	}
}
