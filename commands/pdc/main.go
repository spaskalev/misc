package main

import (
	"fmt"
	iou "github.com/spaskalev/misc/ioutil"
	predictor "github.com/spaskalev/misc/predictor"
	"io"
	"os"
)

func main() {
	var code int
	switch {
	case len(os.Args) == 1:
		code = compress(os.Stdout, os.Stdin)
	case len(os.Args) == 2 && os.Args[1] == "-d":
		code = decompress(os.Stdout, os.Stdin)
	default:
		fmt.Fprintln(os.Stdout, "Usage: pdc [-d]")
	}
	os.Exit(code)
}

// Compress the data from the given io.Reader and write it to the given io.Writer
// I/O is buffered for better performance
func compress(output io.Writer, input io.Reader) int {
	var (
		err        error
		buffer     io.Writer = iou.SizedWriter(output, 4096)
		compressor io.Writer = predictor.Compressor(buffer)
	)

	_, err = io.Copy(compressor, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while compressing.\n", err)
		return 1
	}

	// Flush the compressor
	_, err = compressor.Write(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while flushing compresssor buffer.\n", err)
		return 1
	}

	// Flush the buffer
	_, err = buffer.Write(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while flushing output buffer.\n", err)
		return 1
	}

	return 0
}

// Decompress the data from the given io.Reader and write it to the given io.Writer
// I/O is buffered for better performance
func decompress(output io.Writer, input io.Reader) int {
	var (
		err          error
		decompressor io.Reader = predictor.Decompressor(iou.SizedReader(input, 4096))
	)

	_, err = io.Copy(output, decompressor)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while decompressing.\n", err)
		return 1
	}

	return 0
}
