package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/sshloktiwarii/cipher-cli/internal"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Spins up the full Cipher microservice fleet on the cipher_net network",
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Println(cyan("⚡ [CIPHER CLI] Generating Fleet Deployment Spec..."))

		composePath, err := internal.WriteComposeFile(".")
		if err != nil {
			return fmt.Errorf("failed to write compose configuration: %w", err)
		}
		fmt.Printf("%s Generated %s\n", green("✓"), composePath)

		fmt.Println(yellow("🚀 Booting Cipher Fleet on custom bridge network 'cipher_net'..."))

		// Execute docker compose up -d
		dockerCmd := exec.Command("docker", "compose", "-f", composePath, "up", "-d")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			return fmt.Errorf("docker compose execution failed: %w", err)
		}

		fmt.Println("\n" + green("=================================================="))
		fmt.Println(green("🎉 CIPHER PLATFORM ONLINE — 8 SERVICES RUNNING"))
		fmt.Println(green("=================================================="))
		fmt.Printf("• %-20s %s\n", "Semantic Firewall:", "http://localhost:8079 (cipher-sec)")
		fmt.Printf("• %-20s %s\n", "Context Router:", "http://localhost:8080 (cipher-route)")
		fmt.Printf("• %-20s %s\n", "Execution Sandbox:", "http://localhost:8081 (cipher-exec)")
		fmt.Printf("• %-20s %s\n", "Telemetry Engine:", "http://localhost:8082 (cipher-log)")
		fmt.Printf("• %-20s %s\n", "Obsidian Vector DB:", "http://localhost:8083 (cipher-mem)")
		fmt.Printf("• %-20s %s\n", "MCP Git HITL Server:", "http://localhost:8085 (cipher-mcp)")
		fmt.Printf("• %-20s %s\n", "TimescaleDB:", "localhost:5432")
		fmt.Printf("• %-20s %s\n", "Redis Cluster:", "localhost:6379")
		fmt.Println(green("=================================================="))
		fmt.Println("Run " + cyan("cipher logs") + " to attach to live colored stdout streams.")
		return nil
	},
}
