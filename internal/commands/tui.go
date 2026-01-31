package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tunajam/packs/internal/api"
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

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type viewMode int

const (
	viewList viewMode = iota
	viewSearch
	viewDetail
)

const pageSize = 15

type model struct {
	packs       []pack
	filtered    []pack
	cursor      int
	page        int
	selected    *pack
	mode        viewMode
	searchInput textinput.Model
	searchQuery string
	filter      string // "all", "skill", "context", "prompt"
	message     string
	quitting    bool
	loading     bool
	spinner     spinner.Model
	err         error
}

type pack struct {
	name        string
	version     string
	stars       int
	description string
	packType    string
	author      string
}

// Messages for async operations
type packsLoadedMsg struct {
	packs []pack
}

type packsErrorMsg struct {
	err error
}

func fetchPacks(filter string) tea.Cmd {
	return func() tea.Msg {
		client := api.New()
		ctx := context.Background()
		
		packType := ""
		if filter != "all" {
			packType = filter
		}
		
		results, _, err := client.Search(ctx, api.SearchOpts{
			Query: "",
			Type:  packType,
			Limit: 100,
			Sort:  "stars",
		})
		if err != nil {
			return packsErrorMsg{err: err}
		}

		packs := make([]pack, len(results))
		for i, r := range results {
			packs[i] = pack{
				name:        r.Name,
				version:     r.Version,
				stars:       int(r.Stars),
				description: r.Description,
				packType:    r.Type,
				author:      r.Author,
			}
		}
		return packsLoadedMsg{packs: packs}
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search packs..."
	ti.CharLimit = 50
	ti.Width = 40

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	m := model{
		packs:       []pack{},
		filtered:    []pack{},
		searchInput: ti,
		filter:      "all",
		loading:     true,
		spinner:     s,
	}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchPacks("all"), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case packsLoadedMsg:
		m.packs = msg.packs
		m.filtered = msg.packs
		m.loading = false
		m.err = nil
		return m, nil

	case packsErrorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

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
				// Go to previous page if needed
				if m.cursor < m.page*pageSize {
					m.page--
				}
			}

		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				// Go to next page if needed
				if m.cursor >= (m.page+1)*pageSize {
					m.page++
				}
			}

		case "left", "h", "pgup":
			// Previous page
			if m.page > 0 {
				m.page--
				m.cursor = m.page * pageSize
			}

		case "right", "l", "pgdown":
			// Next page
			maxPage := (len(m.filtered) - 1) / pageSize
			if m.page < maxPage {
				m.page++
				m.cursor = m.page * pageSize
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
			m.loading = true
			m.cursor = 0
			m.page = 0
			return m, tea.Batch(fetchPacks("all"), m.spinner.Tick)
		case "2":
			m.filter = "skill"
			m.loading = true
			m.cursor = 0
			m.page = 0
			return m, tea.Batch(fetchPacks("skill"), m.spinner.Tick)
		case "3":
			m.filter = "context"
			m.loading = true
			m.cursor = 0
			m.page = 0
			return m, tea.Batch(fetchPacks("context"), m.spinner.Tick)
		case "4":
			m.filter = "prompt"
			m.loading = true
			m.cursor = 0
			m.page = 0
			return m, tea.Batch(fetchPacks("prompt"), m.spinner.Tick)
		}
	}
	return m, nil
}

func (m *model) applyFilters() {
	var filtered []pack
	query := strings.ToLower(m.searchQuery)

	for _, p := range m.packs {
		// Search filter (type filter is handled by API)
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
	m.page = 0
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

	// Loading state
	if m.loading {
		s.WriteString(fmt.Sprintf("  %s Loading packs...\n", m.spinner.View()))
		return s.String()
	}

	// Error state
	if m.err != nil {
		s.WriteString(fmt.Sprintf("  %s %v\n", errorStyle.Render("Error:"), m.err))
		s.WriteString(helpStyle.Render("  Press q to quit\n"))
		return s.String()
	}

	// Detail view
	if m.mode == viewDetail && m.selected != nil {
		p := m.selected
		typeIcon := getTypeIcon(p.packType)

		s.WriteString(fmt.Sprintf("  %s %s\n", typeIcon, titleStyle.Render(p.name)))
		s.WriteString(fmt.Sprintf("  %s\n\n", dimStyle.Render(p.author)))
		s.WriteString(fmt.Sprintf("  %s\n\n", p.description))
		s.WriteString(fmt.Sprintf("  ★ %d stars\n\n", p.stars))
		s.WriteString(fmt.Sprintf("  %s\n\n", helpStyle.Render("Press ENTER or 'g' to install, ESC to go back")))
		return s.String()
	}

	// List view
	if len(m.filtered) == 0 {
		s.WriteString("  No packs found.\n")
	} else {
		// Calculate page bounds
		start := m.page * pageSize
		end := start + pageSize
		if end > len(m.filtered) {
			end = len(m.filtered)
		}
		totalPages := (len(m.filtered) + pageSize - 1) / pageSize

		for i := start; i < end; i++ {
			p := m.filtered[i]
			cursor := "  "
			style := normalStyle
			if m.cursor == i {
				cursor = "> "
				style = selectedStyle
			}

			typeIcon := getTypeIcon(p.packType)
			stars := fmt.Sprintf("★ %d", p.stars)
			desc := truncateStr(p.description, 35)

			line := fmt.Sprintf("%s%s %-22s  %-6s  %s",
				cursor,
				typeIcon,
				p.name,
				stars,
				desc)
			s.WriteString(style.Render(line))
			s.WriteString("\n")
		}

		// Page indicator
		s.WriteString("\n")
		pageInfo := fmt.Sprintf("  Page %d/%d (%d packs)", m.page+1, totalPages, len(m.filtered))
		s.WriteString(dimStyle.Render(pageInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("  ↑↓ navigate  ←→ page  ⏎ details  g get  / search  1-4 filter  q quit"))
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
