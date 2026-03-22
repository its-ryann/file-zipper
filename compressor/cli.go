package compressor

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Run is the entry point for all CLI commands.
// args is everything after the binary name — e.g. ["compress", "myfile.txt"]
func Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: zip-it <compress|decompress> <file>")
	}

	command := args[0]
	inputPath := args[1]

	// Get the original file size before doing anything
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("could not read input file: %w", err)
	}

	var outputPath string
	start := time.Now()

	switch command {
	case "compress":
		outputPath = inputPath + ".gz"
		err = Compress(inputPath, outputPath)
		if err != nil {
			return err
		}

	case "decompress":
		outputPath = strings.TrimSuffix(inputPath, ".gz") + ".decompressed.txt"
		err = Decompress(inputPath, outputPath)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown command %q — use compress or decompress", command)
	}

	// Get the output file size
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("could not read output file: %w", err)
	}

	// Print the metrics report
	Result{
		InputPath:    inputPath,
		OutputPath:   outputPath,
		OriginalSize: inputInfo.Size(),
		OutputSize:   outputInfo.Size(),
		Duration:     time.Since(start),
	}.Print()

	return nil
}