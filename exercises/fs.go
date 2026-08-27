package exercises

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dnf0/terralings/internal/manifest"
)

//go:embed all:*
var EmbeddedFS embed.FS

// ExtractAll extracts all embedded exercises to targetDir.
// If targetDir contains files (excluding .git), it returns an error unless force is true.
func ExtractAll(targetDir string, force bool) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	if !force {
		entries, err := os.ReadDir(targetDir)
		if err != nil {
			return fmt.Errorf("failed to read target directory: %w", err)
		}
		for _, e := range entries {
			if e.Name() != ".git" {
				return fmt.Errorf("target directory '%s' is not empty (use --force to overwrite)", targetDir)
			}
		}
	}

	return fs.WalkDir(EmbeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." || strings.HasSuffix(path, ".go") {
			return nil
		}

		targetPath := filepath.Join(targetDir, path)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := EmbeddedFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file '%s': %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}

// GetExerciseContent returns the embedded content of an exercise given its relative path.
func GetExerciseContent(relPath string) ([]byte, error) {
	relPath = strings.TrimPrefix(relPath, "exercises/")
	relPath = strings.TrimPrefix(relPath, "/")
	return EmbeddedFS.ReadFile(relPath)
}

// ResetExercise resets an exercise in baseDir back to its embedded starter code.
func ResetExercise(exerciseName string, baseDir string) error {
	ex := manifest.GetExerciseByName(exerciseName)
	if ex == nil {
		return fmt.Errorf("exercise '%s' not found", exerciseName)
	}

	relPath := strings.TrimPrefix(ex.Path, "exercises/")
	relPath = strings.TrimPrefix(relPath, "/")

	if !ex.IsDirectory() {
		data, err := EmbeddedFS.ReadFile(relPath)
		if err != nil {
			return fmt.Errorf("failed to load original template for %s: %w", exerciseName, err)
		}
		targetFile := filepath.Join(baseDir, relPath)
		if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
			return err
		}
		return os.WriteFile(targetFile, data, 0644)
	}

	return fs.WalkDir(EmbeddedFS, relPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." || strings.HasSuffix(path, ".go") {
			return nil
		}

		targetPath := filepath.Join(baseDir, path)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := EmbeddedFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0644)
	})
}
