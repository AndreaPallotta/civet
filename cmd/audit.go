package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/AndreaPallotta/civet/internal/config"
	"github.com/AndreaPallotta/civet/internal/engine"
	"github.com/AndreaPallotta/civet/internal/parser"
	"github.com/AndreaPallotta/civet/internal/provider"
	"github.com/AndreaPallotta/civet/internal/report"
	_ "github.com/AndreaPallotta/civet/internal/rules/github"
	_ "github.com/AndreaPallotta/civet/internal/rules/gitlab"
	_ "github.com/AndreaPallotta/civet/internal/rules/universal"
)

var failUnder int

var auditCmd = &cobra.Command{
	Use:   "audit [file]",
	Short: "Audit a specific CI/CD pipeline file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		if cfg == nil {
			cfg = config.DefaultConfig()
		}

		// Allow CLI override of AI provider
		if aiProvider != "" {
			cfg.AI.Provider = aiProvider
		}

		pipeline, err := parser.Parse(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse pipeline: %w", err)
		}

		rep := engine.Run(pipeline, cfg)

		// Run AI provider if configured
		if cfg.AI.Provider != "" {
			aiProv, err := provider.NewProvider(cfg)
			if err == nil {
				aiReq := &provider.AnalysisRequest{
					PipelineYAML: rep.RawYAML,
					Platform:     rep.Platform.String(),
					Findings:     rep.Findings,
				}
				aiResp, err := aiProv.Analyze(context.Background(), aiReq)
				if err == nil && aiResp != nil {
					rep.AIAnalysis = &engine.AIAnalysis{
						Summary:         aiResp.Summary,
						Recommendations: aiResp.Recommendations,
						RiskLevel:       aiResp.RiskLevel,
					}
				}
			}
		}

		// Render report based on format
		switch format {
		case "json":
			jsonRep := report.NewJSONReporter()
			if err := jsonRep.Write(os.Stdout, rep); err != nil {
				return err
			}
		case "markdown", "md":
			mdRep := report.NewMarkdownReporter()
			if err := mdRep.Write(os.Stdout, rep); err != nil {
				return err
			}
		default:
			termRep := report.NewTerminalReporter()
			if err := termRep.Write(os.Stdout, rep); err != nil {
				return err
			}
		}

		if failUnder > 0 && rep.OverallScore < failUnder {
			return fmt.Errorf("pipeline score %d is below the threshold of %d", rep.OverallScore, failUnder)
		}

		return nil
	},
}

func init() {
	auditCmd.Flags().IntVar(&failUnder, "fail-under", 0, "Fail if the overall score is below this threshold")
	rootCmd.AddCommand(auditCmd)
}
