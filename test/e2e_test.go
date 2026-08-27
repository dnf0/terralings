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

func TestE2E_FullCurriculumVerification(t *testing.T) {
	m := manifest.GetManifest()
	if len(m.Chapters) != 13 {
		t.Fatalf("Expected 13 chapters in curriculum manifest, got %d", len(m.Chapters))
	}

	allExercises := m.AllExercises()
	const expectedExerciseCount = 56
	if len(allExercises) != expectedExerciseCount {
		t.Fatalf("Expected %d exercises across curriculum, got %d", expectedExerciseCount, len(allExercises))
	}

	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on PATH; skipping E2E exercise execution")
	}

	repoRoot := getRepoRoot()
	r := runner.NewRunner(bin)

	for _, ex := range allExercises {
		ex := ex
		t.Run(ex.Name, func(t *testing.T) {
			// 1. Verify exercise file exists
			exFullPath := filepath.Join(repoRoot, ex.Path)
			if _, statErr := os.Stat(exFullPath); os.IsNotExist(statErr) {
				t.Fatalf("Exercise file does not exist: %s", exFullPath)
			}

			// 2. Initial run must fail deterministically (exercise broken / incomplete)
			testEx := ex
			testEx.Path = exFullPath
			res := r.Run(testEx)
			if res.Passed {
				t.Fatalf("Expected initial run of exercise %s to fail, but passed", ex.Name)
			}

			// 3. Solution file must exist
			solRelPath := ex.SolutionPath()
			solFullPath := filepath.Join(repoRoot, solRelPath)
			if _, statErr := os.Stat(solFullPath); os.IsNotExist(statErr) {
				t.Fatalf("Solution file does not exist: %s", solFullPath)
			}

			if strings.HasPrefix(ex.Name, "tofu") {
				versionInfo, _ := detector.GetBinaryVersion(bin)
				if !strings.Contains(strings.ToLower(versionInfo), "opentofu") {
					t.Skipf("Skipping %s because it requires OpenTofu (detected binary: %s)", ex.Name, versionInfo)
				}
			}

			// 4. Reference solution execution must pass cleanly
			solEx := models.Exercise{
				Name:        ex.Name,
				Title:       ex.Title,
				Path:        solFullPath,
				ChapterName: ex.ChapterName,
				Hints:       ex.Hints,
				Mode:        ex.Mode,
			}
			solRes := r.Run(solEx)
			if !solRes.Passed {
				t.Fatalf("Solution %s failed execution (mode: %s).\nOutput: %s\nError: %s",
					solRelPath, ex.Mode, solRes.Output, solRes.Error)
			}
			if solRes.ExitCode != 0 {
				t.Fatalf("Solution %s returned non-zero exit code %d.\nOutput: %s\nError: %s",
					solRelPath, solRes.ExitCode, solRes.Output, solRes.Error)
			}
		})
	}
}

func TestE2E_CLIFullWorkflow(t *testing.T) {
	// Subtest 1: list command
	t.Run("ListCommand", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "list")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'list', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "TERRALINGS") {
			t.Errorf("Expected banner in 'list' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Chapter 01") && !strings.Contains(stdout, "Chapter 1") {
			t.Errorf("Expected Chapter 01 in 'list' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Chapter 13") {
			t.Errorf("Expected Chapter 13 in 'list' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "primitives01") {
			t.Errorf("Expected primitives01 in 'list' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "gov03") {
			t.Errorf("Expected gov03 in 'list' output, got:\n%s", stdout)
		}
	})

	// Subtest 2: version command
	t.Run("VersionCommand", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "version")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'version', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "terralings v") {
			t.Errorf("Expected 'terralings v' in version output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Detected binary:") {
			t.Errorf("Expected 'Detected binary:' in version output, got:\n%s", stdout)
		}
	})

	// Subtest 3: verify command
	t.Run("VerifyCommand", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "verify")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'verify', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "TERRALINGS") {
			t.Errorf("Expected banner in 'verify' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Progress:") {
			t.Errorf("Expected progress indicator in 'verify' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Chapter 01") && !strings.Contains(stdout, "Chapter 1") {
			t.Errorf("Expected Chapter 01 in 'verify' output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "Chapter 13") {
			t.Errorf("Expected Chapter 13 in 'verify' output, got:\n%s", stdout)
		}
	})

	// Subtest 4: hint command
	t.Run("HintCommand", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "hint", "primitives01")
		if exitCode != 0 {
			t.Fatalf("Expected exit code 0 for 'hint primitives01', got %d. Stderr: %s", exitCode, stderr)
		}
		if !strings.Contains(stdout, "primitives01") {
			t.Errorf("Expected 'primitives01' in hint output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "required_version") {
			t.Errorf("Expected hint text in output, got:\n%s", stdout)
		}
	})

	// Subtest 5: run command on existing exercise
	t.Run("RunCommand_InitialFails", func(t *testing.T) {
		stdout, stderr, exitCode := runCLI(t, "run", "primitives01")
		if exitCode == 0 {
			t.Fatalf("Expected non-zero exit code for initial run of primitives01, got 0. Stdout: %s", stdout)
		}
		combined := stdout + stderr
		if !strings.Contains(combined, "primitives01") {
			t.Errorf("Expected 'primitives01' in output, got:\n%s", combined)
		}
	})
}
