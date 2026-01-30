package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170"))
)

type viewMode int

const (
	viewList viewMode = iota
	viewSearch
	viewDetail
)

type model struct {
	packs       []pack
	filtered    []pack
	cursor      int
	selected    *pack
	mode        viewMode
	searchInput textinput.Model
	searchQuery string
	filter      string // "all", "skill", "context", "prompt"
	message     string
	quitting    bool
}

type pack struct {
	name        string
	version     string
	stars       int
	description string
	packType    string
	author      string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search packs..."
	ti.CharLimit = 50
	ti.Width = 40

	packs := []pack{
		{name: "commit-message", version: "1.0.0", stars: 892, description: "Generate conventional commit messages", packType: "skill", author: "tunajam"},
		{name: "pr-description", version: "1.0.0", stars: 654, description: "Write PR descriptions from branch diff", packType: "skill", author: "tunajam"},
		{name: "humanizer", version: "1.0.0", stars: 543, description: "Remove AI patterns from writing", packType: "skill", author: "blader"},
		{name: "claudeception", version: "1.0.0", stars: 421, description: "Extract learnings into reusable skills", packType: "skill", author: "blader"},
		{name: "test-driven-development", version: "1.0.0", stars: 389, description: "TDD workflow for features and bugfixes", packType: "skill", author: "obra"},
		{name: "brainstorming", version: "1.0.0", stars: 312, description: "Structured ideation and design exploration", packType: "skill", author: "obra"},
		{name: "react-query", version: "2.1.0", stars: 1247, description: "React Query patterns and best practices", packType: "context", author: "tunajam"},
		{name: "drizzle-orm", version: "1.0.0", stars: 876, description: "Drizzle ORM conventions and patterns", packType: "context", author: "tunajam"},
		{name: "changelog-generator", version: "1.0.0", stars: 287, description: "Generate changelogs from git history", packType: "skill", author: "composio"},
		{name: "git-worktrees", version: "1.0.0", stars: 245, description: "Work with isolated git worktrees", packType: "skill", author: "obra"},
		{name: "mcp-builder", version: "1.0.0", stars: 534, description: "Build MCP servers for LLM integrations", packType: "skill", author: "composio"},
		{name: "youtube-transcript", version: "1.0.0", stars: 423, description: "Fetch and summarize YouTube transcripts", packType: "skill", author: "tapestry"},
	}

	m := model{
		packs:       packs,
		filtered:    packs,
		searchInput: ti,
		filter:      "all",
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode
		if m.mode == viewSearch {
			switch msg.String() {
			case "esc":
				m.mode = viewList
				m.searchInput.Blur()
				return m, nil
			case "enter":
				m.searchQuery = m.searchInput.Value()
				m.applyFilters()
				m.mode = viewList
				m.searchInput.Blur()
				return m, nil
			}
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}

		// Handle detail mode
		if m.mode == viewDetail {
			switch msg.String() {
			case "esc", "q", "backspace":
				m.mode = viewList
				m.selected = nil
				return m, nil
			case "enter", "g":
				// Install the pack
				if m.selected != nil {
					m.message = fmt.Sprintf("Installing %s...", m.selected.name)
					m.quitting = true
					return m, tea.Quit
				}
			}
			return m, nil
		}

		// List mode
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.filtered) > 0 {
				p := m.filtered[m.cursor]
				m.selected = &p
				m.mode = viewDetail
			}

		case "g":
			// Quick get
			if len(m.filtered) > 0 {
				p := m.filtered[m.cursor]
				m.selected = &p
				m.message = fmt.Sprintf("Installing %s...", p.name)
				m.quitting = true
				return m, tea.Quit
			}

		case "/":
			m.mode = viewSearch
			m.searchInput.Focus()
			return m, textinput.Blink

		case "1":
			m.filter = "all"
			m.applyFilters()
		case "2":
			m.filter = "skill"
			m.applyFilters()
		case "3":
			m.filter = "context"
			m.applyFilters()
		case "4":
			m.filter = "prompt"
			m.applyFilters()
		}
	}
	return m, nil
}

