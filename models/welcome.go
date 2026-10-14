package models

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/therealphatmike/squeal/components"
	"github.com/therealphatmike/squeal/util/databases"
)

type viewState int

const (
	welcomeView viewState = iota
	newDbForm
	selectDbForm
	queryEditorView
)

var emptySelectDbFormState = SelectDatabase{}

type MainModel struct {
	width             int
	height            int
	databases         []databases.Database
	state             viewState
	selectedOption    string
	newDbFormState    NewDatabase
	selectDbFormState SelectDatabase
}

func InitSqueal() (tea.Model, tea.Cmd) {
	databases, err := databases.ReadDatabaseConfigs()
	if err != nil {
		log.Panic("Error reading databases file. Closing program")
		return nil, tea.Quit
	}

	state := welcomeView

	return MainModel{
		databases:      databases,
		selectedOption: "Yes",
		state:          state,
	}, nil
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.state {
	case welcomeView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
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
				} else {
					return m, tea.Quit
				}
			}
		}
		if len(m.databases) > 0 {
			m.state = selectDbForm
			m.selectDbFormState = NewSelectDatabaseForm(m.width, m.height, m.databases)
			cmds = append(cmds, m.selectDbFormState.Init())
		}
	case newDbForm:
		_, newCmd := m.newDbFormState.Update(msg)
		cmds = append(cmds, newCmd)
	case selectDbForm:
		if m.selectDbFormState.dbSelected {
			m.state = queryEditorView
			cmds = append(cmds, m.selectDbFormState.Init())
		} else if !m.selectDbFormState.ready {
			m.selectDbFormState = NewSelectDatabaseForm(m.width, m.height, m.databases)
			cmds = append(cmds, m.selectDbFormState.Init())
		} else {
			_, newCmd := m.selectDbFormState.Update(msg)
			cmds = append(cmds, newCmd)
		}
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
			newDb := NewDatabaseForm(m.width, m.height)
			m.state = newDbForm
			m.newDbFormState = newDb
			cmds = append(cmds, newDb.Init())
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case welcomeView:
		return m.getNoDatabasesScreen(m.width, m.height)
	case newDbForm:
		return m.newDbFormState.View()
	case selectDbForm:
		if m.selectDbFormState.ready {
			return m.selectDbFormState.View()
		} else {
			return ""
		}
	case queryEditorView:
		return m.getQueryEditorScreen()
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

func (m MainModel) getQueryEditorScreen() string {
	selectedDb := m.selectDbFormState.accessor.Get()
	content := strings.Builder{}
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#874BFD")).
		Width(m.width).
		Align(lipgloss.Center).
		Render("Query Editor — " + selectedDb.ConnectionName)

	content.WriteString(header + "\n\n")
	content.WriteString("Connected to: " + selectedDb.ConnectionName + "\n")
	content.WriteString("Engine: " + selectedDb.Engine + "\n")
	content.WriteString("Host: " + selectedDb.Host + ":" + selectedDb.Port + "\n\n")
	content.WriteString("Query editor coming soon...\n\n")
	content.WriteString("Press Ctrl+N to add a new database, or Ctrl+C to quit.\n")

	quickKeys := components.NewQuickKeys(m.width)
	statusText := "Connected to " + selectedDb.ConnectionName
	status := components.NewStatusBar(m.width, statusText)
	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Bottom,
		quickKeys,
		status,
	))

	return content.String()
}
