package detector

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DetectBinary searches for an OpenTofu or Terraform binary.
// Priority:
// 1. Explicit override argument (if non-empty)
// 2. TERRALINGS_BIN environment variable (if non-empty)
// 3. tofu binary on PATH
// 4. terraform binary on PATH
func DetectBinary(override string) (string, error) {
	if override != "" {
		path, err := exec.LookPath(override)
		if err == nil {
			return path, nil
		}
		return "", fmt.Errorf("specified binary not found: %s", override)
	}

	if env := strings.TrimSpace(os.Getenv("TERRALINGS_BIN")); env != "" {
		path, err := exec.LookPath(env)
		if err == nil {
			return path, nil
		}
		return "", fmt.Errorf("binary from TERRALINGS_BIN not found: %s", env)
	}

	if path, err := exec.LookPath("tofu"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("terraform"); err == nil {
		return path, nil
	}

	return "", errors.New("neither 'tofu' nor 'terraform' was found on your PATH. Please install OpenTofu (https://opentofu.org) or Terraform")
}

// GetBinaryVersion returns the version line reported by the binary.
func GetBinaryVersion(binPath string) (string, error) {
	cmd := exec.Command(binPath, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute '%s version': %w (output: %s)", binPath, err, strings.TrimSpace(string(out)))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) == "" {
		return "", fmt.Errorf("empty version output from '%s version'", binPath)
	}

	return strings.TrimSpace(lines[0]), nil
}
