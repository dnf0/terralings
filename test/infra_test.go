package test

import (
	"os"
	"testing"
)

func testFilesExist(t *testing.T) {
	t.Helper()
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
		if _, err := os.Stat(f); err != nil {
			t.Errorf("Required file error: %s: %v", f, err)
		}
	}
}

func TestInfra(t *testing.T) {
	testFilesExist(t)
}
