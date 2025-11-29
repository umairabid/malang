package steps

import (
  "fmt"

  "github.com/charmbracelet/bubbles/progress"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  services "installer.malang/internal/services"
  types "installer.malang/internal/types"
)

type InstallModel struct {
  drives          types.PartitionConfigMsg
  progress        progress.Model
  currentProgress float64
  currentMessage  string
  done            bool
  err             error
  result          [2]string
  progressChan    chan services.ProgressUpdate
}

type installCompleteMsg struct {
  result [2]string
  err    error
}

type progressUpdateMsg services.ProgressUpdate

func listenForProgress(progressChan <-chan services.ProgressUpdate) tea.Cmd {
  return func() tea.Msg {
    update := <-progressChan
    return progressUpdateMsg(update)
  }
}

func performInstall(drives types.PartitionConfigMsg, progressChan chan services.ProgressUpdate) tea.Cmd {
  return func() tea.Msg {
    result := services.Install([3]string(drives), progressChan)
    close(progressChan)
    return installCompleteMsg{result: result, err: nil}
  }
}

func InitInstallStep(drives types.PartitionConfigMsg) tea.Model {
  prog := progress.New(progress.WithDefaultGradient())
  progressChan := make(chan services.ProgressUpdate, 10)

  return InstallModel{
    drives:          drives,
    progress:        prog,
    currentProgress: 0,
    currentMessage:  "Preparing installation...",
    done:            false,
    progressChan:    progressChan,
  }
}

func (m InstallModel) Init() tea.Cmd {
  return tea.Batch(
    listenForProgress(m.progressChan),
    performInstall(m.drives, m.progressChan),
  )
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    if msg.String() == "ctrl+c" || msg.String() == "q" {
      return m, tea.Quit
    }

  case installCompleteMsg:
    m.done = true
    m.err = msg.err
    m.result = msg.result
    m.currentProgress = 1.0
    if m.err == nil {
      return m, func() tea.Msg {
        return types.InstallCompleteMsg(m.result)
      }
    }
    return m, nil

  case progressUpdateMsg:
    m.currentMessage = msg.Message
    m.currentProgress = float64(msg.Current) / float64(msg.Total)
    cmd := m.progress.SetPercent(m.currentProgress)
    return m, tea.Batch(cmd, listenForProgress(m.progressChan))

  case progress.FrameMsg:
    progressModel, cmd := m.progress.Update(msg)
    m.progress = progressModel.(progress.Model)
    return m, cmd
  }

  return m, nil
}

func (m InstallModel) View() string {
  if m.err != nil {
    errorStyle := lipgloss.NewStyle().
      Foreground(lipgloss.Color("196")).
      Bold(true)
    return errorStyle.Render("Installation failed: " + m.err.Error())
  }

  if m.done {
    successStyle := lipgloss.NewStyle().
      Foreground(lipgloss.Color("42")).
      Bold(true)
    return successStyle.Render("✓ Installation completed successfully!")
  }

  percentage := int(m.currentProgress * 100)
  return fmt.Sprintf("\nInstalling Arch Linux...\n\n%s\n\n[%d%%] %s\n\nPlease wait...",
    m.progress.View(),
    percentage,
    m.currentMessage)
}
