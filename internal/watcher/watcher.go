package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dnf0/terralings/internal/diagnostics"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/ui"
	"github.com/fsnotify/fsnotify"
)

// WatchEvent represents a structured NDJSON event emitted during watch --json mode.
type WatchEvent struct {
	Event            string                   `json:"event"`
	Timestamp        string                   `json:"timestamp"`
	Exercise         *models.Exercise         `json:"exercise,omitempty"`
	Passed           *bool                    `json:"passed,omitempty"`
	HasNotDoneMarker *bool                    `json:"has_not_done_marker,omitempty"`
	Diagnostics      []diagnostics.Diagnostic `json:"diagnostics,omitempty"`
	RawOutput        string                   `json:"raw_output,omitempty"`
	ExitCode         int                      `json:"exit_code,omitempty"`
	CurrentIndex     int                      `json:"current_index,omitempty"`
	TotalCount       int                      `json:"total_count,omitempty"`
}

// RunWatch starts the interactive file watcher using stdout.
func RunWatch(binPath string, watchDir string) error {
	return RunWatchWithContext(context.Background(), binPath, watchDir, os.Stdout)
}

// RunWatchWithContext starts the interactive file watcher with a cancellable context and custom output writer.
func RunWatchWithContext(ctx context.Context, binPath string, watchDir string, out io.Writer) error {
	if watchDir == "" {
		watchDir = "exercises"
	}
	r := runner.NewRunner(binPath)
	store, _ := state.NewStore("")
	m := manifest.GetManifest()
	all := m.AllExercises()
	return RunWatchWithStore(ctx, r, all, store, watchDir, out)
}

// RunWatchWithRunner executes the watch loop with a provided runner against the default curriculum manifest.
func RunWatchWithRunner(ctx context.Context, r *runner.Runner, watchDir string, out io.Writer) error {
	m := manifest.GetManifest()
	all := m.AllExercises()
	store, _ := state.NewStore("")
	return RunWatchWithStore(ctx, r, all, store, watchDir, out)
}

// RunWatchWithExercises executes the watch loop for an explicit list of exercises and watch directory.
func RunWatchWithExercises(ctx context.Context, r *runner.Runner, exercises []models.Exercise, watchDir string, out io.Writer) error {
	return RunWatchWithStore(ctx, r, exercises, nil, watchDir, out)
}

// RunWatchWithStore executes the watch loop with progress tracking via state.Store.
func RunWatchWithStore(ctx context.Context, r *runner.Runner, exercises []models.Exercise, store *state.Store, watchDir string, out io.Writer) error {
	if watchDir == "" {
		watchDir = "exercises"
	}

	if len(exercises) == 0 {
		fmt.Fprintln(out, "No exercises found to watch.")
		return nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to initialize fsnotify watcher: %w", err)
	}
	defer w.Close()

	// Add watch directory and all nested subdirectories
	if info, err := os.Stat(watchDir); err == nil && info.IsDir() {
		_ = w.Add(watchDir)
		_ = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() {
				_ = w.Add(path)
			}
			return nil
		})
	}

	fmt.Fprintln(out, ui.FormatBanner())

	// Find the first incomplete exercise
	currentIdx := 0
	allPassed := true
	for i, ex := range exercises {
		res := r.Run(ex)
		if !res.Passed {
			currentIdx = i
			allPassed = false
			if store != nil {
				_ = store.RecordAttempt(ex.Name, ex.ChapterName, false)
			}
			fmt.Fprint(out, ui.FormatResult(res))
			break
		}
	}

	if allPassed {
		fmt.Fprintln(out, "\n🎉 Congratulations! You have completed all Terralings exercises! 🎉")
		fmt.Fprintln(out, ui.FormatProgress(len(exercises), len(exercises)))
		return nil
	}

	var (
		debounceTimer *time.Timer
		debounceMu    sync.Mutex
		triggerChan   = make(chan struct{}, 1)
	)

	triggerEvaluation := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
			select {
			case triggerChan <- struct{}{}:
			default:
			}
		})
	}

	defer func() {
		debounceMu.Lock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceMu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-w.Events:
			if !ok {
				return nil
			}

			// If a new subdirectory is created, watch it automatically
			if event.Op&fsnotify.Create != 0 {
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = w.Add(event.Name)
				}
			}

			if IsRelevantFile(event.Name) && (event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0) {
				triggerEvaluation()
			}

		case <-triggerChan:
			if currentIdx >= len(exercises) {
				fmt.Fprintln(out, "\n🎉 Congratulations! You have completed all Terralings exercises! 🎉")
				fmt.Fprintln(out, ui.FormatProgress(len(exercises), len(exercises)))
				return nil
			}

			currentEx := exercises[currentIdx]
			res := r.Run(currentEx)
			if store != nil {
				_ = store.RecordAttempt(currentEx.Name, currentEx.ChapterName, res.Passed)
			}
			fmt.Fprint(out, ui.FormatResult(res))

			if res.Passed {
				if currentIdx+1 < len(exercises) {
					currentIdx++
					nextEx := exercises[currentIdx]
					fmt.Fprintf(out, "\n%s\n", ui.FormatProgress(currentIdx, len(exercises)))
					fmt.Fprintf(out, "Advancing to next exercise: %s (%s)\n\n", nextEx.Name, nextEx.Path)
					nextRes := r.Run(nextEx)
					if store != nil {
						_ = store.RecordAttempt(nextEx.Name, nextEx.ChapterName, nextRes.Passed)
					}
					fmt.Fprint(out, ui.FormatResult(nextRes))
					if nextRes.Passed {
						triggerEvaluation()
					}
				} else {
					fmt.Fprintln(out, "\n🎉 Congratulations! You have completed all Terralings exercises! 🎉")
					fmt.Fprintln(out, ui.FormatProgress(len(exercises), len(exercises)))
					return nil
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(out, "Watcher error: %v\n", err)
		}
	}
}

