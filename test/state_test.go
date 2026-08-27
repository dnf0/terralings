package test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/state"
)

func TestState_NewStoreAndAutoInitialize(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize state store: %v", err)
	}

	if store.GetVersion() != "1.0" {
		t.Errorf("Expected version 1.0, got %s", store.GetVersion())
	}

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Errorf("Expected state file %s to be created on disk", statePath)
	}
}

func TestState_AutoGitignore(t *testing.T) {
	tmpDir := t.TempDir()
	// Simulate a git repo root
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	statePath := filepath.Join(tmpDir, ".terralings", "state.json")
	_, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize state store: %v", err)
	}

	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Expected .gitignore to be created in git repo root: %v", err)
	}

	if string(data) != ".terralings/\n" && string(data) != ".terralings\n" {
		t.Errorf("Expected .gitignore to contain .terralings/, got: %q", string(data))
	}
}

func TestState_RecordAttemptAndHints(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	// 1. Initial state check
	exState := store.GetExerciseState("primitives01")
	if exState != nil {
		t.Fatalf("Expected nil for unstarted exercise, got %+v", exState)
	}

	// 2. Record failed attempt
	if err := store.RecordAttempt("primitives01", "01_primitives", false); err != nil {
		t.Fatalf("RecordAttempt failed: %v", err)
	}

	exState = store.GetExerciseState("primitives01")
	if exState == nil {
		t.Fatal("Expected exercise state to exist after recording attempt")
	}
	if exState.Attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", exState.Attempts)
	}
	if exState.Status != state.StatusInProgress {
		t.Errorf("Expected in_progress, got %s", exState.Status)
	}
	if exState.FirstAttemptAt == nil {
		t.Error("Expected FirstAttemptAt to be set")
	}
	if exState.CompletedAt != nil {
		t.Error("Expected CompletedAt to be nil on failed attempt")
	}

	// 3. Record hint view
	if err := store.RecordHint("primitives01", "01_primitives", 1); err != nil {
		t.Fatalf("RecordHint failed: %v", err)
	}
	exState = store.GetExerciseState("primitives01")
	if exState.HintsViewed != 1 {
		t.Errorf("Expected 1 hint viewed, got %d", exState.HintsViewed)
	}

	// 4. Record passed attempt
	if err := store.RecordAttempt("primitives01", "01_primitives", true); err != nil {
		t.Fatalf("RecordAttempt failed: %v", err)
	}
	exState = store.GetExerciseState("primitives01")
	if exState.Attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", exState.Attempts)
	}
	if exState.Status != state.StatusPassed {
		t.Errorf("Expected passed, got %s", exState.Status)
	}
	if exState.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set on pass")
	}

	// 5. Reload store from disk and verify persistence
	reloadedStore, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to reload store: %v", err)
	}
	reloadedState := reloadedStore.GetExerciseState("primitives01")
	if reloadedState == nil {
		t.Fatal("Expected reloaded state to contain primitives01")
	}
	if reloadedState.Attempts != 2 || reloadedState.Status != state.StatusPassed || reloadedState.HintsViewed != 1 {
		t.Errorf("Reloaded state mismatch: %+v", reloadedState)
	}
}

func TestState_ThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = store.RecordAttempt("primitives01", "01_primitives", idx%2 == 0)
			_ = store.RecordHint("primitives01", "01_primitives", 1)
			_ = store.GetExerciseState("primitives01")
			_ = store.GetAllExerciseStates()
			_ = store.GetAnalytics(manifest.GetManifest())
		}(i)
	}
	wg.Wait()

	exState := store.GetExerciseState("primitives01")
	if exState.Attempts != 50 {
		t.Errorf("Expected 50 total attempts across goroutines, got %d", exState.Attempts)
	}
}

func TestState_AnalyticsCalculations(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	_ = store.RecordAttempt("primitives01", "01_primitives", true)
	_ = store.RecordAttempt("primitives02", "01_primitives", false)
	_ = store.RecordHint("primitives02", "01_primitives", 2)

	m := manifest.GetManifest()
	analytics := store.GetAnalytics(m)

	if analytics.TotalExercises != len(m.AllExercises()) {
		t.Errorf("Expected total exercises %d, got %d", len(m.AllExercises()), analytics.TotalExercises)
	}
	if analytics.CompletedCount != 1 {
		t.Errorf("Expected 1 completed, got %d", analytics.CompletedCount)
	}
	if analytics.InProgressCount != 1 {
		t.Errorf("Expected 1 in progress, got %d", analytics.InProgressCount)
	}
	if analytics.TotalAttempts != 2 {
		t.Errorf("Expected 2 total attempts, got %d", analytics.TotalAttempts)
	}
	if analytics.TotalHintsViewed != 2 {
		t.Errorf("Expected 2 total hints viewed, got %d", analytics.TotalHintsViewed)
	}

	// Verify chapter breakdown
	if len(analytics.ChapterSummaries) != len(m.Chapters) {
		t.Fatalf("Expected %d chapter summaries, got %d", len(m.Chapters), len(analytics.ChapterSummaries))
	}

	var ch01 *state.ChapterSummary
	for i := range analytics.ChapterSummaries {
		if analytics.ChapterSummaries[i].ChapterID == "01_primitives" {
			ch01 = &analytics.ChapterSummaries[i]
			break
		}
	}
	if ch01 == nil {
		t.Fatal("Chapter 01 summary not found")
	}
	if ch01.Completed != 1 {
		t.Errorf("Expected 1 completed in chapter 01, got %d", ch01.Completed)
	}
	if ch01.TotalAttempts != 2 {
		t.Errorf("Expected 2 attempts in chapter 01, got %d", ch01.TotalAttempts)
	}
	if ch01.TotalHints != 2 {
		t.Errorf("Expected 2 hints in chapter 01, got %d", ch01.TotalHints)
	}
}

func TestState_AnalyticsWithNilManifest(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, ".terralings", "state.json")

	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	_ = store.RecordAttempt("primitives01", "01_primitives", true)
	_ = store.RecordAttempt("primitives02", "01_primitives", false)

	analytics := store.GetAnalytics(nil)
	if analytics.TotalExercises != 2 {
		t.Errorf("Expected 2 total exercises from store state, got %d", analytics.TotalExercises)
	}
	if analytics.CompletedCount != 1 {
		t.Errorf("Expected 1 completed, got %d", analytics.CompletedCount)
	}
	if analytics.InProgressCount != 1 {
		t.Errorf("Expected 1 in progress, got %d", analytics.InProgressCount)
	}
}

func TestState_CorruptFileRecovery(t *testing.T) {
	tmpDir := t.TempDir()
	stateDir := filepath.Join(tmpDir, ".terralings")
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		t.Fatalf("Failed to create state dir: %v", err)
	}
	statePath := filepath.Join(stateDir, "state.json")

	// Write invalid corrupt JSON
	if err := os.WriteFile(statePath, []byte("INVALID_JSON_CORRUPTED{{{"), 0644); err != nil {
		t.Fatalf("Failed to write corrupt state file: %v", err)
	}

	// NewStore should recover gracefully and create a valid store
	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Expected graceful recovery on corrupt file, got error: %v", err)
	}

	if store.GetVersion() != "1.0" {
		t.Errorf("Expected version 1.0 on recovered store, got %s", store.GetVersion())
	}

	// Verify a backup file was created
	backupPath := statePath + ".bak"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Expected corrupt file backup at %s", backupPath)
	}

	// Verify new state file is valid
	if err := store.RecordAttempt("primitives01", "01_primitives", true); err != nil {
		t.Fatalf("RecordAttempt on recovered store failed: %v", err)
	}
}
