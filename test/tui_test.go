package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
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

func TestTUI_SearchSelectionRestoresListAndCursor(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")
	allExercises := m.AllExercises()
	totalCount := len(allExercises)

	// Activate search
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)

	// Search for "dynamic01"
	targetName := "dynamic01"
	for _, ch := range targetName {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updated.(tui.Model)
	}

	filtered := model.FilteredExercises()
	if len(filtered) != 1 || filtered[0].Name != targetName {
		t.Fatalf("Expected 1 filtered result '%s', got %d", targetName, len(filtered))
	}

	// Confirm selection with Enter
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(tui.Model)

	// Verify search mode exited and search query reset
	if model.IsSearching() {
		t.Error("Expected searching to be false after Enter")
	}
	if model.SearchQuery() != "" {
		t.Errorf("Expected search input to be reset to empty, got %q", model.SearchQuery())
	}

	// Verify all exercises are restored in filteredIndices
	restoredList := model.FilteredExercises()
	if len(restoredList) != totalCount {
		t.Fatalf("Expected full list of %d exercises restored, got %d", totalCount, len(restoredList))
	}

	// Find the expected global index of dynamic01
	expectedGlobalIdx := -1
	for i, ex := range allExercises {
		if ex.Name == targetName {
			expectedGlobalIdx = i
			break
		}
	}
	if expectedGlobalIdx == -1 {
		t.Fatalf("Target exercise %s not found in manifest", targetName)
	}

	// Verify cursor points to the global index
	if model.Cursor() != expectedGlobalIdx {
		t.Errorf("Expected cursor to be global index %d, got %d", expectedGlobalIdx, model.Cursor())
	}
	if model.SelectedExercise().Name != targetName {
		t.Errorf("Expected selected exercise to be %s, got %s", targetName, model.SelectedExercise().Name)
	}

	// Verify subsequent navigation moves relative to the global list
	if expectedGlobalIdx < totalCount-1 {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updated.(tui.Model)
		if model.Cursor() != expectedGlobalIdx+1 {
			t.Errorf("Expected cursor to advance to %d after KeyDown, got %d", expectedGlobalIdx+1, model.Cursor())
		}
		if model.SelectedExercise().Name != allExercises[expectedGlobalIdx+1].Name {
			t.Errorf("Expected selected exercise %s, got %s", allExercises[expectedGlobalIdx+1].Name, model.SelectedExercise().Name)
		}
	}
}

func TestTUI_ViewportScrolling(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	bin, _ := detector.DetectBinary("")
	r := runner.NewRunner(bin)
	m := manifest.GetManifest()

	model := tui.NewModel(r, m, store, "exercises")

	// Set a window size with bounded height to exercise scroll viewport
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 20})
	model = updated.(tui.Model)

	// Toggle hints to add substantial content lines
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	model = updated.(tui.Model)

	// Switch active pane to PaneViewport
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updated.(tui.Model)
	if model.ActivePane() != tui.PaneViewport {
		t.Fatalf("Expected active pane PaneViewport, got %v", model.ActivePane())
	}

	if model.Viewport().YOffset != 0 {
		t.Fatalf("Expected initial viewport YOffset 0, got %d", model.Viewport().YOffset)
	}

	// Scroll down with 'j' and Down arrow
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	model = updated.(tui.Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(tui.Model)

	if model.Viewport().YOffset <= 0 {
		t.Errorf("Expected viewport YOffset > 0 after scrolling down, got %d", model.Viewport().YOffset)
	}

	// Scroll back up with 'k' and Up arrow
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(tui.Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(tui.Model)

	if model.Viewport().YOffset != 0 {
		t.Errorf("Expected viewport YOffset to return to 0 after scrolling up, got %d", model.Viewport().YOffset)
	}

	// First Esc in viewport closes hints
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(tui.Model)
	if model.IsHintsVisible() {
		t.Error("Expected hints to be dismissed after first Esc")
	}

	// Second Esc switches back to sidebar
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(tui.Model)
	if model.ActivePane() != tui.PaneSidebar {
		t.Errorf("Expected active pane PaneSidebar after second Esc, got %v", model.ActivePane())
	}
}

func TestTUI_SpinnerAnimationLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	// Use nil runner for fast, deterministic unit test without subprocess execution
	model := tui.NewModel(nil, m, store, "exercises")

	// 1. Fresh model starts with running = true
	if !model.IsRunning() {
		t.Error("Expected model to be initialized with running = true")
	}

	// 2. While running, spinner.TickMsg returns a non-nil chaining command
	updated, tickCmd := model.Update(spinner.TickMsg{})
	model = updated.(tui.Model)
	if tickCmd == nil {
		t.Error("Expected non-nil tick command when model is running")
	}

	// 3. Execution command finishes and produces result message
	// Run Init command to retrieve execution completion message
	initCmd := model.Init()
	if initCmd != nil {
		msg := initCmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, cmd := range batchMsg {
				if cmd != nil {
					subMsg := cmd()
					if subMsg != nil {
						updated, _ = model.Update(subMsg)
						model = updated.(tui.Model)
					}
				}
			}
		}
	}

	// Verify running transitioned to false
	if model.IsRunning() {
		t.Error("Expected running to be false after exercise evaluation completes")
	}

	// 4. When not running, spinner.TickMsg should return nil (tick chain stopped)
	updated, idleTickCmd := model.Update(spinner.TickMsg{})
	model = updated.(tui.Model)
	if idleTickCmd != nil {
		t.Error("Expected nil tick command when model is not running")
	}

	// 5. Simulate KeyEnter to re-run exercise -> running becomes true and returns non-nil command
	updated, enterCmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(tui.Model)
	if !model.IsRunning() {
		t.Error("Expected model to be running after KeyEnter")
	}
	if enterCmd == nil {
		t.Error("Expected non-nil command returned on KeyEnter")
	}
}

