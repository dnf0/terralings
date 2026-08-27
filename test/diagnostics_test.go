package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/diagnostics"
	"github.com/dnf0/terralings/internal/models"
)

func TestParseDiagnostics_MarkerPresent(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "exercise01.tf")
	content := `# First line comment
# I AM NOT DONE
terraform {
  required_version = ">= 1.6.0"
}
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test exercise file: %v", err)
	}

	ex := models.Exercise{
		Name: "exercise01",
		Path: filePath,
		Mode: models.ModeValidate,
	}

	diags := diagnostics.ParseDiagnostics("", ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic for marker, got %d", len(diags))
	}

	d := diags[0]
	if d.Severity != diagnostics.SeverityWarning {
		t.Errorf("Expected severity %q, got %q", diagnostics.SeverityWarning, d.Severity)
	}
	if d.File != filePath {
		t.Errorf("Expected file %q, got %q", filePath, d.File)
	}
	if d.Line != 2 {
		t.Errorf("Expected line 2, got %d", d.Line)
	}
	if !strings.Contains(d.Summary, "I AM NOT DONE") && !strings.Contains(d.Summary, "not finished") {
		t.Errorf("Expected summary to reference I AM NOT DONE, got %q", d.Summary)
	}
}

func TestParseDiagnostics_StandardCompilerError(t *testing.T) {
	rawOutput := `Error: Missing required argument

  on main.tf line 14, in resource "local_file" "example":
  14: content = "hello"

The argument "filename" is required, but no definition was found.`

	ex := models.Exercise{
		Name: "ex01",
		Path: "exercises/01_primitives/ex01.tf",
		Mode: models.ModeValidate,
	}

	diags := diagnostics.ParseDiagnostics(rawOutput, ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.Severity != diagnostics.SeverityError {
		t.Errorf("Expected severity %q, got %q", diagnostics.SeverityError, d.Severity)
	}
	if d.File != "main.tf" {
		t.Errorf("Expected file %q, got %q", "main.tf", d.File)
	}
	if d.Line != 14 {
		t.Errorf("Expected line 14, got %d", d.Line)
	}
	if d.Summary != "Missing required argument" {
		t.Errorf("Expected summary %q, got %q", "Missing required argument", d.Summary)
	}
	if !strings.Contains(d.Detail, "The argument \"filename\" is required") {
		t.Errorf("Expected detail to contain explanation, got %q", d.Detail)
	}
}

func TestParseDiagnostics_MultipleErrors(t *testing.T) {
	rawOutput := `Error: Missing required argument

  on main.tf line 14, in resource "local_file" "example":
  14: content = "hello"

The argument "filename" is required.

Error: Unsupported attribute

  on outputs.tf line 2, in output "result":
   2:   value = local_file.example.id_wrong

