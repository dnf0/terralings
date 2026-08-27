package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dnf0/terralings/internal/models"
)

// NotDoneMarker is the standard comment marker indicating an incomplete exercise.
const NotDoneMarker = "I AM NOT DONE"

// DefaultTimeout is the default execution timeout per exercise command.
const DefaultTimeout = 30 * time.Second

var markerRegex = regexp.MustCompile(`(?i)i\s+am\s+not\s+done`)

// RunResult contains the execution and evaluation outcome for an exercise.
type RunResult struct {
	Exercise         models.Exercise
	Passed           bool
	HasNotDoneMarker bool
	Output           string
	Error            string
	ExitCode         int
}

// Runner manages subprocess execution of OpenTofu / Terraform commands with plugin caching.
type Runner struct {
	BinaryPath string
	CacheDir   string
}

// PluginCacheDir returns the plugin cache directory path.
// Checks TERRALINGS_PLUGIN_CACHE_DIR environment variable first, falling back to ~/.terralings/plugin-cache.
func PluginCacheDir() string {
	if custom := os.Getenv("TERRALINGS_PLUGIN_CACHE_DIR"); custom != "" {
		return custom
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".terralings", "plugin-cache")
}

// NewRunner creates a new Runner with plugin caching configured in ~/.terralings/plugin-cache.
func NewRunner(binaryPath string) *Runner {
	cacheDir := PluginCacheDir()
	_ = os.MkdirAll(cacheDir, 0755)

	return &Runner{
		BinaryPath: binaryPath,
		CacheDir:   cacheDir,
	}
}

// CheckMarker inspects a file or directory for the 'I AM NOT DONE' marker (case-insensitive with whitespace variations).
// Returns false if the marker is absent or if the path cannot be read.
func CheckMarker(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		hasMarker := false
		_ = filepath.Walk(path, func(p string, i os.FileInfo, walkErr error) error {
			if walkErr == nil && !i.IsDir() && (strings.HasSuffix(p, ".tf") || strings.HasSuffix(p, ".hcl") || strings.HasSuffix(p, ".tftest.hcl")) {
				if data, readErr := os.ReadFile(p); readErr == nil {
					if markerRegex.Match(data) {
						hasMarker = true
						return filepath.SkipAll
					}
				}
			}
			return nil
		})
		return hasMarker
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return markerRegex.Match(data)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath)
	})
}

// Run executes the exercise via init and the appropriate mode (validate, plan, test) in an isolated directory.
func (r *Runner) Run(ex models.Exercise) RunResult {
	hasMarker := CheckMarker(ex.Path)

	stat, err := os.Stat(ex.Path)
	if err != nil {
		return RunResult{
			Exercise:         ex,
			Passed:           false,
			HasNotDoneMarker: hasMarker,
			Output:           "",
			Error:            fmt.Sprintf("Failed to access exercise path: %v", err),
			ExitCode:         1,
		}
	}

	tempDir, err := os.MkdirTemp("", "terralings-stage-*")
	if err != nil {
		return RunResult{
			Exercise:         ex,
			Passed:           false,
			HasNotDoneMarker: hasMarker,
			Output:           "",
			Error:            fmt.Sprintf("Failed to create temporary staging directory: %v", err),
			ExitCode:         1,
		}
	}
	defer os.RemoveAll(tempDir)

	workDir := tempDir
	if stat.IsDir() {
		if copyErr := copyDir(ex.Path, tempDir); copyErr != nil {
			return RunResult{
				Exercise:         ex,
				Passed:           false,
				HasNotDoneMarker: hasMarker,
				Output:           "",
				Error:            fmt.Sprintf("Failed to copy exercise directory: %v", copyErr),
				ExitCode:         1,
			}
		}
	} else {
		targetFile := filepath.Join(tempDir, filepath.Base(ex.Path))
		if copyErr := copyFile(ex.Path, targetFile); copyErr != nil {
			return RunResult{
				Exercise:         ex,
				Passed:           false,
				HasNotDoneMarker: hasMarker,
				Output:           "",
				Error:            fmt.Sprintf("Failed to stage exercise file: %v", copyErr),
				ExitCode:         1,
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	cmdEnv := append(os.Environ(),
		"TF_PLUGIN_CACHE_DIR="+r.CacheDir,
		"TF_IN_AUTOMATION=1",
		"TF_INPUT=0",
	)

	// Step 1: Run init
	initCmd := exec.CommandContext(ctx, r.BinaryPath, "init", "-backend=false", "-no-color", "-input=false")
	initCmd.Dir = workDir
	initCmd.Env = cmdEnv
	initOut, initErr := initCmd.CombinedOutput()

	if initErr != nil {
		exitCode := 1
		if initCmd.ProcessState != nil {
			exitCode = initCmd.ProcessState.ExitCode()
		}
		return RunResult{
			Exercise:         ex,
			Passed:           false,
			HasNotDoneMarker: hasMarker,
			Output:           string(initOut),
			Error:            fmt.Sprintf("Init failed: %s", string(initOut)),
			ExitCode:         exitCode,
		}
	}

	// Step 2: Run verification command (validate / plan / test)
	var cmd *exec.Cmd
	switch ex.Mode {
	case models.ModeTest:
		cmd = exec.CommandContext(ctx, r.BinaryPath, "test", "-no-color")
	case models.ModePlan:
		cmd = exec.CommandContext(ctx, r.BinaryPath, "plan", "-no-color", "-input=false")
	default:
		cmd = exec.CommandContext(ctx, r.BinaryPath, "validate", "-no-color")
	}

	cmd.Dir = workDir
	cmd.Env = cmdEnv
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmdErr := cmd.Run()

	exitCode := 0
	if cmdErr != nil {
		exitCode = 1
		if cmd.ProcessState != nil {
			exitCode = cmd.ProcessState.ExitCode()
		}
	} else if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	outStr := stdout.String()
	errStr := stderr.String()
	if cmdErr != nil && errStr == "" && outStr != "" {
		errStr = outStr
	}

	passed := (cmdErr == nil)

	return RunResult{
		Exercise:         ex,
		Passed:           passed,
		HasNotDoneMarker: hasMarker,
		Output:           outStr,
		Error:            errStr,
		ExitCode:         exitCode,
	}
}
