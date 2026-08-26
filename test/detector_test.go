package test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/detector"
)

func TestDetectBinarySystem(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		if !strings.Contains(err.Error(), "neither 'tofu' nor 'terraform' was found") {
			t.Fatalf("Unexpected error message when binary not found: %v", err)
		}
		return
	}

	base := filepath.Base(bin)
	if base != "tofu" && base != "terraform" && !strings.HasPrefix(base, "tofu") && !strings.HasPrefix(base, "terraform") {
		t.Fatalf("Expected detected binary to be tofu or terraform, got: %s (base: %s)", bin, base)
	}

	// Test fallback error when PATH is empty
	t.Run("EmptyPATH", func(t *testing.T) {
		t.Setenv("PATH", "")
		t.Setenv("TERRALINGS_BIN", "")
		_, emptyErr := detector.DetectBinary("")
		if emptyErr == nil {
			t.Fatal("Expected error when PATH is empty and no binary available")
		}
		if !strings.Contains(emptyErr.Error(), "neither 'tofu' nor 'terraform' was found") {
			t.Fatalf("Expected helpful error message, got: %v", emptyErr)
		}
	})
}

func TestDetectBinaryOverride(t *testing.T) {
	// Valid override with an existing binary
	goBin, err := exec.LookPath("go")
	if err == nil {
		detected, err := detector.DetectBinary(goBin)
		if err != nil {
			t.Fatalf("Failed to detect valid override path %s: %v", goBin, err)
		}
		if detected != goBin {
			t.Fatalf("Expected %s, got %s", goBin, detected)
		}
	}

	// Non-existent override
	_, err = detector.DetectBinary("/non/existent/binary/terralings_test_fake")
	if err == nil {
		t.Fatal("Expected error for non-existent override binary, got nil")
	}
	if !strings.Contains(err.Error(), "specified binary not found") {
		t.Fatalf("Expected 'specified binary not found' error, got: %v", err)
	}
}

func TestDetectBinaryEnv(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err == nil {
		t.Setenv("TERRALINGS_BIN", goBin)
		detected, err := detector.DetectBinary("")
		if err != nil {
			t.Fatalf("Failed to detect binary from TERRALINGS_BIN=%s: %v", goBin, err)
		}
		if detected != goBin {
			t.Fatalf("Expected %s, got %s", goBin, detected)
		}
	}

	// Non-existent binary via env var
	t.Setenv("TERRALINGS_BIN", "non_existent_binary_env_12345")
	_, err = detector.DetectBinary("")
	if err == nil {
		t.Fatal("Expected error for non-existent TERRALINGS_BIN, got nil")
	}
}

func TestGetBinaryVersion(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping version test")
	}

	version, err := detector.GetBinaryVersion(bin)
	if err != nil {
		t.Fatalf("Failed to get binary version for %s: %v", bin, err)
	}

	if version == "" {
		t.Fatal("Expected non-empty version string")
	}

	hasExpectedKeyword := strings.Contains(version, "OpenTofu") ||
		strings.Contains(version, "Terraform") ||
		strings.Contains(version, "v1.") ||
		strings.Contains(version, "v0.")
	if !hasExpectedKeyword {
		t.Fatalf("Unexpected version string: %s", version)
	}

	// Test non-existent binary
	_, err = detector.GetBinaryVersion("/non/existent/bin/path")
	if err == nil {
		t.Fatal("Expected error when getting version of non-existent binary, got nil")
	}
}
