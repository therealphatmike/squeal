package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/therealphatmike/squeal/components"
)

type viewState int

const (
	welcomeView viewState = iota
	newDbForm
)

// MainModel implements tea.Model
type MainModel struct {
	width     int
	height    int
	databases []string
	state     viewState
}

func InitSqueal() (tea.Model, tea.Cmd) {
	return MainModel{
		databases: []string{},
	}, nil
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.state = newDbForm
			return m, nil
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case welcomeView:
		if len(m.databases) <= 0 || m.databases == nil {
			return getNoDatabasesScreen(m.width, m.height)
		}

		return "uh oh"
	case newDbForm:
		return ""
	default:
		return getNoDatabasesScreen(m.width, m.height)
	}
}

func getNoDatabasesScreen(width int, height int) string {
	content := strings.Builder{}
	welcome := components.NewWelcomeDialog(width)
	quickKeys := components.NewQuickKeys(width)
	status := components.NewStatusBar(width)

	content.WriteString(lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		welcome,
	) + "\n\n")

	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		quickKeys,
		status,
	))

	return content.String()
}
