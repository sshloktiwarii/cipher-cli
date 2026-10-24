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
2. Injects a runaway loop to verify cipher-exec Seatbelt process killing (5s timeout).
3. Polls cipher-mem (:8083) until an INC-*-resolution.md incident file is resolved.`,
	Run: func(cmd *cobra.Command, args []string) {
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println("\n" + cyan("🧪 [CIPHER TEST] Initiating Synthetic End-to-End Fleet Integration Harness..."))
		fmt.Println(bold("==================================================================="))

		client := &http.Client{Timeout: 15 * time.Second}
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

		// Phase 2: Test cipher-exec Sandbox Process Kill (5-second timeout trap)
		fmt.Println(bold("\n[GATE 2] Testing Execution Sandbox Process Ceiling (cipher-exec :8081)..."))
		execPayload, _ := json.Marshal(map[string]interface{}{
			"language":        "python",
			"timeout_seconds": 3,
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
				fmt.Printf("   %s Gate 2 PASSED: Runaway compute killed via isolation (%s) in %v (ceiling <= 5s)\n", green("✓"), cid, time.Since(start2))
			} else {
				fmt.Printf("   %s Gate 2 FAILED: Expected timeout and exit code -1, got: %s\n", red("✗"), string(body))
				allPassed = false
			}
		}

		// Phase 3: Poll cipher-mem until INC-*-resolution.md appears
		fmt.Println(bold("\n[GATE 3] Polling Neural Memory & Obsidian Vault (cipher-mem :8083)..."))
		foundIncident := false
		var matchedFile string
		var matchedSim float64
		start3 := time.Now()

		for attempt := 1; attempt <= 10; attempt++ {
			searchPayload, _ := json.Marshal(map[string]interface{}{
				"query": "Execution Timeout Runaway Loop",
				"top_k": 3,
			})

			memResp, err := client.Post("http://localhost:8083/api/search", "application/json", bytes.NewReader(searchPayload))
			if err == nil {
				body, _ := io.ReadAll(memResp.Body)
				memResp.Body.Close()
				var res map[string]interface{}
				_ = json.Unmarshal(body, &res)

				results, _ := res["results"].([]interface{})
				for _, r := range results {
					item, ok := r.(map[string]interface{})
					if !ok {
						continue
					}
					chunk, ok := item["chunk"].(map[string]interface{})
					if !ok {
						continue
					}
					file, _ := chunk["file"].(string)
					sim, _ := item["similarity"].(float64)

					if strings.Contains(file, "INC-") && strings.Contains(file, "-resolution.md") {
						foundIncident = true
						matchedFile = file
						matchedSim = sim
						break
					}
				}
			}

			if foundIncident {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}

		if foundIncident {
			fmt.Printf("   %s Gate 3 PASSED: Verified incident post-mortem [[%s]] resolved (sim: %.3f) in %v\n",
				green("✓"), matchedFile, matchedSim, time.Since(start3))
		} else {
			fmt.Printf("   %s Gate 3 FAILED: Incident resolution file INC-*-resolution.md was not located\n", red("✗"))
			allPassed = false
		}

		// Phase 4: Test cipher-log Cryptographic Merkle Tree Hash Chain
		fmt.Println(bold("\n[GATE 4] Verifying Cryptographic Merkle Hash Chain (cipher-log :8082 / :8080)..."))
		start4 := time.Now()
		auditURL := "http://localhost:8080/api/log/verify-audit-chain"
		logResp, err := client.Get(auditURL)
		if err != nil {
			// Fallback to direct port
			auditURL = "http://localhost:8082/verify-audit-chain"
			logResp, err = client.Get(auditURL)
		}
		if err != nil {
			fmt.Printf("   %s Connection failed to cipher-log audit chain: %v\n", red("✗"), err)
			allPassed = false
		} else {
			defer logResp.Body.Close()
			body, _ := io.ReadAll(logResp.Body)
			var res map[string]interface{}
			_ = json.Unmarshal(body, &res)
			status, _ := res["status"].(string)
			proof, _ := res["cryptographic_proof"].(string)
			tampered, _ := res["tampered"].(bool)

			if status == "VERIFIED_INTEGRITY" && !tampered {
				fmt.Printf("   %s Gate 4 PASSED: Merkle chain verified untampered (%s) in %v\n", green("✓"), proof, time.Since(start4))
			} else {
				fmt.Printf("   %s Gate 4 FAILED: Cryptographic integrity failure! Body: %s\n", red("✗"), string(body))
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
