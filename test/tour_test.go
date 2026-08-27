package test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/tour"
)

func TestTour_DefaultStepsCountAndContent(t *testing.T) {
	steps := tour.DefaultSteps()
	if len(steps) != 5 {
		t.Fatalf("expected 5 default steps, got %d", len(steps))
	}

	expectedTitles := []string{
		"Welcome & Core Philosophy",
		"Anatomy of an Exercise",
		"Continuous Watch & Verification",
		"Interactive TUI, Hints & Analytics",
		"Editor Integration & LSP",
	}

	for i, expected := range expectedTitles {
		if steps[i].Index != i+1 {
			t.Errorf("step %d index mismatch: got %d, want %d", i, steps[i].Index, i+1)
		}
		if steps[i].Title != expected {
			t.Errorf("step %d title mismatch: got %q, want %q", i, steps[i].Title, expected)
		}
		if len(steps[i].Body) == 0 {
			t.Errorf("step %d body should not be empty", i)
		}
	}
}

func TestTour_NonInteractiveAllSteps(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = true

	err := tr.Run(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error running non-interactive tour: %v", err)
	}

	output := out.String()
	for _, title := range []string{
		"Welcome & Core Philosophy",
		"Anatomy of an Exercise",
		"Continuous Watch & Verification",
		"Interactive TUI, Hints & Analytics",
		"Editor Integration & LSP",
	} {
		if !strings.Contains(output, title) {
			t.Errorf("output missing expected step title %q", title)
		}
	}
}

func TestTour_SpecificStepRender(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = true

	err := tr.Run(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Continuous Watch & Verification") {
		t.Errorf("expected step 3 output, got: %s", output)
	}
	if strings.Contains(output, "Welcome & Core Philosophy") {
		t.Errorf("should not contain step 1 output when step 3 requested")
	}
}

func TestTour_ExportJSON(t *testing.T) {
	tr := tour.NewTour(nil, nil)
	jsonData, err := tr.ExportJSON()
	if err != nil {
		t.Fatalf("failed to export tour JSON: %v", err)
	}

	var parsed []tour.Step
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(parsed) != 5 {
		t.Fatalf("expected 5 steps in JSON, got %d", len(parsed))
	}
}

func TestTour_InteractiveNavigation(t *testing.T) {
	var out bytes.Buffer
	// Navigate: next (n) -> next (Enter) -> prev (p) -> jump (5) -> quit (q)
	in := strings.NewReader("n\n\np\n5\nq\n")

	tr := tour.NewTour(&out, in)
	tr.NonInteractive = false

	err := tr.Run(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error in interactive tour: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "Editor Integration & LSP") {
		t.Errorf("expected jump to step 5 in output, got: %s", output)
	}
}
