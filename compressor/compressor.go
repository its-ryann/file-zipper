package compressor

import (
	"compress/gzip"
	"io"
	"os"
)

func Compress(inputPath string, outputPath string) error {
	src, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	writer := gzip.NewWriter(dst)

	if _, err = io.Copy(writer, src); err != nil {
		writer.Close()
		return err
	}

	// Close explicitly so we catch any flush errors
	return writer.Close()
}