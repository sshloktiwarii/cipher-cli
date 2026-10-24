package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ServiceColor struct {
	Name   string
	Prefix string
	Color  *color.Color
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Attaches to and multiplexes colored stdout/stderr streams from all Foundry services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := []ServiceColor{
			{"foundry-sentinel", "[SENTINEL]", color.New(color.FgHiRed, color.Bold)},
			{"foundry-router", "[ROUTER]  ", color.New(color.FgCyan, color.Bold)},
			{"foundry-crucible", "[CRUCIBLE]", color.New(color.FgMagenta, color.Bold)},
			{"foundry-pulse", "[PULSE]   ", color.New(color.FgBlue, color.Bold)},
			{"foundry-graph", "[GRAPH]   ", color.New(color.FgGreen, color.Bold)},
			{"foundry-airgap", "[AIRGAP]  ", color.New(color.FgYellow, color.Bold)},
			{"foundry-timescaledb", "[DB]      ", color.New(color.FgWhite)},
			{"foundry-redis", "[REDIS]   ", color.New(color.FgHiBlack)},
		}

		fmt.Println(color.CyanString("⚡ [FOUNDRY CLI] Multiplexing Fleet Telemetry & Logs (Ctrl+C to detach)..."))

		var wg sync.WaitGroup
		ctxStop := make(chan os.Signal, 1)
		signal.Notify(ctxStop, os.Interrupt, syscall.SIGTERM)

		for _, svc := range services {
			wg.Add(1)
			go func(sc ServiceColor) {
				defer wg.Done()
				tailContainerLogs(sc)
			}(svc)
		}

		<-ctxStop
		fmt.Println(color.YellowString("\nDetached from fleet logs."))
		return nil
	},
}

func tailContainerLogs(sc ServiceColor) {
	cmd := exec.Command("docker", "logs", "-f", "--tail", "20", sc.Name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				_ = err
			}
			break
		}
		prefix := sc.Color.Sprintf("%s", sc.Prefix)
		fmt.Printf("%s %s", prefix, line)
	}
	_ = cmd.Wait()
}
