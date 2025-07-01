package internal

import (
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	selected int
	items    []string
}

// View implements tea.Model.
func (m model) View() string {
	panic("unimplemented")
}

func InitialModel() model {
	return model{
		selected: 0,
		items:    []string{"📡 Repo Service", "📦 Gateway", "📣 Notifier", "💽 Redis", "🔧 Restart Service", "❌ Exit"},
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
			return m, tea.Quit
		case "up":
			if m.selected > 0 {
				m.selected--
			}
		case "down":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		case "enter":
			if m.items[m.selected] == "❌ Exit" {
				return m, tea.Quit
			}
			// Optionally call gRPC methods based on selected item
		}
	}
	return m, nil
}
