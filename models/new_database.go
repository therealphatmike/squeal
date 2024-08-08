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

func NewDatabaseForm(width int, height int) NewDatabase {
	return NewDatabase{
		lg:     lipgloss.DefaultRenderer(),
		width:  width,
		height: height,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Name").
					Key("name"),
				huh.NewSelect[string]().
					Key("engine").
					Options(huh.NewOptions("Postgres")...).
					Title("Database Engine"),
				huh.NewConfirm().
					Key("submit").
					Validate(func(v bool) error {
						if !v {
							return fmt.Errorf("Welp, finish up then")
						}
						return nil
					}).
					Affirmative("Create").
					Negative("Cancel"),
			),
		).
			WithWidth(80).
			WithShowHelp(false).
			WithShowErrors(false),
	}
}

func (m NewDatabase) Init() tea.Cmd {
	return m.form.Init()
}

func (m NewDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = 100
		m.height = 20
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
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m NewDatabase) View() string {
	if m.form.State == huh.StateCompleted {
		name := m.form.GetString("name")
		engine := m.form.GetString("engine")
		return fmt.Sprintf("You selected: %s, Lvl. %s", name, engine)
	}

	form := m.lg.NewStyle().Margin(1, 0).Render(m.form.View())

	subtle := lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	header := lipgloss.
		NewStyle().
		Width(100).
		Height(1).
		Align(lipgloss.Center).
		Render("New Database Connection Form")

	dialog := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		dialogBoxStyle.Render(header),
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	content := strings.Builder{}
	quickKeys := components.NewQuickKeys(m.width)
	status := components.NewStatusBar(m.width)

	content.WriteString(lipgloss.JoinVertical(
		lipgloss.Center,
		dialog,
		form,
		quickKeys,
		status,
	))

	return content.String()
}
