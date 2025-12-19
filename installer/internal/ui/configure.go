package ui

import (
  "github.com/charmbracelet/bubbles/spinner"

  services "installer.malang/internal/services"
  types "installer.malang/internal/types"
  tea "github.com/charmbracelet/bubbletea"
)

type ConfigureModel struct {
  progressChan chan types.ConfigureStream
  spinner      spinner.Model
  mountPoints  [2]string
  progressMsg  types.ConfigureStream
}

func configureSystem(m ConfigureModel) tea.Cmd {
  return func() tea.Msg {
    services.ConfigureSystem(m.mountPoints, m.progressChan)
    return nil
  }
}

func InitConfigureStep(mounts [2]string) tea.Model {
  return ConfigureModel{
    progressChan: make(chan types.ConfigureStream, 10),
    spinner:     spinner.New(),
    mountPoints: mounts,
  }
}

func listenForConfigureProgress(progressChan chan types.ConfigureStream) tea.Cmd {
  return func() tea.Msg {
    update := <-progressChan
    return types.ConfigureStream(update)
  }
}

func (m ConfigureModel) Init() tea.Cmd {
  return tea.Batch(
    configureSystem(m),
    listenForConfigureProgress(m.progressChan),
    m.spinner.Tick,
  )
}

func (m ConfigureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case spinner.TickMsg:
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
  case types.ConfigureStream:
    m.progressMsg = msg
    return m, listenForConfigureProgress(m.progressChan)
  }
  return m, nil
}

func (m ConfigureModel) View() string {
  return m.spinner.View() + " " + m.progressMsg.Line
}
