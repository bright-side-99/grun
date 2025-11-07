package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
func needsRecompile(sourceFile, cachedBinary string) (bool, error) {
	// If cached binary doesn't exist, we need to compile
	if _, err := os.Stat(cachedBinary); os.IsNotExist(err) {
		return true, nil
	}

	// Get modification times
	sourceInfo, err := os.Stat(sourceFile)
	if err != nil {
		return true, err
	}

	cachedInfo, err := os.Stat(cachedBinary)
	if err != nil {
		return true, err
	}

	// If source is newer than cached binary, we need to recompile
	return sourceInfo.ModTime().After(cachedInfo.ModTime()), nil
}

// compileGoFile compiles a Go file to the specified output path
func compileGoFile(sourceFile, outputPath string) error {
	cmd := exec.Command("go", "build", "-o", outputPath, sourceFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

