# grun

A fast Go script runner with intelligent caching. Compiles Go files only when they change, then runs them instantly.

## Features

- **Smart caching**: Only recompiles when the source file changes
- **Fast execution**: Uses cached binaries for instant runs
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
2. **Subsequent runs**: `grun` checks if the source file has been modified
   - If unchanged: Uses the cached binary (instant execution)
   - If changed: Recompiles and updates the cache

The cache key is based on the absolute path of the source file, so different files get different cached binaries.

## Cache Directory Priority

The cache directory is determined in the following order:
1. `-cache-dir` command-line flag
2. `GRUN_CACHE` environment variable
3. Default: `~/.cache/grun`

## Example

The repository includes an example script in `examples/script.go`. Run it:

```bash
$ grun examples/script.go
Hello, world!

$ grun examples/script.go
Hello, world!
```

The second run uses the cached binary. Modify the script and run again - it will automatically recompile:

```bash
$ # Edit examples/script.go to print "Hello, updated!"
$ grun examples/script.go
Hello, updated!
```

## Testing

Run Go unit tests:

```bash
go test ./...
```

## Development

```bash
# Build grun
go build -o grun ./main.go

# Test with example script
./grun examples/script.go

# Install as system tool
./setup.sh install

# Run tests
go test ./...
```
