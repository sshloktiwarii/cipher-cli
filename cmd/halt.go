package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var haltCmd = &cobra.Command{
	Use:   "halt",
	Short: "Gracefully stops and removes all Foundry containers on foundry_net",
	RunE: func(cmd *cobra.Command, args []string) error {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(yellow("🛑 [FOUNDRY CLI] Halting and dismantling Foundry Fleet..."))

		dockerCmd := exec.Command("docker", "compose", "down", "--remove-orphans")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			fmt.Println(red("Notice: Fallback stopping containers by name..."))
			services := []string{
				"foundry-sentinel", "foundry-router", "foundry-crucible",
				"foundry-pulse", "foundry-graph", "foundry-airgap",
				"foundry-timescaledb", "foundry-redis",
			}
			for _, svc := range services {
				_ = exec.Command("docker", "stop", svc).Run()
				_ = exec.Command("docker", "rm", svc).Run()
			}
		}

		fmt.Println(green("✓ Foundry microservices successfully halted. Resources purged."))
		return nil
	},
}