func TestTUI_WatcherErrorStatus(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Initially no watcher warning
	view := model.View()
	if strings.Contains(view, "⚠️ Watcher:") {
		t.Error("Expected no watcher warning initially")
	}

	// Update with WatchErrMsg
	testErr := fmt.Errorf("permission denied")
	updated, cmd := model.Update(tui.WatchErrMsg{Err: testErr})
	model = updated.(tui.Model)

	view = model.View()
	if !strings.Contains(view, "⚠️ Watcher: permission denied") {
		t.Errorf("Expected watcher error in footer view, got:\n%s", view)
	}
	if cmd == nil {
		t.Error("Expected non-nil reset timer command after watcher error")
	}
}

func TestTUI_SearchFilter_MultiMatchResetCursor(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Move cursor down several exercises first
	for i := 0; i < 5; i++ {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updated.(tui.Model)
	}
	if model.Cursor() != 5 {
		t.Fatalf("Expected cursor 5, got %d", model.Cursor())
	}

	// Enter search mode
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)

	// Search for query matching multiple exercises (e.g. 'primitive')
	for _, ch := range "primitive" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updated.(tui.Model)
	}

	// Cursor should reset to 0 (the top match)
	if model.Cursor() != 0 {
		t.Errorf("Expected cursor to reset to 0 after typing search query, got %d", model.Cursor())
	}

	// Press Enter to confirm first match
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(tui.Model)

	// List is restored to full length
	if len(model.FilteredIndices()) != len(m.AllExercises()) {
		t.Errorf("Expected full list restored, got %d / %d", len(model.FilteredIndices()), len(m.AllExercises()))
	}
	// Selected exercise matches the first match (primitives01)
	if model.SelectedExercise().Name != "primitives01" {
		t.Errorf("Expected selected exercise primitives01, got %s", model.SelectedExercise().Name)
	}
}

func TestTUI_RunSkippedWhenNilRunner(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Trigger initial command execution
	initCmd := model.Init()
	if initCmd != nil {
		msg := initCmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, cmd := range batchMsg {
				if cmd != nil {
					subMsg := cmd()
					if subMsg != nil {
						updated, _ := model.Update(subMsg)
						model = updated.(tui.Model)
					}
				}
			}
		}
	}

	// LastResult must remain nil (not a synthetic failed attempt)
	if model.LastResult() != nil {
		t.Errorf("Expected LastResult to be nil when execution was skipped, got %+v", model.LastResult())
	}

	// Store should not contain any bogus recorded attempts
	allStates := store.GetAllExerciseStates()
	if len(allStates) > 0 {
		t.Errorf("Expected 0 exercise states recorded in store for skipped run, got %d", len(allStates))
	}
}

func TestTUI_SearchFilter_UpDownNavigation(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Enter search mode
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)

	// Search for 'primitive'
	for _, ch := range "primitive" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updated.(tui.Model)
	}

	if model.Cursor() != 0 {
		t.Fatalf("Expected cursor 0 after filter, got %d", model.Cursor())
	}

	// Down arrow moves cursor to 1 (primitives02)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(tui.Model)
	if model.Cursor() != 1 {
		t.Errorf("Expected cursor 1 after Down arrow in search mode, got %d", model.Cursor())
	}

	// Down arrow moves cursor to 2 (primitives03)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(tui.Model)
	if model.Cursor() != 2 {
		t.Errorf("Expected cursor 2 after second Down arrow in search mode, got %d", model.Cursor())
	}

	// Up arrow moves cursor back to 1 (primitives02)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updated.(tui.Model)
	if model.Cursor() != 1 {
		t.Errorf("Expected cursor 1 after Up arrow in search mode, got %d", model.Cursor())
	}

	// Press Enter to confirm selection of second match (primitives02)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(tui.Model)

	if model.SelectedExercise().Name != "primitives02" {
		t.Errorf("Expected selected exercise primitives02, got %s", model.SelectedExercise().Name)
	}
	if len(model.FilteredIndices()) != len(m.AllExercises()) {
		t.Errorf("Expected full list restored after Enter, got %d / %d", len(model.FilteredIndices()), len(m.AllExercises()))
	}
}

