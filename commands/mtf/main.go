package main

import (
	"flag"
	fib "github.com/spaskalev/misc/encoding/fibonacci"
	mtf "github.com/spaskalev/misc/encoding/mtf"
	iou "github.com/spaskalev/misc/ioutil"
	"io"
	"os"
)

func main() {
	d := flag.Bool("d", false, "Toggle decode mode.")
	flag.Parse()

	var (
		input  io.Reader = iou.SizedReader(os.Stdin, 4096)
		output io.Writer = iou.SizedWriter(os.Stdout, 4096)
		code   int
	)

	// Exit handler
	defer func() {
		os.Exit(code)
	}()

	// Flush the output buffer
	defer output.Write(nil)

	if *d {
		input = mtf.Decoder(fib.Decoder(input))
	} else {
		input = mtf.Encoder(input)

		// Encode output as fibonacci integers
		output = fib.Encoder(output)
		defer output.Write(nil)
	}

	if _, err := io.Copy(output, input); err != nil {
		os.Stderr.WriteString("Error while transforming data.\n" + err.Error())
		code = 1
	}
}
