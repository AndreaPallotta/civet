package report

import (
	"fmt"
	"io"

	"github.com/AndreaPallotta/civet/internal/engine"
)

type MarkdownReporter struct{}

func NewMarkdownReporter() *MarkdownReporter {
	return &MarkdownReporter{}
}

func (r *MarkdownReporter) Write(w io.Writer, rep *engine.Report) error {
	fmt.Fprintf(w, "# Civet Report\n\n")
	fmt.Fprintf(w, "**File**: `%s`\n", rep.FilePath)
	fmt.Fprintf(w, "**Platform**: `%s`\n", rep.Platform.String())
	fmt.Fprintf(w, "**Overall Score**: %d\n\n", rep.OverallScore)

	fmt.Fprintf(w, "## Category Scores\n\n")
	fmt.Fprintf(w, "| Category | Score |\n")
	fmt.Fprintf(w, "|----------|-------|\n")
	for cat, score := range rep.CategoryScores {
		fmt.Fprintf(w, "| %s | %d |\n", cat, score)
	}
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "## Findings\n\n")
	for _, f := range rep.Findings {
		fmt.Fprintf(w, "### [%s] %s\n\n", f.Severity.String(), f.RuleID)
		fmt.Fprintf(w, "%s\n\n", f.Message)
		if f.Suggestion != "" {
			fmt.Fprintf(w, "**Suggestion**: %s\n\n", f.Suggestion)
		}
		if f.Location.Line > 0 {
			fmt.Fprintf(w, "**Location**: Line %d\n\n", f.Location.Line)
		}
	}

	if rep.AIAnalysis != nil {
		fmt.Fprintf(w, "## AI Analysis\n\n")
		fmt.Fprintf(w, "**Risk Level**: %s\n\n", rep.AIAnalysis.RiskLevel)
		fmt.Fprintf(w, "**Summary**: %s\n\n", rep.AIAnalysis.Summary)
		fmt.Fprintf(w, "### Recommendations\n\n")
		for _, rec := range rep.AIAnalysis.Recommendations {
			fmt.Fprintf(w, "- %s\n", rec)
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}
