package state

import (
	"time"

	"github.com/dnf0/terralings/internal/models"
)

// ChapterSummary holds metrics for a single curriculum chapter.
type ChapterSummary struct {
	ChapterID     string `json:"chapter_id"`
	Title         string `json:"title"`
	Total         int    `json:"total"`
	Completed     int    `json:"completed"`
	TotalAttempts int    `json:"total_attempts"`
	TotalHints    int    `json:"total_hints"`
}

// AnalyticsSummary provides aggregate metrics across all chapters and exercises.
type AnalyticsSummary struct {
	TotalExercises   int              `json:"total_exercises"`
	CompletedCount   int              `json:"completed_count"`
	InProgressCount  int              `json:"in_progress_count"`
	TotalAttempts    int              `json:"total_attempts"`
	TotalHintsViewed int              `json:"total_hints_viewed"`
	TotalTimeSpent   time.Duration    `json:"total_time_spent"`
	ChapterSummaries []ChapterSummary `json:"chapter_summaries"`
}

// GetAnalytics computes progress and learning analytics based on the provided manifest.
func (s *Store) GetAnalytics(m *models.Manifest) AnalyticsSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	summary := AnalyticsSummary{}
	var totalSecs int64

	if m == nil {
		summary.TotalExercises = len(s.data.Exercises)
		for _, ex := range s.data.Exercises {
			if ex == nil {
				continue
			}
			if ex.Status == StatusPassed {
				summary.CompletedCount++
			} else if ex.Status == StatusInProgress {
				summary.InProgressCount++
			}
			summary.TotalAttempts += ex.Attempts
			summary.TotalHintsViewed += ex.HintsViewed
			totalSecs += ex.TimeSpentSeconds
		}
		if s.data.TotalTimeSpentSeconds > totalSecs {
			totalSecs = s.data.TotalTimeSpentSeconds
		}
		summary.TotalTimeSpent = time.Duration(totalSecs) * time.Second
		return summary
	}

	summary.TotalExercises = len(m.AllExercises())
	summary.ChapterSummaries = make([]ChapterSummary, 0, len(m.Chapters))

	for _, ch := range m.Chapters {
		cs := ChapterSummary{
			ChapterID: ch.Name,
			Title:     ch.Title,
			Total:     len(ch.Exercises),
		}

		for _, ex := range ch.Exercises {
			st, ok := s.data.Exercises[ex.Name]
			if ok && st != nil {
				if st.Status == StatusPassed {
					cs.Completed++
					summary.CompletedCount++
				} else if st.Status == StatusInProgress {
					summary.InProgressCount++
				}
				cs.TotalAttempts += st.Attempts
				summary.TotalAttempts += st.Attempts
				cs.TotalHints += st.HintsViewed
				summary.TotalHintsViewed += st.HintsViewed
				totalSecs += st.TimeSpentSeconds
			}
		}

		summary.ChapterSummaries = append(summary.ChapterSummaries, cs)
	}

	if s.data.TotalTimeSpentSeconds > totalSecs {
		totalSecs = s.data.TotalTimeSpentSeconds
	}
	summary.TotalTimeSpent = time.Duration(totalSecs) * time.Second

	return summary
}
