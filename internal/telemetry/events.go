package telemetry

// Event names - keeping them consistent and discoverable
const (
	EventCLIInvoked   = "cli_invoked"   // Any CLI command
	EventPackGet      = "pack_get"      // packs get <name>
	EventPackSearch   = "pack_search"   // packs find <query>
	EventPackInfo     = "pack_info"     // packs info <name>
	EventPackSubmit   = "pack_submit"   // packs submit
	EventTUILaunched  = "tui_launched"  // packs (no args)
	EventTUIInstall   = "tui_install"   // Install from TUI
	EventAuthLogin    = "auth_login"    // packs login
	EventAuthLogout   = "auth_logout"   // packs logout
)

// TrackCLI tracks a CLI command invocation.
func TrackCLI(command string) {
	Track(EventCLIInvoked, map[string]any{
		"command": command,
	})
}

// TrackPackGet tracks a pack installation.
func TrackPackGet(packName string, source string, success bool) {
	Track(EventPackGet, map[string]any{
		"pack":    packName,
		"source":  source, // "registry", "github", "url"
		"success": success,
	})
}

// TrackPackSearch tracks a search query.
func TrackPackSearch(queryLength int, resultsCount int) {
	// We don't track the actual query for privacy - just metrics
	Track(EventPackSearch, map[string]any{
		"query_length":  queryLength,
		"results_count": resultsCount,
	})
}

// TrackPackInfo tracks pack info views.
func TrackPackInfo(packName string) {
	Track(EventPackInfo, map[string]any{
		"pack": packName,
	})
}

// TrackPackSubmit tracks pack submissions.
func TrackPackSubmit(success bool) {
	Track(EventPackSubmit, map[string]any{
		"success": success,
	})
}

// TrackTUI tracks TUI launches and interactions.
func TrackTUI(action string) {
	event := EventTUILaunched
	if action == "install" {
		event = EventTUIInstall
	}
	Track(event, map[string]any{
		"action": action,
	})
}

// TrackAuth tracks authentication events.
func TrackAuth(action string, success bool) {
	event := EventAuthLogin
	if action == "logout" {
		event = EventAuthLogout
	}
	Track(event, map[string]any{
		"success": success,
	})
}
