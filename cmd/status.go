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
	Short: "Displays health, ports, and execution state across all Cipher services",
	Run: func(cmd *cobra.Command, args []string) {
		services := []FleetService{
			{"cipher-sec", 8079, "http://localhost:8079/healthz"},
			{"cipher-route", 8080, "http://localhost:8080/healthz"},
			{"cipher-exec", 8081, "http://localhost:8081/healthz"},
			{"cipher-log", 8082, "http://localhost:8082/healthz"},
			{"cipher-mem", 8083, "http://localhost:8083/health"},
			{"cipher-mcp", 8085, "http://localhost:8085/actions"},
		}

		client := &http.Client{Timeout: 1 * time.Second}

		fmt.Println("\n" + color.CyanString("⚡ CIPHER PLATFORM HEALTH & PORT STATUS"))
		fmt.Printf("%-20s %-8s %-12s %-20s\n", "SERVICE", "PORT", "CONTAINER", "HTTP HEALTH")
		fmt.Println("-------------------------------------------------------------------")

		for _, s := range services {
			// Check docker container
			out, _ := exec.Command("docker", "inspect", "-f", "{{.State.Status}}", s.Name).Output()
			cStatus := strings.TrimSpace(string(out))
			if cStatus == "" {
				cStatus = "offline"
			}

			// Format container color
			var cColor string
			if cStatus == "running" {
				cColor = color.GreenString("%-12s", cStatus)
			} else {
				cColor = color.RedString("%-12s", cStatus)
			}

			// Check HTTP health
			hStatus := "offline"
			resp, err := client.Get(s.URL)
			if err == nil {
				hStatus = fmt.Sprintf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
				_ = resp.Body.Close()
			}

			var hColor string
			if strings.HasPrefix(hStatus, "200") || strings.HasPrefix(hStatus, "202") {
				hColor = color.GreenString(hStatus)
			} else {
				hColor = color.RedString(hStatus)
			}

			fmt.Printf("%-20s %-8d %s %s\n", s.Name, s.Port, cColor, hColor)
		}
		fmt.Println("-------------------------------------------------------------------")
	},
}
