package steps

import (
  "fmt"

  tea "github.com/charmbracelet/bubbletea"
  services "installer.malang/internal/services"
)

type Disks []services.Disk

type model struct {
  msg string
  disks Disks
}

func collectDisks() tea.Msg {
  fmt.Println("Collecting disks...")
  disks := services.CollectDisks()
  return Disks(disks)
}

func InitDiskStep() tea.Model {
  fmt.Println("Initializing disk step model...")
  return model {
    msg: "Select Disk",
    disks: collectDisks().(Disks),
  }
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  fmt.Println("Updating disk step...")
  return m, cmd
}

func (m model) View() string {
  fmt.Println(m.disks)
  return m.msg 
}
