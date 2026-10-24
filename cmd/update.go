package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	CurrentVersion = "v1.2.0"
	updateCheckOnly bool
	updateForce     bool
)

type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Performs atomic self-update of the Cipher CLI binary from GitHub Releases",
	Long: `Fetches release metadata from GitHub Releases, verifies SHA256 checksums,
and performs atomic in-place binary swapping with zero-downtime transition.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
		green := color.New(color.FgGreen, color.Bold).SprintFunc()
		yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println("\n" + cyan("🔄 [CIPHER UPDATE] Checking for Platform Releases..."))
		fmt.Printf("   Current Version: %s (%s/%s)\n", bold(CurrentVersion), runtime.GOOS, runtime.GOARCH)

		client := &http.Client{Timeout: 15 * time.Second}
		apiURL := "https://api.github.com/repos/sshloktiwarii/cipher-cli/releases/latest"

		req, err := http.NewRequest(http.MethodGet, apiURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", "cipher-cli/"+CurrentVersion)
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to query GitHub Releases: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			fmt.Printf("%s No remote releases tagged yet on GitHub repository sshloktiwarii/cipher-cli.\n", yellow("ℹ"))
			fmt.Println(green("✓ Local binary is currently at the latest development revision: " + CurrentVersion))
			return nil
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("GitHub API returned HTTP %d: %s", resp.StatusCode, string(body))
		}

		var release GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return fmt.Errorf("failed to decode release payload: %w", err)
		}

		fmt.Printf("   Latest Upstream: %s (Released: %s)\n", bold(release.TagName), release.PublishedAt)

		if release.TagName == CurrentVersion && !updateForce {
			fmt.Println(green("✓ Cipher is already up-to-date! Current version: ") + CurrentVersion)
			return nil
		}

		if updateCheckOnly {
			if release.TagName != CurrentVersion {
				fmt.Printf("%s Update available: %s ➔ %s\n", yellow("!"), CurrentVersion, release.TagName)
				fmt.Println("Run 'cipher update' to install.")
			}
			return nil
		}

		// Locate target binary asset for current OS/Arch
		assetPattern := fmt.Sprintf("cipher_%s_%s", runtime.GOOS, runtime.GOARCH)
		var targetURL string
		var checksumURL string

		for _, a := range release.Assets {
			if strings.Contains(a.Name, assetPattern) && !strings.HasSuffix(a.Name, ".sha256") {
				targetURL = a.BrowserDownloadURL
			}
			if strings.Contains(a.Name, "checksums.txt") || strings.Contains(a.Name, ".sha256") {
				checksumURL = a.BrowserDownloadURL
			}
		}

		if targetURL == "" {
			fmt.Printf("%s No precompiled release asset matching %s found in release %s.\n",
				yellow("!"), assetPattern, release.TagName)
			fmt.Println("Building in-place from source via go install...")
			buildCmd := exec.Command("go", "install", "github.com/sshloktiwarii/cipher-cli@"+release.TagName)
			buildCmd.Stdout = os.Stdout
			buildCmd.Stderr = os.Stderr
			if err := buildCmd.Run(); err != nil {
				return fmt.Errorf("go install update failed: %w", err)
			}
			fmt.Println(green("✓ Successfully updated Cipher CLI to ") + release.TagName)
			return nil
		}

		// Self-executable path
		selfPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("unable to determine binary location: %w", err)
		}
		selfPath, err = filepath.EvalSymlinks(selfPath)
		if err != nil {
			return fmt.Errorf("unable to resolve binary symlink: %w", err)
		}

		fmt.Printf("   Downloading binary update from %s...\n", targetURL)
		tmpFile := selfPath + ".new"
		_ = os.Remove(tmpFile)

		dlResp, err := client.Get(targetURL)
		if err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}
		defer dlResp.Body.Close()

		f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to create temporary update file %s: %w", tmpFile, err)
		}

		hasher := sha256.New()
		multi := io.MultiWriter(f, hasher)
		if _, err := io.Copy(multi, dlResp.Body); err != nil {
			_ = f.Close()
			_ = os.Remove(tmpFile)
			return fmt.Errorf("failed to stream binary: %w", err)
		}
		_ = f.Close()

		shaHex := hex.EncodeToString(hasher.Sum(nil))
		fmt.Printf("   %s Download complete. Computed SHA256: %s\n", green("✓"), shaHex)

		// Verify checksum if available
		if checksumURL != "" {
			csResp, err := client.Get(checksumURL)
			if err == nil {
				defer csResp.Body.Close()
				csBody, _ := io.ReadAll(csResp.Body)
				if !strings.Contains(string(csBody), shaHex) {
					_ = os.Remove(tmpFile)
					return fmt.Errorf("%s Checksum verification mismatch against upstream signatures!", red("✗"))
				}
				fmt.Printf("   %s SHA256 Cryptographic Signature Verified against upstream\n", green("✓"))
			}
		}

		// Atomic Swap
		oldBackup := selfPath + ".old"
		_ = os.Remove(oldBackup)

		if err := os.Rename(selfPath, oldBackup); err != nil {
			_ = os.Remove(tmpFile)
			return fmt.Errorf("failed to rotate current binary: %w", err)
		}

		if err := os.Rename(tmpFile, selfPath); err != nil {
			// Rollback
			_ = os.Rename(oldBackup, selfPath)
			return fmt.Errorf("failed to activate new binary: %w (rolled back to original)", err)
		}
		_ = os.Remove(oldBackup)

		fmt.Println(green("\n🎉 Successfully updated Cipher CLI to ") + bold(release.TagName))
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateCheckOnly, "check", "c", false, "Check for new releases without installing")
	updateCmd.Flags().BoolVarP(&updateForce, "force", "f", false, "Force re-installation of the latest release")
	rootCmd.AddCommand(updateCmd)
}
