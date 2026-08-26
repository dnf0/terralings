package runner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/dnf0/terralings/internal/models"
)

// NotDoneMarker is the standard comment marker indicating an incomplete exercise.
const NotDoneMarker = "I AM NOT DONE"

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

// NewRunner creates a new Runner with plugin caching configured in ~/.terralings/plugin-cache.
func NewRunner(binaryPath string) *Runner {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	cacheDir := filepath.Join(home, ".terralings", "plugin-cache")
	_ = os.MkdirAll(cacheDir, 0755)

	return &Runner{
		BinaryPath: binaryPath,
		CacheDir:   cacheDir,
	}
}

// CheckMarker inspects a file or path for the 'I AM NOT DONE' marker (case-insensitive with whitespace variations).
// Returns false if the marker is absent or if the file cannot be read.
func CheckMarker(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return markerRegex.Match(data)
}

// Run executes the exercise via init and the appropriate mode (validate, plan, test).
func (r *Runner) Run(ex models.Exercise) RunResult {
	hasMarker := CheckMarker(ex.Path)

	dir := filepath.Dir(ex.Path)
	if info, err := os.Stat(ex.Path); err == nil && info.IsDir() {
		dir = ex.Path
	}

	// Step 1: Run init
	initCmd := exec.Command(r.BinaryPath, "init", "-backend=false", "-no-color", "-input=false")
	initCmd.Dir = dir
	initCmd.Env = append(os.Environ(), "TF_PLUGIN_CACHE_DIR="+r.CacheDir)
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
		cmd = exec.Command(r.BinaryPath, "test", "-no-color")
	case models.ModePlan:
		cmd = exec.Command(r.BinaryPath, "plan", "-no-color", "-input=false")
	default:
		cmd = exec.Command(r.BinaryPath, "validate", "-no-color")
	}

	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TF_PLUGIN_CACHE_DIR="+r.CacheDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmdErr := cmd.Run()

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	outStr := stdout.String()
	errStr := stderr.String()
	if cmdErr != nil && errStr == "" && outStr != "" {
		errStr = outStr
	}

	passed := (cmdErr == nil) && !hasMarker

	return RunResult{
		Exercise:         ex,
		Passed:           passed,
		HasNotDoneMarker: hasMarker,
		Output:           outStr,
		Error:            errStr,
		ExitCode:         exitCode,
	}
}
