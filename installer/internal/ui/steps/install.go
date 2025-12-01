package steps

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	services "installer.malang/internal/services"
	types "installer.malang/internal/types"
)

type InstallModel struct {
	drives       types.PartitionConfigMsg
	progressChan chan types.ProgressUpdate
	streamChan   chan types.InstallPackageStream
  progressMsg types.ProgressUpdate
  streamMsgs []types.InstallPackageStream
}

type progressUpdateMsg types.ProgressUpdate
type streamMsg types.InstallPackageStream

func listenForProgress(progressChan <-chan types.ProgressUpdate) tea.Cmd {
	return func() tea.Msg {
		update := <-progressChan
		return progressUpdateMsg(update)
	}
}

func listenForStream(streamChan <-chan types.InstallPackageStream) tea.Cmd {
	return func() tea.Msg {
		stream := <-streamChan
		return streamMsg(stream)
	}
}

func performInstall(
	drives types.PartitionConfigMsg,
	progressChan chan types.ProgressUpdate,
	streamChan chan types.InstallPackageStream,
) tea.Cmd {
	return func() tea.Msg {
		services.Install([3]string(drives), progressChan, streamChan)
		close(progressChan)
		return progressUpdateMsg(types.ProgressUpdate{
			Message: "Installation complete.",
			Step:    4,
			Success: true,
		})
	}
}

func InitInstallStep(drives types.PartitionConfigMsg) tea.Model {
	progressChan := make(chan types.ProgressUpdate, 10)
	streamChan := make(chan types.InstallPackageStream, 10)

	return InstallModel{
		drives:       drives,
		progressChan: progressChan,
		streamChan:   streamChan,
	}
}

func (m InstallModel) Init() tea.Cmd {
	return tea.Batch(
		listenForProgress(m.progressChan),
		listenForStream(m.streamChan),
		performInstall(m.drives, m.progressChan, m.streamChan),
	)
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m InstallModel) View() string {
	return ""
}
