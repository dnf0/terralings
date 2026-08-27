package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dnf0/terralings/exercises"
)

// Update handles incoming messages and events in the Bubble Tea Elm loop.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Layout sizing
		sidebarWidth := m.width / 3
		if sidebarWidth < 28 {
			sidebarWidth = 28
		}
		if sidebarWidth > 42 {
			sidebarWidth = 42
		}

		contentWidth := m.width - sidebarWidth - 4
		if contentWidth < 20 {
			contentWidth = 20
		}

		contentHeight := m.height - 7
		if contentHeight < 8 {
			contentHeight = 8
		}

		m.viewport.Width = contentWidth
		m.viewport.Height = contentHeight

	case spinner.TickMsg:
		if m.running {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case fileChangedMsg:
		m.running = true
		cmds = append(cmds, m.runExerciseCmd(m.SelectedExercise()))

	case runResultMsg:
		m.running = false
		m.lastResult = &msg.result
		if m.store != nil {
			_ = m.store.RecordAttempt(msg.result.Exercise.Name, msg.result.Exercise.ChapterName, msg.result.Passed)
		}

	case statusResetMsg:
		m.statusMessage = ""

	case tea.KeyMsg:
		if m.searching {
			switch {
			case key.Matches(msg, m.keys.Esc):
				m.searching = false
				m.searchInput.Blur()
				m.searchInput.Reset()
				m.updateFilter()
				return m, nil

			case key.Matches(msg, m.keys.Enter):
				m.searching = false
				m.searchInput.Blur()
				if len(m.filteredIndices) > 0 {
					m.showHints = false
					m.hintIndex = 0
					m.running = true
					return m, m.runExerciseCmd(m.SelectedExercise())
				}
				return m, nil

			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.updateFilter()
				return m, cmd
			}
		}

		// Normal dashboard navigation mode
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Search):
			m.searching = true
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
			}
			return m, nil

		case key.Matches(msg, m.keys.Reset):
			ex := m.SelectedExercise()
			if ex.Name != "" {
				_ = exercises.ResetExercise(ex.Name, m.watchDir)
				m.statusMessage = "🔄 Reset " + ex.Name + " to initial starting code"
				m.running = true
				return m, tea.Batch(
					m.runExerciseCmd(ex),
					tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return statusResetMsg{} }),
				)
			}
			return m, nil

		case key.Matches(msg, m.keys.Enter):
			m.running = true
			return m, m.runExerciseCmd(m.SelectedExercise())

		case key.Matches(msg, m.keys.Esc):
			if m.showHints {
				m.showHints = false
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
					if len(m.filteredIndices) > 0 {
						if m.cursor > 0 {
							m.cursor--
						}
						m.showHints = false
						m.hintIndex = 0
						m.running = true
						return m, m.runExerciseCmd(m.SelectedExercise())
					}
				case key.Matches(msg, m.keys.Down):
					if len(m.filteredIndices) > 0 {
						if m.cursor < len(m.filteredIndices)-1 {
							m.cursor++
						}
						m.showHints = false
						m.hintIndex = 0
						m.running = true
						return m, m.runExerciseCmd(m.SelectedExercise())
					}
				}
			} else if m.activePane == PaneViewport {
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, tea.Batch(cmds...)
}
