package test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/watcher"
)

// safeBuffer is a thread-safe bytes.Buffer for concurrent test assertions
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func TestWatcher_CleanShutdownOnContextCancel(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping watcher test")
	}

	tmpDir := t.TempDir()
	exDir := filepath.Join(tmpDir, "exercises", "01_test")
	_ = os.MkdirAll(exDir, 0755)

	exFile := filepath.Join(exDir, "ex01.tf")
	_ = os.WriteFile(exFile, []byte("# I AM NOT DONE\nterraform {}"), 0644)

	exercises := []models.Exercise{
		{Name: "ex01", Path: exFile, Mode: models.ModeValidate},
	}

	ctx, cancel := context.WithCancel(context.Background())
	var out safeBuffer

	r := runner.NewRunner(bin)
	doneChan := make(chan error, 1)

	go func() {
		doneChan <- watcher.RunWatchWithExercises(ctx, r, exercises, tmpDir, &out)
	}()

	// Allow watcher to initialize and run initial check
	time.Sleep(100 * time.Millisecond)

	// Cancel context to request clean shutdown
	cancel()

	select {
	case err := <-doneChan:
		if err != nil {
			t.Fatalf("Expected nil error on context cancellation, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Watcher did not shut down cleanly within 2 seconds")
	}

	if !strings.Contains(out.String(), "TERRALINGS") {
		t.Fatalf("Expected banner in output, got:\n%s", out.String())
	}
}

func TestWatcher_DirectoryDiscovery(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping watcher test")
	}

	tmpDir := t.TempDir()
	subDir1 := filepath.Join(tmpDir, "01_primitives")
	subDir2 := filepath.Join(tmpDir, "nested", "deep", "dir")
	_ = os.MkdirAll(subDir1, 0755)
	_ = os.MkdirAll(subDir2, 0755)

	ex1 := filepath.Join(subDir1, "ex1.tf")
	_ = os.WriteFile(ex1, []byte("terraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)

	exercises := []models.Exercise{
		{Name: "ex1", Path: ex1, Mode: models.ModeValidate},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var out safeBuffer
	r := runner.NewRunner(bin)
	err = watcher.RunWatchWithExercises(ctx, r, exercises, tmpDir, &out)
	if err != nil {
		t.Fatalf("Unexpected watcher error: %v", err)
	}

	// Should have passed ex1 and completed
	if !strings.Contains(out.String(), "passed") && !strings.Contains(out.String(), "Congratulations") {
		t.Fatalf("Expected pass/congratulations output, got:\n%s", out.String())
	}
}

func TestWatcher_FileDebounceAndAdvancement(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping watcher test")
	}

	tmpDir := t.TempDir()
	ex1File := filepath.Join(tmpDir, "ex01.tf")
	ex2File := filepath.Join(tmpDir, "ex02.tf")

	// ex1 starts with NOT DONE marker
	_ = os.WriteFile(ex1File, []byte("# I AM NOT DONE\nterraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)
	// ex2 starts with NOT DONE marker
	_ = os.WriteFile(ex2File, []byte("# I AM NOT DONE\nterraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)

	exercises := []models.Exercise{
		{Name: "ex01", Path: ex1File, Mode: models.ModeValidate},
		{Name: "ex02", Path: ex2File, Mode: models.ModeValidate},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var out safeBuffer
	r := runner.NewRunner(bin)
	doneChan := make(chan error, 1)

	go func() {
		doneChan <- watcher.RunWatchWithExercises(ctx, r, exercises, tmpDir, &out)
	}()

	// Wait for initial check to run
	var initialReported bool
	deadlineInitial := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadlineInitial) {
		if strings.Contains(out.String(), "I AM NOT DONE") {
			initialReported = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if !initialReported {
		t.Fatalf("Expected initial output to report 'I AM NOT DONE', got:\n%s", out.String())
	}

	// Rapidly edit file multiple times to test debounce
	for i := 0; i < 3; i++ {
		_ = os.WriteFile(ex1File, []byte("terraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounce and execution
	var advanced bool
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		outStr := out.String()
		if strings.Contains(outStr, "Advancing to next exercise") || strings.Contains(outStr, "ex01 passed") {
			advanced = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !advanced {
		t.Fatalf("Watcher did not advance to next exercise after file update within timeout. Output:\n%s", out.String())
	}

	cancel()
	select {
	case <-doneChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Watcher failed to stop after cancel")
	}
}

func TestWatcher_AllExercisesCompleted(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping watcher test")
	}

	tmpDir := t.TempDir()
	ex1File := filepath.Join(tmpDir, "ex01.tf")
	_ = os.WriteFile(ex1File, []byte("terraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)

	exercises := []models.Exercise{
		{Name: "ex01", Path: ex1File, Mode: models.ModeValidate},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var out safeBuffer
	r := runner.NewRunner(bin)
	err = watcher.RunWatchWithExercises(ctx, r, exercises, tmpDir, &out)
	if err != nil {
		t.Fatalf("Unexpected watcher error: %v", err)
	}

	if !strings.Contains(out.String(), "Congratulations") {
		t.Fatalf("Expected congratulations message when all exercises are complete, got:\n%s", out.String())
	}
}
