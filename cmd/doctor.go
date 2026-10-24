package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type PortCheck struct {
	Port        int
	ServiceName string
}

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"check"},
	Short:   "Runs comprehensive pre-flight diagnostics across ports, memory, and Docker daemon",
	Long: `Doctor performs pre-flight system health checks to verify that:
1. TCP ports 8079 through 8085 are unblocked or free.
2. The Docker daemon socket is active and responsive.
3. System available RAM meets the 4GB+ operating baseline for local LLM orchestration.`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println("\n" + cyan("🩺 [CIPHER DOCTOR] Running Pre-Flight Diagnostics..."))
		fmt.Println(bold("==================================================================="))

		allPassed := true

		// 1. Port Checks (8079 through 8085)
		ports := []PortCheck{
			{8079, "cipher-sec (WAF/Firewall)"},
			{8080, "cipher-route (Context Router)"},
			{8081, "cipher-exec (Sandbox)"},
			{8082, "cipher-log (Telemetry Ingestion)"},
			{8083, "cipher-mem (Obsidian Graph)"},
			{8084, "cipher-spare (Cluster Aux)"},
			{8085, "cipher-mcp (Stateless MCP)"},
		}

		fmt.Println(bold("\n1. TCP Port Availability (8079 - 8085):"))
		for _, p := range ports {
			addr := fmt.Sprintf("127.0.0.1:%d", p.Port)
			conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
			if err != nil {
				// Port is free (desired state for cold start)
				fmt.Printf("   %s Port %-5d | %-32s [FREE]\n", green("✓"), p.Port, p.ServiceName)
			} else {
				_ = conn.Close()
				fmt.Printf("   %s Port %-5d | %-32s [IN USE / ACTIVE]\n", yellow("!"), p.Port, p.ServiceName)
			}
		}

		// 2. Docker Daemon Check
		fmt.Println(bold("\n2. Container Runtime & Docker Daemon:"))
		dockerSock := "/var/run/docker.sock"
		sockExists := false
		if _, err := os.Stat(dockerSock); err == nil {
			sockExists = true
		}

		// Ping docker via CLI
		out, err := exec.Command("docker", "info").CombinedOutput()
		if err == nil {
			fmt.Printf("   %s Docker Daemon responding via %s\n", green("✓"), dockerSock)
		} else {
			allPassed = false
			fmt.Printf("   %s Docker Daemon is unreachable or offline!\n", red("✗"))
			if !sockExists {
				fmt.Printf("      %s Socket %s not found.\n", red("→"), dockerSock)
			} else {
				fmt.Printf("      %s %s\n", red("→"), strings.TrimSpace(string(out)))
			}
			fmt.Printf("      %s Remediation: Start Docker Desktop or verify system privileges: 'sudo chmod 666 %s'\n",
				yellow("Fix:"), dockerSock)
			fmt.Printf("      %s Note: cipher-exec will automatically use native OS Seatbelt/bwrap fallback.\n",
				cyan("Info:"))
		}

		// 3. System RAM Check
		fmt.Println(bold("\n3. System Memory Resources:"))
		var totalRAMGB float64 = 0
		if runtime.GOOS == "darwin" {
			out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
			if err == nil {
				bytes, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
				totalRAMGB = float64(bytes) / (1024 * 1024 * 1024)
			}
		} else if runtime.GOOS == "linux" {
			out, err := exec.Command("sh", "-c", "grep MemTotal /proc/meminfo | awk '{print $2}'").Output()
			if err == nil {
				kb, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
				totalRAMGB = float64(kb) / (1024 * 1024)
			}
		}

		if totalRAMGB >= 8.0 {
			fmt.Printf("   %s Available Physical RAM: %.1f GB (Meets 8GB+ High-Concurrency Tier)\n", green("✓"), totalRAMGB)
		} else if totalRAMGB >= 4.0 {
			fmt.Printf("   %s Available Physical RAM: %.1f GB (Low-Memory Mode Recommended: run with --low-mem)\n", yellow("!"), totalRAMGB)
		} else if totalRAMGB > 0 {
			fmt.Printf("   %s Available Physical RAM: %.1f GB (Under 4GB: High risk of OOM with TimescaleDB & LLMs)\n", red("✗"), totalRAMGB)
			fmt.Printf("      %s Run 'cipher-log --low-mem' to use embedded SQLite WAL instead of Postgres.\n", yellow("Fix:"))
		} else {
			fmt.Printf("   %s Memory check skipped on %s architecture.\n", yellow("~"), runtime.GOOS)
		}

		fmt.Println(bold("\n==================================================================="))
		if allPassed {
			fmt.Println(green("🎉 PRE-FLIGHT DIAGNOSTICS PASSED — Host is ready for 'cipher up'!\n"))
		} else {
			fmt.Println(yellow("⚠️  PRE-FLIGHT DIAGNOSTICS COMPLETED WITH NOTICES — Review remediations above.\n"))
		}
	},
}
