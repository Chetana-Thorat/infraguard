package cmd

import (
	"fmt"

	"github.com/Chetana-Thorat/infraguard/internal/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a directory for Terraform and YAML files",
	Args:  cobra.MaximumNArgs(1),
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

		fmt.Println("Found files:")
		for _, f := range files {
			fmt.Println("-", f)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
