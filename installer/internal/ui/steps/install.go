package steps

import (
  "fmt"

  tea "github.com/charmbracelet/bubbletea"
  types "installer.malang/internal/types"
)

type InstallModel struct {
  drives types.PartitionConfigMsg
}

func InitInstallStep(drives types.PartitionConfigMsg) tea.Model {
  return InstallModel{
    drives: drives,
  }
}

func (m InstallModel) Init() tea.Cmd {
  return nil
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  return m, nil
}

func (m InstallModel) View() string {
  return "Installation step placeholder. Drives: " + fmt.Sprintf("%v", m.drives)
}
