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
			Padding(0).
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
	formHeight    int
	selectableDbs []databases.Database
	accessor      *huh.PointerAccessor[databases.Database]
}

func NewSelectDatabaseForm(width int, height int, availableDbs []databases.Database) SelectDatabase {
	return newSelectDatabaseForm(width, height, availableDbs, nil)
}

func newSelectDatabaseForm(width int, height int, availableDbs []databases.Database, focused *databases.Database) SelectDatabase {
	dbOptions := []huh.Option[databases.Database]{}
	for _, database := range availableDbs {
		dbOptions = append(dbOptions, huh.NewOption(database.ConnectionName, database))
	}

	var accessor *huh.PointerAccessor[databases.Database]
	if focused != nil {
		accessor = huh.NewPointerAccessor(focused)
	} else {
		zero := new(databases.Database)
		accessor = huh.NewPointerAccessor(zero)
	}
	formHeight := 24

	return SelectDatabase{
		ready:         true,
		width:         width,
		height:        height,
		formHeight:    formHeight,
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
		).WithHeight(formHeight),
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
		// Recreate form with current selection to handle resize
		focused := m.accessor.Get()
		return newSelectDatabaseForm(m.width, m.height, m.selectableDbs, &focused), m.form.Init()
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
	header := lipgloss.
		NewStyle().
		Width(m.width - 10).
		Padding(0, 0).
		Height(1).
		Align(lipgloss.Center).
		Render("😱😱 Welcome to SQueaL 😱😱\nPlease select a database for connection.")

	highlightedDb := m.accessor.Get()

	formWidth := (m.width - 16) / 2
	dbInfoWidth := m.width - 16 - formWidth

	formContent := dialogBoxStyle.Width(formWidth).Height(m.formHeight + 1).MarginRight(0).Render(m.form.View())
	infoContent := dialogBoxStyle.Width(dbInfoWidth).Height(m.formHeight + 1).Padding(0, 1).Render(getCurrentlyHighlightedDatabaseInfo(highlightedDb))

	content := strings.Builder{}
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
				formContent,
				infoContent,
			),
		),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(subtle),
	))

	quickKeys := components.NewQuickKeys(m.width)
	statusText := "Selecting Database..."
	status := components.NewStatusBar(m.width, statusText)
	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Bottom,
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
			0,
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
