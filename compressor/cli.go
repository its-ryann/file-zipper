package compressor

import (
	"fmt"
	"strings"
)

// Run is the entry point for all CLI commands.
// args is everything after the binary name — e.g. ["compress", "myfile.txt"]
func Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: zip-it <compress|decompress> <file>")
	}

	command := args[0]
	inputPath := args[1]

	switch command {
	case "compress":
		outputPath := inputPath + ".gz"
		return Compress(inputPath, outputPath)

	case "decompress":
		// Strip the .gz extension to get the output path
		outputPath := strings.TrimSuffix(inputPath, ".gz") + ".decompressed.txt"
		return Decompress(inputPath, outputPath)

	default:
		return fmt.Errorf("unknown command %q — use compress or decompress", command)
	}
}