package yonde

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Viewer struct {
	label string
}

func initialMode() Viewer {
	return Viewer{""}
}

func (m Viewer) Init() tea.Cmd {
	return nil
}

func (m Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			m.label = msg.String()
		}
	}

	return m, nil
}

func (m Viewer) View() string {
	s := fmt.Sprintf("%s", m.label)
	return s
}

func Run() {
	p := tea.NewProgram(initialMode(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %+v", err)
		os.Exit(1)
	}
}
