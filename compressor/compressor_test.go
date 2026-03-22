package compressor

import (
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	// Create a temporary input file with known content
	input, err := os.CreateTemp("", "input-*.txt")
	if err != nil {
		t.Fatalf("could not create temp input file: %v", err)
	}
	defer os.Remove(input.Name())

	_, err = input.WriteString("hello hello hello world world world hello hello hello world world world hello hello hello world world world hello hello hello world world world")
	if err != nil {
		t.Fatalf("could not write to temp file: %v", err)
	}
	input.Close()

	// Define where the compressed output will go
	outputPath := input.Name() + ".gz"
	defer os.Remove(outputPath)

	// Call the function we haven't written yet
	err = Compress(input.Name(), outputPath)
	if err != nil {
		t.Fatalf("Compress returned an error: %v", err)
	}

	// Verify the output file was actually created
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file was not created: %v", err)
	}

	// Verify the compressed file is smaller than the input
	inputInfo, _ := os.Stat(input.Name())
	if info.Size() >= inputInfo.Size() {
		t.Errorf("expected compressed file to be smaller, got %d >= %d", info.Size(), inputInfo.Size())
	}
}