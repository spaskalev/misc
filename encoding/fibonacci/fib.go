// Package provides a shifted fibonacci encoding of unsigned integers.
//
// http://en.wikipedia.org/wiki/Fibonacci_coding maps positive integers as
// 1 - 11, 2 - 011, 3 - 0011, 4 - 1011, 5 - 00011
//
// Incrementing input by one to allow for zero gives
// 0 - 11, 1 - 011, 2 - 0011, 3 - 1011, 4 - 00011
//
// The codes are reversed so that they are easily stored in uints,
// effectively avoiding the need to store the number of leading zeroes
// 0 - 11, 1 - 110, 2 - 1100, 3 - 1101, 4 - 11000
package fibonacci // import "github.com/solarsea/misc/encoding/fibonacci"

import (
	"io"
)

// Alias type with methods for encoding and decoding integers
type Numbers []uint64

var (
	// Used for decoding byte values
	codec Numbers
	// Used for encoding byte values
	// The lower 16 bits store the encoded value itself
	// while the remaining upper ones store its length
	lookup [256]uint32
)

func init() {
	codec = New(16)
	for i := uint64(0); i < 256; i++ {
		val, len := codec.Code(i)
		lookup[i] |= uint32(val)
		lookup[i] |= uint32(len) << 16
	}
}

// Returns a slice with fibonacci numbers up to the given length
func New(size int) Numbers {
	var fibs Numbers = make(Numbers, size)
	copy(fibs, []uint64{1, 1})
	for i := 2; i < size; i++ {
		fibs[i] = fibs[i-1] + fibs[i-2]
	}
	return fibs
}

// Returns a fibonacci code for an integer as specified in the package's doc.
func (f Numbers) Code(value uint64) (result uint64, length byte) {
	// Increment to encode zero as one
	value++

	// Find the nearest fibonacci number
	for f[length] <= value {
		length++
	}

	// Leading bit that signals the start of a fibonacci-encoded integer
	result |= 1

	// Find the Zeckendorf's representation by raising a bit for each
	// fibonacci number that is less or equal to the difference
	// between the value and the previous such number
	for i := length - 1; i >= 1; i-- {
		result <<= 1
		if f[i] <= value {
			result |= 1
			value -= f[i]
		}
	}
	return
}

// Returns an integer from a fibonacci code as specified in the package's doc.
func (f Numbers) Decode(value uint64) (result uint64, length byte) {
	length = 1
	// Loop until the lowest two bits are both raised
	for (value & 3) != 3 {
		// Add the fibonacci number for the current bit if it is raised
		if (value & 1) == 1 {
			result += f[length]

			// We know that the next bit cannot be raised by Zeckendorf's theorem
			value >>= 2
			length += 2
			continue
		}

		value >>= 1
		length++
	}
	return result + f[length] - 1, length + 1
}

// Returns a fibonacci encoder over the provided io.Writer
func Encoder(target io.Writer) io.Writer {
	var enc encoder
	enc.target = target
	return &enc
}

type encoder struct {
	target    io.Writer
	buffer    [2]byte
	remaining byte
	length    byte
}

// Implements io.Writer
func (e *encoder) Write(input []byte) (int, error) {
	var (
		total int
		err   error
	)

	// Flush on a nil slice
	if input == nil {
		_, err = e.target.Write([]byte{byte(e.remaining)})
		return 0, err
	}

	for _, currentByte := range input {
		// Get the fibonacci code and bit length for the current byte
		enc, len := uint16(lookup[currentByte]), byte(lookup[currentByte]>>16)

		// Add current bits to higher positions
		e.remaining |= byte(enc << e.length)

		// maximum length of added bits to e.remaining
		added := 8 - e.length

		// Shift the the encoded value and account for its length
		enc >>= added
		e.length += len
		len -= added

		// Not enough bits to write
		if e.length < 8 {
			// Account for the processed input byte
			total++

			continue
		}

		// Clearing e.length is not necessary as it will be overwritten later

		// Stage the complete byte for writing
		buffer := e.buffer[:1]
		buffer[0] = byte(e.remaining)

		// Stage every full byte from the encoded value for writing
		//
		// The bitlength of the largest encoded byte value, 255, is 13.
		// Even with 7 bits already in the buffer this leaves [7+1], [8]
		// and 4 bits remaining => a single if is enough instead of a for.
		//
		// 128 is [1000 0000] in binary. Any value equal or greater than it
		// will be atleast 8 bits in length
		if enc >= 128 {
			buffer = append(buffer, byte(enc))
			enc >>= 8
			len -= 8
		}

		// Store the remaining bits
		e.remaining, e.length = byte(enc), len

		// Write the staged bytes
		_, err = e.target.Write(buffer)

		// Abort write on error
		if err != nil {
			break
		}

		// Account for the processed input byte
		total++
	}
	return total, err
}

// Returns a fibonacci decoder over the provided io.Reader
func Decoder(source io.Reader) io.Reader {
	var dec decoder
	dec.source = source
	return &dec
}

type decoder struct {
	source io.Reader
	buffer uint64
	at     byte
}

// Implements io.Reader
func (d *decoder) Read(output []byte) (int, error) {
	var (
		total int
		err   error
	)

start:
	// While we have suitable buffered data and enough output space
	for (len(output) > 0) && ((d.buffer & (d.buffer >> 1)) > 0) {
		val, len := codec.Decode(d.buffer)

		// Store the decoded byte
		output[0] = byte(val)

		// Advance the internal and output buffers
		output = output[1:]
		d.buffer >>= len
		d.at -= len

		// Account for the processed output byte
		total++
	}

	// Termination condition
	if len(output) == 0 || err != nil {
		return total, err
	}

	// We need to limit the output's size else we could end up with a lot of small values
	// that fit neither in the output slice nor in the internal buffer
	//
	// (63 is [0011 1111] in binary, xor is a substraction and right shift a division)
	free := int((63 ^ d.at) >> 3)
	if free > len(output) {
		free = len(output)
	}

	// Read data and transfer to the internal buffer
	count, err := d.source.Read(output[:free])
	for _, v := range output[:count] {
		d.buffer |= uint64(v) << d.at
		d.at += 8
	}

	// To ensure a tail call :)
	goto start
}
