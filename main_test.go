package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetCacheDir(t *testing.T) {
	// Test default cache directory
	originalCache := *cacheDir
	originalEnv := os.Getenv("GRUN_CACHE")
	defer func() {
		*cacheDir = originalCache
		os.Setenv("GRUN_CACHE", originalEnv)
	}()

	// Reset flag
	*cacheDir = ""
	os.Unsetenv("GRUN_CACHE")

	cachePath := getCacheDir()
	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, defaultCacheDir)

	if cachePath != expected {
		t.Errorf("Expected default cache dir %s, got %s", expected, cachePath)
	}
}

func TestGetCacheDirWithFlag(t *testing.T) {
	originalCache := *cacheDir
	defer func() {
		*cacheDir = originalCache
	}()

	*cacheDir = "/tmp/test-cache"
	cachePath := getCacheDir()

	if cachePath != "/tmp/test-cache" {
		t.Errorf("Expected /tmp/test-cache, got %s", cachePath)
	}
}

func TestGetCacheDirWithEnv(t *testing.T) {
	originalCache := *cacheDir
	originalEnv := os.Getenv("GRUN_CACHE")
	defer func() {
		*cacheDir = originalCache
		os.Setenv("GRUN_CACHE", originalEnv)
	}()

	*cacheDir = ""
	os.Setenv("GRUN_CACHE", "/tmp/env-cache")
	cachePath := getCacheDir()

	if cachePath != "/tmp/env-cache" {
		t.Errorf("Expected /tmp/env-cache, got %s", cachePath)
	}
}

func TestEnsureCacheDir(t *testing.T) {
	testDir := "/tmp/grun-test-cache"
	defer os.RemoveAll(testDir)

	originalCache := *cacheDir
	defer func() {
		*cacheDir = originalCache
	}()

	*cacheDir = testDir
	if err := ensureCacheDir(); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Cache directory was not created: %s", testDir)
	}
}

