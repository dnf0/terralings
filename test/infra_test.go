package test

import (
	"os"
	"testing"
)

func testFilesExist(t *testing.T) {
	requiredFiles := []string{
		"../go.mod",
		"../.gitignore",
		"../Makefile",
		"../LICENSE",
		"../README.md",
		"../CONTRIBUTING.md",
		"../CHANGELOG.md",
		"../.github/workflows/ci.yml",
	}
	for _, f := range requiredFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Fatalf("Required file does not exist: %s", f)
		}
	}
}

func TestInfra(t *testing.T) {
	testFilesExist(t)
}
