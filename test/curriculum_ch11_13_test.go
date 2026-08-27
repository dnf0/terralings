package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
)

func getChapters11To13Exercises() []models.Exercise {
	m := manifest.GetManifest()
	var result []models.Exercise
	for _, ch := range m.Chapters {
		if ch.Number >= 11 && ch.Number <= 13 {
			result = append(result, ch.Exercises...)
		}
	}
	return result
}

func TestCurriculumChapters11To13Count(t *testing.T) {
	exercises := getChapters11To13Exercises()
	expectedCount := 4 + 3 + 3 // ch11: 4, ch12: 3, ch13: 3 = 10
	if len(exercises) != expectedCount {
		t.Fatalf("Expected %d exercises in chapters 11-13, got %d", expectedCount, len(exercises))
	}

	m := manifest.GetManifest()
	ch11 := m.Chapters[10]
	if len(ch11.Exercises) != 4 {
		t.Errorf("Expected Chapter 11 to have 4 exercises, got %d", len(ch11.Exercises))
	}
	ch12 := m.Chapters[11]
	if len(ch12.Exercises) != 3 {
		t.Errorf("Expected Chapter 12 to have 3 exercises, got %d", len(ch12.Exercises))
	}
	ch13 := m.Chapters[12]
	if len(ch13.Exercises) != 3 {
		t.Errorf("Expected Chapter 13 to have 3 exercises, got %d", len(ch13.Exercises))
	}
}

func TestExercisesChapters11To13FailInitial(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping initial exercise run test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters11To13Exercises()

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

func TestSolutionsChapters11To13Pass(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping solution test")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)
	exercises := getChapters11To13Exercises()

	for _, ex := range exercises {
		t.Run(ex.Name, func(t *testing.T) {
			solRelPath := ex.SolutionPath()
			fullPath := filepath.Join(repoRoot, solRelPath)

			if _, statErr := os.Stat(fullPath); os.IsNotExist(statErr) {
				t.Fatalf("Solution file missing: %s", fullPath)
			}

			if strings.HasPrefix(ex.Name, "tofu") {
				versionInfo, _ := detector.GetBinaryVersion(bin)
				if !strings.Contains(strings.ToLower(versionInfo), "opentofu") {
					t.Skipf("Skipping %s because it requires OpenTofu (detected binary: %s)", ex.Name, versionInfo)
				}
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
