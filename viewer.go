package yonde

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type Viewer struct {
	label     string
	size      *size
	in        string
	doc       *document
	rowOffset int
	curChunk  int
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

func (v *Viewer) chunk() *chunk {
	return v.doc.chunks[v.curChunk]
}

func (v Viewer) Init() tea.Cmd {
	return nil
}

func (v *Viewer) moveUp(lines int) {
	for i := 0; i < lines; i++ {
		if v.rowOffset <= 0 {
			break
		}

		ch := v.doc.chunks[v.curChunk]
		if ch.e == v.rowOffset+v.size.rows-2 {
			v.prevChunk()
		}

		v.rowOffset--
	}
}

func (v *Viewer) moveDown(lines int) {
	for i := 0; i < lines; i++ {
		if v.rowOffset+v.size.rows-1 >= v.doc.len() {
			break
		}

		ch := v.doc.chunks[v.curChunk]
		if ch.s == v.rowOffset {
			v.nextChunk()
		}

		v.rowOffset++
	}
}

func (v *Viewer) beginingOfRows() {
	v.rowOffset = 0

	for v.chunk().e >= v.rowOffset+v.size.rows-1 && v.curChunk > 0 {
		v.prevChunk()
	}
}

func (v *Viewer) endOfRows() {
	v.rowOffset = v.doc.len() - v.size.rows + 1

	for v.chunk().s < v.rowOffset && v.curChunk < len(v.doc.chunks)-1 {
		v.nextChunk()
	}
}

func (v *Viewer) prevChunk() {
	if v.curChunk > 0 {
		v.curChunk--
	}
}

func (v *Viewer) nextChunk() {
	if v.curChunk < len(v.doc.chunks)-1 {
		v.curChunk++
	}
}

func (v Viewer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return v, tea.Quit
		case "up", "k":
			v.moveUp(1)
		case "down", "j":
			v.moveDown(1)
		case "pgup", "b":
			v.moveUp(v.size.rows - 1)
		case "pgdown", " ":
			v.moveDown(v.size.rows - 1)
		case "ctrl+u", "u":
			v.moveUp((v.size.rows - 1) / 2)
		case "ctrl+d", "d":
			v.moveDown((v.size.rows - 1) / 2)
		case "g", "<":
			v.beginingOfRows()
		case "G", ">":
			v.endOfRows()
		case "ctrl+n":
			v.nextChunk()
		case "ctrl+p":
			v.prevChunk()
		case "enter":
			s := v.doc.chunks[v.curChunk].s
			e := v.doc.chunks[v.curChunk].e
			b := stripNewLines(v.doc.lines[s : e+1])
			dispatchCmd(b)
			//fmt.Fprintf(os.Stderr, "%s", string(b))
		}
	case tea.WindowSizeMsg:
		v.size.cols = msg.Width
		v.size.rows = msg.Height
	}

	v.scroll()

	return v, nil
}

func (v *Viewer) scroll() {
	ch := v.doc.chunks[v.curChunk]
	if ch.s < v.rowOffset {
		v.rowOffset = ch.s
	}

	if ch.e >= v.rowOffset+v.size.rows-1 {
		v.rowOffset = ch.e - v.size.rows + 2
	}
}

var styleMarker = lg.NewStyle().Foreground(lg.Color("13"))

func (v Viewer) drawRows() string {
	ch := v.doc.chunks[v.curChunk]
	var builder strings.Builder
	for screenY := 0; screenY < v.size.rows-1; screenY++ {
		docY := screenY + v.rowOffset
		if docY < v.doc.len() {
			fringe := "  "
			if ch.containsAt(docY) {
				fringe = styleMarker.Render("❱ ")
			}

			builder.WriteString(fringe)
			builder.WriteString(v.doc.lines[docY])
		} else {
			builder.WriteString("~")
		}

		builder.WriteString("\n")
	}
	return builder.String()
}

var styleStatus = lg.NewStyle().Foreground(lg.AdaptiveColor{Light: "7", Dark: "0"}).Background(lg.AdaptiveColor{Light: "0", Dark: "7"})

func (v Viewer) drawStatus() string {
	filename := v.in
	lines := v.rowOffset + v.size.rows - 1
	if lines > v.doc.len() {
		lines = v.doc.len()
	}
	totalLines := v.doc.len()
	percent := float64(lines) / float64(totalLines) * 100.0
	s := fmt.Sprintf("%s %d/%d [%.0f%%]", filename, lines, totalLines, percent)
	return styleStatus.Render(s)
}

func (v Viewer) View() string {
	s := v.drawRows()
	s += v.drawStatus()
	return s
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
