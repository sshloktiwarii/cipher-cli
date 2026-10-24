package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	flagSynthetic bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Executes synthetic end-to-end integration tests across the Cipher fleet",
	Long: `The test command verifies the full platform security and autonomous healing loop:
1. Injects mock credentials to verify cipher-sec WAF redaction.
2. Injects a runaway loop to verify cipher-exec Seatbelt process killing.
3. Asserts telemetry ingestion in cipher-log.
4. Asserts neural backlink resolution in cipher-mem and Obsidian vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println("\n" + cyan("🧪 [CIPHER TEST] Initiating Synthetic End-to-End Fleet Integration Harness..."))
		fmt.Println(bold("==================================================================="))

		client := &http.Client{Timeout: 10 * time.Second}
		allPassed := true

		// Phase 1: Test cipher-sec WAF Redaction
		fmt.Println(bold("\n[GATE 1] Testing Semantic Firewall (cipher-sec :8079)..."))
		mockKey := "AKIAIOSFODNN7EXAMPLE"
		secPayload, _ := json.Marshal(map[string]string{
			"prompt": fmt.Sprintf("Deploy AWS instance with key %s and execute task", mockKey),
		})

		start1 := time.Now()
		secResp, err := client.Post("http://localhost:8079/inspect", "application/json", bytes.NewReader(secPayload))
		if err != nil {
			fmt.Printf("   %s Connection failed to cipher-sec: %v\n", red("✗"), err)
			fmt.Printf("      %s Ensure cipher-sec is running on :8079\n", yellow("Fix:"))
			allPassed = false
		} else {
			defer secResp.Body.Close()
			body, _ := io.ReadAll(secResp.Body)
			var res map[string]interface{}
			_ = json.Unmarshal(body, &res)

			redacted, _ := res["redacted"].(bool)
			sanitized, _ := res["sanitized"].(string)

			if redacted && strings.Contains(sanitized, "[REDACTED_BY_CIPHER_SEC]") && !strings.Contains(sanitized, mockKey) {
				fmt.Printf("   %s Gate 1 PASSED: Credential %s intercepted & redacted in %v\n", green("✓"), mockKey, time.Since(start1))
			} else {
				fmt.Printf("   %s Gate 1 FAILED: Credential was not properly redacted! Body: %s\n", red("✗"), string(body))
				allPassed = false
			}
		}

		// Phase 2: Test cipher-exec Sandbox Process Kill
		fmt.Println(bold("\n[GATE 2] Testing Execution Sandbox Process Ceiling (cipher-exec :8081)..."))
		execPayload, _ := json.Marshal(map[string]interface{}{
			"language":        "python",
			"timeout_seconds": 2,
			"code":            "import time\nwhile True:\n    pass # Runaway loop",
		})

		start2 := time.Now()
		execResp, err := client.Post("http://localhost:8081/execute", "application/json", bytes.NewReader(execPayload))
		if err != nil {
			fmt.Printf("   %s Connection failed to cipher-exec: %v\n", red("✗"), err)
			allPassed = false
		} else {
			defer execResp.Body.Close()
			body, _ := io.ReadAll(execResp.Body)
			var res map[string]interface{}
			_ = json.Unmarshal(body, &res)

			status, _ := res["status"].(string)
			exitCode, _ := res["exit_code"].(float64)
			cid, _ := res["container_id"].(string)

			if status == "timeout" && exitCode == -1 {
				fmt.Printf("   %s Gate 2 PASSED: Runaway compute killed via isolation (%s) in %v\n", green("✓"), cid, time.Since(start2))
			} else {
				fmt.Printf("   %s Gate 2 FAILED: Expected timeout and exit code -1, got: %s\n", red("✗"), string(body))
				allPassed = false
			}
		}

		// Phase 3: Test cipher-log Telemetry Ingestion
		fmt.Println(bold("\n[GATE 3] Testing Telemetry Engine Ingestion (cipher-log :8082)..."))
		start3 := time.Now()
		logResp, err := client.Get("http://localhost:8082/healthz")
		if err != nil {
			fmt.Printf("   %s Connection failed to cipher-log: %v\n", red("✗"), err)
			allPassed = false
		} else {
			defer logResp.Body.Close()
			body, _ := io.ReadAll(logResp.Body)
			var res map[string]interface{}
			_ = json.Unmarshal(body, &res)

			status, _ := res["status"].(string)
			storage, _ := res["storage"].(string)
			totalIngested, _ := res["total_ingested"].(float64)

			if status == "healthy" && totalIngested >= 1 {
				fmt.Printf("   %s Gate 3 PASSED: Telemetry active [%s] with %d events persisted in %v\n",
					green("✓"), storage, int(totalIngested), time.Since(start3))
			} else {
				fmt.Printf("   %s Gate 3 FAILED: Telemetry unhealthy or zero events ingested: %s\n", red("✗"), string(body))
				allPassed = false
			}
		}

		// Phase 4: Test cipher-mem Obsidian Resolution Backlink
		fmt.Println(bold("\n[GATE 4] Testing Neural Memory Search & Backlinks (cipher-mem :8083)..."))
		searchPayload, _ := json.Marshal(map[string]interface{}{
			"query": "Crucible Execution Timeout",
			"top_k": 1,
		})

		start4 := time.Now()
		memResp, err := client.Post("http://localhost:8083/api/search", "application/json", bytes.NewReader(searchPayload))
		if err != nil {
			fmt.Printf("   %s Connection failed to cipher-mem: %v\n", red("✗"), err)
			allPassed = false
		} else {
			defer memResp.Body.Close()
			body, _ := io.ReadAll(memResp.Body)
			var res map[string]interface{}
			_ = json.Unmarshal(body, &res)

			results, _ := res["results"].([]interface{})
			if len(results) > 0 {
				first := results[0].(map[string]interface{})
				chunk := first["chunk"].(map[string]interface{})
				file := chunk["file"].(string)
				sim := first["similarity"].(float64)
				fmt.Printf("   %s Gate 4 PASSED: Retrieved incident backlink [[%s]] (similarity: %.3f) in %v\n",
					green("✓"), file, sim, time.Since(start4))
			} else {
				fmt.Printf("   %s Gate 4 FAILED: Zero resolution chunks found in Obsidian vector store\n", red("✗"))
				allPassed = false
			}
		}

		fmt.Println(bold("\n==================================================================="))
		if allPassed {
			fmt.Println(green("🎉 SYNTHETIC FLEET INTEGRATION TEST PASSED (All 4 verification gates satisfied)\n"))
			os.Exit(0)
		} else {
			fmt.Println(red("❌ SYNTHETIC FLEET INTEGRATION TEST FAILED — See failed gate details above.\n"))
			os.Exit(1)
		}
	},
}

func init() {
	testCmd.Flags().BoolVarP(&flagSynthetic, "synthetic", "s", true, "Execute full synthetic end-to-end integration harness")
}
