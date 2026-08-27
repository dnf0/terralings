package diagnostics

import (
	"bufio"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/dnf0/terralings/internal/models"
)

// Severity indicates the diagnostic message level.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Diagnostic represents a standardized compiler or exercise diagnostic issue.
type Diagnostic struct {
	File      string   `json:"file"`
	Line      int      `json:"line"`
	Column    int      `json:"column"`
	EndLine   int      `json:"end_line,omitempty"`
	EndColumn int      `json:"end_column,omitempty"`
	Severity  Severity `json:"severity"`
	Summary   string   `json:"summary"`
	Detail    string   `json:"detail"`
}

var (
	ansiRegex       = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	headerRegex     = regexp.MustCompile(`(?i)^(?:[│|]\s*)?(Error|Warning|Info):\s*(.*)$`)
	locationRegex   = regexp.MustCompile(`(?i)\bon\s+([^\s:,]+)\s+line\s+(\d+)(?:,\s*(?:col(?:umn)?\s*(\d+)|in\s+([^:\n]+)))?`)
	pipePrefixRegex = regexp.MustCompile(`^[│|]\s?`)
	zeroErrorsRegex = regexp.MustCompile(`(?i)(?:(?:\b0\s+(?:errors?|failed|failures?|errored)\b)|(?:\bno\s+(?:errors?|failures?)\b)|(?:\b(?:error_count|errors?|failed|failures?)\s*[:=]\s*0\b)|(?:"error_count"\s*:\s*0))`)
)

type rawTFDiagnostic struct {
	Severity string `json:"severity"`
	Summary  string `json:"summary"`
	Detail   string `json:"detail"`
	Range    *struct {
		Filename string `json:"filename"`
		Start    *struct {
			Line   int `json:"line"`
			Column int `json:"column"`
			Byte   int `json:"byte"`
		} `json:"start"`
		End *struct {
			Line   int `json:"line"`
			Column int `json:"column"`
			Byte   int `json:"byte"`
		} `json:"end"`
	} `json:"range"`
}

type rawTFValidateOutput struct {
	FormatVersion string            `json:"format_version"`
	Valid         bool              `json:"valid"`
	ErrorCount    int               `json:"error_count"`
	WarningCount  int               `json:"warning_count"`
	Diagnostics   []rawTFDiagnostic `json:"diagnostics"`
}

type rawTFLogEntry struct {
	Level      string           `json:"@level"`
	Message    string           `json:"@message"`
	Diagnostic *rawTFDiagnostic `json:"diagnostic"`
}

func convertTFDiagnostic(rd rawTFDiagnostic) Diagnostic {
	sev := SeverityError
	switch strings.ToLower(rd.Severity) {
	case "warning":
		sev = SeverityWarning
	case "info", "information":
		sev = SeverityInfo
	}

	d := Diagnostic{
		Severity: sev,
		Summary:  rd.Summary,
		Detail:   rd.Detail,
	}

	if rd.Range != nil {
		d.File = rd.Range.Filename
		if rd.Range.Start != nil {
			d.Line = rd.Range.Start.Line
			d.Column = rd.Range.Start.Column
		}
		if rd.Range.End != nil {
			d.EndLine = rd.Range.End.Line
			d.EndColumn = rd.Range.End.Column
		}
	}

	return d
}

func parseJSONDiagnostics(raw string) ([]Diagnostic, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, false
	}

	// 1. Try Terraform validate JSON output format
	if strings.HasPrefix(trimmed, "{") {
		var valOut rawTFValidateOutput
		if err := json.Unmarshal([]byte(trimmed), &valOut); err == nil {
			if valOut.FormatVersion != "" || valOut.Valid || (valOut.ErrorCount == 0 && valOut.Diagnostics != nil) {
				var result []Diagnostic
				for _, td := range valOut.Diagnostics {
					result = append(result, convertTFDiagnostic(td))
				}
				if result == nil {
					result = []Diagnostic{}
				}
				return result, true
			}
		}

		// Try single rawTFLogEntry
		var logEntry rawTFLogEntry
		if err := json.Unmarshal([]byte(trimmed), &logEntry); err == nil && logEntry.Diagnostic != nil {
			return []Diagnostic{convertTFDiagnostic(*logEntry.Diagnostic)}, true
		}

		// Try single rawTFDiagnostic
		var singleDiag rawTFDiagnostic
		if err := json.Unmarshal([]byte(trimmed), &singleDiag); err == nil && singleDiag.Summary != "" && singleDiag.Severity != "" {
			return []Diagnostic{convertTFDiagnostic(singleDiag)}, true
		}
	}

	// 2. Try JSON Array of rawTFDiagnostic or rawTFLogEntry
	if strings.HasPrefix(trimmed, "[") {
		var diagList []rawTFDiagnostic
		if err := json.Unmarshal([]byte(trimmed), &diagList); err == nil && len(diagList) > 0 && diagList[0].Summary != "" {
			var result []Diagnostic
			for _, td := range diagList {
				result = append(result, convertTFDiagnostic(td))
			}
			return result, true
		}

		var logList []rawTFLogEntry
		if err := json.Unmarshal([]byte(trimmed), &logList); err == nil && len(logList) > 0 {
			var result []Diagnostic
			for _, le := range logList {
				if le.Diagnostic != nil {
					result = append(result, convertTFDiagnostic(*le.Diagnostic))
				}
			}
			if len(result) > 0 {
				return result, true
			}
		}
	}

	// 3. Try NDJSON (line by line)
	scanner := bufio.NewScanner(strings.NewReader(trimmed))
	var ndjsonDiags []Diagnostic
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var logEntry rawTFLogEntry
		if err := json.Unmarshal([]byte(line), &logEntry); err == nil && logEntry.Diagnostic != nil {
			ndjsonDiags = append(ndjsonDiags, convertTFDiagnostic(*logEntry.Diagnostic))
			continue
		}
		var singleDiag rawTFDiagnostic
		if err := json.Unmarshal([]byte(line), &singleDiag); err == nil && singleDiag.Summary != "" && singleDiag.Severity != "" {
			ndjsonDiags = append(ndjsonDiags, convertTFDiagnostic(singleDiag))
		}
	}
	if len(ndjsonDiags) > 0 {
		return ndjsonDiags, true
	}

	return nil, false
}

