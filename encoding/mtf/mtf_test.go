package mtf // import "github.com/solarsea/misc/encoding/mtf"

import (
	"bytes"
	diff "github.com/solarsea/diff"
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

	delta := diff.Diff(diff.WithEqual(len(data), len(processed),
		func(i, j int) bool { return data[i] == processed[j] }))
	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Differences detected ", delta, processed)
	}
}
