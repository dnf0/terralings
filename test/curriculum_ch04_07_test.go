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

func getChapters4To7Exercises() []models.Exercise {
	m := manifest.GetManifest()
	var result []models.Exercise
	for _, ch := range m.Chapters {
		if ch.Number >= 4 && ch.Number <= 7 {
			result = append(result, ch.Exercises...)
		}
	}
	return result
}

func TestCurriculumChapters04To07Count(t *testing.T) {
	exercises := getChapters4To7Exercises()
	expectedCount := 5 + 5 + 4 + 4 // ch4: 5, ch5: 5, ch6: 4, ch7: 4 = 18
	if len(exercises) != expectedCount {
		t.Fatalf("Expected %d exercises in chapters 4-7, got %d", expectedCount, len(exercises))
	}

	m := manifest.GetManifest()
	ch4 := m.Chapters[3]
	if len(ch4.Exercises) != 5 {
		t.Errorf("Expected Chapter 4 to have 5 exercises, got %d", len(ch4.Exercises))
	}
	ch5 := m.Chapters[4]
	if len(ch5.Exercises) != 5 {
		t.Errorf("Expected Chapter 5 to have 5 exercises, got %d", len(ch5.Exercises))
	}
	ch6 := m.Chapters[5]
	if len(ch6.Exercises) != 4 {
		t.Errorf("Expected Chapter 6 to have 4 exercises, got %d", len(ch6.Exercises))
	}
	ch7 := m.Chapters[6]
	if len(ch7.Exercises) != 4 {
		t.Errorf("Expected Chapter 7 to have 4 exercises, got %d", len(ch7.Exercises))
	}
}

func TestExercisesChapters04To07FailInitial(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping initial exercise run test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters4To7Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			fullPath := filepath.Join(repoRoot, ex.Path)
			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Exercise file missing: %s", fullPath)
			}

			testEx := ex
			testEx.Path = fullPath
			res := r.Run(testEx)
			if res.Passed {
				t.Fatalf("Expected exercise %s to fail initial run, but passed", ex.Name)
			}
		})
	}
}

func TestSolutionsChapters04To07Pass(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping solution test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters4To7Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			solRelPath := ex.SolutionPath()
			fullPath := filepath.Join(repoRoot, solRelPath)

			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Solution file missing: %s", fullPath)
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
