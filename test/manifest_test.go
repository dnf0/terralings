package test

import (
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
)

func TestManifestLoadsAllChapters(t *testing.T) {
	m := manifest.GetManifest()
	if len(m.Chapters) != 13 {
		t.Fatalf("Expected 13 chapters, got %d", len(m.Chapters))
	}
	allEx := m.AllExercises()
	if len(allEx) < 50 {
		t.Fatalf("Expected >= 50 exercises, got %d", len(allEx))
	}
	first := allEx[0]
	if first.Name != "primitives01" {
		t.Fatalf("Expected first exercise to be primitives01, got %s", first.Name)
	}
}

func TestGetExerciseByName(t *testing.T) {
	ex := manifest.GetExerciseByName("primitives01")
	if ex == nil {
		t.Fatal("primitives01 not found by short name")
	}
	if ex.Path != "exercises/01_primitives/primitives01.tf" {
		t.Fatalf("Unexpected path: %s", ex.Path)
	}

	exByPath := manifest.GetExerciseByName("exercises/01_primitives/primitives01.tf")
	if exByPath == nil {
		t.Fatal("primitives01 not found by full path")
	}
	if exByPath.Name != "primitives01" {
		t.Fatalf("Unexpected name: %s", exByPath.Name)
	}

	nonExistent := manifest.GetExerciseByName("non_existent_exercise")
	if nonExistent != nil {
		t.Fatalf("Expected nil for non-existent exercise, got %+v", nonExistent)
	}
}

func TestGetNextExercise(t *testing.T) {
	next := manifest.GetNextExercise("primitives01")
	if next == nil {
		t.Fatal("next exercise after primitives01 is nil")
	}
	if next.Name != "primitives02" {
		t.Fatalf("Expected primitives02, got %s", next.Name)
	}

	all := manifest.GetManifest().AllExercises()
	last := all[len(all)-1]
	nextAfterLast := manifest.GetNextExercise(last.Name)
	if nextAfterLast != nil {
		t.Fatalf("Expected nil after last exercise, got %+v", nextAfterLast)
	}
}

func TestExerciseSolutionPath(t *testing.T) {
	ex := models.Exercise{
		Name: "primitives01",
		Path: "exercises/01_primitives/primitives01.tf",
	}
	solPath := ex.SolutionPath()
	if solPath != "solutions/01_primitives/primitives01.tf" {
		t.Fatalf("Expected solutions/01_primitives/primitives01.tf, got %s", solPath)
	}
}

func TestAllExercisesHaveValidFields(t *testing.T) {
	m := manifest.GetManifest()
	allEx := m.AllExercises()

	validModes := map[models.ExerciseMode]bool{
		models.ModeValidate: true,
		models.ModePlan:     true,
		models.ModeTest:     true,
	}

	seenNames := make(map[string]bool)

	for _, ex := range allEx {
		if ex.Name == "" {
			t.Errorf("Exercise has empty Name: %+v", ex)
		}
		if seenNames[ex.Name] {
			t.Errorf("Duplicate exercise name: %s", ex.Name)
		}
		seenNames[ex.Name] = true

		if ex.Title == "" {
			t.Errorf("Exercise %s has empty Title", ex.Name)
		}
		if !strings.HasPrefix(ex.Path, "exercises/") {
			t.Errorf("Exercise %s Path does not start with exercises/: %s", ex.Name, ex.Path)
		}
		if ex.ChapterName == "" {
			t.Errorf("Exercise %s has empty ChapterName", ex.Name)
		}
		if len(ex.Hints) == 0 {
			t.Errorf("Exercise %s has no hints", ex.Name)
		}
		for i, hint := range ex.Hints {
			if strings.TrimSpace(hint) == "" {
				t.Errorf("Exercise %s has empty hint at index %d", ex.Name, i)
			}
		}
		if !validModes[ex.Mode] {
			t.Errorf("Exercise %s has invalid Mode: %s", ex.Name, ex.Mode)
		}
	}
}
