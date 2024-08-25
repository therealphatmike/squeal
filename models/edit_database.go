package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealphatmike/squeal/components"
	"github.com/therealphatmike/squeal/util/databases"
)

type EditDatabase struct {
	width  int
	height int
	form   *huh.Form
	lg     *lipgloss.Renderer
}

func EditDatabaseForm(width int, height int, database databases.Database) EditDatabase {
	return EditDatabase{
		lg:     lipgloss.DefaultRenderer(),
		width:  width,
		height: height,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Connection Name").
					Description("We use this in the connection list so you can easily choose between your DBs.").
					Key("connectionName").
					Value(&database.ConnectionName),

				huh.NewSelect[string]().
					Key("engine").
					Options(
						huh.NewOption("PostgreSQL", "postgres"),
						huh.NewOption("MySQL", "mysql"),
						huh.NewOption("MariaDB", "maria"),
					).
					Title("Database Engine").
					Value(&database.Engine),

				huh.NewSelect[string]().
					Key("mode").
					Options(
						huh.NewOption("Host and Port", "hostAndPort"),
					).
					Title("Connection Mode").
					Value(&database.ConnectionMode),

				huh.NewInput().
					Title("Host").
					Key("host").
					Value(&database.Host),

				huh.NewInput().
					Title("Port").
					Key("port").
					Value(&database.Port),

				huh.NewInput().
					Title("User").
					Key("user").
					Value(&database.Username),

				huh.NewInput().
					Title("Password").
					Key("password").
					Value(&database.Password),

				huh.NewInput().
					Title("Default Database").
					Key("defaultDatabase").
					Value(&database.DefaultDatabase),

				huh.NewConfirm().
					Key("submit").
					Affirmative("Create").
					Negative("Cancel"),
			),
		).
			WithShowHelp(true).
			WithShowErrors(true),
	}
}

func (m EditDatabase) Init() tea.Cmd {
	return m.form.Init()
}

func (m EditDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		db := databases.Database{
			ConnectionName:  m.form.GetString("connectionName"),
			Engine:          m.form.GetString("engine"),
			Username:        m.form.GetString("user"),
			Password:        m.form.GetString("password"),
			Host:            m.form.GetString("host"),
			Port:            m.form.GetString("port"),
			DefaultDatabase: m.form.GetString("defaultDatabase"),
		}
		// remove the old db connection, may want to do this differently?
		if err := databases.RemoveDatabaseConnection(db.ConnectionName); err != nil {
			cmds = append(cmds, tea.Quit)
		}

		if err := databases.AddDatabaseConnection(db); err != nil {
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m EditDatabase) View() string {
	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	if m.form.State == huh.StateCompleted {
		name := m.form.GetString("connectionName")
		engine := m.form.GetString("engine")
		username := m.form.GetString("user")
		password := m.form.GetString("password")
		host := m.form.GetString("host")
		port := m.form.GetString("port")
		db := m.form.GetString("defaultDatabase")

		dialogBoxStyle := m.lg.NewStyle().
			MarginBottom(0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

		return m.lg.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			m.lg.NewStyle().
				Margin(2).
				Render(
					lipgloss.JoinVertical(
						lipgloss.Center,
						dialogBoxStyle.Render("Connection String for "+engine+" Database "+name),
						fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, db),
					),
				),
			lipgloss.WithWhitespaceChars("U+1F631"), // IYKYK
			lipgloss.WithWhitespaceForeground(subtle),
		)
	}

	form := m.lg.NewStyle().
		Align(lipgloss.Left).
		Margin(1, 0).
		Width(80).
		Render(m.form.View())

	dialogBoxStyle := m.lg.NewStyle().
		MarginBottom(0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	header := m.lg.
		NewStyle().
		Width(100).
		Height(1).
		Align(lipgloss.Center).
		Render("Edit Database Connection Form")

	content := strings.Builder{}
	quickKeys := components.NewQuickKeys(m.width)
	statusText := "Editing Existing Database Connection"
	status := components.NewStatusBar(m.width, statusText)

	content.WriteString(m.lg.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			dialogBoxStyle.Render(header),
			form,
		),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(subtle),
	))

	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Bottom,
		// why do I need two of these to get the quick keys bar to show up?
		quickKeys,
		quickKeys,
		status,
	))

	return content.String()
}
