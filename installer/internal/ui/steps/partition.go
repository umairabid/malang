package steps

import (
  tea "github.com/charmbracelet/bubbletea"
  types "installer.malang/internal/types"
)

type PartitionModel struct {
  disk types.Disk
}

func InitPartitionStep(selectedDisk types.Disk) tea.Model {
  return PartitionModel{
    disk: selectedDisk,
  }
}

func (m PartitionModel) Init() tea.Cmd {
  return nil
}

func (m PartitionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  return m, nil
}

func (m PartitionModel) View() string {
  return "Partitioning step for disk: " + m.disk.Name + " (" + m.disk.SizeInGb + ")"
}
