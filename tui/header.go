package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	Width     int
	input     textinput.Model
	isLoading bool
	spinner   spinner.Model
}

func NewHeader() Header {
	ti := textinput.NewModel()
	ti.Prompt = ""
	ti.Placeholder = ""
	ti.PlaceholderStyle = styles.Faint.Copy()
	spinner := spinner.New()
	spinner.Style = styles.Regular.Copy().Padding(0, 1)
	return Header{
		input:   ti,
		spinner: spinner,
	}
}

func (h Header) Init() tea.Cmd {
	if h.isLoading {
		return h.spinner.Tick
	}
	return nil
}

func (h Header) Value() string {
	return h.input.Value()
}

func (h Header) Update(msg tea.Msg) (Header, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	h.input, cmd = h.input.Update(msg)
	cmds = append(cmds, cmd)

	if h.isLoading {
		h.spinner, cmd = h.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return h, tea.Batch(cmds...)
}

func (h *Header) SetIsLoading(isLoading bool) tea.Cmd {
	h.isLoading = isLoading
	if isLoading {
		return h.spinner.Tick
	}
	return nil
}

func (h *Header) Focus() tea.Cmd {
	h.input.Placeholder = "Search..."
	return h.input.Focus()
}

func (c Header) View() string {
	var headerRow string
	if c.isLoading {
		spinner := c.spinner.View()
		textInput := styles.Regular.Copy().Width(c.Width - lipgloss.Width(spinner)).Render(c.input.View())
		headerRow = lipgloss.JoinHorizontal(lipgloss.Top, c.spinner.View(), textInput)
	} else {
		headerRow = styles.Regular.Copy().PaddingLeft(3).Width(c.Width).Render(c.input.View())
	}

	line := strings.Repeat("─", c.Width)
	line = styles.Bold.Render(line)
	return lipgloss.JoinVertical(lipgloss.Left, headerRow, line)
}
