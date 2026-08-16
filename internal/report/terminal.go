package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/AndreaPallotta/civet/internal/engine"
	"github.com/AndreaPallotta/civet/internal/rules"
	"github.com/fatih/color"
)

type TerminalReporter struct{}

func NewTerminalReporter() *TerminalReporter {
	return &TerminalReporter{}
}

func (r *TerminalReporter) Write(w io.Writer, rep *engine.Report) error {
	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Fprintf(w, "File: %s | Platform: %s\n", rep.FilePath, rep.Platform.String())

	var scoreColor *color.Color
	if rep.OverallScore >= 80 {
		scoreColor = color.New(color.FgGreen, color.Bold)
	} else if rep.OverallScore >= 60 {
		scoreColor = color.New(color.FgYellow, color.Bold)
	} else {
		scoreColor = color.New(color.FgRed, color.Bold)
	}

	fmt.Fprintf(w, "Overall Score: ")
	scoreColor.Fprintf(w, "%d\n\n", rep.OverallScore)

	fmt.Fprintf(w, "Category Scores:\n")
	for cat, score := range rep.CategoryScores {
		filled := score / 10
		empty := 10 - filled
		if filled < 0 {
			filled = 0
		}
		if empty < 0 {
			empty = 0
		}
		if filled > 10 {
			filled = 10
			empty = 0
		}
		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		fmt.Fprintf(w, "  %-15s [%s] %d\n", cat, bar, score)
	}
	fmt.Fprintf(w, "\n")

	errors := []rules.Finding{}
	warnings := []rules.Finding{}
	infos := []rules.Finding{}

	for _, f := range rep.Findings {
		switch f.Severity {
		case rules.SeverityError:
			errors = append(errors, f)
		case rules.SeverityWarning:
			warnings = append(warnings, f)
		default:
			infos = append(infos, f)
		}
	}

	errColor := color.New(color.FgRed, color.Bold)
	warnColor := color.New(color.FgYellow, color.Bold)
	infoColor := color.New(color.FgBlue, color.Bold)

	printFindings := func(fs []rules.Finding, c *color.Color, tag string) {
		for _, f := range fs {
			c.Fprintf(w, "[%s]", tag)
			fmt.Fprintf(w, " %s: %s\n", f.RuleID, f.Message)
			if f.Suggestion != "" {
				fmt.Fprintf(w, "  Suggestion: %s\n", f.Suggestion)
			}
			if f.Location.Line > 0 {
				fmt.Fprintf(w, "  Location: line %d\n", f.Location.Line)
			}
			fmt.Fprintf(w, "\n")
		}
	}

	printFindings(errors, errColor, "ERROR")
	printFindings(warnings, warnColor, "WARN")
	printFindings(infos, infoColor, "INFO")

	if rep.AIAnalysis != nil {
		aiColor := color.New(color.FgMagenta, color.Bold)
		aiColor.Fprintf(w, "AI Analysis\n")
		fmt.Fprintf(w, "Risk Level: %s\n", rep.AIAnalysis.RiskLevel)
		fmt.Fprintf(w, "Summary: %s\n", rep.AIAnalysis.Summary)
		fmt.Fprintf(w, "Recommendations:\n")
		for _, rec := range rep.AIAnalysis.Recommendations {
			fmt.Fprintf(w, "  - %s\n", rec)
		}
		fmt.Fprintf(w, "\n")
	}

	return nil
}
