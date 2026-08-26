package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
)

func getChapters8To10Exercises() []models.Exercise {
	m := manifest.GetManifest()
	var result []models.Exercise
	for _, ch := range m.Chapters {
		if ch.Number >= 8 && ch.Number <= 10 {
			result = append(result, ch.Exercises...)
		}
	}
	return result
}

func TestCurriculumChapters08To10Count(t *testing.T) {
	exercises := getChapters8To10Exercises()
	expectedCount := 5 + 4 + 4 // ch8: 5, ch9: 4, ch10: 4 = 13
	if len(exercises) != expectedCount {
		t.Fatalf("Expected %d exercises in chapters 8-10, got %d", expectedCount, len(exercises))
	}

	m := manifest.GetManifest()
	ch8 := m.Chapters[7]
	if len(ch8.Exercises) != 5 {
		t.Errorf("Expected Chapter 8 to have 5 exercises, got %d", len(ch8.Exercises))
	}
	ch9 := m.Chapters[8]
	if len(ch9.Exercises) != 4 {
		t.Errorf("Expected Chapter 9 to have 4 exercises, got %d", len(ch9.Exercises))
	}
	ch10 := m.Chapters[9]
	if len(ch10.Exercises) != 4 {
		t.Errorf("Expected Chapter 10 to have 4 exercises, got %d", len(ch10.Exercises))
	}
}

func TestExercisesChapters08To10FailInitial(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping initial exercise run test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters8To10Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			fullPath := filepath.Join(repoRoot, ex.Path)
			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Exercise file/dir missing: %s", fullPath)
			}

			if !runner.CheckMarker(fullPath) {
				t.Fatalf("Exercise %s MUST have '%s' marker near top", ex.Name, runner.NotDoneMarker)
			}

			testEx := ex
			testEx.Path = fullPath
			res := r.Run(testEx)
			if res.Passed {
				t.Fatalf("Expected exercise %s to fail initial run, but passed", ex.Name)
			}
			if !res.HasNotDoneMarker {
				t.Fatalf("Expected exercise %s HasNotDoneMarker to be true", ex.Name)
			}
		})
	}
}

func TestSolutionsChapters08To10Pass(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping solution test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters8To10Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			solRelPath := ex.SolutionPath()
			fullPath := filepath.Join(repoRoot, solRelPath)

			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Solution file/dir missing: %s", fullPath)
			}

			if runner.CheckMarker(fullPath) {
				t.Fatalf("Solution %s MUST NOT contain '%s' marker", solRelPath, runner.NotDoneMarker)
			}

			solEx := models.Exercise{
				Name:        ex.Name,
				Title:       ex.Title,
				Path:        fullPath,
				ChapterName: ex.ChapterName,
				Hints:       ex.Hints,
				Mode:        ex.Mode,
			}

			res := r.Run(solEx)
			if !res.Passed {
				t.Fatalf("Solution %s failed execution (mode: %s).\nOutput: %s\nError: %s",
					solRelPath, ex.Mode, res.Output, res.Error)
			}
			if res.ExitCode != 0 {
				t.Fatalf("Solution %s returned non-zero exit code %d", solRelPath, res.ExitCode)
			}
		})
	}
}
