package compressor

import (
	"testing"
)

func TestCalculateRatio(t *testing.T) {
	tests := []struct {
		name       string
		original   int64
		compressed int64
		want       float64
	}{
		{"half size", 100, 50, 50.0},
		{"no reduction", 100, 100, 0.0},
		{"typical text", 1000, 300, 70.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateRatio(tt.original, tt.compressed)
			if got != tt.want {
				t.Errorf("CalculateRatio(%d, %d) = %.1f, want %.1f",
					tt.original, tt.compressed, got, tt.want)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatSize(tt.bytes)
			if got != tt.want {
				t.Errorf("FormatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}