package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/version"
)

var semverRegex = regexp.MustCompile(`^\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`)

func TestVersion_SingleSourceOfTruth(t *testing.T) {
	if version.Number == "" {
		t.Fatal("version.Number must not be empty")
	}
	if !semverRegex.MatchString(version.Number) {
		t.Fatalf("version.Number %q does not match valid SemVer format", version.Number)
	}
	if version.String != "v"+version.Number {
		t.Fatalf("version.String %q must match 'v' + version.Number %q", version.String, version.Number)
	}
}

func TestVersion_VSCodeExtensionMatches(t *testing.T) {
	pkgPath, err := filepath.Abs("../extensions/vscode/package.json")
	if err != nil {
		t.Fatalf("Failed to resolve package.json path: %v", err)
	}
	pkgData, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("Failed to read package.json at %s: %v", pkgPath, err)
	}

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(pkgData, &pkg); err != nil {
		t.Fatalf("Failed to parse package.json: %v", err)
	}

	if pkg.Version != version.Number {
		t.Fatalf("VS Code extension version %q in package.json does not match Go version.Number %q", pkg.Version, version.Number)
	}

	// Also check package-lock.json if present
	lockPath, err := filepath.Abs("../extensions/vscode/package-lock.json")
	if err == nil {
		if lockData, err := os.ReadFile(lockPath); err == nil {
			var lock struct {
				Version string `json:"version"`
			}
			if err := json.Unmarshal(lockData, &lock); err == nil && lock.Version != "" {
				if lock.Version != version.Number {
					t.Fatalf("VS Code extension version %q in package-lock.json does not match Go version.Number %q", lock.Version, version.Number)
				}
			}
		}
	}
}

func TestVersion_ChangelogContainsCurrentVersion(t *testing.T) {
	changelogPath, err := filepath.Abs("../CHANGELOG.md")
	if err != nil {
		t.Fatalf("Failed to resolve CHANGELOG.md path: %v", err)
	}
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Fatalf("Failed to read CHANGELOG.md: %v", err)
	}

	expectedHeading := "## [" + version.Number + "]"
	if !strings.Contains(string(content), expectedHeading) {
		t.Fatalf("CHANGELOG.md does not contain release heading %q", expectedHeading)
	}
}

func TestVersion_CLIVersionCommandOutput(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "version")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'version', got %d. Stderr: %s", exitCode, stderr)
	}
	expectedOutput := "terralings " + version.String
	if !strings.Contains(stdout, expectedOutput) {
		t.Fatalf("Expected CLI version output to contain %q, got:\n%s", expectedOutput, stdout)
	}
}
