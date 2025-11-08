# grun Examples

This directory contains example scripts demonstrating different use cases for `grun`.

## Examples

### 1. Simple Script (`scripts/`)

A basic Go script using only the standard library. Perfect for quick utilities and throwaway scripts.

**Key features:**
- No `go.mod` required
- Zero setup
- Instant compilation and caching

**Run it:**
```bash
grun scripts/script.go
```

### 2. Script with Dependencies (`with-deps/`)

A Go script that uses external packages (demonstrates `github.com/fatih/color`). Shows how `grun` automatically handles Go modules and dependencies.

**Key features:**
- Uses `go.mod` for dependency management
- Automatic detection and building
- Tracks changes to `go.mod` and `go.sum`

**Run it:**
```bash
grun with-deps/script.go
```

## Creating Your Own Scripts

### Simple Scripts (No Dependencies)

Just write a Go file and run it:

```go
// hello.go
package main
import "fmt"
func main() { fmt.Println("Hello!") }
```

```bash
grun hello.go
```

### Scripts with Dependencies

1. Create a directory for your script
2. Initialize a Go module
3. Add dependencies
4. Run with `grun`

```bash
mkdir myscript && cd myscript
go mod init myscript
go get github.com/some/package
# Write your script.go
grun script.go
```

## Testing All Examples

From the grun root directory:

```bash
# Build grun first
go build -o grun ./main.go

# Test simple script
./grun examples/scripts/script.go

# Test script with dependencies
./grun examples/with-deps/script.go

# Try with arguments
./grun examples/scripts/script.go arg1 arg2
./grun examples/with-deps/script.go hello world
```

