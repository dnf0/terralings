package test

import (
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/tui"
)

func TestTUI_InitialModel(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to create state store: %v", err)
	}

	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	// Scenario 1: Clean state - first exercise selected
	model := tui.NewModel(r, m, store, "exercises")
	if model.Cursor() != 0 {
		t.Errorf("Expected cursor at 0 for fresh store, got %d", model.Cursor())
	}
	if model.SelectedExercise().Name != "primitives01" {
		t.Errorf("Expected selected exercise primitives01, got %s", model.SelectedExercise().Name)
	}
	if model.ActivePane() != tui.PaneSidebar {
		t.Errorf("Expected active pane PaneSidebar, got %v", model.ActivePane())
	}

	// Scenario 2: First exercise passed - second exercise selected
	if err := store.RecordAttempt("primitives01", "01_primitives", true); err != nil {
		t.Fatalf("Failed to record passed attempt: %v", err)
	}

	model2 := tui.NewModel(r, m, store, "exercises")
	if model2.Cursor() != 1 {
		t.Errorf("Expected cursor at 1 when first exercise passed, got %d", model2.Cursor())
	}
	if model2.SelectedExercise().Name != "primitives02" {
		t.Errorf("Expected selected exercise primitives02, got %s", model2.SelectedExercise().Name)
	}
}

func TestTUI_KeyNavigation(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")

	// 1. Move Down with Down arrow
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(tui.Model)
	if model.Cursor() != 1 {
		t.Errorf("Expected cursor 1 after KeyDown, got %d", model.Cursor())
	}

	// 2. Move Down with 'j'
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(tui.Model)
	if model.Cursor() != 2 {
		t.Errorf("Expected cursor 2 after 'j', got %d", model.Cursor())
	}

	// 3. Move Up with Up arrow
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(tui.Model)
	if model.Cursor() != 1 {
		t.Errorf("Expected cursor 1 after KeyUp, got %d", model.Cursor())
	}

	// 4. Move Up with 'k'
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(tui.Model)
	if model.Cursor() != 0 {
		t.Errorf("Expected cursor 0 after 'k', got %d", model.Cursor())
	}

	// 5. Boundary check at top: Move Up at 0 clamps or stays at 0
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(tui.Model)
	if model.Cursor() < 0 || model.Cursor() >= len(m.AllExercises()) {
		t.Errorf("Cursor out of bounds after KeyUp at 0: %d", model.Cursor())
	}
}

func TestTUI_TabSwitchPane(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")
	if model.ActivePane() != tui.PaneSidebar {
		t.Fatalf("Expected initial pane to be PaneSidebar")
	}

	// Press Tab -> Switch to PaneViewport
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updated.(tui.Model)
	if model.ActivePane() != tui.PaneViewport {
		t.Errorf("Expected active pane PaneViewport after Tab, got %v", model.ActivePane())
	}

	// Press Tab again -> Switch back to PaneSidebar
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updated.(tui.Model)
	if model.ActivePane() != tui.PaneSidebar {
		t.Errorf("Expected active pane PaneSidebar after second Tab, got %v", model.ActivePane())
	}
}

func TestTUI_HintToggle(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")
	ex := model.SelectedExercise()
	totalHints := len(ex.Hints)
	if totalHints < 2 {
		t.Fatalf("Expected at least 2 hints for %s", ex.Name)
	}

	if model.IsHintsVisible() {
		t.Error("Expected hints to be hidden initially")
	}

	// Press 'h' -> Show hint 0
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model = updated.(tui.Model)
	if !model.IsHintsVisible() {
		t.Error("Expected hints to be visible after pressing 'h'")
	}
	if model.HintIndex() != 0 {
		t.Errorf("Expected hint index 0, got %d", model.HintIndex())
	}

	// Press 'h' again -> Show hint 1
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model = updated.(tui.Model)
	if model.HintIndex() != 1 {
		t.Errorf("Expected hint index 1, got %d", model.HintIndex())
	}

	// Cycle through remaining hints and check wrap around
	for i := 2; i < totalHints; i++ {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		model = updated.(tui.Model)
	}
	// One more 'h' wraps back to 0
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model = updated.(tui.Model)
	if model.HintIndex() != 0 {
		t.Errorf("Expected hint index to wrap to 0, got %d", model.HintIndex())
	}

	// Verify state store recorded hint consultation
	exState := store.GetExerciseState(ex.Name)
	if exState == nil || exState.HintsViewed < 1 {
		t.Error("Expected hint consultation to be recorded in state store")
	}
}

func TestTUI_WindowResize(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")

	updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updated.(tui.Model)

	if model.Width() != 120 {
		t.Errorf("Expected width 120, got %d", model.Width())
	}
	if model.Height() != 40 {
		t.Errorf("Expected height 40, got %d", model.Height())
	}
}

func TestTUI_ViewRendering(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = updated.(tui.Model)

	view := model.View()

	if !strings.Contains(view, "TERRALINGS") && !strings.Contains(view, "Terralings") {
		t.Fatalf("Expected header with 'Terralings' in View output, got:\n%s", view)
	}
	if !strings.Contains(view, "primitives01") {
		t.Fatalf("Expected exercise 'primitives01' in View output, got:\n%s", view)
	}
	if !strings.Contains(view, "Navigate") || !strings.Contains(view, "Quit") {
		t.Fatalf("Expected keymap help bar in View output, got:\n%s", view)
	}
}

func TestTUI_SearchFilter(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")

	// 1. Activate search mode with '/'
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)
	if !model.IsSearching() {
		t.Fatal("Expected search mode to be active after '/'")
	}

	// 2. Type "dynamic"
	for _, char := range "dynamic" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
		model = updated.(tui.Model)
	}

	filtered := model.FilteredExercises()
	if len(filtered) == 0 {
		t.Fatal("Expected filtered exercises matching 'dynamic', got 0")
	}
	for _, ex := range filtered {
		if !strings.Contains(ex.Name, "dynamic") && !strings.Contains(ex.ChapterName, "dynamic") && !strings.Contains(strings.ToLower(ex.Title), "dynamic") {
			t.Errorf("Filtered exercise %s does not match 'dynamic'", ex.Name)
		}
	}

	// 3. Press Enter to confirm search selection
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(tui.Model)
	if model.IsSearching() {
		t.Error("Expected search mode to be deactivated after Enter")
	}
	if !strings.Contains(model.SelectedExercise().Name, "dynamic") && !strings.Contains(model.SelectedExercise().ChapterName, "dynamic") && !strings.Contains(strings.ToLower(model.SelectedExercise().Title), "dynamic") {
		t.Errorf("Expected selected exercise to match 'dynamic', got %s", model.SelectedExercise().Name)
	}

	// 4. Test Esc to cancel search
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)
	if !model.IsSearching() {
		t.Fatal("Expected search mode to be active after '/'")
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(tui.Model)
	if model.IsSearching() {
		t.Error("Expected search mode to be deactivated after Esc")
	}
}
