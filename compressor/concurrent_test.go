package compressor

import (
	"os"
	"testing"
)

func TestCompressConcurrent(t *testing.T) {
	// Create three temp input files
	paths := make([]string, 3)
	for i := range paths {
		f, err := os.CreateTemp("", "concurrent-*.txt")
		if err != nil {
			t.Fatalf("could not create temp file: %v", err)
		}
		_, err = f.WriteString("hello hello hello world world world hello hello hello world world world")
		if err != nil {
			t.Fatalf("could not write temp file: %v", err)
		}
		f.Close()
		paths[i] = f.Name()
		defer os.Remove(f.Name())
		defer os.Remove(f.Name() + ".gz")
	}

	// Compress all three concurrently
	results, errors := CompressConcurrent(paths)

	// Check no errors occurred
	for _, err := range errors {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	// Check we got a result for every file
	if len(results) != len(paths) {
		t.Errorf("expected %d results, got %d", len(paths), len(results))
	}

	// Check every output file exists and is smaller
	for _, r := range results {
		info, err := os.Stat(r.OutputPath)
		if err != nil {
			t.Errorf("output file %s not found", r.OutputPath)
			continue
		}
		if info.Size() >= r.OriginalSize {
			t.Errorf("expected %s to be smaller than original", r.OutputPath)
		}
	}
}
