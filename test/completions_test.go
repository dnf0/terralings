package test

import (
	"strings"
	"testing"
)

func TestCLI_CompletionsCommand(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, sh := range shells {
		t.Run(sh, func(t *testing.T) {
			stdout, stderr, exitCode := runCLI(t, "completions", sh)
			if exitCode != 0 {
				t.Fatalf("Expected exit code 0 for 'completions %s', got %d. Stderr: %s", sh, exitCode, stderr)
			}
			if len(stdout) < 100 {
				t.Fatalf("Expected completion script for %s, got:\n%s", sh, stdout)
			}
			if !strings.Contains(strings.ToLower(stdout), "terralings") {
				t.Fatalf("Expected completion script to reference 'terralings', got:\n%s", stdout)
			}
		})
	}
}

func TestCLI_CompletionsInvalidShell(t *testing.T) {
	_, stderr, exitCode := runCLI(t, "completions", "unknown_shell_xyz")
	if exitCode == 0 {
		t.Fatal("Expected non-zero exit code for invalid shell name")
	}
	if !strings.Contains(stderr, "unsupported shell") && !strings.Contains(stderr, "invalid") {
		t.Fatalf("Expected error message for invalid shell, got: %s", stderr)
	}
}

func TestCLI_SearchCommand(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "search", "dynamic")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'search dynamic', got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "dynamic01") {
		t.Fatalf("Expected 'dynamic01' in search output, got:\n%s", stdout)
	}
}
