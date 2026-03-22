# zip-it — A Gzip CLI Compressor

> A high-performance command-line tool for compressing and decompressing files using the Gzip format. Built with streaming I/O to handle files far larger than available system RAM.

---

## Why I Built This

Most developers use compression tools daily without ever thinking about what happens underneath. I built `zip-it` to answer that question for myself: *How does data actually shrink?*

This project taught me how computers move data through memory in chunks using buffer streams, how the DEFLATE algorithm exploits repeating patterns in text, and how a professional Go project is structured for real-world use. The result is a tool I actually use — and a codebase I'm proud to show.

---

## Features

| Feature | Description |
|---|---|
| Streaming I/O | Processes files in chunks — handles multi-GB files without crashing |
| Performance Metrics | Reports original size, compressed size, and ratio in real time |
| Robust Error Handling | Validates file paths, permissions, and flags before execution |
| Clean CLI Interface | Simple, intuitive argument-driven commands |

---

## Installation

**Prerequisites:** Go 1.21+

```bash
# Clone the repository
git clone https://github.com/YOUR_USERNAME/file-zipper.git
cd file-zipper

# Build the binary
make build

# (Optional) Install globally
make install
```

---

## Usage

```bash
# Compress a file
zip-it compress myfile.txt

# Decompress a file
zip-it decompress myfile.txt.gz

# Compress with verbose output
zip-it compress myfile.txt --verbose
```

**Example output:**

```
✔ Compressed: myfile.txt → myfile.txt.gz
  Original size:    2.40 MB
  Compressed size:  612.00 KB
  Compression ratio: 74.5% reduction
  Time elapsed:     0.03s
```

---

## Compression Benchmarks

Results on a MacBook Pro M2 / Ubuntu 22.04 (your results may vary):

| File Type | Original Size | Compressed Size | Reduction |
|---|---|---|---|
| Plain text `.txt` | 2.4 MB | 612 KB | **74.5%** |
| JSON data `.json` | 1.1 MB | 98 KB | **91.1%** |
| Go source `.go` | 480 KB | 142 KB | **70.4%** |
| PNG image `.png` | 3.2 MB | 3.1 MB | **3.1%** |
| Binary `.exe` | 8.0 MB | 7.2 MB | **10.0%** |

> **Why do images compress poorly?** PNG is already compressed. Gzip excels at text and structured data where repeating patterns are abundant.

---

## Project Structure

```
/file-zipper
├── main.go              # Entry point — parses CLI arguments
├── compressor/
│   ├── compress.go      # Gzip compression logic
│   └── decompress.go    # Gzip decompression logic
├── go.mod               # Module and dependency management
├── Makefile             # Build, clean, and install automation
└── README.md            # Documentation and learning log
```

---

## How It Works

### 1. Streaming I/O with `io.Reader` & `io.Writer`

Instead of loading an entire file into RAM, `zip-it` uses a **bucket brigade** approach:

```
[Source File] → [Buffered Reader] → [Gzip Writer] → [Output File]
```

The `io.Copy` function moves data in 32KB chunks. A 4GB file uses only ~32KB of working memory at any given moment.

```go
// The core of the compressor — deceptively simple
src, _ := os.Open(inputPath)
dst, _ := os.Create(outputPath)
writer := gzip.NewWriter(dst)

io.Copy(writer, src) // The magic happens here
writer.Close()
```

### 2. The DEFLATE Algorithm

Gzip uses **DEFLATE**, which combines two techniques:

- **LZ77** — Finds repeating sequences and replaces them with back-references: *"repeat what I said 50 bytes ago, for 12 bytes."*
- **Huffman Coding** — Assigns shorter bit patterns to more frequent bytes, like Morse code assigns `·` to the letter E.

This is why text compresses so well (lots of repeated words) and why already-compressed images do not.

---

## Makefile Commands

```bash
make build    # Compiles the binary to ./bin/zip-it
make run      # Builds and runs with default test args
make clean    # Removes build artifacts
make install  # Installs binary to /usr/local/bin
make test     # Runs the test suite
```

---

## Roadmap

- [x] Single-file compression and decompression
- [x] Real-time performance metrics
- [ ] **Concurrent compression** — compress multiple files simultaneously using goroutines
- [ ] Progress bar for large files
- [ ] Directory compression (recursive)
- [ ] Custom compression level flag (`--level 1-9`)

---

## What I Learned

> *This section is my engineering learning log.*

**Streams vs. Buffers:** I initially tried to read entire files into memory with `os.ReadFile`. It worked for small files but would have collapsed on anything large. Switching to streaming `io.Copy` was the first real systems-thinking decision I made in Go.

**The cost of `Close()`:** A Gzip writer buffers its final bytes and only flushes them when you call `.Close()`. Forgetting this produces a corrupt file. Silent bugs like this taught me to always use `defer writer.Close()`.

**Why compression ratios matter for the web:** HTTP responses use Gzip by default. Understanding this project made me realize why a 1MB JSON API response might only cost 80KB of bandwidth.

---

## License

MIT — see [LICENSE](LICENSE) for details.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

```bash
# Run tests before submitting a PR
make test
```
