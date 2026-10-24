package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type FleetService struct {
	Name string
	Port int
	URL  string
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays health, ports, and execution state across all Foundry services",
	Run: func(cmd *cobra.Command, args []string) {
		services := []FleetService{
			{"foundry-sentinel", 8079, "http://localhost:8079/healthz"},
			{"foundry-router", 8080, "http://localhost:8080/healthz"},
			{"foundry-crucible", 8081, "http://localhost:8081/healthz"},
			{"foundry-pulse", 8082, "http://localhost:8082/healthz"},
			{"foundry-graph", 8083, "http://localhost:8083/health"},
			{"foundry-airgap", 8085, "http://localhost:8085/actions"},
		}

		client := &http.Client{Timeout: 1 * time.Second}

		fmt.Println("\n" + color.CyanString("⚡ FOUNDRY PLATFORM HEALTH & PORT STATUS"))
		fmt.Printf("%-20s %-8s %-12s %-20s\n", "SERVICE", "PORT", "CONTAINER", "HTTP HEALTH")
		fmt.Println("-------------------------------------------------------------------")

		for _, s := range services {
			// Check docker container
			out, _ := exec.Command("docker", "inspect", "-f", "{{.State.Status}}", s.Name).Output()
			cStatus := strings.TrimSpace(string(out))
			if cStatus == "" {
				cStatus = "offline"
			}

			// Check HTTP
			hStatus := color.RedString("UNREACHABLE")
			resp, err := client.Get(s.URL)
			if err == nil {
				if resp.StatusCode == 200 || resp.StatusCode == 204 {
					hStatus = color.GreenString("200 OK")
				} else {
					hStatus = color.YellowString("HTTP %d", resp.StatusCode)
				}
				_ = resp.Body.Close()
			}

			cColor := color.GreenString(cStatus)
			if cStatus != "running" {
				cColor = color.RedString(cStatus)
			}

			fmt.Printf("%-20s :%-7d %-22s %s\n", s.Name, s.Port, cColor, hStatus)
		}
		fmt.Println("-------------------------------------------------------------------\n")
	},
}
