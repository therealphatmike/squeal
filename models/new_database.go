package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealphatmike/squeal/components"
)

type NewDatabase struct {
	width  int
	height int
	form   *huh.Form
	lg     *lipgloss.Renderer
}

type DatabaseConnection struct {
	ConnectionName  string `toml:"connectionName"`
	Engine          string `toml:"engine"`
	ConnectionMode  string `toml:"connectionMode"`
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	Username        string `toml:"username"`
	Password        string `toml:"password"`
	DefaultDatabase string `toml:"defaultDatabase"`
}

func NewDatabaseForm(width int, height int) NewDatabase {
	return NewDatabase{
		lg:     lipgloss.DefaultRenderer(),
		width:  width,
		height: height,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Connection Name").
					Description("We use this in the connection list so you can easily choose between your DBs.").
					Key("connectionName"),

				huh.NewSelect[string]().
					Key("engine").
					Options(
						huh.NewOption("PostgreSQL", "postgres"),
						huh.NewOption("MySQL", "mysql"),
						huh.NewOption("MariaDB", "maria"),
					).
					Title("Database Engine"),

				huh.NewSelect[string]().
					Key("mode").
					Options(
						huh.NewOption("Host and Port", "hostAndPort"),
					).
					Title("Connection Mode"),

				huh.NewInput().
					Title("Host").
					Key("host"),

				huh.NewInput().
					Title("Port").
					Key("port"),

				huh.NewInput().
					Title("User").
					Key("user"),

				huh.NewInput().
					Title("Password").
					Key("password"),

				huh.NewInput().
					Title("Default Database").
					Key("defaultDatabase"),

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

func (m NewDatabase) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	return m, tea.Batch(cmds...)
}

func (m NewDatabase) View() string {
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
		Render("New Database Connection Form")

	content := strings.Builder{}
	quickKeys := components.NewQuickKeys(m.width)
	statusText := "Configuring New Connection"
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