func parseTextDiagnostics(raw string) []Diagnostic {
	cleaned := ansiRegex.ReplaceAllString(raw, "")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return nil
	}

	lines := strings.Split(cleaned, "\n")
	var diags []Diagnostic

	var curDiag *Diagnostic
	var detailLines []string

	finalizeCurrent := func() {
		if curDiag != nil {
			curDiag.Detail = strings.TrimSpace(strings.Join(detailLines, "\n"))
			diags = append(diags, *curDiag)
			curDiag = nil
			detailLines = nil
		}
	}

	for _, rawLine := range lines {
		line := strings.TrimRight(rawLine, "\r")
		cleanLine := strings.TrimSpace(line)

		// Strip decorative box drawing characters (e.g. ╷, ╵, │, |)
		if cleanLine == "╷" || cleanLine == "╵" {
			continue
		}
		cleanLine = strings.TrimPrefix(cleanLine, "│")
		cleanLine = strings.TrimPrefix(cleanLine, "|")
		cleanLine = strings.TrimSpace(cleanLine)

		if hm := headerRegex.FindStringSubmatch(cleanLine); hm != nil {
			finalizeCurrent()

			sev := SeverityError
			switch strings.ToLower(hm[1]) {
			case "warning":
				sev = SeverityWarning
			case "info":
				sev = SeverityInfo
			}

			curDiag = &Diagnostic{
				Severity: sev,
				Summary:  strings.TrimSpace(hm[2]),
			}
			continue
		}

		if curDiag != nil {
			if curDiag.Summary == "" && cleanLine != "" && !strings.HasPrefix(cleanLine, "on ") {
				curDiag.Summary = cleanLine
				continue
			}

			if locMatch := locationRegex.FindStringSubmatch(cleanLine); locMatch != nil {
				if curDiag.File == "" {
					curDiag.File = locMatch[1]
					if lineNum, err := strconv.Atoi(locMatch[2]); err == nil {
						curDiag.Line = lineNum
					}
					if len(locMatch) > 3 && locMatch[3] != "" {
						if colNum, err := strconv.Atoi(locMatch[3]); err == nil {
							curDiag.Column = colNum
						}
					}
					continue
				}
			}

			cleanedDetailLine := pipePrefixRegex.ReplaceAllString(line, "")
			detailLines = append(detailLines, cleanedDetailLine)
		}
	}

	finalizeCurrent()

	// Fallback: If no structured diagnostics were extracted but raw output clearly indicates an error/failure
	if len(diags) == 0 {
		lower := strings.ToLower(cleaned)
		lowerWithoutZeros := zeroErrorsRegex.ReplaceAllString(lower, "")
		if strings.Contains(lowerWithoutZeros, "error") ||
			strings.Contains(lowerWithoutZeros, "failed") ||
			strings.Contains(lowerWithoutZeros, "failure") ||
			strings.Contains(lowerWithoutZeros, "fatal") {
			nonEmpty := []string{}
			for _, l := range lines {
				tl := strings.TrimSpace(l)
				if tl != "" && tl != "╷" && tl != "╵" {
					nonEmpty = append(nonEmpty, tl)
				}
			}
			if len(nonEmpty) > 0 {
				summary := nonEmpty[0]
				detail := ""
				if len(nonEmpty) > 1 {
					detail = strings.Join(nonEmpty[1:], "\n")
				}
				diags = append(diags, Diagnostic{
					Severity: SeverityError,
					Summary:  summary,
					Detail:   detail,
				})
			}
		}
	}

	return diags
}

// ParseDiagnostics parses compiler output (text or JSON) into a normalized slice of Diagnostics.
func ParseDiagnostics(rawOutput string, ex models.Exercise) []Diagnostic {
	var results []Diagnostic

	trimmed := strings.TrimSpace(rawOutput)
	if trimmed != "" {
		if jsonDiags, ok := parseJSONDiagnostics(trimmed); ok {
			results = append(results, jsonDiags...)
		} else {
			textDiags := parseTextDiagnostics(rawOutput)
			results = append(results, textDiags...)
		}
	}

	if results == nil {
		results = []Diagnostic{}
	}
	return results
}

// Parse is an alias for ParseDiagnostics.
func Parse(rawOutput string, ex models.Exercise) []Diagnostic {
	return ParseDiagnostics(rawOutput, ex)
}
