package main

// This example demonstrates various Lip Gloss style and layout features.

import (
  "log"
  "time"

  tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/progress"
  "github.com/charmbracelet/lipgloss"
)

var (

	// General.

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		SetString("•").
		Padding(0, 1).
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

  pageStyle = lipgloss.NewStyle().
      Border(lipgloss.RoundedBorder()).
      BorderForeground(lipgloss.Color("#26ff00")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
	// Dialog.

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

type Model struct {
  width    int
  height   int
  loaded   bool
  progress progress.Model
}

type tickMsg time.Time

func (m Model) Init() tea.Cmd {
  return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
  case tea.KeyMsg:
    switch msg.String() {
      case "ctrl+c":
        return m, tea.Quit
    }
  case tickMsg:
		if m.progress.Percent() == 1.0 {
		  m.loaded = true
      return m, nil
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
  }
  return m, nil
}

func (m Model) getWelcomeDialog() string {
    okButton := activeButtonStyle.Render("Hell yeah, brother!")
		cancelButton := buttonStyle.Render("Naw")

		question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Welcom to SQueaL, the TUI database manager")
		buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
		ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

		dialog := lipgloss.Place(m.width, 9,
			lipgloss.Center, lipgloss.Center,
			dialogBoxStyle.Render(ui),
			lipgloss.WithWhitespaceChars("//"),
			lipgloss.WithWhitespaceForeground(subtle),
		) 

  return dialog
}

func (m Model) View() string {
  content := ""
  if !m.loaded {
    content = lipgloss.JoinVertical(lipgloss.Center, m.progress.View())
  } else {
    content = m.getWelcomeDialog()
  }

  return lipgloss.Place(
    m.width,
    m.height,
    lipgloss.Center,
    lipgloss.Center,
    content,
  )
}

func main() {
  m := Model{
		progress: progress.New(progress.WithDefaultGradient()),
    loaded:  false,
	}

  f, err := tea.LogToFile("debug.log", "debug")
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  p := tea.NewProgram(m, tea.WithAltScreen())
  if _, err := p.Run(); err != nil {
    log.Fatal(err)
  }
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
