package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
)

func TestCheckMarker(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "HashCommentExact",
			content:  "# I AM NOT DONE\nterraform {}",
			expected: true,
		},
		{
			name:     "SlashCommentExact",
			content:  "// I AM NOT DONE\nterraform {}",
			expected: true,
		},
		{
			name:     "BlockCommentLowerCase",
			content:  "/* i am not done */\nterraform {}",
			expected: true,
		},
		{
			name:     "MultipleSpacesAndMixedCase",
			content:  "#   i   AM   not   Done   \nterraform {}",
			expected: true,
		},
		{
			name:     "TabsAndWhitespace",
			content:  "//\tI\tAM\tNOT\tDONE\nterraform {}",
			expected: true,
		},
		{
			name:     "MarkerAbsentDone",
			content:  "# I AM DONE\nterraform {}",
			expected: false,
		},
		{
			name:     "PlainValidConfigNoMarker",
			content:  "terraform {\n  required_version = \">= 1.6.0\"\n}\n",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, tc.name+".tf")
			if err := os.WriteFile(filePath, []byte(tc.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}
			result := runner.CheckMarker(filePath)
			if result != tc.expected {
				t.Fatalf("Expected CheckMarker to be %v, got %v for content:\n%s", tc.expected, result, tc.content)
			}
		})
	}

	// Test non-existent file returns false
	t.Run("NonExistentFile", func(t *testing.T) {
		result := runner.CheckMarker(filepath.Join(tmpDir, "does_not_exist.tf"))
		if result != false {
			t.Fatalf("Expected CheckMarker to return false for non-existent file, got %v", result)
		}
	})
}

func TestRunnerDetectsNotDoneMarker(t *testing.T) {
	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "ex01.tf")
	content := `# I AM NOT DONE
terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "example" {
  input = "marker test"
}
`
	if err := os.WriteFile(exFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write exercise file: %v", err)
	}

	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping runner test")
	}

	r := runner.NewRunner(bin)
	ex := models.Exercise{Name: "ex01", Path: exFile, Mode: models.ModeValidate}
	res := r.Run(ex)

	if res.Passed {
		t.Fatal("Expected exercise with NOT DONE marker to fail (res.Passed == false)")
	}
	if !res.HasNotDoneMarker {
		t.Fatal("Expected HasNotDoneMarker to be true")
	}
}

func TestRunnerExecutesValidConfiguration(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping runner execution test")
	}

	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "valid.tf")
	content := `terraform {
  required_version = ">= 1.6.0"
}

resource "terraform_data" "greeting" {
  input = "Hello, Terralings!"
}
`
	if err := os.WriteFile(exFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write valid exercise file: %v", err)
	}

	r := runner.NewRunner(bin)

	// Mode: Validate
	t.Run("ModeValidate", func(t *testing.T) {
		ex := models.Exercise{Name: "valid_validate", Path: exFile, Mode: models.ModeValidate}
		res := r.Run(ex)
		if !res.Passed {
			t.Fatalf("Expected valid exercise to pass in validate mode, but failed.\nOutput: %s\nError: %s", res.Output, res.Error)
		}
		if res.HasNotDoneMarker {
			t.Fatal("Expected HasNotDoneMarker to be false")
		}
		if res.ExitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d", res.ExitCode)
		}
	})

	// Mode: Plan
	t.Run("ModePlan", func(t *testing.T) {
		ex := models.Exercise{Name: "valid_plan", Path: exFile, Mode: models.ModePlan}
		res := r.Run(ex)
		if !res.Passed {
			t.Fatalf("Expected valid exercise to pass in plan mode, but failed.\nOutput: %s\nError: %s", res.Output, res.Error)
		}
		if res.HasNotDoneMarker {
			t.Fatal("Expected HasNotDoneMarker to be false")
		}
		if res.ExitCode != 0 {
			t.Fatalf("Expected exit code 0, got %d", res.ExitCode)
		}
	})
}

func TestRunnerFailsOnInvalidConfiguration(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping runner test")
	}

	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "invalid.tf")
	content := `terraform {
  required_version = ">= 1.6.0"
}

resource "invalid_unknown_resource_xyz" "bad" {
  unknown_attribute = "bad"
}
`
	if err := os.WriteFile(exFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write invalid exercise file: %v", err)
	}

	r := runner.NewRunner(bin)
	ex := models.Exercise{Name: "invalid_ex", Path: exFile, Mode: models.ModeValidate}
	res := r.Run(ex)

	if res.Passed {
		t.Fatal("Expected invalid exercise to fail (res.Passed == false)")
	}
	if res.HasNotDoneMarker {
		t.Fatal("Expected HasNotDoneMarker to be false for file without marker")
	}
	if res.ExitCode == 0 && res.Error == "" && res.Output == "" {
		t.Fatal("Expected failure exit code or diagnostic message")
	}
}

func TestRunnerPluginCacheDirCreated(t *testing.T) {
	customCache := filepath.Join(t.TempDir(), "custom-cache")
	r := &runner.Runner{
		BinaryPath: "tofu",
		CacheDir:   customCache,
	}
	if err := os.MkdirAll(r.CacheDir, 0755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}

	if info, err := os.Stat(customCache); err != nil || !info.IsDir() {
		t.Fatalf("Expected custom cache dir to exist as a directory: %v", err)
	}
}
