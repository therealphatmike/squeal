package main

// This example demonstrates various Lip Gloss style and layout features.

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	models "github.com/therealphatmike/squeal/models"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m, _ := models.InitSqueal()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
