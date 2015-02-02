package fibonacci // import "github.com/solarsea/misc/fibonacci"

import (
	"bytes"
	diff "github.com/solarsea/misc/diff"
	"io"
	"strings"
	"testing"
)

func TestNumbers(t *testing.T) {
	n := New(32)

	expected := []uint64{1, 1, 2, 3, 5, 8, 13, 21}

	for i, v := range expected {

		if f := n[i]; f != v {
			t.Error("Unexpected value for", i, f, "expected", v)
		}
	}
}

func TestCoding(t *testing.T) {
	n := New(32)

	for i := uint64(0); i < 4096; i++ {
		enc, encLen := n.Code(i)
		dec, decLen := n.Decode(enc)

		if i != dec {
			t.Errorf("Unexpected value for %d - enc is %b, dec is %d\n", i, enc, dec)
		}
		if encLen != decLen {
			t.Errorf("Unexpected difference between encoded and decoded lengths.", encLen, decLen)
		}
	}
}

func TestWriterReader(t *testing.T) {
	var (
		buf   bytes.Buffer
		w     io.Writer = Encoder(&buf)
		input []byte    = make([]byte, 256)
		fib   Numbers   = New(16)
	)

	for i := uint64(0); i < 256; i++ {
		input[i] = byte(i)
	}

	// Write the input
	count, err := w.Write(input)
	if count != len(input) {
		t.Error("Unexpected write count", count)
	}
	if err != nil {
		t.Error("Unexpected write error", err.Error())
	}

	// Flush remaining bits
	count, err = w.Write(nil)
	if count != 0 {
		t.Error("Unexpected write count while flushing", count)
	}
	if err != nil {
		t.Error("Unexpected write error while flushing", err.Error())
	}

	var output string
	for _, v := range buf.Bytes() {
		output += u2s(uint64(v), 8)
	}

	for i, v := range input {
		c, l := fib.Code(uint64(v))
		vs := u2s(c, l)
		if loc := strings.Index(output, vs); loc != 0 {
			t.Fatal("Unexpected location for", i, "value", vs)
		}
		output = output[len(vs):]
	}

	var (
		in  *bytes.Reader = bytes.NewReader(buf.Bytes())
		r   io.Reader     = Decoder(in)
		out bytes.Buffer
	)
	io.Copy(&out, r)
	decoded := out.Bytes()

	delta := diff.Diff(diff.D{Len1: len(decoded), Len2: len(input), EqualFunc: func(i, j int) bool {
		return decoded[i] == input[j]
	}})

	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Differences detected ", delta)
	}
}

func u2s(b uint64, l byte) (result string) {
	for i := byte(0); i < l; i++ {
		if (b & 1) > 0 {
			result += "1"
		} else {
			result += "0"
		}
		b >>= 1
	}
	return
}
