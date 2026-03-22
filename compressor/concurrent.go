package compressor

import (
	"os"
	"sync"
	"time"
)

type job struct {
	inputPath string
}

type jobResult struct {
	result Result
	err    error
}

func CompressConcurrent(inputPaths []string) ([]Result, []error) {
	var wg sync.WaitGroup
	resultsCh := make(chan jobResult, len(inputPaths))

	// Launch one goroutine per file
	for _, path := range inputPaths {
		wg.Add(1)
		go func(inputPath string) {
			defer wg.Done()

			inputInfo, err := os.Stat(inputPath)
			if err != nil {
				resultsCh <- jobResult{err: err}
				return
			}

			outputPath := inputPath + ".gz"
			start := time.Now()

			if err := Compress(inputPath, outputPath); err != nil {
				resultsCh <- jobResult{err: err}
				return
			}

			outputInfo, err := os.Stat(outputPath)
			if err != nil {
				resultsCh <- jobResult{err: err}
				return
			}

			resultsCh <- jobResult{
				result: Result{
					InputPath:    inputPath,
					OutputPath:   outputPath,
					OriginalSize: inputInfo.Size(),
					OutputSize:   outputInfo.Size(),
					Duration:     time.Since(start),
				},
			}
		}(path)
	}

	// Close the channel once all goroutines finish
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results from the channel
	var results []Result
	var errors []error
	for jr := range resultsCh {
		if jr.err != nil {
			errors = append(errors, jr.err)
		} else {
			results = append(results, jr.result)
		}
	}

	return results, errors
}