// RunWatchJSON starts the watch loop emitting structured NDJSON events to out.
func RunWatchJSON(ctx context.Context, r *runner.Runner, exercises []models.Exercise, store *state.Store, watchDir string, out io.Writer) error {
	if watchDir == "" {
		watchDir = "exercises"
	}

	var writeMu sync.Mutex
	emit := func(evt WatchEvent) {
		writeMu.Lock()
		defer writeMu.Unlock()
		if evt.Timestamp == "" {
			evt.Timestamp = time.Now().UTC().Format(time.RFC3339)
		}
		data, err := json.Marshal(evt)
		if err == nil {
			_, _ = fmt.Fprintf(out, "%s\n", string(data))
		}
	}

	if len(exercises) == 0 {
		emit(WatchEvent{
			Event:      "completed",
			TotalCount: 0,
		})
		return nil
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to initialize fsnotify watcher: %w", err)
	}
	defer w.Close()

	if info, err := os.Stat(watchDir); err == nil && info.IsDir() {
		_ = w.Add(watchDir)
		_ = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() {
				_ = w.Add(path)
			}
			return nil
		})
	}

	evaluate := func(idx int) runner.RunResult {
		ex := exercises[idx]
		emit(WatchEvent{
			Event:        "exercise_start",
			Exercise:     &ex,
			CurrentIndex: idx,
			TotalCount:   len(exercises),
		})

		res := r.Run(ex)
		if store != nil {
			_ = store.RecordAttempt(ex.Name, ex.ChapterName, res.Passed)
		}

		rawOut := res.Output
		if res.Error != "" {
			rawOut = rawOut + "\n" + res.Error
		}
		diags := diagnostics.Parse(rawOut, ex)

		passed := res.Passed
		hasMarker := res.HasNotDoneMarker

		emit(WatchEvent{
			Event:            "exercise_result",
			Exercise:         &ex,
			Passed:           &passed,
			HasNotDoneMarker: &hasMarker,
			Diagnostics:      diags,
			RawOutput:        rawOut,
			ExitCode:         res.ExitCode,
			CurrentIndex:     idx,
			TotalCount:       len(exercises),
		})

		return res
	}

	// Initial evaluation
	currentIdx := 0
	allPassed := true
	for i := range exercises {
		res := evaluate(i)
		if !res.Passed {
			currentIdx = i
			allPassed = false
			break
		}
	}

	if allPassed {
		emit(WatchEvent{
			Event:      "completed",
			TotalCount: len(exercises),
		})
		return nil
	}

	var (
		debounceTimer *time.Timer
		debounceMu    sync.Mutex
		triggerChan   = make(chan struct{}, 1)
	)

	triggerEvaluation := func() {
		debounceMu.Lock()
		defer debounceMu.Unlock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceTimer = time.AfterFunc(50*time.Millisecond, func() {
			select {
			case triggerChan <- struct{}{}:
			default:
			}
		})
	}

	defer func() {
		debounceMu.Lock()
		if debounceTimer != nil {
			debounceTimer.Stop()
		}
		debounceMu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Create != 0 {
				if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
					_ = w.Add(event.Name)
				}
			}
			if IsRelevantFile(event.Name) && (event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0) {
				triggerEvaluation()
			}

		case <-triggerChan:
			if currentIdx >= len(exercises) {
				emit(WatchEvent{
					Event:      "completed",
					TotalCount: len(exercises),
				})
				return nil
			}

			res := evaluate(currentIdx)
			if res.Passed {
				for currentIdx+1 < len(exercises) {
					currentIdx++
					nextRes := evaluate(currentIdx)
					if !nextRes.Passed {
						break
					}
				}
				if currentIdx == len(exercises)-1 {
					lastEx := exercises[currentIdx]
					lastRes := r.Run(lastEx)
					if lastRes.Passed {
						emit(WatchEvent{
							Event:      "completed",
							TotalCount: len(exercises),
						})
						return nil
					}
				}
			}

		case _, ok := <-w.Errors:
			if !ok {
				return nil
			}
		}
	}
}

// IsRelevantFile returns true if the file extension or suffix matches Terraform/OpenTofu files.
func IsRelevantFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".tf" || ext == ".hcl" || strings.HasSuffix(path, ".tftest.hcl") || strings.HasSuffix(path, ".tfvars")
}
