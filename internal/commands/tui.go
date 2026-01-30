package commands

import (
	"fmt"
	"os"

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

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type model struct {
	packs    []pack
	cursor   int
	selected string
	quitting bool
}

type pack struct {
	name        string
	version     string
	stars       int
	description string
	packType    string
}

func initialModel() model {
	// Demo packs for now
	return model{
		packs: []pack{
			{name: "commit-message", version: "1.0.0", stars: 892, description: "Conventional commit generator", packType: "skill"},
			{name: "pr-creator", version: "1.0.0", stars: 654, description: "Create PRs with gh CLI", packType: "skill"},
			{name: "pr-description", version: "1.0.0", stars: 543, description: "Generate PR descriptions", packType: "skill"},
			{name: "react-query", version: "2.1.0", stars: 1247, description: "React Query patterns and best practices", packType: "context"},
			{name: "code-review", version: "1.1.0", stars: 421, description: "Review code for issues", packType: "skill"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.packs)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.packs[m.cursor].name
			return m, tea.Quit
		case "g":
			// Get the selected pack
			m.selected = m.packs[m.cursor].name
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		if m.selected != "" {
			return fmt.Sprintf("\n  Run: packs get %s\n\n", m.selected)
		}
		return ""
	}

	s := "\n"
	s += titleStyle.Render("  🎒 packs") + "                                      "
	s += helpStyle.Render("[?] help [q] quit") + "\n"
	s += "  ────────────────────────────────────────────────────\n"
	s += "  " + subtitleStyle.Render("[All]  Skills   Contexts   Prompts") + "\n"
	s += "  ────────────────────────────────────────────────────\n\n"

	for i, p := range m.packs {
		cursor := "  "
		style := normalStyle
		if m.cursor == i {
			cursor = "> "
			style = selectedStyle
		}

		stars := fmt.Sprintf("★ %d", p.stars)
		line := fmt.Sprintf("%s📦 %-20s %s  %-8s %s",
			cursor,
			p.name,
			p.version,
			stars,
			p.description)
		s += style.Render(line) + "\n"
	}

	s += "\n"
	s += helpStyle.Render("  ↑↓ navigate  ⏎ select  g get  / search  q quit") + "\n"

	return s
}

func RunTUI() {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// If a pack was selected with 'g', get it
	if model, ok := m.(model); ok && model.selected != "" && model.quitting {
		// Run the get command
		fmt.Printf("\n📦 Getting %s...\n\n", model.selected)
	}
}
