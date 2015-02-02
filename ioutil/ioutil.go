// Package ioutil contains various constructs for io operations.
package ioutil // import "github.com/solarsea/misc/ioutil"

import (
	"io"
)

// An function alias type that implements io.Writer.
type WriterFunc func([]byte) (int, error)

// Delegates the call to the WriterFunc while implementing io.Writer.
func (w WriterFunc) Write(b []byte) (int, error) {
	return w(b)
}

// An function alias type that implements io.Reader.
type ReaderFunc func([]byte) (int, error)

// Delegates the call to the WriterFunc while implementing io.Reader.
func (r ReaderFunc) Read(b []byte) (int, error) {
	return r(b)
}

// Reads a single byte from the provided io.Reader
func ReadByte(reader io.Reader) (byte, error) {
	var (
		arr [1]byte
		err error
	)
	_, err = reader.Read(arr[:])
	return arr[0], err
}

// Returns a writer that delegates calls to Write(...) while ensuring
// that it is never called with less bytes than the specified amount.
//
// Calls with fewer bytes are buffered while a call with a nil slice
// causes the buffer to be flushed to the underlying writer.
func SizedWriter(writer io.Writer, size int) io.Writer {
	var sw sizedWriter
	sw.writer = writer
	sw.buffer = make([]byte, 0, size)
	sw.size = size
	return &sw
}

type sizedWriter struct {
	writer io.Writer
	buffer []byte
	size   int
}

func (sw *sizedWriter) Write(input []byte) (int, error) {
	var (
		count int
		err   error
	)

	// Flush the buffer when called with no bytes to write
	if input == nil {
		// Call the writer with whatever we have in store..
		count, err = sw.writer.Write(sw.buffer)

		// Advance the buffer
		sw.buffer = sw.buffer[:copy(sw.buffer, sw.buffer[count:])]

		return 0, err
	}

	// Delegate to the writer if the size is right
	if len(sw.buffer) == 0 && len(input) >= sw.size {
		reduced := (len(input) / sw.size) * sw.size
		count, err = sw.writer.Write(input[:reduced])
		if count < reduced || err != nil {
			return count, err
		}

		// Stage any remaining data in the buffer
		sw.buffer = append(sw.buffer, input[count:]...)
		return len(input), nil
	}

	// Append data to the buffer
	count = copy(sw.buffer[len(sw.buffer):sw.size], input)
	sw.buffer = sw.buffer[:len(sw.buffer)+count]

	// Return if we don't have enough bytes to write
	if len(sw.buffer) < sw.size {
		return len(input), nil
	}

	// Flush the buffer as it is filled
	_, err = sw.Write(nil)
	if err != nil {
		return count, err
	}

	// Handle the rest of the input
	return sw.Write(input[count:])
}

// Returns a reader that delegates calls to Read(...) while ensuring
// that the output buffer is never smaller than the required size
// and is downsized to a multiple of the required size if larger.
func SizedReader(reader io.Reader, size int) io.Reader {
	var sr sizedReader
	sr.reader = reader
	sr.buffer = make([]byte, size)
	sr.size, sr.from, sr.to = size, 0, 0
	return &sr
}

type sizedReader struct {
	reader         io.Reader
	buffer         []byte
	from, to, size int
}

func (sr *sizedReader) Read(output []byte) (int, error) {
	var (
		count int
		err   error
	)

start:
	// Reply with the buffered data if there is any
	if sr.to > 0 {
		count = copy(output, sr.buffer[sr.from:sr.to])

		// Advance the data in the buffer
		sr.from += count

		// Check whether we have reached the end of the buffer
		if sr.from == sr.to {
			// Reset the buffer
			sr.from, sr.to = 0, 0

			return count, err
		}

		// Do not propagate an error until the buffer is exhausted
		return count, nil
	}

	// Delegate if the buffer is empty and the destination buffer is large enough
	if len(output) >= sr.size {
		return sr.reader.Read(output[:(len(output)/sr.size)*sr.size])
	}

	// Perform a read into the buffer
	count, err = sr.reader.Read(sr.buffer)

	// Size the buffer down to the read data size
	// and restart if we have successfully read some bytes
	sr.from, sr.to = 0, count
	if sr.to > 0 {
		goto start
	}

	// Returning on err/misbehaving noop reader
	return 0, err
}