This object has no argument, attribute, or block named "id_wrong".`

	ex := models.Exercise{
		Name: "ex02",
		Path: "exercises/01_primitives/ex02.tf",
		Mode: models.ModeValidate,
	}

	diags := diagnostics.ParseDiagnostics(rawOutput, ex)
	if len(diags) != 2 {
		t.Fatalf("Expected 2 diagnostics, got %d", len(diags))
	}

	// First diagnostic
	if diags[0].Severity != diagnostics.SeverityError {
		t.Errorf("Expected diag[0] severity error, got %q", diags[0].Severity)
	}
	if diags[0].File != "main.tf" {
		t.Errorf("Expected diag[0] file main.tf, got %q", diags[0].File)
	}
	if diags[0].Line != 14 {
		t.Errorf("Expected diag[0] line 14, got %d", diags[0].Line)
	}
	if diags[0].Summary != "Missing required argument" {
		t.Errorf("Expected diag[0] summary 'Missing required argument', got %q", diags[0].Summary)
	}

	// Second diagnostic
	if diags[1].Severity != diagnostics.SeverityError {
		t.Errorf("Expected diag[1] severity error, got %q", diags[1].Severity)
	}
	if diags[1].File != "outputs.tf" {
		t.Errorf("Expected diag[1] file outputs.tf, got %q", diags[1].File)
	}
	if diags[1].Line != 2 {
		t.Errorf("Expected diag[1] line 2, got %d", diags[1].Line)
	}
	if diags[1].Summary != "Unsupported attribute" {
		t.Errorf("Expected diag[1] summary 'Unsupported attribute', got %q", diags[1].Summary)
	}
}

func TestParseDiagnostics_JSONDiagnostic(t *testing.T) {
	rawJSON := `{"@level":"error","diagnostic":{"severity":"error","summary":"Missing required argument","detail":"The argument is required","range":{"filename":"main.tf","start":{"line":10,"column":5},"end":{"line":10,"column":20}}}}`

	ex := models.Exercise{
		Name: "ex03",
		Path: "exercises/01_primitives/ex03.tf",
		Mode: models.ModeValidate,
	}

	diags := diagnostics.ParseDiagnostics(rawJSON, ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.Severity != diagnostics.SeverityError {
		t.Errorf("Expected severity %q, got %q", diagnostics.SeverityError, d.Severity)
	}
	if d.File != "main.tf" {
		t.Errorf("Expected file %q, got %q", "main.tf", d.File)
	}
	if d.Line != 10 {
		t.Errorf("Expected line 10, got %d", d.Line)
	}
	if d.Column != 5 {
		t.Errorf("Expected column 5, got %d", d.Column)
	}
	if d.EndLine != 10 {
		t.Errorf("Expected end line 10, got %d", d.EndLine)
	}
	if d.EndColumn != 20 {
		t.Errorf("Expected end column 20, got %d", d.EndColumn)
	}
	if d.Summary != "Missing required argument" {
		t.Errorf("Expected summary 'Missing required argument', got %q", d.Summary)
	}
	if d.Detail != "The argument is required" {
		t.Errorf("Expected detail 'The argument is required', got %q", d.Detail)
	}
}

func TestParseDiagnostics_EmptyOrClean(t *testing.T) {
	diagsEmpty := diagnostics.ParseDiagnostics("", models.Exercise{})
	if len(diagsEmpty) != 0 {
		t.Fatalf("Expected 0 diagnostics for empty output, got %d", len(diagsEmpty))
	}

	diagsClean := diagnostics.ParseDiagnostics("Success! The configuration is valid.", models.Exercise{})
	if len(diagsClean) != 0 {
		t.Fatalf("Expected 0 diagnostics for clean output, got %d", len(diagsClean))
	}
}

func TestParseDiagnostics_TerraformValidateJSON(t *testing.T) {
	validateJSON := `{
  "format_version": "1.0",
  "valid": false,
  "error_count": 1,
  "warning_count": 1,
  "diagnostics": [
    {
      "severity": "error",
      "summary": "Invalid resource type",
      "detail": "The resource type is invalid.",
      "range": {
        "filename": "infra.tf",
        "start": {
          "line": 7,
          "column": 1
        },
        "end": {
          "line": 7,
          "column": 30
        }
      }
    },
    {
      "severity": "warning",
      "summary": "Deprecated feature",
      "detail": "Use new feature instead.",
      "range": {
        "filename": "infra.tf",
        "start": {
          "line": 12,
          "column": 3
        }
      }
    }
  ]
}`

	ex := models.Exercise{
		Name: "validate_json_test",
		Path: "exercises/01_primitives/test.tf",
	}

	diags := diagnostics.ParseDiagnostics(validateJSON, ex)
	if len(diags) != 2 {
		t.Fatalf("Expected 2 diagnostics from validate JSON, got %d", len(diags))
	}

	if diags[0].Severity != diagnostics.SeverityError || diags[0].Line != 7 || diags[0].File != "infra.tf" {
		t.Errorf("Unexpected diag[0]: %+v", diags[0])
	}
	if diags[1].Severity != diagnostics.SeverityWarning || diags[1].Line != 12 || diags[1].File != "infra.tf" {
		t.Errorf("Unexpected diag[1]: %+v", diags[1])
	}
}

func TestParseDiagnostics_MarkerInDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "main.tf")
	f2 := filepath.Join(tmpDir, "variables.tf")

	_ = os.WriteFile(f1, []byte("terraform {}\n"), 0644)
	_ = os.WriteFile(f2, []byte("// line 1\n// line 2\n// I AM NOT DONE\nvariable \"foo\" {}\n"), 0644)

	ex := models.Exercise{
		Name: "dir_marker_test",
		Path: tmpDir,
	}

	diags := diagnostics.ParseDiagnostics("", ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic for marked directory, got %d", len(diags))
	}
	if diags[0].File != f2 {
		t.Errorf("Expected file %q, got %q", f2, diags[0].File)
	}
	if diags[0].Line != 3 {
		t.Errorf("Expected line 3, got %d", diags[0].Line)
	}
	if diags[0].Severity != diagnostics.SeverityWarning {
		t.Errorf("Expected severity warning, got %q", diags[0].Severity)
	}
}

func TestParseDiagnostics_AnsiEscapeCodes(t *testing.T) {
	ansiOutput := "\x1b[31mError:\x1b[0m \x1b[1mMissing required argument\x1b[0m\n\n  on main.tf line 5, in resource \"foo\" \"bar\":\n   5: content = \"test\"\n\nMissing filename argument."

	ex := models.Exercise{
		Name: "ansi_test",
		Path: "exercises/01_primitives/ansi.tf",
	}

	diags := diagnostics.ParseDiagnostics(ansiOutput, ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Summary != "Missing required argument" {
		t.Errorf("Expected clean summary without ANSI codes, got %q", diags[0].Summary)
	}
	if diags[0].Line != 5 {
		t.Errorf("Expected line 5, got %d", diags[0].Line)
	}
}

func TestParseDiagnostics_ValidValidateJSON_NoFalsePositive(t *testing.T) {
	validJSON := `{
  "format_version": "1.0",
  "valid": true,
  "error_count": 0,
  "warning_count": 0,
  "diagnostics": []
}`

	ex := models.Exercise{
		Name: "valid_json_test",
		Path: "exercises/01_primitives/valid.tf",
	}

	diags := diagnostics.ParseDiagnostics(validJSON, ex)
	if len(diags) != 0 {
		t.Fatalf("Expected 0 diagnostics for valid validate JSON, got %d (%+v)", len(diags), diags)
	}
}

func TestParseDiagnostics_SecondaryLocations(t *testing.T) {
	rawOutput := `Error: Error in function call

  on main.tf line 14, in locals:
  14:   config = templatefile("config.tpl", {})
    ├────────────────
    │ while calling templatefile(path, vars)
    │ on config.tpl line 4:
    │  4: ${nonexistent_var}
    ├────────────────

