package ioutil // import "github.com/spaskalev/misc/ioutil"

import (
	"bytes"
	"errors"
	diff "github.com/solarsea/diff"
	"io"
	"testing"
)

func TestWriterFunc(t *testing.T) {
	var (
		input  []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7}
		output []byte

		reader *bytes.Reader = bytes.NewReader(input)
		buffer bytes.Buffer
	)

	reader.WriteTo(WriterFunc(buffer.Write))
	output = buffer.Bytes()

	// Diff the result against the initial input
	delta := diff.Diff(diff.WithEqual(len(input), len(output),
		func(i, j int) bool { return input[i] == output[j] }))
	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Differences detected ", delta)
	}
}

func TestReaderFunc(t *testing.T) {
	var (
		input  []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7}
		output []byte

		reader *bytes.Reader = bytes.NewReader(input)
		buffer bytes.Buffer
	)

	buffer.ReadFrom(ReaderFunc(reader.Read))
	output = buffer.Bytes()

	// Diff the result against the initial input
	delta := diff.Diff(diff.WithEqual(len(input), len(output),
		func(i, j int) bool { return input[i] == output[j] }))
	if len(delta.Added) > 0 || len(delta.Removed) > 0 {
		t.Error("Differences detected ", delta)
	}
}

func TestReadByte(t *testing.T) {
	var (
		input  []byte        = []byte{255}
		reader *bytes.Reader = bytes.NewReader(input)
	)

	result, err := ReadByte(reader)
	if result != input[0] {
		t.Error("Unexpected read result from ReadByte", result)
	}
	if err != nil {
		t.Error("Unexpected error from ReadByte", err)
	}

	result, err = ReadByte(reader)
	if err != io.EOF {
		t.Error("Unexpected nil error from ReadByte, read value:", result)
	}
}

func TestSizedWriter(t *testing.T) {
	var (
		buffer bytes.Buffer
		writer io.Writer = SizedWriter(&buffer, 4)
	)

	count, err := writer.Write([]byte("12"))
	if count != 2 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}

	count, err = writer.Write([]byte("3456"))
	if count != 2 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}
	if buffer.String() != "1234" {
		t.Error("Unexpected value in wrapped writer", buffer.String())
	}

	// Flush the buffer
	count, err = writer.Write(nil)
	if count != 0 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}
	if buffer.String() != "123456" {
		t.Error("Unexpected value in wrapped writer", buffer.String())
	}

	count, err = writer.Write([]byte("7890"))
	if count != 4 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}
	if buffer.String() != "1234567890" {
		t.Error("Unexpected value in wrapped writer", buffer.String())
	}
}

func TestSizeWriterLarger(t *testing.T) {
	var (
		input  []byte = []byte("0123456789AB")
		buffer bytes.Buffer
		writer = SizedWriter(&buffer, 8)
	)

	count, err := writer.Write(input)
	if count != 12 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}
	if buffer.String() != "01234567" {
		t.Error("Unexpected value in wrapped writer", buffer.String())
	}

	count, err = writer.Write(nil)
	if count != 0 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}
	if buffer.String() != "0123456789AB" {
		t.Error("Unexpected value in wrapped writer", buffer.String())
	}
}

func TestSizedWriterError1(t *testing.T) {
	var (
		errorWriter io.Writer = WriterFunc(func([]byte) (int, error) {
			return 1, errors.New("Invalid write")
		})
		writer io.Writer = SizedWriter(errorWriter, 2)
	)

	count, err := writer.Write([]byte("1"))
	if count != 1 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedWriter", err)
	}

	count, err = writer.Write([]byte("2"))
	if count != 1 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err == nil {
		t.Error("Unexpected lack of error from SizedWriter")
	}
}

func TestSizedWriterError2(t *testing.T) {
	var (
		errorWriter io.Writer = WriterFunc(func([]byte) (int, error) {
			return 1, errors.New("Invalid write")
		})
		writer io.Writer = SizedWriter(errorWriter, 1)
	)

	count, err := writer.Write([]byte("12"))
	if count != 1 {
		t.Error("Unexpected write count from SizedWriter", count)
	}
	if err == nil {
		t.Error("Unexpected lack of error from SizedWriter")
	}
}

func TestSizedReader(t *testing.T) {
	var (
		input  []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7}
		output []byte = make([]byte, 16)

		reader *bytes.Reader = bytes.NewReader(input)
		min    io.Reader     = SizedReader(reader, 4)
	)

	// Expecting a read count of 2
	count, err := min.Read(output[:2])
	if count != 2 {
		t.Error("Invalid read count from SizedReader", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedReader", err)
	}

	// Expecting a read count of 2 as it should have 2 bytes in its buffer
	count, err = min.Read(output[:3])
	if count != 2 {
		t.Error("Invalid read count from SizedReader", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedReader", err)
	}

	// Expecting a read count of 4 as the buffer should be empty
	count, err = min.Read(output[:4])
	if count != 4 {
		t.Error("Invalid read count from SizedReader", count)
	}
	if err != nil {
		t.Error("Unexpected error from SizedReader", err)
	}

	// Expecting a read count of 0 with an EOF as the buffer should be empty
	count, err = min.Read(output[:1])
	if count != 0 {
		t.Error("Invalid read count from SizedReader", count)
	}
	if err != io.EOF {
		t.Error("Unexpected error from SizedReader", err)
	}
}
