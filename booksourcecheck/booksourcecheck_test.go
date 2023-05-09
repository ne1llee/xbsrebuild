package booksourcecheck

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func GetTestFilePath(relativePath string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Failed to get current file information")
	}
	currentDir := filepath.Dir(currentFile)
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		}
		currentDir = filepath.Dir(currentDir)
	}
	projectRoot, err := filepath.Abs(currentDir)
	if err != nil {
		return "", fmt.Errorf("Failed to get project root directory: %v", err)
	}
	return filepath.Join(projectRoot, relativePath), nil
}

func TestBookSourceCheckCheck(t *testing.T) {
	testFilePath, err := GetTestFilePath("testdata/test.json")
	if err != nil {
		t.Fatalf("Failed to get test file path: %v", err)
	}
	data, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to ReadFile test file: %v", err)
	}
	bsc, err := NewBookSourceCheck(data)
	if err != nil {
		t.Fatalf("Failed to NewBookSourceCheck: %v", err)
	}
	bsc.Check()
}
