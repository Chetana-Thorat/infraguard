package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Chetana-Thorat/infraguard/internal/checks"
	"github.com/Chetana-Thorat/infraguard/internal/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:           "scan [path]",
	Short:         "Scan a directory for Terraform and YAML files",
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}

		files, err := scanner.ScanDirectory(path)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			fmt.Println("No Terraform or YAML files found.")
			return nil
		}

		rules := []checks.Check{
			checks.SSHOpenToWorldCheck{},
			checks.K8sResourceLimitsCheck{},
		}

		hasFindings := false

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			for _, rule := range rules {
				findings, err := rule.Run(file, content)
				if err != nil {
					return err
				}

				for _, finding := range findings {
					hasFindings = true
					printFinding(
						finding.RuleID,
						string(finding.Severity),
						finding.FilePath,
						finding.Message,
						finding.Recommendation,
					)
				}
			}
		}

		if hasFindings {
			return fmt.Errorf("policy violations found")
		}

		fmt.Println("No policy violations found.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

func printFinding(ruleID, severity, filePath, message, recommendation string) {
	fmt.Println()
	fmt.Printf("❌ %s\n", strings.ToUpper(severity))
	fmt.Printf("Rule: %s\n", ruleID)
	fmt.Printf("File: %s\n", filePath)
	fmt.Printf("%s\n", message)
	fmt.Printf("Recommendation: %s\n", recommendation)
}
