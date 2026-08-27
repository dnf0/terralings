package test

import (
	"strings"
	"testing"

	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/search"
	"github.com/dnf0/terralings/internal/ui"
)

func TestSearch_ExactNameMatch(t *testing.T) {
	m := manifest.GetManifest()
	results := search.SearchExercises(m, "primitives01")
	if len(results) == 0 {
		t.Fatal("Expected results for 'primitives01', got none")
	}

	if results[0].Exercise.Name != "primitives01" {
		t.Errorf("Expected first match to be primitives01, got %s", results[0].Exercise.Name)
	}
}

func TestSearch_ChapterAndTopicMatch(t *testing.T) {
	m := manifest.GetManifest()
	results := search.SearchExercises(m, "dynamic")
	if len(results) < 4 {
		t.Errorf("Expected at least 4 dynamic results, got %d", len(results))
	}

	for _, r := range results {
		nameMatches := strings.Contains(strings.ToLower(r.Exercise.Name), "dynamic")
		titleMatches := strings.Contains(strings.ToLower(r.Exercise.Title), "dynamic")
		chapMatches := strings.Contains(strings.ToLower(r.Exercise.ChapterName), "dynamic")
		if !nameMatches && !titleMatches && !chapMatches {
			t.Errorf("Result %s does not contain 'dynamic'", r.Exercise.Name)
		}
	}
}

func TestSearch_HintKeywordMatch(t *testing.T) {
	m := manifest.GetManifest()
	results := search.SearchExercises(m, "encryption")
	if len(results) == 0 {
		t.Fatal("Expected results searching for hint keyword 'encryption', got 0")
	}

	foundTofu := false
	for _, r := range results {
		if strings.HasPrefix(r.Exercise.Name, "tofu") {
			foundTofu = true
		}
	}
	if !foundTofu {
		t.Errorf("Expected tofu exercise in encryption search results")
	}
}

func TestSearch_NoMatch(t *testing.T) {
	m := manifest.GetManifest()
	results := search.SearchExercises(m, "xyznonexistentquery123")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent query, got %d", len(results))
	}
}

func TestSearch_FormatSearchResults(t *testing.T) {
	m := manifest.GetManifest()
	results := search.SearchExercises(m, "variables")
	output := ui.FormatSearchResults("variables", results)

	if !strings.Contains(output, "variables") {
		t.Errorf("Expected formatted output to contain 'variables', got:\n%s", output)
	}
	if !strings.Contains(output, "matches found") && !strings.Contains(output, "Search Results") {
		t.Errorf("Expected header in search output, got:\n%s", output)
	}
}
