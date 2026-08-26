# Terralings Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `terralings`, a high-performance interactive CLI learning tool and 13-chapter, 50+ exercise hands-on curriculum (with solutions and tests) for mastering Terraform and OpenTofu from scratch.

**Architecture:** A lightweight Go CLI engine built on Cobra and Lipgloss with an instantaneous file watcher (`fsnotify`), shared local provider caching for sub-100ms execution, declarative curriculum manifest, automated solutions validator, and complete repo infrastructure.

**Tech Stack:** Go 1.22+, OpenTofu >= 1.6 / Terraform >= 1.6, Cobra, Lipgloss, Fsnotify, HCL, .tftest.hcl.

---

### File Structure Map

```
terralings/
├── .github/workflows/ci.yml
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── CONTRIBUTING.md
├── CHANGELOG.md
├── LICENSE
├── cmd/
│   └── terralings/
│       └── main.go
├── internal/
│   ├── detector/
│   │   └── detector.go
│   ├── models/
│   │   └── models.go
│   ├── manifest/
│   │   └── manifest.go
│   ├── runner/
│   │   └── runner.go
│   ├── ui/
│   │   └── ui.go
│   └── watcher/
│       └── watcher.go
├── exercises/
│   ├── 01_primitives/ (primitives01.tf .. primitives06.tf)
│   ├── 02_variables/ (variables01.tf .. variables05.tf)
│   ├── 03_outputs_locals/ (outputs01.tf, locals01.tf, expr01.tf, expr02.tf)
│   ├── 04_functions/ (func01.tf .. func05.tf)
│   ├── 05_meta_arguments/ (meta01.tf .. meta05.tf)
│   ├── 06_dynamic_blocks/ (dynamic01.tf .. dynamic04.tf)
│   ├── 07_data_sources/ (data01.tf .. data04.tf)
│   ├── 08_modules/ (module01 .. module05)
│   ├── 09_state_refactoring/ (state01.tf .. state04.tf)
│   ├── 10_testing/ (test01.tftest.hcl .. test04.tftest.hcl)
│   ├── 11_patterns/ (pattern01.tf .. pattern04.tf)
│   ├── 12_opentofu/ (tofu01.tf .. tofu03.tf)
│   └── 13_governance/ (gov01.tf .. gov03.tf)
├── solutions/
│   └── (mirrors exercises/ 1:1 with working HCL)
└── test/
    ├── detector_test.go
    ├── manifest_test.go
    ├── runner_test.go
    ├── cli_test.go
    └── solutions_test.go
```

---

### Task 1: Project Setup, Go Module, Infrastructure & Gitignore

**Files:**
- Create: `go.mod`
- Create: `.gitignore`
- Create: `Makefile`
- Create: `LICENSE`
- Create: `README.md`
- Create: `CONTRIBUTING.md`
- Create: `CHANGELOG.md`
- Create: `.github/workflows/ci.yml`
- Test: `test/infra_test.go`

- [ ] **Step 1: Write infrastructure test**

```go
// test/infra_test.go
package test

import (
	"os"
	"testing"
)

func testFilesExist(t *testing.T) {
	requiredFiles := []string{
		"../go.mod",
		"../.gitignore",
		"../Makefile",
		"../LICENSE",
		"../README.md",
		"../CONTRIBUTING.md",
		"../CHANGELOG.md",
		"../.github/workflows/ci.yml",
	}
	for _, f := range requiredFiles {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			t.Fatalf("Required file does not exist: %s", f)
		}
	}
}

func TestInfra(t *testing.T) {
	testFilesExist(t)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/infra_test.go`  
Expected: FAIL (go.mod missing)

- [ ] **Step 3: Create go.mod, .gitignore, Makefile, CI workflow and metadata**

```
// go.mod
module github.com/dnf0/terralings

go 1.22.0

require (
	github.com/charmbracelet/lipgloss v0.10.0
	github.com/fsnotify/fsnotify v1.7.0
	github.com/spf13/cobra v1.8.0
)
```

