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
	label     string
	size      *size
	in        string
	doc       *document
	rowOffset int
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

func (v *Viewer) moveUp() {
	if v.rowOffset > 0 {
		v.rowOffset--
	}
}

func (v *Viewer) moveDown() {
	if v.rowOffset+v.size.rows < v.doc.len() {
		v.rowOffset++
	}
}

func (v Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		case "up", "k":
			v.moveUp()
		case "down", "j":
			v.moveDown()
		}
	case tea.WindowSizeMsg:
		v.size.cols = msg.Width
		v.size.rows = msg.Height
	}

	return v, nil
}

func (v Viewer) drawRows() string {
	var builder strings.Builder
	for screenY := 0; screenY < v.size.rows; screenY++ {
		docY := screenY + v.rowOffset
		if docY < v.doc.len() {
			builder.WriteString(v.doc.lines[docY])
		} else {
			builder.WriteString("~")
		}

		if screenY < v.size.rows-1 {
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
