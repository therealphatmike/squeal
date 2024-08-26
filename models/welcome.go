package models

import (
	"database/sql"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/therealphatmike/squeal/components"
	"github.com/therealphatmike/squeal/util/databases"
)

type (
	viewState int
	errMsg    struct {
		error
	}
)

const (
	welcomeView viewState = iota
	newDbForm
	selectDbForm
	connectedDbView
)

var emptySelectDbFormState = SelectDatabase{}

type MainModel struct {
	width                int
	height               int
	databases            []databases.Database
	state                viewState
	selectedOption       string
	newDbFormState       NewDatabase
	selectDbFormState    SelectDatabase
	connectedDbViewState ConnectedDatabaseView
	connectedDatabase    databases.Database
	databaseConnection   *sql.DB
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
	}, func() tea.Msg { return errMsg{err} }
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
		if m.newDbFormState.form.State == huh.StateCompleted {
			databases, err := databases.ReadDatabaseConfigs()
			if err != nil {
				log.Panic("Error reading databases file. Closing program")
				return nil, tea.Quit
			}

			m.state = selectDbForm
			m.selectDbFormState = NewSelectDatabaseForm(m.width, m.height, databases)
			return m, tea.Batch(cmds...)
		}
	case selectDbForm:
		if !m.selectDbFormState.ready {
			m.selectDbFormState = NewSelectDatabaseForm(m.width, m.height, m.databases)
			cmds = append(cmds, m.selectDbFormState.Init())
		} else {
			if m.selectDbFormState.form.State == huh.StateCompleted {
				m.state = connectedDbView
				m.connectedDatabase = m.selectDbFormState.accessor.Get()

				// TODO handle this error elegantly. If we can't connect to DB, we should stay on
				// select db view and place an error somewhere visible for the user
				m.connectedDbViewState, _ = NewConnectedDatabaseView(m.width, m.height, m.connectedDatabase)

				// cmds = append(cmds, m.connectedDbViewState.Init())
				return m, tea.Batch(cmds...)
			}

			_, newCmd := m.selectDbFormState.Update(msg)
			cmds = append(cmds, newCmd)
		}
	case connectedDbView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				m.connectedDbViewState.activeTab = min(m.connectedDbViewState.activeTab+1, len(m.connectedDbViewState.tabs)-1)
			case "shift+tab":
				m.connectedDbViewState.activeTab = max(m.connectedDbViewState.activeTab-1, 0)
			case "ctrl+c":
				return m, tea.Quit
			}
		}
		_, newCmd := m.connectedDbViewState.Update(msg)
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
	case connectedDbView:
		return m.connectedDbViewState.View()
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
