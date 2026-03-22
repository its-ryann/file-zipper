#  zip-it

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Build](https://img.shields.io/badge/build-passing-brightgreen?style=flat)
![TDD](https://img.shields.io/badge/built%20with-TDD-blue?style=flat)

> A high-performance CLI tool for Gzip compression and decompression, built in Go. Processes files of any size using streaming I/O, reports real-time performance metrics, and compresses multiple files simultaneously using goroutines.

---

## Why I Built This

Most developers use compression tools every day without thinking about what happens underneath. I built `zip-it` to answer that question for myself — *how does data actually shrink?*

This project taught me how computers move data through memory in chunks using buffer streams, how the DEFLATE algorithm exploits repeating byte patterns, and how to write concurrent Go programs that do real work in parallel. Every feature was built using **Test-Driven Development** and managed with the **Gitflow** branching strategy.

---

## Features

| | Feature | Description |
|---|---|---|
| 1. | Streaming I/O | Processes files in chunks — handles multi-GB files on minimal RAM |
| 2. | Concurrent compression | Compresses multiple files simultaneously using goroutines |
| 3. | Real-time metrics | Reports original size, compressed size, ratio, and elapsed time |
| 4. | Robust error handling | Wraps errors with context so failures are always actionable |
| 5. | Clean CLI interface | Simple argument-driven commands with no dependencies |

---

## Installation

**Prerequisites:** Go 1.22+

```bash
git clone https://github.com/its-ryann/file-zipper.git
cd file-zipper
make build
```

The binary is output to `./bin/zip-it`.

---

## Usage

**Compress a single file:**
```bash
./bin/zip-it compress myfile.txt
```

**Decompress a file:**
```bash
./bin/zip-it decompress myfile.txt.gz
```

**Compress multiple files concurrently:**
```bash
./bin/zip-it compress file1.txt file2.txt file3.txt
```

**Example output:**
```
✔  file1.txt → file1.txt.gz
   original:   72 B
   compressed: 42 B
   ratio:      41.7% reduction
   time:       0.001s

✔  file2.txt → file2.txt.gz
   original:   88 B
   compressed: 69 B
   ratio:      21.6% reduction
   time:       0.001s

✔  file3.txt → file3.txt.gz
   original:   95 B
   compressed: 97 B
   ratio:      -2.1% reduction
   time:       0.001s
```

> Note: a negative ratio means the compressed output is larger than the input. This is expected for very small or already-compressed files where Gzip's header overhead exceeds the savings. `zip-it` reports this honestly rather than hiding it.

---

## How It Works

### Streaming I/O

Instead of loading an entire file into RAM, `zip-it` streams data through the gzip encoder in chunks:

```
[Source File] ──► [Gzip Writer] ──► [Output File]
```

`io.Copy` moves data continuously without buffering the whole file. A 4 GB file uses the same working memory as a 4 KB file.

```go
writer := gzip.NewWriter(dst)
if _, err = io.Copy(writer, src); err != nil {
    writer.Close()
    return fmt.Errorf("compress: failed during streaming: %w", err)
}
return writer.Close()
```

Note that `writer.Close()` is called explicitly rather than deferred — this is intentional. The gzip writer flushes its final compressed bytes on close, and that flush can fail. Deferring it silently swallows that error and can produce corrupt output.

### Concurrent Compression

When multiple files are passed, `zip-it` spawns one goroutine per file and collects results through a channel:

```go
for _, path := range inputPaths {
    wg.Add(1)
    go func(inputPath string) {
        defer wg.Done()
        // compress and send result to channel
    }(path)
}
```

A `sync.WaitGroup` tracks completion, and the channel is closed once all goroutines finish — allowing the main goroutine to range over results cleanly without a race condition.

### The DEFLATE Algorithm

Gzip uses DEFLATE under the hood, which combines two techniques:

- **LZ77** — identifies repeating byte sequences and replaces them with compact back-references
- **Huffman coding** — assigns shorter bit patterns to more frequent bytes

This is why plain text and JSON compress aggressively (many repeated words and characters), while PNG images compress poorly (already encoded, minimal repetition).

---

## Benchmarks

Tested on Ubuntu 22.04, Go 1.22.2:

| File type | Original | Compressed | Reduction |
|---|---|---|---|
| Repetitive plain text | 72 B | 42 B | **41.7%** |
| Natural language text | 88 B | 69 B | **21.6%** |
| Low-repetition text | 95 B | 97 B | -2.1% |

For larger, real-world files the gains are more pronounced — a 1 MB JSON API response typically compresses to under 100 KB. This is why HTTP uses Gzip by default.

---

## Project Structure

```
file-zipper/
├── main.go                    # Entry point — passes os.Args to Run()
├── compressor/
│   ├── compressor.go          # Compress() and Decompress() — streaming gzip
│   ├── compressor_test.go     # Tests for core compression logic
│   ├── cli.go                 # Run() — argument parsing and dispatch
│   ├── cli_test.go            # Tests for CLI behaviour
│   ├── metrics.go             # CalculateRatio(), FormatSize(), Result.Print()
│   ├── metrics_test.go        # Tests for metrics calculations
│   ├── concurrent.go          # CompressConcurrent() — goroutine pool
│   └── concurrent_test.go     # Tests for concurrent behaviour
├── go.mod
├── Makefile
└── README.md
```

---

## Makefile

```bash
make build   # compile binary to ./bin/zip-it
make test    # run full test suite
make clean   # remove build artifacts
```

---

## Development Approach

This project was built using **Test-Driven Development** — every function was tested before it was implemented. The Red → Green → Refactor cycle was followed for each feature:

1. Write a failing test that defines the expected behaviour
2. Write the minimum code to make it pass
3. Refactor for clarity and correctness without breaking the test

Branch management followed the **Gitflow** workflow. Each feature lived on its own `feature/*` branch, was merged into `develop` via a no-fast-forward merge, and `main` was only updated on release.

---

## What I Learned

**Streams vs buffers:** loading a file with `os.ReadFile` works until the file is larger than available RAM. Streaming with `io.Copy` removes that ceiling entirely — memory usage stays flat regardless of file size.

**The cost of `Close()`:** a gzip writer holds compressed bytes in a buffer and only flushes them when closed. Deferring close and ignoring its return value silently produces truncated, corrupt output. Explicit close with error handling is the correct pattern.

**Variable scoping in Go:** using `:=` inside a switch case declares a new local variable — it does not assign to an outer one. This caused a real bug during development where metrics were computed against an empty path. The fix (`=` instead of `:=`) is one character, but catching it required understanding Go's scoping rules precisely.

**Goroutines and channels:** launching concurrent work is easy in Go. The discipline is in collecting results safely — a buffered channel sized to the number of jobs and a `WaitGroup` to signal completion is the standard pattern, and it composes cleanly.

---

## Roadmap

- [x] Single-file compression
- [x] Single-file decompression
- [x] Real-time performance metrics
- [x] Concurrent multi-file compression
- [ ] Progress bar for large files
- [ ] Recursive directory compression
- [ ] Custom compression level flag (`--level 1-9`)
- [ ] Decompress multiple files concurrently

---

## License

MIT — see [LICENSE](LICENSE) for details.