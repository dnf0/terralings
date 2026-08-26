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

func getRepoRoot() string {
	if _, err := os.Stat("exercises"); err == nil {
		return "."
	}
	if _, err := os.Stat("../exercises"); err == nil {
		return ".."
	}
	return "."
}

func getChapters1To3Exercises() []models.Exercise {
	m := manifest.GetManifest()
	var result []models.Exercise
	for _, ch := range m.Chapters {
		if ch.Number >= 1 && ch.Number <= 3 {
			result = append(result, ch.Exercises...)
		}
	}
	return result
}

func TestCurriculumChapters01To03Count(t *testing.T) {
	exercises := getChapters1To3Exercises()
	expectedCount := 6 + 5 + 4 // ch1: 6, ch2: 5, ch3: 4 = 15
	if len(exercises) != expectedCount {
		t.Fatalf("Expected %d exercises in chapters 1-3, got %d", expectedCount, len(exercises))
	}

	m := manifest.GetManifest()
	ch1 := m.Chapters[0]
	if len(ch1.Exercises) != 6 {
		t.Errorf("Expected Chapter 1 to have 6 exercises, got %d", len(ch1.Exercises))
	}
	ch2 := m.Chapters[1]
	if len(ch2.Exercises) != 5 {
		t.Errorf("Expected Chapter 2 to have 5 exercises, got %d", len(ch2.Exercises))
	}
	ch3 := m.Chapters[2]
	if len(ch3.Exercises) != 4 {
		t.Errorf("Expected Chapter 3 to have 4 exercises, got %d", len(ch3.Exercises))
	}
}

func TestExercisesChapters01To03FailInitial(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping initial exercise run test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters1To3Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			fullPath := filepath.Join(repoRoot, ex.Path)
			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Exercise file missing: %s", fullPath)
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

func TestSolutionsChapters01To03Pass(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping solution test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters1To3Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			solRelPath := ex.SolutionPath()
			fullPath := filepath.Join(repoRoot, solRelPath)

			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Solution file missing: %s", fullPath)
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
