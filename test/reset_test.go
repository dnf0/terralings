package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/exercises"
	"github.com/dnf0/terralings/internal/manifest"
)

func TestReset_SingleFileExercise(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-reset-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := exercises.ExtractAll(tmpDir, true); err != nil {
		t.Fatalf("Failed to extract exercises: %v", err)
	}

	ex := manifest.GetExerciseByName("primitives01")
	if ex == nil {
		t.Fatal("primitives01 not found in manifest")
	}

	rel := strings.TrimPrefix(ex.Path, "exercises/")
	targetFile := filepath.Join(tmpDir, rel)

	// Tamper with the file (remove marker, change content)
	tamperedContent := []byte("// completely modified content without marker\n")
	if err := os.WriteFile(targetFile, tamperedContent, 0644); err != nil {
		t.Fatalf("Failed to overwrite file: %v", err)
	}

	// Reset the exercise
	if err := exercises.ResetExercise("primitives01", tmpDir); err != nil {
		t.Fatalf("ResetExercise failed: %v", err)
	}

	// Read restored file
	restored, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("Failed to read restored file: %v", err)
	}

	if !strings.Contains(string(restored), "primitives01") {
		t.Errorf("Expected reset exercise to contain 'primitives01', got:\n%s", string(restored))
	}
}

func TestReset_DirectoryModuleExercise(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-reset-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := exercises.ExtractAll(tmpDir, true); err != nil {
		t.Fatalf("Failed to extract exercises: %v", err)
	}

	ex := manifest.GetExerciseByName("module04")
	if ex == nil {
		t.Fatal("module04 not found in manifest")
	}

	rel := strings.TrimPrefix(ex.Path, "exercises/")
	targetDir := filepath.Join(tmpDir, rel)

	// Tamper with main.tf inside module04
	mainTf := filepath.Join(targetDir, "main.tf")
	if err := os.WriteFile(mainTf, []byte("// tampered"), 0644); err != nil {
		t.Fatalf("Failed to tamper main.tf: %v", err)
	}

	// Reset exercise
	if err := exercises.ResetExercise("module04", tmpDir); err != nil {
		t.Fatalf("ResetExercise for module04 failed: %v", err)
	}

	// Verify main.tf restored
	restored, err := os.ReadFile(mainTf)
	if err != nil {
		t.Fatalf("Failed to read restored main.tf: %v", err)
	}
	if !strings.Contains(string(restored), "module04") && !strings.Contains(string(restored), "terraform") {
		t.Errorf("Expected reset module04/main.tf to be restored, got:\n%s", string(restored))
	}
}

func TestReset_NonExistentExercise(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-reset-nonexist-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = exercises.ResetExercise("non_existent_exercise", tmpDir)
	if err == nil {
		t.Fatal("Expected error resetting non-existent exercise, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected error to mention 'not found', got: %v", err)
	}
}
