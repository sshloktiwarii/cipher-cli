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
	Short: "Attaches to and multiplexes colored stdout/stderr streams from all Cipher services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services := []ServiceColor{
			{"cipher-sec", "[CIPHER-SEC]  ", color.New(color.FgHiRed, color.Bold)},
			{"cipher-route", "[CIPHER-ROUTE]", color.New(color.FgCyan, color.Bold)},
			{"cipher-exec", "[CIPHER-EXEC] ", color.New(color.FgMagenta, color.Bold)},
			{"cipher-log", "[CIPHER-LOG]  ", color.New(color.FgBlue, color.Bold)},
			{"cipher-mem", "[CIPHER-MEM]  ", color.New(color.FgGreen, color.Bold)},
			{"cipher-mcp", "[CIPHER-MCP]  ", color.New(color.FgYellow, color.Bold)},
			{"cipher-timescaledb", "[DB]          ", color.New(color.FgWhite)},
			{"cipher-redis", "[REDIS]       ", color.New(color.FgHiBlack)},
		}

		fmt.Println(color.CyanString("⚡ [CIPHER CLI] Multiplexing Fleet Telemetry & Logs (Ctrl+C to detach)..."))

		var wg sync.WaitGroup
		ctxStop := make(chan os.Signal, 1)
		signal.Notify(ctxStop, os.Interrupt, syscall.SIGTERM)

		for _, svc := range services {
			wg.Add(1)
			go func(s ServiceColor) {
				defer wg.Done()
				dockerCmd := exec.Command("docker", "logs", "-f", "--tail", "10", s.Name)

				stdout, err := dockerCmd.StdoutPipe()
				if err != nil {
					return
				}
				stderr, err := dockerCmd.StderrPipe()
				if err != nil {
					return
				}

				if err := dockerCmd.Start(); err != nil {
					return
				}

				readLines := func(r io.Reader) {
					scanner := bufio.NewScanner(r)
					for scanner.Scan() {
						line := scanner.Text()
						fmt.Printf("%s %s\n", s.Color.Sprint(s.Prefix), line)
					}
				}

				go readLines(stdout)
				go readLines(stderr)

				<-ctxStop
				if dockerCmd.Process != nil {
					_ = dockerCmd.Process.Kill()
				}
			}(svc)
		}

		<-ctxStop
		fmt.Println("\n" + color.YellowString("Detached from Cipher log stream."))
		return nil
	},
}
