package search

import (
	"sort"
	"strings"

	"github.com/dnf0/terralings/internal/models"
)

// SearchResult contains an exercise and relevance match metadata.
type SearchResult struct {
	Exercise  models.Exercise
	MatchedIn string
	Score     int
}

// SearchExercises searches across manifest exercises matching a query string.
func SearchExercises(m *models.Manifest, query string) []SearchResult {
	if m == nil {
		return nil
	}
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return nil
	}

	var results []SearchResult

	for _, ch := range m.Chapters {
		chName := strings.ToLower(ch.Name)
		chTitle := strings.ToLower(ch.Title)
		chDesc := strings.ToLower(ch.Description)

		for _, ex := range ch.Exercises {
			exName := strings.ToLower(ex.Name)
			exTitle := strings.ToLower(ex.Title)

			score := 0
			var matchedFields []string

			// Exact name match -> highest score
			if exName == q {
				score += 100
				matchedFields = append(matchedFields, "exact name")
			} else if strings.Contains(exName, q) {
				score += 50
				matchedFields = append(matchedFields, "name")
			}

			if strings.Contains(exTitle, q) {
				score += 30
				matchedFields = append(matchedFields, "title")
			}

			if strings.Contains(chTitle, q) || strings.Contains(chName, q) || strings.Contains(chDesc, q) {
				score += 20
				matchedFields = append(matchedFields, "chapter")
			}

			for _, hint := range ex.Hints {
				if strings.Contains(strings.ToLower(hint), q) {
					score += 10
					matchedFields = append(matchedFields, "hint")
					break
				}
			}

			if score > 0 {
				results = append(results, SearchResult{
					Exercise:  ex,
					MatchedIn: strings.Join(matchedFields, ", "),
					Score:     score,
				})
			}
		}
	}

	// Sort results descending by score, then ascending by exercise name
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Exercise.Name < results[j].Exercise.Name
	})

	return results
}
