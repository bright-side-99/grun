# Example: Scripts with External Dependencies

This directory demonstrates using `grun` with Go scripts that require external dependencies.

## What's Here

- `script.go` - A Go script using `github.com/fatih/color` (call with `grun script.go`)
- `color-demo.go` - An executable script with shebang and dependencies (call with `./color-demo.go`)
- `go.mod` - Module file defining dependencies
- `go.sum` - Checksums for dependencies

## How to Run

**Traditional approach:**
```bash
# From the grun root directory
grun examples/with-deps/script.go

# Or from this directory
grun script.go arg1 arg2
```

**Executable with shebang:**
```bash
# Make it executable (already done in repo)
chmod +x color-demo.go

# Run directly!
./color-demo.go arg1 arg2
```

## What This Demonstrates

1. **Automatic dependency detection**: `grun` detects the `go.mod` file and builds with dependencies
2. **Dependency management**: All external packages are properly resolved
3. **Shebang with modules**: Executable scripts work even with external dependencies
4. **Smart caching**: Binary is cached and reused until source or dependencies change
5. **Transparent experience**: Works just like a simple script, but with full module support

## Behind the Scenes

When you run `grun` on these scripts:
1. Detects `go.mod` in the script's directory
2. For scripts with `#!/usr/bin/env grun`, strips the shebang before compilation
3. Creates temporary file in `/tmp` (keeps workspace clean)
4. Builds using `go build` with module support
5. Caches the compiled binary
6. Cleans up temporary files automatically
7. Recompiles if `script.go`, `go.mod`, or `go.sum` changes