func (m *model) applyFilters() {
	var filtered []pack
	query := strings.ToLower(m.searchQuery)

	for _, p := range m.packs {
		// Type filter
		if m.filter != "all" && p.packType != m.filter {
			continue
		}
		// Search filter
		if query != "" {
			if !strings.Contains(strings.ToLower(p.name), query) &&
				!strings.Contains(strings.ToLower(p.description), query) {
				continue
			}
		}
		filtered = append(filtered, p)
	}

	m.filtered = filtered
	m.cursor = 0
}

func (m model) View() string {
	if m.quitting {
		if m.selected != nil {
			return fmt.Sprintf("\n  %s\n  Run: packs get %s\n\n", 
				successStyle.Render("✓"), m.selected.name)
		}
		return ""
	}

	var s strings.Builder

	// Header
	s.WriteString("\n")
	s.WriteString(titleStyle.Render("  🎒 packs"))
	s.WriteString("                                      ")
	s.WriteString(helpStyle.Render("[?] help [q] quit"))
	s.WriteString("\n")
	s.WriteString("  ────────────────────────────────────────────────────\n")

	// Search bar or filter tabs
	if m.mode == viewSearch {
		s.WriteString("  ")
		s.WriteString(m.searchInput.View())
		s.WriteString("\n")
	} else {
		s.WriteString("  ")
		tabs := []string{"[1] All", "[2] Skills", "[3] Contexts", "[4] Prompts"}
		filters := []string{"all", "skill", "context", "prompt"}
		for i, tab := range tabs {
			if filters[i] == m.filter {
				s.WriteString(accentStyle.Render(tab))
			} else {
				s.WriteString(dimStyle.Render(tab))
			}
			s.WriteString("  ")
		}
		if m.searchQuery != "" {
			s.WriteString(dimStyle.Render(fmt.Sprintf(" 🔍 \"%s\"", m.searchQuery)))
		}
		s.WriteString("\n")
	}
	s.WriteString("  ────────────────────────────────────────────────────\n\n")

	// Detail view
	if m.mode == viewDetail && m.selected != nil {
		p := m.selected
		typeIcon := getTypeIcon(p.packType)

		s.WriteString(fmt.Sprintf("  %s %s\n", typeIcon, titleStyle.Render(p.name)))
		s.WriteString(fmt.Sprintf("  %s\n\n", dimStyle.Render(p.version+" · "+p.author)))
		s.WriteString(fmt.Sprintf("  %s\n\n", p.description))
		s.WriteString(fmt.Sprintf("  ★ %d stars\n\n", p.stars))
		s.WriteString(fmt.Sprintf("  %s\n\n", helpStyle.Render("Press ENTER or 'g' to install, ESC to go back")))
		return s.String()
	}

	// List view
	if len(m.filtered) == 0 {
		s.WriteString("  No packs found.\n")
	} else {
		for i, p := range m.filtered {
			cursor := "  "
			style := normalStyle
			if m.cursor == i {
				cursor = "> "
				style = selectedStyle
			}

			typeIcon := getTypeIcon(p.packType)
			stars := fmt.Sprintf("★ %d", p.stars)
			desc := truncateStr(p.description, 35)

			line := fmt.Sprintf("%s%s %-22s %s  %-6s  %s",
				cursor,
				typeIcon,
				p.name,
				p.version,
				stars,
				desc)
			s.WriteString(style.Render(line))
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(helpStyle.Render("  ↑↓ navigate  ⏎ details  g get  / search  1-4 filter  q quit"))
	s.WriteString("\n")

	return s.String()
}

func getTypeIcon(t string) string {
	switch t {
	case "context":
		return "📚"
	case "prompt":
		return "💬"
	default:
		return "📦"
	}
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// If a pack was selected, install it
	if model, ok := m.(model); ok && model.selected != nil {
		fmt.Printf("\n📦 Getting %s...\n\n", model.selected.name)
		err := runGet(model.selected.name, "", false, false)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}
