package models

import (
	"strings"
)

// ExerciseStatus represents the completion state of an exercise.
type ExerciseStatus string

const (
	StatusNotStarted ExerciseStatus = "not_started"
	StatusInProgress ExerciseStatus = "in_progress"
	StatusCompleted  ExerciseStatus = "completed"
	StatusFailed     ExerciseStatus = "failed"
)

// ExerciseMode represents the execution/evaluation mode for an exercise.
type ExerciseMode string

const (
	ModeValidate ExerciseMode = "validate"
	ModePlan     ExerciseMode = "plan"
	ModeTest     ExerciseMode = "test"
)

// Exercise represents a single hands-on curriculum exercise.
type Exercise struct {
	Name        string
	Title       string
	Path        string
	ChapterName string
	Hints       []string
	Mode        ExerciseMode
}

// SolutionPath returns the corresponding solution file path in solutions/.
func (e Exercise) SolutionPath() string {
	return strings.Replace(e.Path, "exercises/", "solutions/", 1)
}

// Chapter represents a curriculum chapter containing a sequence of exercises.
type Chapter struct {
	Number      int
	Name        string
	Title       string
	Description string
	Exercises   []Exercise
}

// Manifest represents the complete curriculum catalogue.
type Manifest struct {
	Chapters []Chapter
}

// AllExercises returns a flattened slice of all exercises across all chapters.
func (m Manifest) AllExercises() []Exercise {
	var result []Exercise
	for _, ch := range m.Chapters {
		result = append(result, ch.Exercises...)
	}
	return result
}
