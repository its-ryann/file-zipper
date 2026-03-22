package compressor

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func Compress(inputPath string, outputPath string) error {
	src, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("compress: could not open input file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("compress: could not create output file: %w", err)
	}
	defer dst.Close()

	writer := gzip.NewWriter(dst)

	if _, err = io.Copy(writer, src); err != nil {
		writer.Close()
		return fmt.Errorf("compress: failed during streaming: %w", err)
	}

	// Close explicitly so we catch any flush errors
	return writer.Close()
}

func Decompress(inputPath string, outputPath string) error {
	src, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("decompress: could not open input file: %w", err)
	}
	defer src.Close()

	// Wrap the source in a gzip reader to decode the compressed stream
	reader, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("decompress: input is not a valid gzip file: %w", err)
	}
	defer reader.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("decompress: could not create output file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, reader); err != nil {
		return fmt.Errorf("decompress: failed during streaming: %w", err)
	}

	return nil
}