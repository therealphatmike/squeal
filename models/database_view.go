package models

import (
	"database/sql"
	"strings"

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
	tabs           []string
	tabContent     []string
	activeTab      int
	status         string
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Margin(0).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.RoundedBorder()).UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func NewConnectedDatabaseView(width int, height int, squealDb databases.Database) (ConnectedDatabaseView, error) {
	sqlDb, err := postgres.OpenDb(squealDb)
	if err != nil {
		return ConnectedDatabaseView{}, err
	}

	tabs := []string{"SQL Editor"}
	tabContent := []string{"You can write SQL queries here!"}

	return ConnectedDatabaseView{
		width:          width,
		height:         height,
		database:       sqlDb,
		squealDbConfig: squealDb,
		tabs:           tabs,
		tabContent:     tabContent,
		activeTab:      0,
		status:         "Connected",
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
	}

	return m, tea.Batch(cmds...)
}

func (m ConnectedDatabaseView) View() string {
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A550DF")).
		Bold(true).
		Padding(0).
		Margin(0).
		Align(lipgloss.Center).
		Render(m.squealDbConfig.ConnectionName)

	navContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		// db tables, etc here later
	)

	navWidth := m.width / 5
	availableWidth := m.width - (navWidth + 4)
	tableContent := lipgloss.Place(
		availableWidth,
		m.height-4,
		lipgloss.Top,
		lipgloss.Left,
		lipgloss.NewStyle().MarginBottom(0).Render(m.getTabLine(availableWidth-2)),
	)

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
				lipgloss.NewStyle().
					Margin(0).
					Padding(0).
					Render(tableContent),
			),
			components.NewQuickKeys(m.width),
			components.NewStatusBar(m.width, m.status),
		),
	)
}

func (m ConnectedDatabaseView) getTabLine(availableWidth int) string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(availableWidth).Height(m.height - 6).Render(m.tabContent[m.activeTab]))
	return docStyle.Render(doc.String())
}
