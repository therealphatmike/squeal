package models

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealphatmike/squeal/components"
	"github.com/therealphatmike/squeal/postgres"
	"github.com/therealphatmike/squeal/util/databases"
)

type ConnectedDatabaseView struct {
	width          int
	height         int
	squealDbConfig databases.Database
	database       *sql.DB
}

func NewConnectedDatabaseView(width int, height int, squealDb databases.Database) (ConnectedDatabaseView, error) {
	sqlDb, err := postgres.OpenDb(squealDb)
	if err != nil {
		return ConnectedDatabaseView{}, err
	}

	return ConnectedDatabaseView{
		width:          width,
		height:         height,
		database:       sqlDb,
		squealDbConfig: squealDb,
	}, nil
}

func (m ConnectedDatabaseView) Init() tea.Cmd {
	return nil
}

func (m ConnectedDatabaseView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m ConnectedDatabaseView) View() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A550DF")).
		Bold(true).
		Align(lipgloss.Right).
		Render(m.squealDbConfig.ConnectionName)
	navContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		// db tables, etc here later
	)

	tableContent := lipgloss.JoinVertical(
		lipgloss.Right,
		// tabs,
		"", // db area,
	)

	navWidth := m.width / 5
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Top,
		lipgloss.Left,
		lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				dialogBoxStyle.
					Margin(0).
					Padding(0, 1).
					Width(navWidth).
					Height(m.height-4).
					Render(navContent),
				dialogBoxStyle.
					Margin(0).
					Padding(0, 1).
					Width(m.width-(navWidth+4)).
					Height(m.height-4).
					Render(tableContent),
			),
			components.NewQuickKeys(m.width),
			components.NewStatusBar(m.width, "Connected"),
		),
	)
}
