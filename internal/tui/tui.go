package tui

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/fsnotify/fsnotify"
)

// RunTUI launches the full-screen interactive Bubble Tea terminal dashboard.
func RunTUI(ctx context.Context, binPath string, watchDir string, statePath string, in io.Reader, out io.Writer) error {
	if watchDir == "" {
		watchDir = "exercises"
	}

	store, err := state.NewStore(statePath)
	if err != nil {
		return err
	}

	m := manifest.GetManifest()
	r := runner.NewRunner(binPath)
	model := NewModel(r, m, store, watchDir)

	var opts []tea.ProgramOption
	opts = append(opts, tea.WithAltScreen())
	if in != nil {
		opts = append(opts, tea.WithInput(in))
	}
	if out != nil {
		opts = append(opts, tea.WithOutput(out))
	}
	if ctx != nil {
		opts = append(opts, tea.WithContext(ctx))
	}

	p := tea.NewProgram(model, opts...)

	// Start background fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		defer watcher.Close()

		_ = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info != nil && info.IsDir() {
				_ = watcher.Add(path)
			}
			return nil
		})

		var mu sync.Mutex
		var lastEvent time.Time

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
						if strings.HasSuffix(event.Name, ".tf") || strings.HasSuffix(event.Name, ".hcl") {
							mu.Lock()
							now := time.Now()
							if now.Sub(lastEvent) > 150*time.Millisecond {
								lastEvent = now
								mu.Unlock()
								p.Send(fileChangedMsg{path: event.Name})
							} else {
								mu.Unlock()
							}
						}
					}
				case <-watcher.Errors:
				}
			}
		}()
	}

	_, err = p.Run()
	return err
}
