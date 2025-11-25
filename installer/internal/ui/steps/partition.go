package steps

import (
  "fmt"
  "strconv"
  "github.com/charmbracelet/lipgloss"

  tea "github.com/charmbracelet/bubbletea"
  types "installer.malang/internal/types"
  services "installer.malang/internal/services"
)

type PartitionModel struct {
  disk         types.Disk
  percentages  [3]string
  focusedField int
  errorMsg     string
}

func InitPartitionStep(selectedDisk types.Disk) tea.Model {
  return PartitionModel{
    disk:         selectedDisk,
    percentages:  [3]string{"20", "10", "70"},
    focusedField: 0,
    errorMsg:     "",
  }
}

func (m PartitionModel) Init() tea.Cmd {
  return nil
}

func (m PartitionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c", "q":
      return m, tea.Quit

    case "tab", "down":
      m.focusedField = (m.focusedField + 1) % 3
      m.errorMsg = ""

    case "shift+tab", "up":
      m.focusedField = (m.focusedField - 1 + 3) % 3
      m.errorMsg = ""

    case "enter":
      values, err := validatePercentages(m.percentages)
      if err != nil {
        m.errorMsg = err.Error()
        return m, nil
      }
      
      driveNames := services.PartitionDisk(m.disk, values)
      return m, func() tea.Msg {
        return types.PartitionConfigMsg(driveNames)
      }

    case "backspace":
      if len(m.percentages[m.focusedField]) > 0 {
        m.percentages[m.focusedField] = m.percentages[m.focusedField][:len(m.percentages[m.focusedField])-1]
      }
      m.errorMsg = ""

    default:
      if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
        if len(m.percentages[m.focusedField]) < 3 {
          m.percentages[m.focusedField] += msg.String()
        }
        m.errorMsg = ""
      }
    }
  }

  return m, nil
}

func (m PartitionModel) View() string {
  title := fmt.Sprintf("Configure Partitions for %s (%s)\n\n", m.disk.Name, m.disk.SizeInGb)

  labels := [3]string{"Boot (EFI): ", "Swap:       ", "Root:       "}

  var lines string
  for i, label := range labels {
    value := m.percentages[i] + "%"
    lines += label + applyStyle(value, m.focusedField == i) + "\n"
  }

  totalLabel := totalLabel(m)
  help := "\n\n(↑/↓/tab: navigate • 0-9: enter value • enter: confirm • q: quit)"

  errorMsg := ""
  if m.errorMsg != "" {
    errorMsg = "\n" + errorStyle().Render("Error: "+m.errorMsg)
  }

  return title + lines + totalLabel + errorMsg + help
}

func focusedStyle() lipgloss.Style {
  return lipgloss.NewStyle().
  Foreground(lipgloss.Color("229")).
  Background(lipgloss.Color("57")).
  Bold(true)
}

func normalStyle() lipgloss.Style {
  return lipgloss.NewStyle().
  Foreground(lipgloss.Color("240"))
}

func errorStyle() lipgloss.Style {
  return lipgloss.NewStyle().
  Foreground(lipgloss.Color("196")).
  Bold(true)
}

func applyStyle(s string, focused bool) string {
  if focused {
    return focusedStyle().Render(s)
  }
  return normalStyle().Render(s)
}

func totalLabel(m PartitionModel) string {
  total := 0
  for _, percent := range m.percentages {
    val, _ := strconv.Atoi(percent)
    total += val
  }

  label := fmt.Sprintf("Total: %d%%", total)
  if total > 100 {
    label = errorStyle().Render(label + " (exceeds 100%)")
  } else if total < 100 {
    label = errorStyle().Render(label + " (less than 100%)")
  }
  return label
}

func validatePercentages(percentages [3]string) ([3]int, error) {
  labels := [3]string{"boot", "swap", "root"}
  var values [3]int
  total := 0

  for i, percent := range percentages {
    val, err := strconv.Atoi(percent)
    if err != nil || val < 0 {
      return [3]int{}, fmt.Errorf("%s percentage must be a valid positive number", labels[i])
    }
    values[i] = val
    total += val
  }

  if total > 100 {
    return [3]int{}, fmt.Errorf("total percentage (%d%%) exceeds 100%%", total)
  }

  if total < 100 {
    return [3]int{}, fmt.Errorf("total percentage (%d%%) is less than 100%%", total)
  }

  return values, nil
}
