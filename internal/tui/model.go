package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
)

// Pane represents the actively focused split-pane in the dashboard.
type Pane int

const (
	PaneSidebar Pane = iota
	PaneViewport
)

// Internal message types
type fileChangedMsg struct {
	path string
}

// FileChangedMsgForTest constructs a fileChangedMsg for test packages.
func FileChangedMsgForTest(path string) tea.Msg {
	return fileChangedMsg{path: path}
}

type runResultMsg struct {
	result  runner.RunResult
	skipped bool
	gen     int
}

type statusResetMsg struct {
	seq int
}

// WatchErrMsg represents an error emitted by the file watcher.
type WatchErrMsg struct {
	Err error
}

// Model is the main Bubble Tea model for the Terralings terminal dashboard.
type Model struct {
	runner          *runner.Runner
	manifest        *models.Manifest
	store           *state.Store
	watchDir        string
	exercises       []models.Exercise
	filteredIndices []int
	cursor          int
	activePane      Pane
	viewport        viewport.Model
	spinner         spinner.Model
	searchInput     textinput.Model
	searching       bool
	preSearchCursor int
	showHints       bool
	hintIndex       int
	statusMessage   string
	statusSeq       int
	lastResult      *runner.RunResult
	width           int
	height          int
	ready           bool
	running         bool
	evalGen         int
	keys            KeyMap
}

// NewModel constructs a new interactive dashboard Model.
func NewModel(r *runner.Runner, m *models.Manifest, s *state.Store, watchDir string) Model {
	if m == nil {
		m = manifest.GetManifest()
	}
	if watchDir == "" {
		watchDir = "exercises"
	}

	exercises := m.AllExercises()
	filtered := make([]int, len(exercises))
	for i := range exercises {
		filtered[i] = i
	}

	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Prompt = "🔍 "
	ti.CharLimit = 64

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7D7"))

	width := 100
	height := 30
	layout := computeLayout(width, height)
	vp := viewport.New(layout.ViewportWidth, layout.ViewportHeight)

	// Determine starting cursor
	startingCursor := 0
	if s != nil {
		for i, ex := range exercises {
			st := s.GetExerciseState(ex.Name)
			if st == nil || st.Status != state.StatusPassed {
				startingCursor = i
				break
			}
		}
	}

	mod := Model{
		runner:          r,
		manifest:        m,
		store:           s,
		watchDir:        watchDir,
		exercises:       exercises,
		cursor:          startingCursor,
		activePane:      PaneSidebar,
		viewport:        vp,
		spinner:         sp,
		running:         true,
		showHints:       false,
		hintIndex:       0,
		searching:       false,
		searchInput:     ti,
		filteredIndices: filtered,
		keys:            DefaultKeyMap(),
		width:           width,
		height:          height,
		ready:           false,
		evalGen:         1,
	}

	mod.updateViewportContent()
	return mod
}

// Viewport returns the viewport model for inspecting scroll state.
func (m Model) Viewport() viewport.Model {
	return m.viewport
}

// IsRunning returns whether an exercise evaluation is actively running.
func (m Model) IsRunning() bool {
	return m.running
}

// Cursor returns the index of the currently selected exercise in the active list.
func (m Model) Cursor() int {
	return m.cursor
}

// SelectedExercise returns the currently selected exercise model.
func (m Model) SelectedExercise() models.Exercise {
	if len(m.filteredIndices) == 0 || m.cursor < 0 || m.cursor >= len(m.filteredIndices) {
		if len(m.exercises) > 0 {
			return m.exercises[0]
		}
		return models.Exercise{}
	}
	return m.exercises[m.filteredIndices[m.cursor]]
}

// ActivePane returns the currently focused Pane.
func (m Model) ActivePane() Pane {
	return m.activePane
}

// IsHintsVisible returns whether the hint drawer is visible.
func (m Model) IsHintsVisible() bool {
	return m.showHints
}

// HintIndex returns the active hint index.
func (m Model) HintIndex() int {
	return m.hintIndex
}

// IsSearching returns whether search mode is active.
func (m Model) IsSearching() bool {
	return m.searching
}

// SearchQuery returns current search input text.
func (m Model) SearchQuery() string {
	return m.searchInput.Value()
}

// FilteredExercises returns the currently visible exercises matching filter.
func (m Model) FilteredExercises() []models.Exercise {
	res := make([]models.Exercise, len(m.filteredIndices))
	for i, idx := range m.filteredIndices {
		res[i] = m.exercises[idx]
	}
	return res
}

// FilteredIndices returns the indices of currently visible exercises.
func (m Model) FilteredIndices() []int {
	res := make([]int, len(m.filteredIndices))
	copy(res, m.filteredIndices)
	return res
}

// Width returns current terminal width.
func (m Model) Width() int {
	return m.width
}

// Height returns current terminal height.
func (m Model) Height() int {
	return m.height
}

// LastResult returns the most recent execution result.
func (m Model) LastResult() *runner.RunResult {
	return m.lastResult
}

// Init starts initial commands when the TUI boots.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		makeRunCmd(m.runner, m.SelectedExercise(), m.evalGen),
	)
}

func makeRunCmd(r *runner.Runner, ex models.Exercise, gen int) tea.Cmd {
	return func() tea.Msg {
		if r == nil || ex.Name == "" {
			return runResultMsg{result: runner.RunResult{Exercise: ex}, skipped: true, gen: gen}
		}
		res := r.Run(ex)
		return runResultMsg{result: res, skipped: false, gen: gen}
	}
}

func (m *Model) startEvaluationCmd(ex models.Exercise) tea.Cmd {
	m.evalGen++
	gen := m.evalGen
	r := m.runner
	wasRunning := m.running
	m.running = true

	var cmds []tea.Cmd
	if !wasRunning {
		cmds = append(cmds, m.spinner.Tick)
	}
	cmds = append(cmds, makeRunCmd(r, ex, gen))
	return tea.Batch(cmds...)
}

func (m *Model) updateFilter() {
	q := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	if q == "" {
		m.filteredIndices = make([]int, len(m.exercises))
		for i := range m.exercises {
			m.filteredIndices[i] = i
		}
		if m.cursor >= len(m.filteredIndices) {
			if len(m.filteredIndices) > 0 {
				m.cursor = len(m.filteredIndices) - 1
			} else {
				m.cursor = 0
			}
		}
		return
	}

	m.filteredIndices = nil
	for i, ex := range m.exercises {
		if strings.Contains(strings.ToLower(ex.Name), q) ||
			strings.Contains(strings.ToLower(ex.Title), q) ||
			strings.Contains(strings.ToLower(ex.ChapterName), q) {
			m.filteredIndices = append(m.filteredIndices, i)
			continue
		}
		matchedHint := false
		for _, h := range ex.Hints {
			if strings.Contains(strings.ToLower(h), q) {
				matchedHint = true
				break
			}
		}
		if matchedHint {
			m.filteredIndices = append(m.filteredIndices, i)
		}
	}

	if m.cursor >= len(m.filteredIndices) {
		if len(m.filteredIndices) > 0 {
			m.cursor = len(m.filteredIndices) - 1
		} else {
			m.cursor = 0
		}
	}
}
