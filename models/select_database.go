package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/therealphatmike/squeal/util/databases"
)

type SelectDatabase struct {
	width     int
	height    int
	form      *huh.Form
	databases []databases.Database
}

func NewSelectDatabaseForm(width int, height int) SelectDatabase {
	return SelectDatabase{
		width:  width,
		height: height,
		form:   huh.NewForm(),
	}
}

func (m SelectDatabase) Init() tea.Cmd {
	return m.form.Init()
}

func (m SelectDatabase) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m SelectDatabase) View() string {
	return "We're really doing it"
}
