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
	width          int
	height         int
	databases      []string
	state          viewState
	selectedOption string
}

func InitSqueal() (tea.Model, tea.Cmd) {
	return MainModel{
		databases:      []string{},
		selectedOption: "Yes",
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
		case "left":
			if m.selectedOption == "Yes" {
				m.selectedOption = "No"
			} else {
				m.selectedOption = "Yes"
			}
		case "right":
			if m.selectedOption == "Yes" {
				m.selectedOption = "No"
			} else {
				m.selectedOption = "Yes"
			}
		case "enter":
			if m.selectedOption == "Yes" {
				m.state = newDbForm
				return m, nil
			} else {
				return m, tea.Quit
			}
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case welcomeView:
		if len(m.databases) <= 0 || m.databases == nil {
			return m.getNoDatabasesScreen(m.width, m.height)
		}

		return "uh oh"
	case newDbForm:
		return NewDatabaseForm(m.width, m.height).View()
	default:
		return m.getNoDatabasesScreen(m.width, m.height)
	}
}

func (m MainModel) getNoDatabasesScreen(width int, height int) string {
	content := strings.Builder{}
	welcome := components.NewWelcomeDialog(width, m.selectedOption)
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

// TODO This is a placeholder. Not intended to be the actual form implementation
func getNewDbForm(width int, height int) string {
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 0).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	newDb := lipgloss.NewStyle().Width(75).Height(25).Align(lipgloss.Center).Render("New Database Connection")
	dialog := lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		dialogBoxStyle.Render(newDb),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(subtle),
	)
	content := strings.Builder{}
	quickKeys := components.NewQuickKeys(width)
	status := components.NewStatusBar(width)

	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		dialog,
		quickKeys,
		status,
	))

	return content.String()
}
