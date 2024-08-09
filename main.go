package main

// This example demonstrates various Lip Gloss style and layout features.

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	models "github.com/therealphatmike/squeal/models"
	"github.com/therealphatmike/squeal/util/bootstrap"
)

func main() {
	err := bootstrap.BootstrapSqueal()
	if err != nil {
		log.Fatal(err)
	}

	m, _ := models.InitSqueal()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
