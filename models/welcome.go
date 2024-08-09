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
	newDbFormState NewDatabase
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

	switch m.state {
	case welcomeView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
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
					newDb := NewDatabaseForm(m.width, m.height)
					m.state = newDbForm
					m.newDbFormState = newDb
					cmds = append(cmds, newDb.Init())
					return m, nil
				} else {
					return m, tea.Quit
				}
			}
		}
	case newDbForm:
		_, newCmd := m.newDbFormState.Update(msg)
		cmds = append(cmds, newCmd)
	}

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
			return m.getNoDatabasesScreen(m.width, m.height)
		}

		return "uh oh"
	case newDbForm:
		return m.newDbFormState.View()
	default:
		return m.getNoDatabasesScreen(m.width, m.height)
	}
}

func (m MainModel) getNoDatabasesScreen(width int, height int) string {
	content := strings.Builder{}
	welcome := components.NewWelcomeDialog(width, m.selectedOption)
	quickKeys := components.NewQuickKeys(width)
	status := components.NewStatusBar(width, "No Databases Configured")

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
