package mtf

import (
	diff "0dev.org/diff"
	"bytes"
	"io"
	"testing"
)

func TestMTF(t *testing.T) {
	var (
		data []byte = []byte{1, 1, 0, 0}

		input   *bytes.Reader = bytes.NewReader(data)
		encoder io.Reader     = Encoder(input)
		decoder io.Reader     = Decoder(encoder)

		output bytes.Buffer
	)

	io.Copy(&output, decoder)
	processed := output.Bytes()

	delta := diff.Diff(diff.D{Len1: len(data), Len2: len(processed),
		EqualFunc: func(i, j int) bool { return data[i] == processed[j] }})
	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Differences detected ", delta, processed)
	}
}
