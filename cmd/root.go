package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "foundry",
	Short: "Foundry CLI — Fleet Commander and Control Plane for the Foundry AI PaaS",
	Long: `Foundry CLI orchestrates the full lifecycle of the Foundry Local-First AI PaaS.
It programmatically controls networking, container execution, and unified log aggregation
across all 8 platform services on the custom foundry_net bridge network.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(haltCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(statusCmd)
}
