package predictor // import "github.com/spaskalev/misc/predictor"

import (
	"bytes"
	"fmt"
	diff "github.com/spaskalev/diff"
	"io/ioutil"
	"testing"
)

// Sample input from RFC1978 - PPP Predictor Compression Protocol
var input = []byte{0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0a,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0a,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0a,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0a,
	0x41, 0x42, 0x41, 0x42, 0x41, 0x42, 0x41, 0x0a,
	0x42, 0x41, 0x42, 0x41, 0x42, 0x41, 0x42, 0x0a,
	0x78, 0x78, 0x78, 0x78, 0x78, 0x78, 0x78, 0x0a}

// Sample output from RFC1978 - PPP Predictor Compression Protocol
var output = []byte{0x60, 0x41, 0x41, 0x41, 0x41, 0x41, 0x0a, 0x60,
	0x41, 0x41, 0x41, 0x41, 0x41, 0x0a, 0x6f, 0x41,
	0x0a, 0x6f, 0x41, 0x0a, 0x41, 0x42, 0x41, 0x42,
	0x41, 0x42, 0x0a, 0x60, 0x42, 0x41, 0x42, 0x41,
	0x42, 0x0a, 0x60, 0x78, 0x78, 0x78, 0x78, 0x78, 0x0a}

func TestCompressorSample(t *testing.T) {
	var (
		buf bytes.Buffer
		err error
	)

	out := Compressor(&buf)
	_, err = out.Write(input)
	if err != nil {
		t.Error(err)
	}

	_, err = out.Write(nil)
	if err != nil {
		t.Error(err)
	}

	result := buf.Bytes()
	delta := diff.Diff(diff.WithEqual(len(result), len(output), func(i, j int) bool { return result[i] == output[j] }))

	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Unexpected compressed output", delta)
	}
}

func TestDecompressorSample(t *testing.T) {
	in := Decompressor(bytes.NewReader(output))
	result, err := ioutil.ReadAll(in)
	if err != nil {
		t.Error("Unexpected error while decompressing", err)
	}

	delta := diff.Diff(diff.WithEqual(len(result), len(input),
		func(i, j int) bool { return result[i] == input[j] }))

	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Unexpected decompressed output", delta)
	}
}

var testData = [][]byte{
	[]byte{},
	[]byte{0, 1, 2, 3},
	[]byte{0, 1, 2, 3, 4, 5, 6, 7},
	[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
}

func TestCycle(t *testing.T) {
	for i := 0; i < len(testData); i++ {
		if err := cycle(testData[i], len(testData[i])); err != nil {
			t.Error(err)
		}
	}
}

func TestStepCycle(t *testing.T) {
	for i := 0; i < len(testData); i++ {
		for j := 1; j < len(testData[i]); j++ {
			if err := cycle(testData[i], j); err != nil {
				t.Error("Error for testData[", i, "], step[", j, "] ", err)
			}
		}
	}
}

func cycle(input []byte, step int) error {
	var (
		buf bytes.Buffer
		err error
	)

	if step > len(input) {
		return nil
	}

	// Create a compressor and write the given data
	compressor := Compressor(&buf)

	var data []byte = input
	var trace []byte = make([]byte, 0)

	for len(data) > 0 {
		if step <= len(data) {

			trace = append(trace, data[:step]...)

			_, err = compressor.Write(data[:step])
			if err != nil {
				return err
			}

			data = data[step:]
		} else {
			step = len(data)
		}
	}

	// Flush the compressor
	_, err = compressor.Write(nil)
	if err != nil {
		return err
	}

	// Attempt to decompress the data
	compressed := buf.Bytes()
	decompressed, err := ioutil.ReadAll(Decompressor(bytes.NewReader(compressed)))
	if err != nil {
		return err
	}

	// Diff the result against the initial input
	delta := diff.Diff(diff.WithEqual(len(input), len(decompressed),
		func(i, j int) bool { return input[i] == decompressed[j] }))

	// Return a well-formated error if any differences are found
	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		return fmt.Errorf("Unexpected decompressed output for step %d, delta %v\ninput:  (%d) %#x\ntrace:  (%d) %#x\noutput: (%d) %#x\n",
			step, delta, len(input), input, len(trace), trace, len(decompressed), decompressed)
	}

	// All is good :)
	return nil
}
