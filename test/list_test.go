package test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/ui"
)

func TestListJSONStructure(t *testing.T) {
	m := manifest.GetManifest()
	if len(m.Chapters) != 13 {
		t.Fatalf("expected 13 chapters, got %d", len(m.Chapters))
	}

	all := m.AllExercises()
	if len(all) != 56 {
		t.Fatalf("expected 56 exercises, got %d", len(all))
	}

	type jsonExercise struct {
		Name    string `json:"name"`
		Title   string `json:"title"`
		Chapter string `json:"chapter"`
		Path    string `json:"path"`
		Mode    string `json:"mode"`
		Status  string `json:"status"`
	}

	var list []jsonExercise
	for _, ch := range m.Chapters {
		for _, ex := range ch.Exercises {
			list = append(list, jsonExercise{
				Name:    ex.Name,
				Title:   ex.Title,
				Chapter: ex.ChapterName,
				Path:    ex.Path,
				Mode:    string(ex.Mode),
				Status:  "not_started",
			})
		}
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("failed to marshal json list: %v", err)
	}

	var parsed []jsonExercise
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal json list: %v", err)
	}

	if len(parsed) != 56 {
		t.Fatalf("expected 56 parsed exercises, got %d", len(parsed))
	}
}

func TestListStateAwareness(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "state.json")

	store, err := state.NewStore(stateFile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	_ = store.RecordAttempt("primitives01", "01_primitives", true)
	_ = store.RecordAttempt("primitives02", "01_primitives", false)

	m := manifest.GetManifest()
	statuses := make(map[string]models.ExerciseStatus)
	for _, ch := range m.Chapters {
		for _, ex := range ch.Exercises {
			if exState := store.GetExerciseState(ex.Name); exState != nil {
				switch exState.Status {
				case state.StatusPassed:
					statuses[ex.Name] = models.StatusCompleted
				case state.StatusInProgress:
					statuses[ex.Name] = models.StatusInProgress
				default:
					statuses[ex.Name] = models.StatusNotStarted
				}
			} else {
				statuses[ex.Name] = models.StatusNotStarted
			}
		}
	}

	if statuses["primitives01"] != models.StatusCompleted {
		t.Fatalf("expected primitives01 to be completed, got %v", statuses["primitives01"])
	}
	if statuses["primitives02"] != models.StatusInProgress {
		t.Fatalf("expected primitives02 to be in_progress, got %v", statuses["primitives02"])
	}
	if statuses["primitives03"] != models.StatusNotStarted {
		t.Fatalf("expected primitives03 to be not_started, got %v", statuses["primitives03"])
	}

	output := ui.FormatChapterList(m, statuses)
	if len(output) == 0 {
		t.Fatal("expected formatted output, got empty string")
	}
}
