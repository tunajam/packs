// Package telemetry provides anonymous usage tracking for packs.
// All tracking is opt-out via PACKS_NO_TELEMETRY=1 or config setting.
// No PII is collected - only command usage and pack install counts.
package telemetry

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	posthogEndpoint = "https://us.i.posthog.com/capture"
	posthogAPIKey   = "phc_pdBQMVMnYlUt7vycs418ZDRoS5nTJ5OcaGnuLrrZpGC"
	timeout         = 2 * time.Second
)

var (
	// version is set at build time via ldflags
	version = "dev"
	// distinctID is computed once per session
	distinctID     string
	distinctIDOnce sync.Once
)

// SetVersion sets the CLI version for telemetry events.
func SetVersion(v string) {
	version = v
}

// IsEnabled returns true if telemetry is enabled.
func IsEnabled() bool {
	// Check environment variable first (highest priority)
	if os.Getenv("PACKS_NO_TELEMETRY") == "1" {
		return false
	}
	// TODO: Check config file setting
	return true
}

// Track sends an event to PostHog asynchronously.
// This never blocks the CLI - events are fire-and-forget.
func Track(event string, properties map[string]any) {
	if !IsEnabled() {
		return
	}

	// Fire and forget - don't block the CLI
	go func() {
		_ = send(event, properties)
	}()
}

// TrackSync sends an event and waits for completion.
// Use sparingly - only when you need to ensure delivery before exit.
func TrackSync(event string, properties map[string]any) error {
	if !IsEnabled() {
		return nil
	}
	return send(event, properties)
}

func send(event string, properties map[string]any) error {
	if properties == nil {
		properties = make(map[string]any)
	}

	// Add standard properties
	properties["$lib"] = "packs-cli"
	properties["version"] = version
	properties["os"] = runtime.GOOS
	properties["arch"] = runtime.GOARCH

	payload := map[string]any{
		"api_key":     posthogAPIKey,
		"event":       event,
		"distinct_id": getDistinctID(),
		"properties":  properties,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("POST", posthogEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// getDistinctID returns a stable anonymous identifier for this machine.
// It's a hash of system identifiers - no PII.
func getDistinctID() string {
	distinctIDOnce.Do(func() {
		distinctID = computeDistinctID()
	})
	return distinctID
}

func computeDistinctID() string {
	// Try to get a stable machine ID
	var machineID string

	switch runtime.GOOS {
	case "darwin":
		// macOS: Use hardware UUID
		out, err := exec.Command("ioreg", "-rd1", "-c", "IOPlatformExpertDevice").Output()
		if err == nil {
			// Parse IOPlatformUUID from output
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if strings.Contains(line, "IOPlatformUUID") {
					parts := strings.Split(line, "=")
					if len(parts) == 2 {
						machineID = strings.TrimSpace(strings.Trim(parts[1], `"`))
						break
					}
				}
			}
		}
	case "linux":
		// Linux: Try /etc/machine-id
		data, err := os.ReadFile("/etc/machine-id")
		if err == nil {
			machineID = strings.TrimSpace(string(data))
		}
	}

	// Fallback: Use hostname + home directory
	if machineID == "" {
		hostname, _ := os.Hostname()
		home, _ := os.UserHomeDir()
		machineID = hostname + home
	}

	// Hash it for privacy
	hash := sha256.Sum256([]byte("packs:" + machineID))
	return hex.EncodeToString(hash[:])[:16]
}
