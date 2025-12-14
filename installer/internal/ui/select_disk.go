package ui

import (
  "fmt"

  "github.com/charmbracelet/bubbles/table"
  "github.com/charmbracelet/lipgloss"

  services "installer.malang/internal/services"
  tea "github.com/charmbracelet/bubbletea"
  types "installer.malang/internal/types"
)

type SelectDiskModel struct {
  msg string
  table table.Model
  disks []types.Disk
}

func createDiskTable(disks []types.Disk) table.Model {
  columns := []table.Column{
    {Title: "Name", Width: 15},
    {Title: "Size", Width: 10},
    {Title: "Type", Width: 10},
  }

  rows := []table.Row{}
  for _, disk := range disks {
    rows = append(rows, table.Row{
      disk.Name, 
      disk.SizeInGb, 
      disk.Type,
    })
  }

  t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(7),
  )

  s := table.DefaultStyles()
  s.Header = s.Header.
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("240")).
    BorderBottom(true).
    Bold(false)
  s.Selected = s.Selected.
    Foreground(lipgloss.Color("229")).
    Background(lipgloss.Color("57")).
    Bold(false)
  t.SetStyles(s)

  return t
}

func InitDiskStep() tea.Model {
  disks := services.CollectDisks()

  return SelectDiskModel{
    msg: "Select Disk",
    disks: disks,
    table: createDiskTable(disks),
  }
}

func (m SelectDiskModel) Init() tea.Cmd {
  return nil
}

func findDiskByName(disks []types.Disk, name string) types.Disk {
  for _, disk := range disks {
    if disk.Name == name {
      return disk
    }
  }
  panic("Disk not found: " + name)
}

func (m SelectDiskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "q", "ctrl+c":
      return m, tea.Quit
    case "enter":
      fmt.Println("Selected disk:", m.table.SelectedRow())
      return m, func() tea.Msg {
        selectedDiskName := m.table.SelectedRow()[0]
        selectedDisk := findDiskByName(m.disks, selectedDiskName)
        return types.SelectedDiskMsg(selectedDisk)
      }
    }
  }

  m.table, cmd = m.table.Update(msg)
  return m, cmd
}

func (m SelectDiskModel) View() string {
  return fmt.Sprintf("%s\n\n%s\n\n%s",
    m.msg,
    m.table.View(),
    "(↑/↓: navigate • enter: select • q: quit)",
  )
}
