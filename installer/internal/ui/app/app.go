package app

import (
  "github.com/charmbracelet/bubbles/textarea"

  tea "github.com/charmbracelet/bubbletea"
)

type model struct {
  textarea textarea.Model
}

func App() {
  tea.NewProgram(initApp()).Run()
}

func initApp() tea.Model {
  ta := textarea.New()
  ta.SetWidth(50)
  ta.SetHeight(10)
  ta.Placeholder = "Type here..."
  ta.Focus()

  return model {
    textarea: ta,
  }
}

func (m model) Init() tea.Cmd {
  return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  m.textarea, cmd = m.textarea.Update(msg)
  return m, cmd
}

func (m model) View() string {
  return "Welcome to Installer" + m.textarea.View()
}

