package test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Completion_Bash(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "bash")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion bash failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "bash completion") && !strings.Contains(out, "__terralings_") && !strings.Contains(out, "complete -o default -F") {
		t.Errorf("expected bash completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_Zsh(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "zsh")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion zsh failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "compdef") && !strings.Contains(out, "zsh completion") {
		t.Errorf("expected zsh completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_Fish(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "fish")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion fish failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "fish completion") && !strings.Contains(out, "complete -c terralings") {
		t.Errorf("expected fish completion script, got:\n%s", out)
	}
}

func TestCLI_Completion_PowerShell(t *testing.T) {
	cmd := exec.Command("go", "run", "../cmd/terralings", "completion", "powershell")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("completion powershell failed: %v, stderr: %s", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "Register-ArgumentCompleter") && !strings.Contains(out, "powershell completion") {
		t.Errorf("expected powershell completion script, got:\n%s", out)
	}
}

func TestCLI_DynamicExerciseCompletion(t *testing.T) {
	commands := []string{"run", "hint", "reset", "search"}
	for _, subcmd := range commands {
		t.Run(subcmd, func(t *testing.T) {
			cmd := exec.Command("go", "run", "../cmd/terralings", "__complete", subcmd, "prim")
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err != nil {
				t.Fatalf("__complete %s prim failed: %v, stderr: %s", subcmd, err, stderr.String())
			}
			out := stdout.String()
			if !strings.Contains(out, "primitives01") {
				t.Errorf("expected dynamic completion for %s to return 'primitives01', got:\n%s", subcmd, out)
			}
		})
	}
}
