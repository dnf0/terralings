package test

import (
	"strings"
	"testing"
	"time"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/ui"
)

func TestFormatBanner(t *testing.T) {
	banner := ui.FormatBanner()
	if banner == "" {
		t.Fatal("FormatBanner returned empty string")
	}
	if !strings.Contains(banner, "TERRALINGS") {
		t.Fatalf("Expected banner to contain 'TERRALINGS', got:\n%s", banner)
	}
}

func TestFormatResult_Passed(t *testing.T) {
	ex := models.Exercise{Name: "primitives01", Path: "exercises/01_primitives/primitives01.tf"}
	res := runner.RunResult{
		Exercise: ex,
		Passed:   true,
	}

	rendered := ui.FormatResult(res)
	if !strings.Contains(rendered, "primitives01") {
		t.Fatalf("Expected output to contain exercise name 'primitives01', got:\n%s", rendered)
	}
	if !strings.Contains(rendered, "✓") && !strings.Contains(rendered, "passed") {
		t.Fatalf("Expected output to contain checkmark or 'passed', got:\n%s", rendered)
	}
}

func TestFormatSuccess_And_Prompt(t *testing.T) {
	succ := ui.FormatSuccess("✓ primitives01 passed!")
	if !strings.Contains(succ, "primitives01 passed!") {
		t.Fatalf("Expected success message to contain 'primitives01 passed!', got:\n%s", succ)
	}

	prompt := ui.FormatInteractivePrompt()
	if !strings.Contains(prompt, "Next exercise") || !strings.Contains(prompt, "Quit") {
		t.Fatalf("Expected prompt to contain navigation controls, got:\n%s", prompt)
	}
}

func TestFormatResult_Error(t *testing.T) {
	ex := models.Exercise{Name: "primitives03", Path: "exercises/01_primitives/primitives03.tf"}
	errMsg := "Error: Unsupported block type 'invalid_block'"
	res := runner.RunResult{
		Exercise: ex,
		Passed:   false,
		Error:    errMsg,
	}

	rendered := ui.FormatResult(res)
	if !strings.Contains(rendered, "primitives03") {
		t.Fatalf("Expected output to contain exercise name 'primitives03', got:\n%s", rendered)
	}
	if !strings.Contains(rendered, errMsg) {
		t.Fatalf("Expected output to contain error message %q, got:\n%s", errMsg, rendered)
	}
}

func TestFormatHint(t *testing.T) {
	ex := &models.Exercise{
		Name:  "primitives01",
		Hints: []string{"First hint for primitives01", "Second hint for primitives01"},
	}

	t.Run("ValidHintIndex", func(t *testing.T) {
		h0 := ui.FormatHint(ex, 0)
		if !strings.Contains(h0, "First hint for primitives01") {
			t.Fatalf("Expected first hint content, got:\n%s", h0)
		}
		if !strings.Contains(h0, "primitives01") {
			t.Fatalf("Expected exercise name in hint, got:\n%s", h0)
		}

		h1 := ui.FormatHint(ex, 1)
		if !strings.Contains(h1, "Second hint for primitives01") {
			t.Fatalf("Expected second hint content, got:\n%s", h1)
		}
	})

	t.Run("OutOfBoundsClampsToLast", func(t *testing.T) {
		hOut := ui.FormatHint(ex, 99)
		if !strings.Contains(hOut, "Second hint for primitives01") {
			t.Fatalf("Expected clamp to last hint, got:\n%s", hOut)
		}
	})

	t.Run("NegativeIndexClampsToFirst", func(t *testing.T) {
		hNeg := ui.FormatHint(ex, -5)
		if !strings.Contains(hNeg, "First hint for primitives01") {
			t.Fatalf("Expected clamp to first hint, got:\n%s", hNeg)
		}
	})

	t.Run("NilExercise", func(t *testing.T) {
		hNil := ui.FormatHint(nil, 0)
		if !strings.Contains(strings.ToLower(hNil), "no hints available") {
			t.Fatalf("Expected 'no hints available' for nil exercise, got:\n%s", hNil)
		}
	})

	t.Run("EmptyHints", func(t *testing.T) {
		exEmpty := &models.Exercise{Name: "empty", Hints: []string{}}
		hEmpty := ui.FormatHint(exEmpty, 0)
		if !strings.Contains(strings.ToLower(hEmpty), "no hints available") {
			t.Fatalf("Expected 'no hints available' for empty hints, got:\n%s", hEmpty)
		}
	})
}

