package main

import (
	"fmt"
	"os"

	"github.com/its-ryann/file-zipper/compressor"
)

func main() {
	// os.Args[0] is the binary name, so we pass everything after it
	if err := compressor.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}