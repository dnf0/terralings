package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	cliBinaryPath string
	buildOnce     sync.Once
	buildErr      error
)

func getCLIBinary(t *testing.T) string {
	t.Helper()
	buildOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "terralings-cli-build-*")
		if err != nil {
			buildErr = err
			return
		}
		binPath := filepath.Join(tmpDir, "terralings")
		cmd := exec.Command("go", "build", "-o", binPath, "../cmd/terralings")
		out, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = err
			return
		}
		cliBinaryPath = binPath
		_ = out
	})

	if buildErr != nil {
		t.Fatalf("Failed to build terralings CLI binary: %v", buildErr)
	}
	return cliBinaryPath
}

func runCLI(t *testing.T, args ...string) (stdout string, stderr string, exitCode int) {
	t.Helper()
	bin := getCLIBinary(t)
	cmd := exec.Command(bin, args...)
	repoRoot, err := filepath.Abs("..")
	if err == nil {
		cmd.Dir = repoRoot
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func TestCLI_ListCommand(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "list")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'list', got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "TERRALINGS") {
		t.Fatalf("Expected banner in 'list' output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Chapter 01") && !strings.Contains(stdout, "Chapter 1") {
		t.Fatalf("Expected Chapter 01 in 'list' output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Chapter 13") {
		t.Fatalf("Expected Chapter 13 in 'list' output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "primitives01") {
		t.Fatalf("Expected 'primitives01' in 'list' output, got:\n%s", stdout)
	}
}

func TestCLI_HintCommand(t *testing.T) {
	t.Run("ValidExercise", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "hint", "primitives01")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'hint primitives01', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "Hint") {
			t.Fatalf("Expected 'Hint' in output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "primitives01") {
			t.Fatalf("Expected 'primitives01' in hint output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "required_version") {
			t.Fatalf("Expected hint text to mention required_version, got:\n%s", stdout)
		}
	})

	t.Run("ValidExerciseWithIndexFlag", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "hint", "primitives01", "--index", "1")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'hint primitives01 --index 1', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "required_providers") {
			t.Fatalf("Expected second hint to mention required_providers, got:\n%s", stdout)
		}
	})

	t.Run("NonExistentExercise", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "hint", "nonexistent_exercise_12345")
		if exitCode == 0 {
			t.Fatal("Expected non-zero exit code for nonexistent exercise hint")
		}
		combined := stdout + stderr
		if !strings.Contains(combined, "not found") {
			t.Fatalf("Expected 'not found' message, got:\nStdout: %s\nStderr: %s", stdout, stderr)
		}
	})
}

func TestCLI_RunCommand(t *testing.T) {
	t.Run("NonExistentExercise", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "run", "nonexistent_exercise_12345")
		if exitCode == 0 {
			t.Fatal("Expected non-zero exit code for nonexistent exercise run")
		}
		combined := stdout + stderr
		if !strings.Contains(combined, "not found") {
			t.Fatalf("Expected 'not found' message, got:\nStdout: %s\nStderr: %s", stdout, stderr)
		}
	})

	t.Run("ExistingExerciseInvocation", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "run", "primitives01")
		combined := stdout + stderr
		if !strings.Contains(combined, "primitives01") {
			t.Fatalf("Expected exercise name 'primitives01' in output, got exitCode=%d:\nStdout: %s\nStderr: %s", exitCode, stdout, stderr)
		}
	})
}

func TestCLI_VersionCommand(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "version")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'version', got %d. Stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "terralings") {
		t.Fatalf("Expected 'terralings' in version output, got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Detected binary:") {
		t.Fatalf("Expected 'Detected binary:' in version output, got:\n%s", stdout)
	}
}

func TestCLI_VerifyCommand(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "verify")
	if !strings.Contains(stdout, "TERRALINGS") {
		t.Fatalf("Expected banner in 'verify' output, got exitCode=%d:\nStdout: %s\nStderr: %s", exitCode, stdout, stderr)
	}
	if !strings.Contains(stdout, "Progress:") {
		t.Fatalf("Expected 'Progress:' in verify output, got:\n%s", stdout)
	}
}

func TestCLI_FlagsAndHelp(t *testing.T) {
	t.Run("HelpFlag", func(t *testing.T) {
		stdout, _, exitCode := runCLI(t, "--help")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for '--help', got %d", exitCode)
		}
		for _, cmdName := range []string{"list", "hint", "run", "watch", "verify", "version"} {
			if !strings.Contains(stdout, cmdName) {
				t.Fatalf("Expected command %q to be listed in --help output, got:\n%s", cmdName, stdout)
			}
		}
	})

	t.Run("InvalidBinaryOverride", func(t *testing.T) {
		_, stderr, exitCode := runCLI(t, "--bin", "/non/existent/bin_xyz", "run", "primitives01")
		if exitCode == 0 {
			t.Fatal("Expected non-zero exit code when invalid --bin override is passed")
		}
		if !strings.Contains(stderr, "specified binary not found") && !strings.Contains(stderr, "not found") {
			t.Fatalf("Expected binary not found error in stderr, got:\n%s", stderr)
		}
	})
}
