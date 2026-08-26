package test

import (
	"os"
	"os/exec"
	"strings"
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
	}
	for _, f := range requiredFiles {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("required file %s: %v", f, err)
		}
	}
}

func TestEntrypointSmoke(t *testing.T) {
	bin := getCLIBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run terralings binary: %v\noutput: %s", err, string(out))
	}
	if !strings.Contains(string(out), "terralings") {
		t.Fatalf("expected output to contain 'terralings', got: %s", string(out))
	}
}
