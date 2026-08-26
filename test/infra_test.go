package test

import (
	"os"
	"testing"
)

func TestInfra(t *testing.T) {
	requiredFiles := []string{
		"../go.mod",
		"../.gitignore",
		"../Makefile",
		"../LICENSE",
		"../README.md",
		"../CONTRIBUTING.md",
		"../CHANGELOG.md",
		"../.github/workflows/ci.yml",
		"../cmd/terralings/main.go",
	}
	for _, f := range requiredFiles {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("required file %s: %v", f, err)
		}
	}
}
