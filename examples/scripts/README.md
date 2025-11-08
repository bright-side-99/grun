# Example: Simple Scripts (No Dependencies)

This directory demonstrates using `grun` with simple Go scripts that only use the standard library.

## What's Here

- `script.go` - A simple Go script with shebang (call with `./script.go` or `grun script.go`)
- `hello-executable.go` - An executable script with shebang (call with `./hello-executable.go` or `grun hello-executable.go`)

## How to Run

**Run directly (recommended):**
```bash
# All scripts have shebangs and are executable
./script.go arg1 arg2
./hello-executable.go arg1 arg2
```

**Or use grun explicitly:**
```bash
# From the grun root directory
grun examples/scripts/script.go

# Or from this directory
grun script.go arg1 arg2
```

## What This Demonstrates

1. **Zero setup required**: No `go.mod` needed for simple scripts
2. **Instant execution**: First run compiles, subsequent runs use cached binary
3. **Shebang support**: Go scripts can be executable like shell scripts
4. **Standard library only**: Perfect for quick utilities and throwaway scripts

## Behind the Scenes

When you run `grun` on these scripts:
1. Detects no `go.mod` file in the directory
2. For scripts with `#!/usr/bin/env grun`, strips the shebang before compilation
3. Creates temporary file in `/tmp` (not in your workspace)
4. Builds using `go build` (single file compilation)
5. Caches the compiled binary
6. Cleans up temporary files automatically
7. Recompiles only if script changes

