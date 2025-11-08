package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultCacheDir = ".cache/grun"
)

var (
	cacheDir = flag.String("cache-dir", "", "Directory to use for caching (default: .cache/grun)")
)

func getCacheDir() string {
	// Check if cache-dir flag is set
	if *cacheDir != "" {
		return *cacheDir
	}

	// Check if GRUN_CACHE environment variable is set
	if envCache := os.Getenv("GRUN_CACHE"); envCache != "" {
		return envCache
	}

	// Use default cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home directory
		return defaultCacheDir
	}

	return filepath.Join(homeDir, defaultCacheDir)
}

func ensureCacheDir() error {
	cachePath := getCacheDir()
	return os.MkdirAll(cachePath, 0755)
}

// getFileHash computes a hash of the file content
func getFileHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// getCachedBinaryPath returns the path where the cached binary should be stored
func getCachedBinaryPath(sourceFile string) (string, error) {
	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return "", err
	}

	// Use hash of absolute path to create unique cache key
	hash := sha256.Sum256([]byte(absPath))
	hashStr := hex.EncodeToString(hash[:])[:16]

	cacheDir := getCacheDir()
	binaryName := filepath.Base(sourceFile)
	// Remove .go extension if present
	if ext := filepath.Ext(binaryName); ext == ".go" {
		binaryName = binaryName[:len(binaryName)-len(ext)]
	}

	return filepath.Join(cacheDir, fmt.Sprintf("%s-%s", binaryName, hashStr)), nil
}

// needsRecompile checks if the source file needs to be recompiled
// It considers the source file, go.mod, and go.sum timestamps
func needsRecompile(sourceFile, cachedBinary string) (bool, error) {
	// If cached binary doesn't exist, we need to compile
	if _, err := os.Stat(cachedBinary); os.IsNotExist(err) {
		return true, nil
	}

	// Get cached binary modification time
	cachedInfo, err := os.Stat(cachedBinary)
	if err != nil {
		return true, err
	}
	cachedTime := cachedInfo.ModTime()

	// Get source file modification time
	sourceInfo, err := os.Stat(sourceFile)
	if err != nil {
		return true, err
	}

	// If source is newer than cached binary, we need to recompile
	if sourceInfo.ModTime().After(cachedTime) {
		return true, nil
	}

	// Check if go.mod or go.sum changed
	absSourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		return true, err
	}

	sourceDir := filepath.Dir(absSourceFile)

	// Check go.mod
	goModPath := filepath.Join(sourceDir, "go.mod")
	if modInfo, err := os.Stat(goModPath); err == nil {
		if modInfo.ModTime().After(cachedTime) {
			return true, nil
		}
	}

	// Check go.sum
	goSumPath := filepath.Join(sourceDir, "go.sum")
	if sumInfo, err := os.Stat(goSumPath); err == nil {
		if sumInfo.ModTime().After(cachedTime) {
			return true, nil
		}
	}

	return false, nil
}

// hasShebang checks if a file starts with a shebang line
func hasShebang(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read first two bytes
	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		// EOF means file is shorter than 2 bytes (no shebang possible)
		// Other errors should be propagated
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}

	return header[0] == '#' && header[1] == '!', nil
}

// createTempFileWithoutShebang creates a temporary copy of the source file without the shebang line
// Returns the temp file path and a cleanup function
func createTempFileWithoutShebang(sourceFile string) (string, func(), error) {
	data, err := os.ReadFile(sourceFile)
	if err != nil {
		return "", nil, err
	}

	// Find the first newline and skip the shebang line
	lines := string(data)
	if idx := strings.Index(lines, "\n"); idx != -1 {
		lines = lines[idx+1:]
	}

	// Create temporary file in system temp directory to avoid polluting script directory
	// Use the original filename to make error messages clearer
	baseName := filepath.Base(sourceFile)
	tmpFile, err := os.CreateTemp("", "grun_"+baseName+"_*.go")
	if err != nil {
		return "", nil, err
	}

	if _, err := tmpFile.WriteString(lines); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, err
	}

	tmpFile.Close()
	tmpPath := tmpFile.Name()

	// Return cleanup function
	cleanup := func() {
		os.Remove(tmpPath)
	}

	return tmpPath, cleanup, nil
}

// compileGoFile compiles a Go file to the specified output path
// It detects if a go.mod exists and builds the package accordingly
// It also handles shebang lines by creating a temporary file without them
func compileGoFile(sourceFile, outputPath string) error {
	absSourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	sourceDir := filepath.Dir(absSourceFile)
	goModPath := filepath.Join(sourceDir, "go.mod")

	// Check if file has a shebang
	hasShebangLine, err := hasShebang(absSourceFile)
	if err != nil {
		return fmt.Errorf("failed to check for shebang: %w", err)
	}

	var fileToCompile string
	var cleanup func()

	if hasShebangLine {
		// Create temporary file without shebang in system temp directory
		tmpFile, cleanupFunc, err := createTempFileWithoutShebang(absSourceFile)
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %w", err)
		}
		fileToCompile = tmpFile
		cleanup = cleanupFunc
		// Ensure cleanup happens even if compilation fails
		defer cleanup()
	} else {
		fileToCompile = absSourceFile
		cleanup = nil
	}

	var cmd *exec.Cmd
	hasGoMod := false

	// Check if go.mod exists in the source directory
	if _, err := os.Stat(goModPath); err == nil {
		hasGoMod = true

		// For module-based builds with shebang, we build the specific file
		// The temp file is in a different directory, so we provide its absolute path
		if hasShebangLine {
			cmd = exec.Command("go", "build", "-o", outputPath, fileToCompile)
			cmd.Dir = sourceDir
		} else {
			// No shebang - build the package normally
			cmd = exec.Command("go", "build", "-o", outputPath)
			cmd.Dir = sourceDir
		}
	} else {
		// No go.mod - build single file (original behavior)
		cmd = exec.Command("go", "build", "-o", outputPath, fileToCompile)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil && !hasGoMod {
		// If build failed without go.mod, provide helpful hint
		return fmt.Errorf("build failed: %w\n\nHint: If your script uses external dependencies, initialize a Go module:\n  cd %s\n  go mod init <module-name>\n  go get <dependencies>", err, sourceDir)
	}

	return err
}

// runBinary executes the binary with the given arguments
func runBinary(binaryPath string, args []string) error {
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <go-file> [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s script.go\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s script.go arg1 arg2\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := flag.Arg(0)

	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}

	// Ensure cache directory exists
	if err := ensureCacheDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating cache directory: %v\n", err)
		os.Exit(1)
	}

	// Get cached binary path
	cachedBinary, err := getCachedBinaryPath(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining cache path: %v\n", err)
		os.Exit(1)
	}

	// Check if we need to recompile
	needsCompile, err := needsRecompile(sourceFile, cachedBinary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking file status: %v\n", err)
		os.Exit(1)
	}

	if needsCompile {
		if err := compileGoFile(sourceFile, cachedBinary); err != nil {
			fmt.Fprintf(os.Stderr, "Error compiling: %v\n", err)
			os.Exit(1)
		}
	}

	// Run the binary with remaining arguments
	args := flag.Args()[1:]
	if err := runBinary(cachedBinary, args); err != nil {
		// If it's an exit error, use its exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running binary: %v\n", err)
		os.Exit(1)
	}
}
