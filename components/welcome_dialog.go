package components

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

func NewWelcomeDialog(width int) string {
	okButton := activeButtonStyle.Render("Yes")
	cancelButton := buttonStyle.Render("No")

	welcome := lipgloss.NewStyle().Width(56).Align(lipgloss.Center).Render("Welcom to SQueaL, the TUI database manager")
	noDatabases := lipgloss.NewStyle().Align(lipgloss.Center).Render("It looks like you don't have any databases configured.")
	setupPrompt := lipgloss.NewStyle().Align(lipgloss.Center).Render("Would you like to set one up?")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, welcome, "\n\n", noDatabases, setupPrompt, buttons)

	dialog := lipgloss.Place(width, 14,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("//"),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	return dialog
}
