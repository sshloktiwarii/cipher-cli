package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cipher",
	Short: "Cipher CLI — Fleet Commander and Control Plane for the Cipher AI PaaS",
	Long: `Cipher CLI orchestrates the full lifecycle of the Cipher Local-First AI PaaS.
It programmatically controls networking, container execution, and unified log aggregation
across all 8 platform services on the custom cipher_net bridge network.`,
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
