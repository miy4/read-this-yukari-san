package yonde

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Viewer struct {
	label string
	size  *size
}

type size struct {
	cols int
	rows int
}

func initialMode() Viewer {
	return Viewer{label: "", size: &size{}}
}

func (v Viewer) Init() tea.Cmd {
	return nil
}

func (v Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		}
	case tea.WindowSizeMsg:
		v.size.cols = msg.Width
		v.size.rows = msg.Height
	}

	return v, nil
}

func (v Viewer) drawRows() string {
	var builder strings.Builder
	for y := 0; y < v.size.rows; y++ {
		builder.WriteString("~\n")
	}
	return builder.String()
}

func (v Viewer) View() string {
	return v.drawRows()
}

func Run() {
	p := tea.NewProgram(initialMode(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Alas, there's been an error: %+v", err)
		os.Exit(1)
	}
}
