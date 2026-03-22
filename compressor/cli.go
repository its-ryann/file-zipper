package compressor

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: zip-it <compress|decompress> <file> [file2 file3 ...]")
	}

	command := args[0]
	inputPaths := args[1:]

	if command == "compress" && len(inputPaths) > 1 {
		results, errors := CompressConcurrent(inputPaths)
		for _, r := range results {
			r.Print()
		}
		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			return fmt.Errorf("%d file(s) failed to compress", len(errors))
		}
		return nil
	}

	inputPath := inputPaths[0]

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

	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("could not read output file: %w", err)
	}

	Result{
		InputPath:    inputPath,
		OutputPath:   outputPath,
		OriginalSize: inputInfo.Size(),
		OutputSize:   outputInfo.Size(),
		Duration:     time.Since(start),
	}.Print()

	return nil
}