```gitignore
# .gitignore
# Binaries
bin/
dist/
terralings

# Terraform & OpenTofu
.terraform/
.terraform.lock.hcl
*.tfstate
*.tfstate.*
*.tfplan
crash.log
crash.*.log
.cache/
.terralings/

# Agent rules & AI tooling
.agents/
.agent-state/
.superpowers/
.roborev/
.claude/
.gemini/
.cursor/
graphify-out/
.smellcheck-cache/
```

```makefile
# Makefile
.PHONY: build test run clean lint

build:
	go build -o bin/terralings cmd/terralings/main.go

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	rm -rf bin/ .terraform/ .cache/
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./test`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add go.mod .gitignore Makefile LICENSE README.md CONTRIBUTING.md CHANGELOG.md .github/ test/infra_test.go
git commit -m "chore: setup project infrastructure, packaging and CI"
```

---

### Task 2: Models & Manifest Engine

**Files:**
- Create: `internal/models/models.go`
- Create: `internal/manifest/manifest.go`
- Test: `test/manifest_test.go`

- [ ] **Step 1: Write test for models and manifest**

```go
// test/manifest_test.go
package test

import (
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
)

func TestManifestLoadsAllChapters(t *testing.T) {
	m := manifest.GetManifest()
	if len(m.Chapters) != 13 {
		t.Fatalf("Expected 13 chapters, got %d", len(m.Chapters))
	}
	allEx := m.AllExercises()
	if len(allEx) < 50 {
		t.Fatalf("Expected >= 50 exercises, got %d", len(allEx))
	}
	first := allEx[0]
	if first.Name != "primitives01" {
		t.Fatalf("Expected first exercise to be primitives01, got %s", first.Name)
	}
}

func TestGetExerciseByName(t *testing.T) {
	ex := manifest.GetExerciseByName("primitives01")
	if ex == nil {
		t.Fatal("primitives01 not found")
	}
	if ex.Path != "exercises/01_primitives/primitives01.tf" {
		t.Fatalf("Unexpected path: %s", ex.Path)
	}
}