func TestTUI_StaleEvaluationDropped(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Start at exercise 0 (evalGen = 1)
	initCmd := model.Init()
	var gen1Msg tea.Msg
	if initCmd != nil {
		msg := initCmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, cmd := range batchMsg {
				if cmd != nil {
					res := cmd()
					if res != nil {
						gen1Msg = res
					}
				}
			}
		}
	}

	// Move cursor down to exercise 1 -> triggers evalGen = 2
	updated, enterCmd := model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(tui.Model)

	var gen2Msg tea.Msg
	if enterCmd != nil {
		msg := enterCmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, cmd := range batchMsg {
				if cmd != nil {
					res := cmd()
					if res != nil {
						gen2Msg = res
					}
				}
			}
		}
	}

	// Now send the stale gen1Msg (simulating out-of-order arrival)
	updated, _ = model.Update(gen1Msg)
	model = updated.(tui.Model)

	// Since gen1Msg is stale, running state should STILL be true (waiting for gen2)
	if !model.IsRunning() {
		t.Error("Expected model to remain running when stale generation result arrives")
	}

	// Now send gen2Msg
	updated, _ = model.Update(gen2Msg)
	model = updated.(tui.Model)

	// Now running state transitions to false
	if model.IsRunning() {
		t.Error("Expected model to stop running after current generation result arrives")
	}
}

func TestTUI_SearchFilter_QueryWithJK(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Enter search mode
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)

	// Search for query containing 'j' and 'k', e.g. "json" or "block"
	query := "block"
	for _, ch := range query {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updated.(tui.Model)
	}

	if model.SearchQuery() != "block" {
		t.Errorf("Expected search query 'block', got '%s'", model.SearchQuery())
	}

	// Verify 'k' can also be typed into query
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	model = updated.(tui.Model)
	if model.SearchQuery() != "blockk" {
		t.Errorf("Expected search query 'blockk', got '%s'", model.SearchQuery())
	}
}

func TestTUI_SearchFilter_EscRestoresPreviousSelection(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Move cursor down to index 4
	for i := 0; i < 4; i++ {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		model = updated.(tui.Model)
	}
	if model.Cursor() != 4 {
		t.Fatalf("Expected cursor 4, got %d", model.Cursor())
	}
	expectedExercise := model.SelectedExercise().Name

	// Enter search mode
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	model = updated.(tui.Model)

	// Type query that changes filter and sets cursor to 0
	for _, ch := range "primitive" {
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updated.(tui.Model)
	}

	// Cancel search with Esc
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updated.(tui.Model)

	// Selection and cursor must be restored to 4
	if model.Cursor() != 4 {
		t.Errorf("Expected cursor restored to 4 after Esc, got %d", model.Cursor())
	}
	if model.SelectedExercise().Name != expectedExercise {
		t.Errorf("Expected selected exercise %s, got %s", expectedExercise, model.SelectedExercise().Name)
	}
}

func TestTUI_FileChangedMsg_Reevaluates(t *testing.T) {
	tmpDir := t.TempDir()
	statePath := filepath.Join(tmpDir, "state.json")
	store, _ := state.NewStore(statePath)
	m := manifest.GetManifest()

	model := tui.NewModel(nil, m, store, "exercises")

	// Finish initial run
	initCmd := model.Init()
	if initCmd != nil {
		msg := initCmd()
		if batchMsg, ok := msg.(tea.BatchMsg); ok {
			for _, cmd := range batchMsg {
				if cmd != nil {
					res := cmd()
					if res != nil {
						updated, _ := model.Update(res)
						model = updated.(tui.Model)
					}
				}
			}
		}
	}
	if model.IsRunning() {
		t.Fatal("Expected model to be idle")
	}

	// Send fileChangedMsg
	updated, cmd := model.Update(tui.FileChangedMsgForTest("exercises/01_syntax/syntax01.tf"))
	model = updated.(tui.Model)

	if !model.IsRunning() {
		t.Error("Expected model to transition to running on fileChangedMsg")
	}
	if cmd == nil {
		t.Error("Expected evaluation command returned on fileChangedMsg")
	}
}