Invalid variable reference.`

	ex := models.Exercise{
		Name: "secondary_loc_test",
		Path: "exercises/01_primitives/secondary.tf",
	}

	diags := diagnostics.ParseDiagnostics(rawOutput, ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.File != "main.tf" {
		t.Errorf("Expected primary file main.tf, got %q", d.File)
	}
	if d.Line != 14 {
		t.Errorf("Expected primary line 14, got %d", d.Line)
	}
	if !strings.Contains(d.Detail, "config.tpl") {
		t.Errorf("Expected secondary location to be preserved in detail, got %q", d.Detail)
	}
}

func TestParseDiagnostics_MarkerInDirectory_SkipsHidden(t *testing.T) {
	tmpDir := t.TempDir()

	// Root exercise file without marker
	mainFile := filepath.Join(tmpDir, "main.tf")
	_ = os.WriteFile(mainFile, []byte("terraform {}\n"), 0644)

	// Hidden directories (.terraform, .git, .terralings) with marker
	tfDir := filepath.Join(tmpDir, ".terraform", "modules", "child")
	_ = os.MkdirAll(tfDir, 0755)
	_ = os.WriteFile(filepath.Join(tfDir, "sub.tf"), []byte("// I AM NOT DONE\n"), 0644)

	gitDir := filepath.Join(tmpDir, ".git")
	_ = os.MkdirAll(gitDir, 0755)
	_ = os.WriteFile(filepath.Join(gitDir, "hook.tf"), []byte("// I AM NOT DONE\n"), 0644)

	terralingsDir := filepath.Join(tmpDir, ".terralings")
	_ = os.MkdirAll(terralingsDir, 0755)
	_ = os.WriteFile(filepath.Join(terralingsDir, "cache.tf"), []byte("// I AM NOT DONE\n"), 0644)

	ex := models.Exercise{
		Name: "skip_hidden_test",
		Path: tmpDir,
	}

	diags := diagnostics.ParseDiagnostics("", ex)
	if len(diags) != 0 {
		t.Fatalf("Expected 0 diagnostics because markers are in hidden directories, got %d: %+v", len(diags), diags)
	}
}

func TestParseDiagnostics_BoxDrawingDetail(t *testing.T) {
	rawOutput := `╷
