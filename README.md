# grun

A fast Go script runner with intelligent caching. Compiles Go files only when they change, then runs them instantly.

## Features

- **Smart caching**: Only recompiles when the source file changes
- **Fast execution**: Uses cached binaries for instant runs
- **Shebang support**: Make Go scripts directly executable like shell scripts
- **Dependency support**: Automatically handles scripts with external dependencies via `go.mod`
- **Hybrid approach**: Works with both simple single-file scripts and complex modules
- **Intelligent recompilation**: Tracks changes to source files, `go.mod`, and `go.sum`
- **Clean temp handling**: Temporary files are managed in system temp directory, never polluting your workspace
- **Configurable cache directory**: Default `~/.cache/grun`, override via flag or env var
- **Easy installation**: Install as a system tool with one command

## Installation

Install `grun` as a system tool:

```bash
./setup.sh install
```

This will:
- Build the binary
- Install it to `~/.local/bin` (or `/usr/local/bin` if user directory is not writable)
- Make it available system-wide

**Uninstall:**

```bash
./setup.sh uninstall
```

**Note:** Make sure `~/.local/bin` is in your PATH. If not, add this to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.):
```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Usage

### Basic usage

Run a Go file:

```bash
grun script.go
```

The first run will compile the file and cache the binary. Subsequent runs will use the cached version unless the source file changes.

### Passing arguments

Pass arguments to your script:

```bash
grun script.go arg1 arg2 arg3
```

### Making scripts executable (shebang)

Make your Go scripts directly executable like shell scripts:

```bash
# Add shebang to your script
cat > hello.go << 'EOF'
#!/usr/bin/env grun
package main

import "fmt"

func main() {
    fmt.Println("Hello from executable Go script!")
}
EOF

# Make it executable
chmod +x hello.go

# Run it directly!
./hello.go
```

**How it works:**
- The shebang `#!/usr/bin/env grun` tells the OS to run the script with `grun`
- Make the script executable: `chmod +x hello.go`
- `grun` automatically strips the shebang before compilation
- Temporary files are created in `/tmp`, keeping your workspace clean
- Caching works normally - scripts execute instantly after first compilation

**Note:** Don't forget to make your script executable with `chmod +x`! Without execute permissions, you'll get a "permission denied" error.

**Works with dependencies too:**
```bash
#!/usr/bin/env grun
package main

import "github.com/fatih/color"

func main() {
    color.Green("Executable with dependencies!")
}
```

Just ensure your script directory has a `go.mod` file!

### Using external dependencies

For scripts that need external packages, initialize a Go module in the script's directory:

```bash
# Create your script
mkdir my-script && cd my-script
cat > script.go << 'EOF'
package main

import (
    "fmt"
    "github.com/fatih/color"
)

func main() {
    color.Green("Hello with dependencies!")
}
EOF

# Initialize module and add dependencies
go mod init my-script
go get github.com/fatih/color

# Run with grun
grun script.go
```

Once set up, `grun` will automatically:
- Detect the `go.mod` file
- Build the package with all dependencies
- Recompile when `go.mod` or `go.sum` changes
- Cache the binary for fast subsequent runs

### Override cache directory

**Using command-line flag:**
```bash
grun -cache-dir /tmp/my-cache script.go
```

**Using environment variable:**
```bash
GRUN_CACHE=/tmp/my-cache grun script.go
```

## How it works

1. **First run**: `grun` compiles your Go file and stores the binary in the cache directory
2. **Subsequent runs**: `grun` checks if the source file, `go.mod`, or `go.sum` has been modified
   - If unchanged: Uses the cached binary (instant execution)
   - If changed: Recompiles and updates the cache

The cache key is based on the absolute path of the source file, so different files get different cached binaries.

### Dependency Handling

`grun` intelligently detects how to build your script:

- **With `go.mod`**: If a `go.mod` file exists in the script's directory, `grun` builds the entire package (with all dependencies)
- **Without `go.mod`**: Builds as a single file (standard library only)
- **Build failure hint**: If building fails without `go.mod`, you'll get helpful suggestions to initialize a module

This hybrid approach means you can use `grun` for:
- Quick throwaway scripts (no setup needed)
- Production-quality scripts with external dependencies (just add `go.mod`)

## Cache Directory Priority

The cache directory is determined in the following order:
1. `-cache-dir` command-line flag
2. `GRUN_CACHE` environment variable
3. Default: `~/.cache/grun`

## Examples

The repository includes two example scripts demonstrating different use cases:

### Simple Script (No Dependencies)

Located in `examples/scripts/script.go` - demonstrates a simple script with shebang using only the standard library:

```bash
$ ./examples/scripts/script.go
=== Simple grun Example ===

This is a simple Go script with no external dependencies.
It uses only the standard library and runs instantly with grun!

Arguments received: test, args
```

Or use `grun` explicitly: `grun examples/scripts/script.go`

### Script with External Dependencies

Located in `examples/with-deps/script.go` - demonstrates using external packages with shebang:

```bash
$ ./examples/with-deps/script.go
=== grun Example with Dependencies ===

✓ Successfully imported and used github.com/fatih/color
✓ This script requires go.mod to work
✓ grun automatically detects go.mod and builds accordingly
```

Or use `grun` explicitly: `grun examples/with-deps/script.go`

Both examples have shebangs (`#!/usr/bin/env grun`) and can be run directly. They cache their binaries and only recompile when files change.

## Testing

Run Go unit tests:

```bash
go test ./...
```

## Development

```bash
# Build grun
go build -o grun ./main.go

# Test with simple example
./grun examples/scripts/script.go
# Or run directly (has shebang)
./examples/scripts/script.go

# Test with dependencies example
./grun examples/with-deps/script.go
# Or run directly (has shebang)
./examples/with-deps/script.go

# Install as system tool
./setup.sh install

# Run tests
go test ./...
```
