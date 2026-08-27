package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/doctor"
)

func TestDoctor_RunDiagnostics(t *testing.T) {
	tmpDir := t.TempDir()

	report := doctor.RunDiagnostics(tmpDir, "", "")
	if len(report.Checks) < 4 {
		t.Fatalf("expected at least 4 diagnostic checks, got %d", len(report.Checks))
	}

	// Should contain binary, exercises, cache, state checks
	var foundBinary, foundExercises, foundCache, foundState bool
	for _, check := range report.Checks {
		switch check.ID {
		case "binary":
			foundBinary = true
		case "exercises":
			foundExercises = true
		case "cache":
			foundCache = true
		case "state":
			foundState = true
		}
	}

	if !foundBinary || !foundExercises || !foundCache || !foundState {
		t.Errorf("missing core checks: binary=%v, exercises=%v, cache=%v, state=%v",
			foundBinary, foundExercises, foundCache, foundState)
	}
}

func TestDoctor_FormatReport(t *testing.T) {
	rep := doctor.DiagnosticReport{
		EngineDetected: "opentofu",
		EngineVersion:  "1.8.0",
		WorkspaceDir:   "/tmp/test",
		Checks: []doctor.CheckResult{
			{
				ID:      "binary",
				Name:    "IaC Engine Binary",
				Status:  doctor.StatusPass,
				Message: "Found OpenTofu v1.8.0",
			},
			{
				ID:         "exercises",
				Name:       "Curriculum Scaffold",
				Status:     doctor.StatusWarn,
				Message:    "Missing exercises directory",
				Resolution: "Run `terralings init` to scaffold",
			},
			{
				ID:         "git",
				Name:       "Git Ignore Integration",
				Status:     doctor.StatusFail,
				Message:    "Critical git issue",
				Resolution: "Fix git configuration",
			},
		},
		Passed:       false,
		FailureCount: 1,
		WarningCount: 1,
	}

	formatted := doctor.FormatReport(rep)
	if !strings.Contains(formatted, "Terralings Doctor Diagnostic Report") {
		t.Errorf("expected header in formatted report, got: %s", formatted)
	}
	if !strings.Contains(formatted, "IaC Engine Binary") {
		t.Errorf("expected binary check in formatted report")
	}
	if !strings.Contains(formatted, "Run `terralings init`") {
		t.Errorf("expected resolution message in formatted report")
	}
	if !strings.Contains(formatted, "✓") || !strings.Contains(formatted, "!") || !strings.Contains(formatted, "✗") {
		t.Errorf("expected status icons (✓, !, ✗) in formatted report, got: %s", formatted)
	}
}

func TestDoctor_JSONSerialization(t *testing.T) {
	tmpDir := t.TempDir()
	rep := doctor.RunDiagnostics(tmpDir, "", "")

	jsonData, err := json.Marshal(rep)
	if err != nil {
		t.Fatalf("failed to serialize report to JSON: %v", err)
	}

	var decoded doctor.DiagnosticReport
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("failed to deserialize report JSON: %v", err)
	}

	if decoded.WorkspaceDir != tmpDir {
		t.Errorf("workspace mismatch: got %q, want %q", decoded.WorkspaceDir, tmpDir)
	}
	if len(decoded.Checks) != len(rep.Checks) {
		t.Errorf("checks count mismatch: got %d, want %d", len(decoded.Checks), len(rep.Checks))
	}
}

func TestDoctor_MissingBinary(t *testing.T) {
	tmpDir := t.TempDir()
	invalidBin := "/path/to/non_existent_tofu_binary_12345"

	report := doctor.RunDiagnostics(tmpDir, invalidBin, "")
	if report.Passed {
		t.Fatal("expected report to fail when binary is missing")
	}
	if report.FailureCount == 0 {
		t.Errorf("expected failure count > 0, got %d", report.FailureCount)
	}

	var binaryCheck *doctor.CheckResult
	for i := range report.Checks {
		if report.Checks[i].ID == "binary" {
			binaryCheck = &report.Checks[i]
			break
		}
	}

	if binaryCheck == nil {
		t.Fatal("binary check result not found")
	}
	if binaryCheck.Status != doctor.StatusFail {
		t.Errorf("expected binary status FAIL, got %s", binaryCheck.Status)
	}
	if binaryCheck.Resolution == "" {
		t.Error("expected non-empty resolution instructions for missing binary")
	}
}

