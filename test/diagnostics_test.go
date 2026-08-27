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