func TestFormatProgress(t *testing.T) {
	t.Run("PartialProgress", func(t *testing.T) {
		p := ui.FormatProgress(25, 50)
		if !strings.Contains(p, "25/50") {
			t.Fatalf("Expected '25/50' in progress output, got:\n%s", p)
		}
		if !strings.Contains(p, "50%") {
			t.Fatalf("Expected '50%%' in progress output, got:\n%s", p)
		}
	})

	t.Run("ZeroProgress", func(t *testing.T) {
		p := ui.FormatProgress(0, 50)
		if !strings.Contains(p, "0/50") {
			t.Fatalf("Expected '0/50' in progress output, got:\n%s", p)
		}
		if !strings.Contains(p, "0%") {
			t.Fatalf("Expected '0%%' in progress output, got:\n%s", p)
		}
	})

	t.Run("CompleteProgress", func(t *testing.T) {
		p := ui.FormatProgress(50, 50)
		if !strings.Contains(p, "50/50") {
			t.Fatalf("Expected '50/50' in progress output, got:\n%s", p)
		}
		if !strings.Contains(p, "100%") {
			t.Fatalf("Expected '100%%' in progress output, got:\n%s", p)
		}
	})

	t.Run("ZeroTotalHandled", func(t *testing.T) {
		p := ui.FormatProgress(0, 0)
		if p == "" {
			t.Fatal("Expected non-empty progress string for 0 total")
		}
	})
}

func TestFormatChapterList(t *testing.T) {
	m := manifest.GetManifest()
	statuses := map[string]models.ExerciseStatus{
		"primitives01": models.StatusCompleted,
		"primitives02": models.StatusInProgress,
		"primitives03": models.StatusFailed,
	}

	t.Run("RenderWithStatuses", func(t *testing.T) {
		list := ui.FormatChapterList(m, statuses)
		if !strings.Contains(list, "Chapter 1") && !strings.Contains(list, "Chapter 01") {
			t.Fatalf("Expected chapter header in list, got:\n%s", list)
		}
		if !strings.Contains(list, "primitives01") {
			t.Fatalf("Expected exercise 'primitives01' in list, got:\n%s", list)
		}
		if !strings.Contains(list, "primitives02") {
			t.Fatalf("Expected exercise 'primitives02' in list, got:\n%s", list)
		}
	})

	t.Run("RenderWithNilStatuses", func(t *testing.T) {
		list := ui.FormatChapterList(m, nil)
		if !strings.Contains(list, "primitives01") {
			t.Fatalf("Expected exercise 'primitives01' in list, got:\n%s", list)
		}
	})

	t.Run("RenderNilManifest", func(t *testing.T) {
		list := ui.FormatChapterList(nil, statuses)
		if list == "" {
			t.Fatal("Expected non-empty response or message for nil manifest")
		}
	})
}

func TestFormatAnalytics(t *testing.T) {
	t.Run("EmptyAnalytics", func(t *testing.T) {
		summary := state.AnalyticsSummary{
			TotalExercises: 51,
		}
		rendered := ui.FormatAnalytics(summary)
		if !strings.Contains(rendered, "TERRALINGS LEARNING ANALYTICS") {
			t.Fatalf("Expected title banner in analytics, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "0%") {
			t.Fatalf("Expected 0%% progress, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "Time Invested:") {
			t.Fatalf("Expected Time Invested in analytics, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "Chapter Breakdown:") {
			t.Fatalf("Expected Chapter Breakdown in analytics, got:\n%s", rendered)
		}
	})

	t.Run("PopulatedAnalytics", func(t *testing.T) {
		summary := state.AnalyticsSummary{
			TotalExercises:   50,
			CompletedCount:   10,
			InProgressCount:  5,
			TotalAttempts:    25,
			TotalHintsViewed: 8,
			TotalTimeSpent:   75 * time.Minute,
			ChapterSummaries: []state.ChapterSummary{
				{
					ChapterID:     "01_primitives",
					Title:         "HCL Foundations & Core Primitives",
					Total:         6,
					Completed:     6,
					TotalAttempts: 8,
					TotalHints:    2,
				},
				{
					ChapterID:     "02_variables",
					Title:         "Input Variables, Types & Validations",
					Total:         5,
					Completed:     4,
					TotalAttempts: 17,
					TotalHints:    6,
				},
			},
		}
		rendered := ui.FormatAnalytics(summary)
		if !strings.Contains(rendered, "20%") {
			t.Fatalf("Expected 20%% progress, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "10/50") {
			t.Fatalf("Expected 10/50 completed count, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "1h 15m") {
			t.Fatalf("Expected '1h 15m' time invested, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "25") {
			t.Fatalf("Expected total attempts 25, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "8") {
			t.Fatalf("Expected total hints 8, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "01_primitives") || !strings.Contains(rendered, "100%") {
			t.Fatalf("Expected 01_primitives 100%% in chapter breakdown, got:\n%s", rendered)
		}
		if !strings.Contains(rendered, "02_variables") || !strings.Contains(rendered, "80%") {
			t.Fatalf("Expected 02_variables 80%% in chapter breakdown, got:\n%s", rendered)
		}
	})
}

func TestFormatFirstRunWelcome(t *testing.T) {
	msg := ui.FormatFirstRunWelcome()
	if !strings.Contains(msg, "Welcome to Terralings") {
		t.Errorf("Expected welcome title, got:\n%s", msg)
	}
	if !strings.Contains(msg, "terralings tour") {
		t.Errorf("Expected mention of 'terralings tour', got:\n%s", msg)
	}
	if !strings.Contains(msg, "terralings doctor") {
		t.Errorf("Expected mention of 'terralings doctor', got:\n%s", msg)
	}
	if !strings.Contains(msg, "terralings watch") {
		t.Errorf("Expected mention of 'terralings watch', got:\n%s", msg)
	}
}
