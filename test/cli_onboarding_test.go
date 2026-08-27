package test

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/doctor"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/tour"
)

func TestCLI_TourCommand_NonInteractive(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "tour", "--non-interactive")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'tour --non-interactive', got %d. Stderr: %s", exitCode, stderr)
	}

	expectedTitles := []string{
		"Welcome & Core Philosophy",
		"Anatomy of an Exercise",
		"Continuous Watch & Verification",
		"Interactive TUI, Hints & Analytics",
		"Editor Integration & LSP",
	}

	for _, title := range expectedTitles {
		if !strings.Contains(stdout, title) {
			t.Errorf("Expected tour output to contain %q, but was missing. Output:\n%s", title, stdout)
		}
	}

	for i := 1; i <= 5; i++ {
		stepBadge := "STEP " + string(rune('0'+i)) + " OF 5"
		if !strings.Contains(stdout, stepBadge) {
			t.Errorf("Expected tour output to contain badge %q", stepBadge)
		}
	}
}

func TestCLI_TourCommand_SingleStep(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "tour", "--step", "2", "--non-interactive")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'tour --step 2 --non-interactive', got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "Anatomy of an Exercise") {
		t.Errorf("Expected output to contain step 2 title 'Anatomy of an Exercise', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "STEP 2 OF 5") {
		t.Errorf("Expected output to contain 'STEP 2 OF 5', got:\n%s", stdout)
	}

	// Should not contain step 1 or step 3
	if strings.Contains(stdout, "STEP 1 OF 5") || strings.Contains(stdout, "Welcome & Core Philosophy") {
		t.Errorf("Expected output NOT to contain step 1, but found it in:\n%s", stdout)
	}
	if strings.Contains(stdout, "STEP 3 OF 5") || strings.Contains(stdout, "Continuous Watch & Verification") {
		t.Errorf("Expected output NOT to contain step 3, but found it in:\n%s", stdout)
	}
}

func TestCLI_TourCommand_JSON(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "tour", "--json")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'tour --json', got %d. Stderr: %s", exitCode, stderr)
	}

	var steps []tour.Step
	if err := json.Unmarshal([]byte(stdout), &steps); err != nil {
		t.Fatalf("Failed to parse JSON output from 'tour --json': %v\nOutput:\n%s", err, stdout)
	}

	if len(steps) != 5 {
		t.Fatalf("Expected 5 tour steps in JSON, got %d", len(steps))
	}

	if steps[0].Index != 1 || steps[0].Title != "Welcome & Core Philosophy" {
		t.Errorf("Unexpected step 1: %+v", steps[0])
	}
	if steps[1].Index != 2 || steps[1].Title != "Anatomy of an Exercise" {
		t.Errorf("Unexpected step 2: %+v", steps[1])
	}
}

func TestCLI_DoctorCommand_Text(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "doctor")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'doctor', got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "Terralings Doctor Diagnostic Report") {
		t.Errorf("Expected report header in doctor output, got:\n%s", stdout)
	}

	expectedChecks := []string{
		"IaC Engine Binary",
		"Curriculum Scaffold",
		"Provider Plugin Cache",
		"Progress Persistence Store",
	}

	for _, chk := range expectedChecks {
		if !strings.Contains(stdout, chk) {
			t.Errorf("Expected doctor output to contain check %q, but was missing:\n%s", chk, stdout)
		}
	}
}

func TestCLI_DoctorCommand_JSON(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "doctor", "--json")
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0 for 'doctor --json', got %d. Stderr: %s", exitCode, stderr)
	}

	var report doctor.DiagnosticReport
	if err := json.Unmarshal([]byte(stdout), &report); err != nil {
		t.Fatalf("Failed to parse JSON output from 'doctor --json': %v\nOutput:\n%s", err, stdout)
	}

	if len(report.Checks) < 4 {
		t.Errorf("Expected at least 4 diagnostic checks, got %d", len(report.Checks))
	}

	var foundBinary, foundScaffold, foundCache bool
	for _, c := range report.Checks {
		switch c.ID {
		case "binary":
			foundBinary = true
		case "exercises":
			foundScaffold = true
		case "cache":
			foundCache = true
		}
	}

	if !foundBinary || !foundScaffold || !foundCache {
		t.Errorf("Missing expected check IDs: binary=%v, exercises=%v, cache=%v", foundBinary, foundScaffold, foundCache)
	}
}

func TestCLI_FirstRunWelcomeGuidance(t *testing.T) {
	t.Run("InitCommandDisplaysGuidance", func(t *testing.T) {
		tmpDir := t.TempDir()
		initTarget := filepath.Join(tmpDir, "exercises")

		stdout, stderr, exitCode := runCLI(t, "init", initTarget)
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for init, got %d. Stderr: %s", exitCode, stderr)
		}

		if !strings.Contains(stdout, "terralings tour") {
			t.Errorf("Expected init output to recommend 'terralings tour', got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "terralings doctor") {
			t.Errorf("Expected init output to recommend 'terralings doctor', got:\n%s", stdout)
		}
	})

	t.Run("RootNoArgsDisplaysWelcomeGuidance", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t)
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for terralings with no args, got %d. Stderr: %s", exitCode, stderr)
		}

		if !strings.Contains(stdout, "terralings tour") {
			t.Errorf("Expected root command without args to mention 'terralings tour', got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "terralings doctor") {
			t.Errorf("Expected root command without args to mention 'terralings doctor', got:\n%s", stdout)
		}
	})
}

func TestRunner_PluginCacheDirEnvOverride(t *testing.T) {
	customDir := "/custom/terralings/plugin-cache-dir"
	t.Setenv("TERRALINGS_PLUGIN_CACHE_DIR", customDir)

	got := runner.PluginCacheDir()
	if got != customDir {
		t.Errorf("PluginCacheDir() = %q, want %q", got, customDir)
	}
}
