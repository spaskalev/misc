package main

import (
	"bufio"
	"fmt"
	diff "github.com/solarsea/diff"
	"hash"
	"hash/fnv"
	"os"
	"path/filepath"
)

const usage = "Usage: plaindiff <file1> <file2>\n"

func main() {
	if len(os.Args) <= 3 {
		os.Stderr.WriteString(usage)
		os.Exit(1)
	}

	var hd hashDiff
	hd.hash = fnv.New64a()
	hd.first = read(os.Args[2], hd.hash)
	hd.second = read(os.Args[1], hd.hash)

	var result diff.Delta = diff.Diff(diff.WithEqual(len(hd.first), len(hd.second), hd.Equal))

	gen := source(result)
	out := bufio.NewWriter(os.Stdout)
	for have, added, mark := gen(); have; have, added, mark = gen() {
		var from []line = hd.line(!added)

		fmt.Fprintln(out)
		for i := mark.From; i < mark.Length; i++ {
			fmt.Fprint(out, i+1) // Line numbers start from 1 for most people :)
			if added {
				fmt.Fprint(out, " > ")
			} else {
				fmt.Fprint(out, " < ")
			}
			fmt.Fprintln(out, from[i].text)
		}
	}
	out.Flush()
}

// Returns a closure over the provided diff.Delta
// that returns diff.Mark entries while prioritizing removals when possible
func source(d diff.Delta) func() (bool, bool, diff.Mark) {
	var addedAt, removedAt int = 0, 0
	return func() (bool, bool, diff.Mark) {
		var addsOver bool = addedAt == len(d.Added)
		var removesOver bool = removedAt == len(d.Removed)

		var add, remove diff.Mark

		// Check whether both mark slices have been exhausted
		if addsOver && removesOver {
			return false, false, diff.Mark{}
		}

		// Return an add if removes are over
		if removesOver {
			add = d.Added[addedAt]
			addedAt++
			return true, true, add
		}

		// Return a remove if the adds are over
		if addsOver {
			remove = d.Removed[removedAt]
			removedAt++
			return true, false, remove
		}

		add = d.Added[addedAt]
		remove = d.Removed[removedAt]

		// Prioritize a remove if it happens before an add
		if remove.From <= add.From {
			removedAt++
			return true, false, remove
		}

		// Else
		addedAt++
		return true, true, add

	}
}

// A line-based diff.Interface implementation
type hashDiff struct {
	first, second []line
	hash          hash.Hash64
}

// Required per diff.Interface
func (h *hashDiff) Equal(i, j int) bool {
	if h.first[i].hash != h.second[j].hash {
		return false
	}
	return h.first[i].text == h.second[j].text
}

// A helper method for getting a line slice
func (h *hashDiff) line(first bool) []line {
	if first {
		return h.first
	}
	return h.second
}

// Holds a text line and its hash
type line struct {
	hash uint64
	text string
}

// Reads all lines in a file and returns a line entry for each
func read(name string, h hash.Hash64) []line {
	abs, err := filepath.Abs(name)
	fatal(err)
	f, err := os.Open(abs)
	fatal(err)
	scanner := bufio.NewScanner(f)
	result := make([]line, 0)
	for scanner.Scan() {
		h.Reset()
		h.Write(scanner.Bytes())
		result = append(result, line{hash: h.Sum64(), text: scanner.Text()})
	}
	fatal(f.Close())
	return result
}

// Write an error to stderr and exit
func fatal(e error) {
	if e != nil {
		os.Stderr.WriteString(e.Error())
		os.Exit(1)
	}
}
