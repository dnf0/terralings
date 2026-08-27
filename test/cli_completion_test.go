package test

import (
	"strings"
	"testing"
)

func TestCLI_Completion_Bash(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "completion", "bash")
	if exitCode != 0 {
		t.Fatalf("completion bash failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "bash completion V2") || !strings.Contains(stdout, "__start_terralings") {
		t.Errorf("expected bash completion V2 script with __start_terralings, got:\n%s", stdout)
	}
}

func TestCLI_Completion_Zsh(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "completion", "zsh")
	if exitCode != 0 {
		t.Fatalf("completion zsh failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "compdef") && !strings.Contains(stdout, "zsh completion") {
		t.Errorf("expected zsh completion script, got:\n%s", stdout)
	}
}

func TestCLI_Completion_Fish(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "completion", "fish")
	if exitCode != 0 {
		t.Fatalf("completion fish failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "fish completion") && !strings.Contains(stdout, "complete -c terralings") {
		t.Errorf("expected fish completion script, got:\n%s", stdout)
	}
}

func TestCLI_Completion_PowerShell(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "completion", "powershell")
	if exitCode != 0 {
		t.Fatalf("completion powershell failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "Register-ArgumentCompleter") && !strings.Contains(stdout, "powershell completion") {
		t.Errorf("expected powershell completion script, got:\n%s", stdout)
	}
}

func TestCLI_Completion_AliasCompletions(t *testing.T) {
	canonicalOut, canonicalErr, canonicalExit := runCLI(t, "completion", "bash")
	if canonicalExit != 0 {
		t.Fatalf("canonical completion failed: %s", canonicalErr)
	}
	aliasOut, aliasErr, aliasExit := runCLI(t, "completions", "bash")
	if aliasExit != 0 {
		t.Fatalf("completions alias failed with exit code %d, stderr: %s", aliasExit, aliasErr)
	}
	if aliasOut != canonicalOut {
		t.Errorf("expected 'completions bash' output to match 'completion bash', got mismatch")
	}
}

func TestCLI_Completion_UnsupportedShell(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "completion", "invalid_shell_name")
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for unsupported shell, got 0. Stdout: %s", stdout)
	}
	if !strings.Contains(stderr, "unsupported shell") {
		t.Errorf("expected 'unsupported shell' error message in stderr, got:\n%s", stderr)
	}
}

func TestCLI_DynamicExerciseCompletion(t *testing.T) {
	commands := []string{"run", "hint", "reset", "search"}
	for _, subcmd := range commands {
		t.Run(subcmd, func(t *testing.T) {
			stdout, stderr, exitCode := runCLI(t, "__complete", subcmd, "prim")
			if exitCode != 0 {
				t.Fatalf("__complete %s prim failed with exit code %d, stderr: %s", subcmd, exitCode, stderr)
			}
			if !strings.Contains(stdout, "primitives01") {
				t.Errorf("expected dynamic completion for %s to return 'primitives01', got:\n%s", subcmd, stdout)
			}
			// Verify tab description is included (e.g. primitives01\t...)
			if !strings.Contains(stdout, "primitives01\t") {
				t.Errorf("expected dynamic completion to include tab description annotation, got:\n%s", stdout)
			}
		})
	}
}

func TestCLI_DynamicSearchCompletion_Chapters(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "__complete", "search", "01_")
	if exitCode != 0 {
		t.Fatalf("__complete search 01_ failed with exit code %d, stderr: %s", exitCode, stderr)
	}
	if !strings.Contains(stdout, "01_primitives") {
		t.Errorf("expected dynamic completion for search to return chapter '01_primitives', got:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Chapter:") {
		t.Errorf("expected dynamic completion description to mention 'Chapter:', got:\n%s", stdout)
	}
}
