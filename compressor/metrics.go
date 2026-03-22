package compressor

import (
	"fmt"
	"time"
)

func CalculateRatio(original, compressed int64) float64 {
	if original == 0 {
		return 0.0
	}
	return float64(original-compressed) / float64(original) * 100
}

func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

type Result struct {
	InputPath    string
	OutputPath   string
	OriginalSize int64
	OutputSize   int64
	Duration     time.Duration
}

func (r Result) Print() {
	ratio := CalculateRatio(r.OriginalSize, r.OutputSize)
	fmt.Printf("✔  %s → %s\n", r.InputPath, r.OutputPath)
	fmt.Printf("   original:   %s\n", FormatSize(r.OriginalSize))
	fmt.Printf("   compressed: %s\n", FormatSize(r.OutputSize))
	fmt.Printf("   ratio:      %.1f%% reduction\n", ratio)
	fmt.Printf("   time:       %.3fs\n", r.Duration.Seconds())
}