func TestDoctor_ScaffoldIntegrity(t *testing.T) {
	// 1. Missing exercises directory
	tmpDirMissing := t.TempDir()
	reportMissing := doctor.RunDiagnostics(tmpDirMissing, "", "")

	var exCheckMissing *doctor.CheckResult
	for i := range reportMissing.Checks {
		if reportMissing.Checks[i].ID == "exercises" {
			exCheckMissing = &reportMissing.Checks[i]
			break
		}
	}
	if exCheckMissing == nil {
		t.Fatal("exercises check result not found for missing dir")
	}
	if exCheckMissing.Status != doctor.StatusWarn {
		t.Errorf("expected StatusWarn for missing exercises dir, got %s", exCheckMissing.Status)
	}
	if !strings.Contains(exCheckMissing.Resolution, "terralings init") {
		t.Errorf("expected resolution to recommend 'terralings init', got: %s", exCheckMissing.Resolution)
	}

	// 2. Complete exercises directory
	tmpDirComplete := t.TempDir()
	exDir := filepath.Join(tmpDirComplete, "exercises")
	if err := os.MkdirAll(filepath.Join(exDir, "01_test"), 0755); err != nil {
		t.Fatalf("failed to create exercises dir: %v", err)
	}
	testFile := filepath.Join(exDir, "01_test", "test.tf")
	if err := os.WriteFile(testFile, []byte("terraform {}"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	reportComplete := doctor.RunDiagnostics(tmpDirComplete, "", "")
	var exCheckComplete *doctor.CheckResult
	for i := range reportComplete.Checks {
		if reportComplete.Checks[i].ID == "exercises" {
			exCheckComplete = &reportComplete.Checks[i]
			break
		}
	}
	if exCheckComplete == nil {
		t.Fatal("exercises check result not found for complete dir")
	}
	if exCheckComplete.Status != doctor.StatusPass {
		t.Errorf("expected StatusPass for complete exercises dir, got %s", exCheckComplete.Status)
	}
	if !strings.Contains(exCheckComplete.Message, "1") {
		t.Errorf("expected exercise count in message, got: %s", exCheckComplete.Message)
	}
}

func TestDoctor_GitIgnore(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}

	// Without .gitignore or without .terralings in .gitignore
	report1 := doctor.RunDiagnostics(tmpDir, "", "")
	var gitCheck1 *doctor.CheckResult
	for i := range report1.Checks {
		if report1.Checks[i].ID == "git" {
			gitCheck1 = &report1.Checks[i]
			break
		}
	}
	if gitCheck1 == nil {
		t.Fatal("git check result not found when .git exists")
	}
	if gitCheck1.Status != doctor.StatusWarn {
		t.Errorf("expected StatusWarn when .terralings not ignored, got %s", gitCheck1.Status)
	}

	// With .gitignore containing .terralings
	tmpDir2 := t.TempDir()
	gitDir2 := filepath.Join(tmpDir2, ".git")
	if err := os.MkdirAll(gitDir2, 0755); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}
	gitignorePath := filepath.Join(tmpDir2, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(".terralings/\n"), 0644); err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	report2 := doctor.RunDiagnostics(tmpDir2, "", "")
	var gitCheck2 *doctor.CheckResult
	for i := range report2.Checks {
		if report2.Checks[i].ID == "git" {
			gitCheck2 = &report2.Checks[i]
			break
		}
	}
	if gitCheck2 == nil {
		t.Fatal("git check result not found when .git exists")
	}
	if gitCheck2.Status != doctor.StatusPass {
		t.Errorf("expected StatusPass when .terralings is ignored, got %s", gitCheck2.Status)
	}
}
