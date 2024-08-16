package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/therealphatmike/squeal/components"
	"github.com/therealphatmike/squeal/util/databases"
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Margin(0, 5).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

type SelectDatabase struct {
	ready         bool
	width         int
	height        int
	form          *huh.Form
	selectableDbs []databases.Database
	accessor      *huh.PointerAccessor[databases.Database]
}

func NewSelectDatabaseForm(width int, height int, availableDbs []databases.Database) SelectDatabase {
	dbOptions := []huh.Option[databases.Database]{}
	for _, database := range availableDbs {
		dbOptions = append(dbOptions, huh.NewOption(database.ConnectionName, database))
	}

	focused := new(databases.Database)
	accessor := huh.NewPointerAccessor(focused)

	return SelectDatabase{
		ready:         true,
		width:         width,
		height:        height,
		selectableDbs: availableDbs,
		accessor:      accessor,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[databases.Database]().
					Title("Which database would you like to connect to?").
					Key("database").
					Accessor(accessor).
					Options(dbOptions...),
			),
		).WithHeight(24),
	}
}

func (m SelectDatabase) Init() tea.Cmd {
	return m.form.Init()
}

func (m SelectDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func (m SelectDatabase) View() string {
	if m.form.State == huh.StateCompleted {
		return ""
	}

	content := strings.Builder{}

	header := lipgloss.
		NewStyle().
		Width(102).
		Padding(0, 0).
		Height(1).
		Align(lipgloss.Center).
		Render("😱😱 Welcome to SQueaL 😱😱\nPlease select a database for connection.")

	highlightedDb := m.accessor.Get()

	content.WriteString(lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			dialogBoxStyle.Render(header),
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				dialogBoxStyle.Width(50).Height(25).MarginRight(0).Render(m.form.View()),
				dialogBoxStyle.Width(50).Height(25).MarginLeft(0).Padding(0).Render(getCurrentlyHighlightedDatabaseInfo(highlightedDb)),
			),
		),
		lipgloss.WithWhitespaceChars("#!"),
		lipgloss.WithWhitespaceForeground(subtle),
	))

	quickKeys := components.NewQuickKeys(m.width)
	statusText := "Selecting Database..."
	status := components.NewStatusBar(m.width, statusText)
	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Bottom,
		// why do I need two of these to get the quick keys bar to show up?
		quickKeys,
		quickKeys,
		status,
	))

	return content.String()
}

func getCurrentlyHighlightedDatabaseInfo(highlighted databases.Database) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#874BFD"))

	infoKeyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#32a852")).
		Align(lipgloss.Left)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.Place(
			50,
			1,
			lipgloss.Center,
			lipgloss.Center,
			headerStyle.Render(highlighted.ConnectionName+" Connection Details\n"),
		),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Engine: "), highlighted.Engine),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Host: "), highlighted.Host),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("ConnectionMode: "), highlighted.ConnectionMode),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Port: "), highlighted.Port),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Username: "), highlighted.Username),
		lipgloss.JoinHorizontal(lipgloss.Left, infoKeyStyle.Render("Default Database: "), highlighted.DefaultDatabase),
	)
}
