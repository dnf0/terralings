package tui

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/watcher"
	"github.com/fsnotify/fsnotify"
)

// RunTUI launches the full-screen interactive Bubble Tea terminal dashboard.
func RunTUI(ctx context.Context, binPath string, watchDir string, statePath string, in io.Reader, out io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}
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
	opts = append(opts, tea.WithContext(ctx))

	p := tea.NewProgram(model, opts...)

	// Start background fsnotify watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err == nil {
		defer fsWatcher.Close()

		_ = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info != nil && info.IsDir() {
				_ = fsWatcher.Add(path)
			}
			return nil
		})

		var mu sync.Mutex
		var debounceTimer *time.Timer

		go func() {
			defer func() {
				mu.Lock()
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				mu.Unlock()
			}()

			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-fsWatcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Create != 0 {
						if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
							_ = fsWatcher.Add(event.Name)
						}
					}
					if watcher.IsRelevantFile(event.Name) && (event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0) {
						eventName := event.Name
						mu.Lock()
						if debounceTimer != nil {
							debounceTimer.Stop()
						}
						debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
							p.Send(fileChangedMsg{path: eventName})
						})
						mu.Unlock()
					}
				case err, ok := <-fsWatcher.Errors:
					if !ok {
						return
					}
					if err != nil {
						p.Send(WatchErrMsg{Err: err})
					}
				}
			}
		}()
	}

	_, err = p.Run()
	return err
}
