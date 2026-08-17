package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreaPallotta/civet/internal/config"
	"github.com/AndreaPallotta/civet/internal/engine"
	"github.com/AndreaPallotta/civet/internal/parser"
	"github.com/AndreaPallotta/civet/internal/report"
	"github.com/AndreaPallotta/civet/internal/scanner"
)

var scanFailUnder int

var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Short: "Scan a directory for pipeline files and audit them",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		files, err := scanner.Scan(dir)
		if err != nil {
			return fmt.Errorf("failed to scan directory: %w", err)
		}

		if len(files) == 0 {
			fmt.Println("No pipeline files found.")
			return nil
		}

		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		type result struct {
			path     string
			platform string
			score    int
		}
		var results []result
		lowestScore := 100

		for _, file := range files {
			pipeline, err := parser.Parse(file)
			if err != nil {
				fmt.Printf("Failed to parse %s: %v\n", file, err)
				continue
			}

			rep := engine.Run(pipeline, cfg)

			results = append(results, result{
				path:     file,
				platform: rep.Platform.String(),
				score:    rep.OverallScore,
			})

			if rep.OverallScore < lowestScore {
				lowestScore = rep.OverallScore
			}
		}

		if llmMode || format == "llm" {
			type scanLLMOutput struct {
				Header    string                 `json:"system_prompt_context"`
				TotalDone int                    `json:"total_pipelines_scanned"`
				Pipelines []report.LLMContext    `json:"pipelines"`
			}

			var contexts []report.LLMContext
			for _, file := range files {
				pipeline, err := parser.Parse(file)
				if err != nil {
					continue
				}
				rep := engine.Run(pipeline, cfg)
				
				grade := "HEALTHY"
				if rep.OverallScore < 60 {
					grade = "CRITICAL_ATTENTION_REQUIRED"
				} else if rep.OverallScore < 85 {
					grade = "NEEDS_IMPROVEMENT"
				}

				var findings []report.LLMFinding
				var actionItems []string
				for _, f := range rep.Findings {
					findings = append(findings, report.LLMFinding{
						RuleID:      f.RuleID,
						RuleName:    f.RuleName,
						Severity:    f.Severity,
						Category:    f.Category,
						Location:    f.Location,
						Description: f.Message,
						ProposedFix: f.Suggestion,
					})
					actionItems = append(actionItems, fmt.Sprintf("[%s] Resolve %s: %s", f.Severity.String(), f.RuleID, f.Suggestion))
				}

				contexts = append(contexts, report.LLMContext{
					Header: "Civet Deterministic CI/CD Quality Audit",
					Pipeline: report.LLMPipelineInfo{
						FilePath: rep.FilePath,
						Platform: rep.Platform.String(),
					},
					Scores:      rep.CategoryScores,
					TotalScore:  rep.OverallScore,
					HealthGrade: grade,
					Findings:    findings,
					ActionItems: actionItems,
				})
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(scanLLMOutput{
				Header:    "Civet Workspace Pipeline Scan. Review these pipelines and prioritize refactoring based on health_grade.",
				TotalDone: len(contexts),
				Pipelines: contexts,
			})
		}

		fmt.Println("Scan Results:")
		fmt.Printf("%-50s %-18s %s\n", "File", "Platform", "Score")
		for _, r := range results {
			fmt.Printf("%-50s %-18s %d/100\n", r.path, r.platform, r.score)
		}

		if scanFailUnder > 0 && lowestScore < scanFailUnder {
			return fmt.Errorf("lowest pipeline score %d is below the threshold of %d", lowestScore, scanFailUnder)
		}

		return nil
	},
}

func init() {
	scanCmd.Flags().IntVar(&scanFailUnder, "fail-under", 0, "Fail if the lowest scoring pipeline is below this threshold")
	rootCmd.AddCommand(scanCmd)
}
