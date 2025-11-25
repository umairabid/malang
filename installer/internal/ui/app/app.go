package app

import (
  "fmt"
  "github.com/charmbracelet/lipgloss"

  tea "github.com/charmbracelet/bubbletea"
  steps  "installer.malang/internal/ui/steps"
  types  "installer.malang/internal/types"
)

type model struct {
  currentStep tea.Model
  selectedDisk types.Disk
  drives      types.PartitionConfigMsg
}

func App() {
  tea.NewProgram(initApp(), tea.WithAltScreen()).Run()
}

func initApp() tea.Model {
  model := model{}
  model.currentStep = steps.InitDiskStep()
  return model
}

func (m model) Init() tea.Cmd {
  return tea.SetWindowTitle("Malang Installer")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  m.currentStep, cmd = m.currentStep.Update(msg)
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "q", "esc", "ctrl+c":
      return m, tea.Quit
    }
  case types.SelectedDiskMsg:
    fmt.Println("Selected disk for installation:", msg)
    m.selectedDisk = types.Disk(msg)
    m.currentStep = steps.InitPartitionStep(m.selectedDisk)
    return m, nil
  case types.PartitionConfigMsg:
    fmt.Println("Partition configuration received:", msg)
    m.drives = msg
    return m, nil
  }
  return m, cmd
}

func (m model) View() string {
  return headerStyle().Render("Welcome to Malang Installer") + "\n\n" + bodyStyle().Render(m.currentStep.View()) 
}

func headerStyle() lipgloss.Style {
  return lipgloss.NewStyle().Bold(true).
  Foreground(lipgloss.Color("#FAFAFA")).
  Background(lipgloss.Color("#7D56F4")).
  Padding(1, 2)
}

func bodyStyle() lipgloss.Style {
  return lipgloss.NewStyle().
  Padding(1, 2)
}
