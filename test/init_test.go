package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/exercises"
	"github.com/dnf0/terralings/internal/manifest"
)

func TestInit_ExtractsAllExercises(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-init-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = exercises.ExtractAll(tmpDir, false)
	if err != nil {
		t.Fatalf("ExtractAll failed: %v", err)
	}

	m := manifest.GetManifest()
	all := m.AllExercises()
	if len(all) == 0 {
		t.Fatal("Expected manifest to contain exercises, got 0")
	}

	for _, ex := range all {
		// Exercise paths in manifest are like "exercises/01_primitives/primitives01.tf"
		// When extracted to tmpDir, relPath is relative to "exercises/"
		rel := strings.TrimPrefix(ex.Path, "exercises/")
		targetPath := filepath.Join(tmpDir, rel)

		info, err := os.Stat(targetPath)
		if err != nil {
			t.Errorf("Expected exercise file %s to exist in %s, got error: %v", rel, tmpDir, err)
			continue
		}

		if !ex.IsDirectory() && info.IsDir() {
			t.Errorf("Expected %s to be a file, but got directory", rel)
		}
	}
}

func TestInit_FailsOnNonEmptyDirectoryWithoutForce(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-init-nonempty-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dummyFile := filepath.Join(tmpDir, "dummy.txt")
	if err := os.WriteFile(dummyFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	err = exercises.ExtractAll(tmpDir, false)
	if err == nil {
		t.Fatal("Expected ExtractAll to fail on non-empty directory without force, but it succeeded")
	}
	if !strings.Contains(err.Error(), "not empty") {
		t.Errorf("Expected error message to mention 'not empty', got: %v", err)
	}
}

func TestInit_SucceedsOnNonEmptyDirectoryWithForce(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "terralings-init-force-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dummyFile := filepath.Join(tmpDir, "dummy.txt")
	if err := os.WriteFile(dummyFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to write dummy file: %v", err)
	}

	err = exercises.ExtractAll(tmpDir, true)
	if err != nil {
		t.Fatalf("Expected ExtractAll to succeed with force=true, got error: %v", err)
	}

	// Verify dummy file still exists alongside extracted files
	if _, err := os.Stat(dummyFile); err != nil {
		t.Errorf("Expected dummy file to still exist, got error: %v", err)
	}
}