│ Error: Invalid expression
│ 
│   on main.tf line 8:
│    8: count = var.something
│ 
│ A single value is expected.
╵`

	ex := models.Exercise{
		Name: "box_drawing_test",
		Path: "exercises/01_primitives/box.tf",
	}

	diags := diagnostics.ParseDiagnostics(rawOutput, ex)
	if len(diags) != 1 {
		t.Fatalf("Expected 1 diagnostic, got %d", len(diags))
	}

	d := diags[0]
	if d.Summary != "Invalid expression" {
		t.Errorf("Expected summary 'Invalid expression', got %q", d.Summary)
	}
	if d.File != "main.tf" {
		t.Errorf("Expected file 'main.tf', got %q", d.File)
	}
	if d.Line != 8 {
		t.Errorf("Expected line 8, got %d", d.Line)
	}
	for _, l := range strings.Split(d.Detail, "\n") {
		if strings.HasPrefix(l, "│") || strings.HasPrefix(l, "|") {
			t.Errorf("Detail line still contains leading pipe: %q", l)
		}
	}
	if !strings.Contains(d.Detail, "A single value is expected.") {
		t.Errorf("Expected detail to contain explanation, got %q", d.Detail)
	}
}

func TestParseDiagnostics_StrictFallback(t *testing.T) {
	ex := models.Exercise{
		Name: "fallback_test",
		Path: "exercises/01_primitives/test.tf",
	}

	// Clean outputs with "0 failed" or "0 errors" or "errors: 0"
	testCasesClean := []string{
		"0 failed, 0 errors",
		"Test suite finished: 10 passed, 0 failed, 0 errors",
		"Apply complete! Resources: 1 added, 0 changed, 0 destroyed. (errors: 0)",
		"error_count: 0, warning_count: 0",
		"no errors found",
	}

	for _, text := range testCasesClean {
		diags := diagnostics.ParseDiagnostics(text, ex)
		if len(diags) != 0 {
			t.Errorf("Expected 0 diagnostics for clean text %q, got %d (%+v)", text, len(diags), diags)
		}
	}

	// Real unstructured error output
	realErrorText := "fatal: could not connect to registry: connection timed out"
	diagsErr := diagnostics.ParseDiagnostics(realErrorText, ex)
	if len(diagsErr) != 1 {
		t.Fatalf("Expected 1 fallback diagnostic for real error, got %d", len(diagsErr))
	}
	if diagsErr[0].Severity != diagnostics.SeverityError {
		t.Errorf("Expected severity error, got %q", diagsErr[0].Severity)
	}
	if diagsErr[0].Summary != realErrorText {
		t.Errorf("Expected summary %q, got %q", realErrorText, diagsErr[0].Summary)
	}
}