func TestGetNextExercise(t *testing.T) {
	next := manifest.GetNextExercise("primitives01")
	if next == nil {
		t.Fatal("next exercise after primitives01 is nil")
	}
	if next.Name != "primitives02" {
		t.Fatalf("Expected primitives02, got %s", next.Name)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/manifest_test.go`  
Expected: FAIL (packages missing)

- [ ] **Step 3: Implement internal/models/models.go and internal/manifest/manifest.go**

```go
// internal/models/models.go
package models

import (
	"strings"
)

type ExerciseStatus string

const (
	StatusNotStarted ExerciseStatus = "not_started"
	StatusInProgress ExerciseStatus = "in_progress"
	StatusCompleted  ExerciseStatus = "completed"
	StatusFailed     ExerciseStatus = "failed"
)

type ExerciseMode string

const (
	ModeValidate ExerciseMode = "validate"
	ModePlan     ExerciseMode = "plan"
	ModeTest     ExerciseMode = "test"
)

type Exercise struct {
	Name        string
	Title       string
	Path        string
	ChapterName string
	Hints       []string
	Mode        ExerciseMode
}

func (e Exercise) SolutionPath() string {
	return strings.Replace(e.Path, "exercises/", "solutions/", 1)
}

type Chapter struct {
	Number      int
	Name        string
	Title       string
	Description string
	Exercises   []Exercise
}

type Manifest struct {
	Chapters []Chapter
}

func (m Manifest) AllExercises() []Exercise {
	var result []Exercise
	for _, ch := range m.Chapters {
		result = append(result, ch.Exercises...)
	}
	return result
}
```

```go
// internal/manifest/manifest.go
package manifest

import (
	"sync"

	"github.com/dnf0/terralings/internal/models"
)

var (
	manifestInstance *models.Manifest
	once             sync.Once
)

func buildManifest() *models.Manifest {
	return &models.Manifest{
		Chapters: []models.Chapter{
			{
				Number:      1,
				Name:        "01_primitives",
				Title:       "HCL Foundations & Core Primitives",
				Description: "Blocks, attributes, provider requirements and first resources",
				Exercises: []models.Exercise{
					{Name: "primitives01", Title: "Terraform Configuration Block", Path: "exercises/01_primitives/primitives01.tf", ChapterName: "01_primitives", Hints: []string{"Use required_version = \">= 1.6.0\"", "Specify required_providers with local source"}, Mode: models.ModeValidate},
					{Name: "primitives02", Title: "First Resource Declaration", Path: "exercises/01_primitives/primitives02.tf", ChapterName: "01_primitives", Hints: []string{"Declare resource local_file with filename and content"}, Mode: models.ModePlan},
					{Name: "primitives03", Title: "Resource Dependencies", Path: "exercises/01_primitives/primitives03.tf", ChapterName: "01_primitives", Hints: []string{"Reference local_file.first.content in second resource"}, Mode: models.ModePlan},
					{Name: "primitives04", Title: "String Interpolation & Heredoc", Path: "exercises/01_primitives/primitives04.tf", ChapterName: "01_primitives", Hints: []string{"Use <<-EOT for heredoc multi-line strings"}, Mode: models.ModePlan},
					{Name: "primitives05", Title: "Syntax & Formatting", Path: "exercises/01_primitives/primitives05.tf", ChapterName: "01_primitives", Hints: []string{"Align equals signs and use double quotes for strings"}, Mode: models.ModeValidate},
					{Name: "primitives06", Title: "Lifecycle Mechanics", Path: "exercises/01_primitives/primitives06.tf", ChapterName: "01_primitives", Hints: []string{"Ensure terraform_data resource has input and triggers_replace"}, Mode: models.ModePlan},
				},
			},
			// Chapters 2 through 13 fully populated...
		},
	}
}

func GetManifest() *models.Manifest {
	once.Do(func() {
		manifestInstance = buildManifest()
	})
	return manifestInstance
}

func GetExerciseByName(name string) *models.Exercise {
	for _, ex := range GetManifest().AllExercises() {
		if ex.Name == name || ex.Path == name {
			e := ex
			return &e
		}
	}
	return nil
}

func GetNextExercise(currentName string) *models.Exercise {
	all := GetManifest().AllExercises()
	for i, ex := range all {
		if ex.Name == currentName || ex.Path == currentName {
			if i+1 < len(all) {
				next := all[i+1]
				return &next
			}
		}
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./test/manifest_test.go`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/models/ internal/manifest/ test/manifest_test.go
git commit -m "feat: implement curriculum manifest and data models"
```

---

### Task 3: Binary Detector & Subprocess Execution Engine with Plugin Caching

**Files:**
- Create: `internal/detector/detector.go`
- Create: `internal/runner/runner.go`
- Test: `test/detector_test.go`
- Test: `test/runner_test.go`

- [ ] **Step 1: Write detector and runner tests**

```go
// test/detector_test.go
package test

import (
	"testing"

	"github.com/dnf0/terralings/internal/detector"
)

func TestDetectBinary(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping")
	}
	if bin != "tofu" && bin != "terraform" {
		t.Fatalf("Unexpected binary detected: %s", bin)
	}
}
```

```go
// test/runner_test.go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
)

func TestRunnerDetectsNotDoneMarker(t *testing.T) {
	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "ex01.tf")
	os.WriteFile(exFile, []byte("# I AM NOT DONE\nterraform {}"), 0644)

	r := runner.NewRunner("tofu")
	ex := models.Exercise{Name: "ex01", Path: exFile, Mode: models.ModeValidate}
	res := r.Run(ex)
	if res.Passed {
		t.Fatal("Expected exercise with NOT DONE marker to fail")
	}
	if !res.HasNotDoneMarker {
		t.Fatal("Expected HasNotDoneMarker to be true")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/detector_test.go ./test/runner_test.go`  
Expected: FAIL

- [ ] **Step 3: Implement internal/detector and internal/runner**

```go
// internal/detector/detector.go
package detector

import (
	"errors"
	"os"
	"os/exec"
)

func DetectBinary(override string) (string, error) {
	if override != "" {
		if path, err := exec.LookPath(override); err == nil {
			return path, nil
		}
		return "", errors.New("specified binary not found: " + override)
	}
	if env := os.Getenv("TERRALINGS_BIN"); env != "" {
		if path, err := exec.LookPath(env); err == nil {
			return path, nil
		}
	}
	if path, err := exec.LookPath("tofu"); err == nil {
		return path, nil
	}
	if path, err := exec.LookPath("terraform"); err == nil {
		return path, nil
	}
	return "", errors.New("neither 'tofu' nor 'terraform' was found on your PATH. Please install OpenTofu (https://opentofu.org) or Terraform")
}
```

```go
// internal/runner/runner.go
package runner

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dnf0/terralings/internal/models"
)

const NotDoneMarker = "I AM NOT DONE"

type RunResult struct {
	Exercise         models.Exercise
	Passed           bool
	HasNotDoneMarker bool
	Output           string
	Error            string
	ExitCode         int
}

type Runner struct {
	BinaryPath string
	CacheDir   string
}

func NewRunner(binaryPath string) *Runner {
	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".terralings", "plugin-cache")
	os.MkdirAll(cacheDir, 0755)
	return &Runner{
		BinaryPath: binaryPath,
		CacheDir:   cacheDir,
	}
}

func (r *Runner) CheckMarker(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), NotDoneMarker)
}

func (r *Runner) Run(ex models.Exercise) RunResult {
	hasMarker := r.CheckMarker(ex.Path)
	dir := filepath.Dir(ex.Path)
	
	// Step 1: Run init
	initCmd := exec.Command(r.BinaryPath, "init", "-backend=false", "-no-color")
	initCmd.Dir = dir
	initCmd.Env = append(os.Environ(), "TF_PLUGIN_CACHE_DIR="+r.CacheDir)
	initOut, err := initCmd.CombinedOutput()
	if err != nil && !hasMarker {
		return RunResult{
			Exercise:         ex,
			Passed:           false,
			HasNotDoneMarker: false,
			Output:           string(initOut),
			Error:            "Init failed: " + string(initOut),
			ExitCode:         1,
		}
	}

	// Step 2: Run verification command (validate / plan / test)
	var cmd *exec.Cmd
	switch ex.Mode {
	case models.ModeTest:
		cmd = exec.Command(r.BinaryPath, "test", "-no-color")
	case models.ModePlan:
		cmd = exec.Command(r.BinaryPath, "plan", "-no-color")
	default:
		cmd = exec.Command(r.BinaryPath, "validate", "-no-color")
	}

	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TF_PLUGIN_CACHE_DIR="+r.CacheDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmdErr := cmd.Run()

	passed := (cmdErr == nil) && !hasMarker
	return RunResult{
		Exercise:         ex,
		Passed:           passed,
		HasNotDoneMarker: hasMarker,
		Output:           stdout.String(),
		Error:            stderr.String(),
		ExitCode:         cmd.ProcessState.ExitCode(),
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./test/detector_test.go ./test/runner_test.go`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/detector/ internal/runner/ test/detector_test.go test/runner_test.go
git commit -m "feat: implement binary detector and runner engine with plugin caching"
```

---

### Task 4: Terminal UI, Diagnostics & Progressive Hints

**Files:**
- Create: `internal/ui/ui.go`
- Test: `test/ui_test.go`

- [ ] **Step 1: Write UI tests**

```go
// test/ui_test.go
package test

import (
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/ui"
)

func TestRenderResult(t *testing.T) {
	ex := models.Exercise{Name: "primitives01", Path: "exercises/01_primitives/primitives01.tf"}
	res := runner.RunResult{Exercise: ex, Passed: true}
	rendered := ui.FormatResult(res)
	if !strings.Contains(rendered, "primitives01 passed") {
		t.Fatalf("Expected output to contain 'primitives01 passed', got: %s", rendered)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./test/ui_test.go`  
Expected: FAIL

- [ ] **Step 3: Implement internal/ui/ui.go**

```go
// internal/ui/ui.go
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7D7")).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700"))

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF5F87")).
			Padding(1).
			MarginTop(1)

	hintBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5FAFFF")).
			Padding(1).
			MarginTop(1)
)

func FormatBanner() string {
	return headerStyle.Render("⚡ TERRALINGS: Master Terraform & OpenTofu from Scratch ⚡\n")
}

func FormatResult(res runner.RunResult) string {
	var b strings.Builder
	if res.Passed {
		b.WriteString(successStyle.Render(fmt.Sprintf("✓ Exercise %s passed!\n", res.Exercise.Name)))
	} else {
		if res.HasNotDoneMarker {
			b.WriteString(warningStyle.Render(fmt.Sprintf("⌛ %s still contains '%s' marker. Keep going!\n", res.Exercise.Name, runner.NotDoneMarker)))
		}
		if res.Error != "" {
			b.WriteString(errorBoxStyle.Render(fmt.Sprintf("Error in %s:\n%s", res.Exercise.Name, res.Error)))
		} else if res.Output != "" {
			b.WriteString(fmt.Sprintf("\n%s\n", res.Output))
		}
	}
	return b.String()
}

func FormatHint(ex *models.Exercise, hintIdx int) string {
	if ex == nil || len(ex.Hints) == 0 {
		return warningStyle.Render("No hints available for this exercise.")
	}
	idx := hintIdx
	if idx >= len(ex.Hints) {
		idx = len(ex.Hints) - 1
	}
	return hintBoxStyle.Render(fmt.Sprintf("💡 Hint for %s:\n%s", ex.Name, ex.Hints[idx]))
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./test/ui_test.go`  
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/ui/ test/ui_test.go
git commit -m "feat: implement rich terminal UI and formatting"
```

---

### Task 5: File Watcher Engine & CLI Commands

**Files:**
- Create: `internal/watcher/watcher.go`
- Create: `cmd/terralings/main.go`
- Test: `test/cli_test.go`

- [ ] **Step 1: Write CLI tests**

```go
// test/cli_test.go
package test

import (
	"bytes"
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/ui"
)

func TestListChapters(t *testing.T) {
	m := manifest.GetManifest()
	if len(m.Chapters) == 0 {
		t.Fatal("Chapters list should not be empty")
	}
}
```

- [ ] **Step 2: Run test to verify it passes / fails**

Run: `go test -v ./test/cli_test.go`  
Expected: PASS

- [ ] **Step 3: Implement internal/watcher and cmd/terralings/main.go**

```go
// internal/watcher/watcher.go
package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/ui"
)

func RunWatch(binPath string) error {
	r := runner.NewRunner(binPath)
	m := manifest.GetManifest()
	all := m.AllExercises()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	// Watch exercises dir
	w.Add("exercises")
	filepath.Walk("exercises", func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			w.Add(path)
		}
		return nil
	})

	fmt.Println(ui.FormatBanner())
	
	// Find first incomplete exercise
	currentIdx := 0
	for i, ex := range all {
		res := r.Run(ex)
		if !res.Passed {
			currentIdx = i
			fmt.Print(ui.FormatResult(res))
			break
		}
	}

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				time.Sleep(50 * time.Millisecond) // debounce
				currentEx := all[currentIdx]
				res := r.Run(currentEx)
				fmt.Print(ui.FormatResult(res))
				if res.Passed {
					if currentIdx+1 < len(all) {
						currentIdx++
						nextEx := all[currentIdx]
						fmt.Printf("\nAdvancing to next exercise: %s (%s)\n", nextEx.Name, nextEx.Path)
					} else {
						fmt.Println("\n🎉 Congratulations! You have completed all Terralings exercises! 🎉")
					}
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
		}
	}
}
```

```go
// cmd/terralings/main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/ui"
	"github.com/dnf0/terralings/internal/watcher"
)

var binOverride string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "terralings",
		Short: "Terralings - Learn Terraform & OpenTofu interactively",
	}

	rootCmd.PersistentFlags().StringVar(&binOverride, "bin", "", "Custom path to tofu or terraform binary")

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all curriculum chapters and exercises",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ui.FormatBanner())
			m := manifest.GetManifest()
			for _, ch := range m.Chapters {
				fmt.Printf("Chapter %02d: %s - %s\n", ch.Number, ch.Title, ch.Description)
				for _, ex := range ch.Exercises {
					fmt.Printf("  • %-16s : %s (%s)\n", ex.Name, ex.Title, ex.Path)
				}
			}
		},
	}

	var hintCmd = &cobra.Command{
		Use:   "hint [exercise_name]",
		Short: "Show progressive hints for an exercise",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ex := manifest.GetExerciseByName(args[0])
			if ex == nil {
				fmt.Fprintf(os.Stderr, "Exercise '%s' not found.\n", args[0])
				os.Exit(1)
			}
			fmt.Println(ui.FormatHint(ex, 0))
		},
	}

	var runCmd = &cobra.Command{
		Use:   "run [exercise_name]",
		Short: "Run verification on a single exercise",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			ex := manifest.GetExerciseByName(args[0])
			if ex == nil {
				fmt.Fprintf(os.Stderr, "Exercise '%s' not found.\n", args[0])
				os.Exit(1)
			}
			r := runner.NewRunner(bin)
			res := r.Run(*ex)
			fmt.Print(ui.FormatResult(res))
		},
	}

	var watchCmd = &cobra.Command{
		Use:   "watch",
		Short: "Start continuous interactive watch mode",
		Run: func(cmd *cobra.Command, args []string) {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			if err := watcher.RunWatch(bin); err != nil {
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(listCmd, hintCmd, runCmd, watchCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 4: Build binary and run test**

Run: `go build -o bin/terralings cmd/terralings/main.go && ./bin/terralings list`  
Expected: PASS (prints curriculum table)

- [ ] **Step 5: Commit**

```bash
git add internal/watcher/ cmd/terralings/ test/cli_test.go
git commit -m "feat: implement CLI commands and interactive watch loop"
```

---

### Task 6: Chapters 1 to 3 Curriculum & Reference Solutions

**Files:**
- Create: `exercises/01_primitives/` (primitives01.tf to primitives06.tf)
- Create: `solutions/01_primitives/` (primitives01.tf to primitives06.tf)
- Create: `exercises/02_variables/` (variables01.tf to variables05.tf)
- Create: `solutions/02_variables/` (variables01.tf to variables05.tf)
- Create: `exercises/03_outputs_locals/` (outputs01.tf, locals01.tf, expr01.tf, expr02.tf)
- Create: `solutions/03_outputs_locals/` (outputs01.tf, locals01.tf, expr01.tf, expr02.tf)
- Test: `test/chapters_1_3_test.go`

- [ ] **Step 1: Write test for Chapters 1-3 solutions**
- [ ] **Step 2: Author exercises and solutions for Chapters 1, 2, and 3**
- [ ] **Step 3: Run test to verify all solutions pass**
- [ ] **Step 4: Commit**

```bash
git add exercises/01_primitives exercises/02_variables exercises/03_outputs_locals solutions/ test/chapters_1_3_test.go
git commit -m "feat: add curriculum and reference solutions for chapters 1 to 3"
```

---

### Task 7: Chapters 4 to 7 Curriculum & Reference Solutions

**Files:**
- Create: `exercises/04_functions/` (func01.tf to func05.tf)
- Create: `solutions/04_functions/` (func01.tf to func05.tf)
- Create: `exercises/05_meta_arguments/` (meta01.tf to meta05.tf)
- Create: `solutions/05_meta_arguments/` (meta01.tf to meta05.tf)
- Create: `exercises/06_dynamic_blocks/` (dynamic01.tf to dynamic04.tf)
- Create: `solutions/06_dynamic_blocks/` (dynamic01.tf to dynamic04.tf)
- Create: `exercises/07_data_sources/` (data01.tf to data04.tf)
- Create: `solutions/07_data_sources/` (data01.tf to data04.tf)
- Test: `test/chapters_4_7_test.go`

- [ ] **Step 1: Write test for Chapters 4-7 solutions**
- [ ] **Step 2: Author exercises and solutions for Chapters 4, 5, 6, and 7**
- [ ] **Step 3: Run test to verify all solutions pass**
- [ ] **Step 4: Commit**

```bash
git add exercises/04_functions exercises/05_meta_arguments exercises/06_dynamic_blocks exercises/07_data_sources solutions/ test/chapters_4_7_test.go
git commit -m "feat: add curriculum and reference solutions for chapters 4 to 7"
```

---

### Task 8: Chapters 8 to 10 Curriculum & Reference Solutions

**Files:**
- Create: `exercises/08_modules/` (module01 to module05)
- Create: `solutions/08_modules/` (module01 to module05)
- Create: `exercises/09_state_refactoring/` (state01.tf to state04.tf)
- Create: `solutions/09_state_refactoring/` (state01.tf to state04.tf)
- Create: `exercises/10_testing/` (test01.tftest.hcl to test04.tftest.hcl)
- Create: `solutions/10_testing/` (test01.tftest.hcl to test04.tftest.hcl)
- Test: `test/chapters_8_10_test.go`

- [ ] **Step 1: Write test for Chapters 8-10 solutions**
- [ ] **Step 2: Author exercises and solutions for Chapters 8, 9, and 10**
- [ ] **Step 3: Run test to verify all solutions pass**
- [ ] **Step 4: Commit**

```bash
git add exercises/08_modules exercises/09_state_refactoring exercises/10_testing solutions/ test/chapters_8_10_test.go
git commit -m "feat: add curriculum and reference solutions for chapters 8 to 10"
```

---

### Task 9: Chapters 11 to 13 Curriculum & Reference Solutions

**Files:**
- Create: `exercises/11_patterns/` (pattern01.tf to pattern04.tf)
- Create: `solutions/11_patterns/` (pattern01.tf to pattern04.tf)
- Create: `exercises/12_opentofu/` (tofu01.tf to tofu03.tf)
- Create: `solutions/12_opentofu/` (tofu01.tf to tofu03.tf)
- Create: `exercises/13_governance/` (gov01.tf to gov03.tf)
- Create: `solutions/13_governance/` (gov01.tf to gov03.tf)
- Test: `test/chapters_11_13_test.go`

- [ ] **Step 1: Write test for Chapters 11-13 solutions**
- [ ] **Step 2: Author exercises and solutions for Chapters 11, 12, and 13**
- [ ] **Step 3: Run test to verify all solutions pass**
- [ ] **Step 4: Commit**

```bash
git add exercises/11_patterns exercises/12_opentofu exercises/13_governance solutions/ test/chapters_11_13_test.go
git commit -m "feat: add curriculum and reference solutions for chapters 11 to 13"
```

---

### Task 10: End-to-End Test Suite, Comprehensive Validation & CI Gates

**Files:**
- Create: `test/solutions_test.go`
- Update: `README.md`
- Update: `Makefile`

- [ ] **Step 1: Implement global solution verification test in `test/solutions_test.go`**
- [ ] **Step 2: Run full test suite (`go test -v ./test`)**
- [ ] **Step 3: Run `tofu fmt -check -recursive solutions/`**
- [ ] **Step 4: Commit**

```bash
git add test/ README.md Makefile
git commit -m "test: implement comprehensive test suite and CI verification harness"
```
