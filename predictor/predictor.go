// Package predictor implements the predictor compression/decompression algorithm
// as specified by RFC1978 - PPP Predictor Compression Protocol
package predictor // import "github.com/spaskalev/misc/predictor"

import (
	bits "github.com/spaskalev/bits"
	iou "github.com/spaskalev/misc/ioutil"
	"io"
)

// The context struct contains the predictor's algorithm guess table
// and the current value of its input/output hash
type context struct {
	table [1 << 16]byte
	hash  uint16
}

// The following hash code is the heart of the algorithm:
// It builds a sliding hash sum of the previous 3-and-a-bit
// characters which will be used to index the guess table.
// A better hash function would result in additional compression,
// at the expense of time.
func (ctx *context) update(val byte) {
	ctx.hash = (ctx.hash << 4) ^ uint16(val)
}

// Returns an io.Writer implementation that wraps the provided io.Writer
// and compresses data according to the predictor algorithm
//
// It can buffer data as the predictor mandates 8-byte blocks with a header.
// A call with no data will force a flush.
func Compressor(writer io.Writer) io.Writer {
	var cmp compressor
	cmp.target = writer
	return iou.SizedWriter(&cmp, 8)
}

type compressor struct {
	context
	target io.Writer
}

// Note: this method does not implement the full io.Writer's Write() semantics
func (ctx *compressor) Write(data []byte) (int, error) {
	var (
		blockSize  int = 8
		datalength int = len(data)
	)

	if datalength == 0 {
		return 0, nil
	}

	if datalength < blockSize {
		blockSize = datalength
	}

	var buf []byte = make([]byte, 1, blockSize+1)
	for block := 0; block < datalength/blockSize; block++ {
		for i := 0; i < blockSize; i++ {
			var current byte = data[(block*blockSize)+i]
			if ctx.table[ctx.hash] == current {
				// Guess was right - don't output
				buf[0] |= 1 << uint(i)
			} else {
				// Guess was wrong, output char
				ctx.table[ctx.hash] = current
				buf = append(buf, current)
			}
			ctx.update(current)
		}

		if c, err := ctx.target.Write(buf); err != nil {
			return (block * blockSize) + c, err
		}

		// Reset the flags and buffer for the next iteration
		buf, buf[0] = buf[:1], 0
	}

	return datalength, nil
}

// Returns an io.Reader implementation that wraps the provided io.Reader
// and decompresses data according to the predictor algorithm
func Decompressor(reader io.Reader) io.Reader {
	var dcmp decompressor
	dcmp.source = reader
	return iou.SizedReader(&dcmp, 8)
}

type decompressor struct {
	context
	source io.Reader
}

// Note: this method does not implement the full io.Reader's Read() semantics
func (ctx *decompressor) Read(output []byte) (int, error) {
	var (
		err                          error
		buffer                       []byte = make([]byte, 8)
		flags                        byte
		predicted, rc, total, copied int
	)

	// Read the next prediction header
readHeader:
	rc, err = ctx.source.Read(buffer[:1])
	// Fail on error unless it is EOF
	if err != nil && err != io.EOF {
		return total, err
	} else if rc == 0 {
		return total, err
	}

	// Copy the prediction header and calculate the number of subsequent bytes to read
	flags = buffer[0]
	predicted = bits.Hamming(flags)

	// Read the non-predicted bytes and place them in the end of the buffer
	rc, err = ctx.source.Read(buffer[predicted:])
retryData:
	if (rc < (8 - predicted)) && err == nil {
		// Retry the read if we have fewer bytes than what the prediction header indicates
		var r int
		r, err = ctx.source.Read(buffer[predicted+rc:])
		rc += r
		goto retryData
	} // Continue on any error, try to decompress and return it along the result

	// rc now contains the amount of actual bytes in this cycle (usually 8)
	rc += predicted

	// Walk the buffer, filling in the predicted blanks,
	// relocating read bytes and and updating the guess table
	for i, a := 0, predicted; i < rc; i++ {
		if (flags & (1 << uint(i))) > 0 {
			// Guess succeeded, fill in from the table
			buffer[i] = ctx.table[ctx.hash]
		} else {
			// Relocate a read byte and advance the read byte index
			buffer[i], a = buffer[a], a+1
			// Guess failed, update the table
			ctx.table[ctx.hash] = buffer[i]
		}
		// Update the hash
		ctx.update(buffer[i])
	}

	// Copy the decompressed data to the output and accumulate the count
	copied = copy(output, buffer[:rc])
	total += copied

	// Loop for another pass if there is available space in the output
	output = output[copied:]
	if len(output) > 0 && err == nil {
		goto readHeader
	}

	return total, err
}
