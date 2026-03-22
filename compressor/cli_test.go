package compressor

import (
	"os"
	"testing"
)

func TestRunCompress(t *testing.T) {
	// Create a temp input file
	input, err := os.CreateTemp("", "cli-input-*.txt")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(input.Name())

	_, err = input.WriteString("hello hello hello world world world hello hello hello world world world")
	if err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	input.Close()

	outputPath := input.Name() + ".gz"
	defer os.Remove(outputPath)

	// Simulate: zip-it compress <input>
	err = Run([]string{"compress", input.Name()})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	// Output file should exist
	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file %s to exist", outputPath)
	}
}

func TestRunDecompress(t *testing.T) {
	// Create and compress a temp file first
	input, err := os.CreateTemp("", "cli-input-*.txt")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer os.Remove(input.Name())

	_, err = input.WriteString("hello hello hello world world world hello hello hello world world world")
	if err != nil {
		t.Fatalf("could not write temp file: %v", err)
	}
	input.Close()

	compressedPath := input.Name() + ".gz"
	defer os.Remove(compressedPath)

	err = Compress(input.Name(), compressedPath)
	if err != nil {
		t.Fatalf("setup Compress failed: %v", err)
	}

	outputPath := input.Name() + ".decompressed.txt"
	defer os.Remove(outputPath)

	// Simulate: zip-it decompress <input.gz>
	err = Run([]string{"decompress", compressedPath})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Errorf("expected output file %s to exist", outputPath)
	}
}

func TestRunInvalidCommand(t *testing.T) {
	err := Run([]string{"explode", "somefile.txt"})
	if err == nil {
		t.Error("expected an error for unknown command, got nil")
	}
}

func TestRunMissingArgs(t *testing.T) {
	err := Run([]string{"compress"})
	if err == nil {
		t.Error("expected an error when no file argument given, got nil")
	}
}