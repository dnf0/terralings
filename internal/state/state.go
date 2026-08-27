package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ExerciseStatus represents the completion status of an individual exercise.
type ExerciseStatus string

const (
	StatusNotStarted ExerciseStatus = "not_started"
	StatusInProgress ExerciseStatus = "in_progress"
	StatusPassed     ExerciseStatus = "passed"
)

// ExerciseState holds progress metrics and timestamps for a single exercise.
type ExerciseState struct {
	Name             string         `json:"name"`
	Chapter          string         `json:"chapter"`
	Status           ExerciseStatus `json:"status"`
	Attempts         int            `json:"attempts"`
	HintsViewed      int            `json:"hints_viewed"`
	FirstAttemptAt   *time.Time     `json:"first_attempt_at,omitempty"`
	CompletedAt      *time.Time     `json:"completed_at,omitempty"`
	TimeSpentSeconds int64          `json:"time_spent_seconds"`
}

// StateData is the root JSON structure persisted to state.json.
type StateData struct {
	Version               string                    `json:"version"`
	CreatedAt             time.Time                 `json:"created_at"`
	LastActiveAt          time.Time                 `json:"last_active_at"`
	TotalTimeSpentSeconds int64                     `json:"total_time_spent_seconds"`
	Exercises             map[string]*ExerciseState `json:"exercises"`
}

// Store provides thread-safe access to persistent state data.
type Store struct {
	mu       sync.RWMutex
	filePath string
	data     StateData
}

// NewStore initializes a Store at the given file path.
// It handles directory creation, corrupt file backup/recovery, and auto-gitignore.
func NewStore(filePath string) (*Store, error) {
	if filePath == "" {
		filePath = filepath.Join(".terralings", "state.json")
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	ensureGitignore(dir)

	store := &Store{
		filePath: filePath,
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		now := time.Now().UTC()
		store.data = StateData{
			Version:      "1.0",
			CreatedAt:    now,
			LastActiveAt: now,
			Exercises:    make(map[string]*ExerciseState),
		}
		if err := store.saveLocked(); err != nil {
			return nil, err
		}
		return store, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data StateData
	if err := json.Unmarshal(content, &data); err != nil {
		// Backup corrupt file and initialize a clean store
		backupPath := filePath + ".bak"
		_ = os.WriteFile(backupPath, content, 0644)

		now := time.Now().UTC()
		store.data = StateData{
			Version:      "1.0",
			CreatedAt:    now,
			LastActiveAt: now,
			Exercises:    make(map[string]*ExerciseState),
		}
		if err := store.saveLocked(); err != nil {
			return nil, err
		}
		return store, nil
	}

	if data.Exercises == nil {
		data.Exercises = make(map[string]*ExerciseState)
	}
	if data.Version == "" {
		data.Version = "1.0"
	}
	store.data = data

	return store, nil
}

func ensureGitignore(stateDir string) {
	if filepath.Base(stateDir) != ".terralings" {
		return
	}

	workspaceDir := filepath.Dir(stateDir)
	gitDir := filepath.Join(workspaceDir, ".git")
	gitignorePath := filepath.Join(workspaceDir, ".gitignore")

	gitExists := false
	if fi, err := os.Stat(gitDir); err == nil && fi.IsDir() {
		gitExists = true
	}

	gitignoreExists := false
	if _, err := os.Stat(gitignorePath); err == nil {
		gitignoreExists = true
	}

	if !gitExists && !gitignoreExists {
		return
	}

	if gitignoreExists {
		data, err := os.ReadFile(gitignorePath)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == ".terralings" || trimmed == ".terralings/" || trimmed == "/.terralings" || trimmed == "/.terralings/" {
					return
				}
			}
			// Append entry
			content := string(data)
			if len(content) > 0 && !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += ".terralings/\n"
			_ = os.WriteFile(gitignorePath, []byte(content), 0644)
		}
	} else if gitExists {
		_ = os.WriteFile(gitignorePath, []byte(".terralings/\n"), 0644)
	}
}

// GetVersion returns the data schema version.
func (s *Store) GetVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.Version
}

// GetFilePath returns the file path on disk.
func (s *Store) GetFilePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filePath
}

// GetExerciseState returns a copy of the state for the named exercise, or nil if not tracked.
func (s *Store) GetExerciseState(name string) *ExerciseState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ex, ok := s.data.Exercises[name]
	if !ok || ex == nil {
		return nil
	}
	cp := *ex
	return &cp
}

// GetAllExerciseStates returns a copy of all exercise states.
func (s *Store) GetAllExerciseStates() map[string]ExerciseState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]ExerciseState, len(s.data.Exercises))
	for k, v := range s.data.Exercises {
		if v != nil {
			result[k] = *v
		}
	}
	return result
}

// RecordAttempt updates attempt count, status, timestamps, and persists changes.
func (s *Store) RecordAttempt(name string, chapter string, passed bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ex, ok := s.data.Exercises[name]
	if !ok || ex == nil {
		ex = &ExerciseState{
			Name:    name,
			Chapter: chapter,
			Status:  StatusInProgress,
		}
		s.data.Exercises[name] = ex
	}

	now := time.Now().UTC()
	if ex.FirstAttemptAt == nil {
		ex.FirstAttemptAt = &now
	}
	ex.Attempts++
	s.data.LastActiveAt = now

	if passed {
		ex.Status = StatusPassed
		if ex.CompletedAt == nil {
			ex.CompletedAt = &now
		}
	} else if ex.Status != StatusPassed {
		ex.Status = StatusInProgress
	}

	return s.saveLocked()
}

// RecordHint updates hint viewed count and persists changes.
func (s *Store) RecordHint(name string, chapter string, hintIndex int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ex, ok := s.data.Exercises[name]
	if !ok || ex == nil {
		ex = &ExerciseState{
			Name:    name,
			Chapter: chapter,
			Status:  StatusInProgress,
		}
		s.data.Exercises[name] = ex
	}

	if hintIndex > ex.HintsViewed {
		ex.HintsViewed = hintIndex
	} else if ex.HintsViewed == 0 && hintIndex >= 0 {
		if hintIndex == 0 {
			ex.HintsViewed = 1
		} else {
			ex.HintsViewed = hintIndex
		}
	}

	s.data.LastActiveAt = time.Now().UTC()
	return s.saveLocked()
}

// AddTimeSpent adds active duration to an exercise and total time spent.
func (s *Store) AddTimeSpent(name string, chapter string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	secs := int64(duration.Seconds())
	if secs <= 0 {
		return nil
	}

	s.data.TotalTimeSpentSeconds += secs
	s.data.LastActiveAt = time.Now().UTC()

	ex, ok := s.data.Exercises[name]
	if !ok || ex == nil {
		ex = &ExerciseState{
			Name:    name,
			Chapter: chapter,
			Status:  StatusInProgress,
		}
		s.data.Exercises[name] = ex
	}
	ex.TimeSpentSeconds += secs

	return s.saveLocked()
}

// Save flushes current state to disk atomically.
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	tmpPath := s.filePath + ".tmp"
	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s.data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmpPath)
		return err
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, s.filePath)
}
