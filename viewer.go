package yonde

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Viewer struct {
	label string
	size  *size
	in    string
	doc   *document
}

type size struct {
	cols int
	rows int
}

func NewViewer() (Viewer, error) {
	var r io.Reader
	var filename string

	// TODO read from stdin
	if len(os.Args) > 1 {
		filename = os.Args[1]
		var err error
		r, err = os.Open(filename)
		if err != nil {
			return Viewer{}, nil
		}
	} else {
		return Viewer{}, errors.New("Missing filename")
	}

	document, err := open(r)
	if err != nil {
		return Viewer{}, err
	}

	return Viewer{
		size: &size{},
		in:   filename,
		doc:  document,
	}, nil
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
		if y < len(v.doc.lines) {
			builder.WriteString(v.doc.lines[y])
		} else {
			builder.WriteString("~")
		}

		if y < v.size.rows-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func (v Viewer) View() string {
	return v.drawRows()
}

func Run() {
	viewer, err := NewViewer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(viewer, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
