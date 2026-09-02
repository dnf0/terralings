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

func getChapters14To15Exercises() []models.Exercise {
	m := manifest.GetManifest()
	var result []models.Exercise
	for _, ch := range m.Chapters {
		if ch.Number >= 14 && ch.Number <= 15 {
			result = append(result, ch.Exercises...)
		}
	}
	return result
}

func TestCurriculumChapters14To15Count(t *testing.T) {
	exercises := getChapters14To15Exercises()
	expectedCount := 6 + 6 // ch14: 6, ch15: 6 = 12
	if len(exercises) != expectedCount {
		t.Fatalf("Expected %d exercises in chapters 14-15, got %d", expectedCount, len(exercises))
	}

	m := manifest.GetManifest()
	ch14 := m.Chapters[13]
	if len(ch14.Exercises) != 6 {
		t.Errorf("Expected Chapter 14 to have 6 exercises, got %d", len(ch14.Exercises))
	}
	ch15 := m.Chapters[14]
	if len(ch15.Exercises) != 6 {
		t.Errorf("Expected Chapter 15 to have 6 exercises, got %d", len(ch15.Exercises))
	}
}

func TestExercisesChapters14To15FailInitial(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping initial exercise run test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters14To15Exercises()

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

func TestSolutionsChapters14To15Pass(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping solution test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters14To15Exercises()

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
