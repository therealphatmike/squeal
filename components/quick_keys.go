package components

import (
    "github.com/charmbracelet/lipgloss"
)

var (
    quickKeysStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
        Background(lipgloss.Color("#353533"))
    
    keyStyle = lipgloss.NewStyle().
        Padding(0, 0).
        MarginRight(1).
        Inherit(quickKeysStyle).
        Foreground(lipgloss.Color("#A550DF"))

    labelStyle = lipgloss.NewStyle().
        Padding(0, 1).
        MarginRight(2).
        Inherit(quickKeysStyle).
        Foreground(lipgloss.Color("#FFFDF5"))
)

func NewQuickKeys(width int) string {
    quitKey := keyStyle.Render("^c")
    quitLabel := labelStyle.Render("Quit")
    newKey := keyStyle.Render("^n")
    newLabel := labelStyle.Render("New")
    disconnectKey := keyStyle.Render("^d")
    disconnectLabel := labelStyle.Render("Disconnect")

    bar := lipgloss.JoinHorizontal(
        lipgloss.Top,
        quitKey,
        quitLabel,
        newKey,
        newLabel,
        disconnectKey,
        disconnectLabel,
    )

    return quickKeysStyle.Width(width).Render(bar)
}
