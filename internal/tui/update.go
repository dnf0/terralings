package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dnf0/terralings/exercises"
)

func (m *Model) updateViewportContent() {
	m.viewport.SetContent(m.mainContent())
	// Re-clamps viewport scroll offset against newly set content lines.
	m.viewport.SetYOffset(m.viewport.YOffset)
}

// Update handles incoming messages and events in the Bubble Tea Elm loop.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		layout := computeLayout(m.width, m.height)
		m.viewport.Width = layout.ViewportWidth
		m.viewport.Height = layout.ViewportHeight
		m.updateViewportContent()
		return m, nil

	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			m.updateViewportContent()
			return m, cmd
		}
		return m, nil

	case fileChangedMsg:
		ex := m.SelectedExercise()
		m.showHints = false
		m.hintIndex = 0
		m.viewport.GotoTop()
		cmd := m.startEvaluationCmd(ex)
		m.updateViewportContent()
		return m, cmd

	case runResultMsg:
		// Drop stale evaluation results from outdated generations
		if msg.gen != m.evalGen {
			return m, nil
		}
		m.running = false
		if !msg.skipped && msg.result.Exercise.Name != "" {
			m.lastResult = &msg.result
			if m.store != nil {
				_ = m.store.RecordAttempt(msg.result.Exercise.Name, msg.result.Exercise.ChapterName, msg.result.Passed)
			}
		}
		m.updateViewportContent()
		return m, nil

	case WatchErrMsg:
		if msg.Err != nil {
			m.statusSeq++
			seq := m.statusSeq
			m.statusMessage = "⚠️ Watcher: " + msg.Err.Error()
			m.updateViewportContent()
			return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusResetMsg{seq: seq} })
		}
		return m, nil

	case statusResetMsg:
		if msg.seq == m.statusSeq {
			m.statusMessage = ""
			m.updateViewportContent()
		}
		return m, nil

	case tea.KeyMsg:
		if m.searching {
			switch {
			case key.Matches(msg, m.keys.Esc):
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.Reset()
				m.updateFilter()
				m.cursor = m.preSearchCursor
				m.updateViewportContent()
				return m, nil

			case key.Matches(msg, m.keys.Enter):
				if len(m.filteredIndices) > 0 {
					selectedIdx := m.filteredIndices[m.cursor]
					m.searching = false
					m.searchInput.Blur()
					m.searchInput.Reset()
					m.updateFilter()
					m.cursor = selectedIdx
					m.showHints = false
					m.hintIndex = 0
					m.viewport.GotoTop()
					cmd := m.startEvaluationCmd(m.SelectedExercise())
					m.updateViewportContent()
					return m, cmd
				}
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.Reset()
				m.updateFilter()
				m.cursor = m.preSearchCursor
				m.updateViewportContent()
				return m, nil

			case msg.Type == tea.KeyUp:
				if m.cursor > 0 {
					m.cursor--
					m.updateViewportContent()
				}
				return m, nil

			case msg.Type == tea.KeyDown:
				if m.cursor < len(m.filteredIndices)-1 {
					m.cursor++
					m.updateViewportContent()
				}
				return m, nil

			default:
				oldVal := m.searchInput.Value()
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				if m.searchInput.Value() != oldVal {
					m.cursor = 0
					m.updateFilter()
				}
				m.updateViewportContent()
				return m, cmd
			}
		}

		// Normal dashboard navigation mode
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Search):
			m.searching = true
			m.preSearchCursor = m.cursor
			m.searchInput.Focus()
			return m, textinput.Blink

		case key.Matches(msg, m.keys.Tab):
			if m.activePane == PaneSidebar {
				m.activePane = PaneViewport
			} else {
				m.activePane = PaneSidebar
			}
			return m, nil

		case key.Matches(msg, m.keys.Hint):
			ex := m.SelectedExercise()
			if len(ex.Hints) > 0 {
				if !m.showHints {
					m.showHints = true
					m.hintIndex = 0
				} else {
					m.hintIndex = (m.hintIndex + 1) % len(ex.Hints)
				}
				if m.store != nil {
					_ = m.store.RecordHint(ex.Name, ex.ChapterName, m.hintIndex+1)
				}
				m.updateViewportContent()
			}
			return m, nil

		case key.Matches(msg, m.keys.Reset):
			ex := m.SelectedExercise()
			if ex.Name != "" {
				_ = exercises.ResetExercise(ex.Name, m.watchDir)
				m.statusSeq++
				seq := m.statusSeq
				m.statusMessage = "🔄 Reset " + ex.Name + " to initial starting code"
				m.viewport.GotoTop()
				cmd := m.startEvaluationCmd(ex)
				m.updateViewportContent()
				return m, tea.Batch(
					cmd,
					tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusResetMsg{seq: seq} }),
				)
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			m.viewport.GotoTop()
			cmd := m.startEvaluationCmd(m.SelectedExercise())
			m.updateViewportContent()
			return m, cmd

		case key.Matches(msg, m.keys.Esc):
			if m.showHints {
				m.showHints = false
				m.updateViewportContent()
				return m, nil
			}
			if m.activePane == PaneViewport {
				m.activePane = PaneSidebar
				return m, nil
			}
			return m, nil

		default:
			if m.activePane == PaneSidebar {
				switch {
				case key.Matches(msg, m.keys.Up):
					if len(m.filteredIndices) > 0 && m.cursor > 0 {
						m.cursor--
						m.showHints = false
						m.hintIndex = 0
						m.viewport.GotoTop()
						cmd := m.startEvaluationCmd(m.SelectedExercise())
						m.updateViewportContent()
						return m, cmd
					}
					return m, nil
				case key.Matches(msg, m.keys.Down):
					if len(m.filteredIndices) > 0 && m.cursor < len(m.filteredIndices)-1 {
						m.cursor++
						m.showHints = false
						m.hintIndex = 0
						m.viewport.GotoTop()
						cmd := m.startEvaluationCmd(m.SelectedExercise())
						m.updateViewportContent()
						return m, cmd
					}
					return m, nil
				}
			} else if m.activePane == PaneViewport {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}
