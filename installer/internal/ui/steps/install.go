package steps

import (
	"fmt"
	"strings"

	spin "github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
	services "installer.malang/internal/services"
	types "installer.malang/internal/types"
)

type InstallModel struct {
	drives       types.PartitionConfigMsg
	progressChan chan types.ProgressUpdate
	streamChan   chan types.InstallPackageStream
	progressMsg  types.ProgressUpdate
	streamMsgs   []types.InstallPackageStream
	spinner      spin.Model
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
		close(streamChan)
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

	sp := spin.New()
	sp.Style = lipgloss.NewStyle()

	return InstallModel{
		drives:       drives,
		progressChan: progressChan,
		streamChan:   streamChan,
		spinner:      sp,
	}
}

func (m InstallModel) Init() tea.Cmd {
	return tea.Batch(
		listenForProgress(m.progressChan),
		listenForStream(m.streamChan),
		performInstall(m.drives, m.progressChan, m.streamChan),
		m.spinner.Tick,
	)
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	fmt.Println("InstallModel received message:", msg)
	switch msg := msg.(type) {
	case progressUpdateMsg:
		m.progressMsg = types.ProgressUpdate(msg)
		// ignore zero-value reads from closed channel
		if m.progressMsg.Message == "" && m.progressMsg.Step == 0 {
			return m, nil
		}
		// when installation completes, stop listening; otherwise keep listening
		if m.progressMsg.Step != 4 {
			return m, listenForProgress(m.progressChan)
		}
		return m, nil

	case streamMsg:
		stream := types.InstallPackageStream(msg)
		// ignore zero-value reads from closed channel
		if stream.Line == "" && stream.Source == "" {
			return m, m.spinner.Tick
		}
		m.streamMsgs = append(m.streamMsgs, stream)
		// continue listening for more stream messages unless installation is complete
		if m.progressMsg.Step != 4 {
			return m, tea.Batch(listenForStream(m.streamChan), m.spinner.Tick)
		}
		return m, m.spinner.Tick

	case spin.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m InstallModel) View() string {
	// always show the latest progress message
	progressLine := m.progressMsg.Message
	if progressLine == "" {
		progressLine = "Starting..."
	}

	out := strings.Builder{}
	out.WriteString(progressLine + "\n")

	// show stream window when packages are being installed or when any stream messages exist
	if m.progressMsg.Step == 2 || len(m.streamMsgs) > 0 {
		// spinner + label
		out.WriteString("\n")
		out.WriteString(fmt.Sprintf("%s Installing packages...\n\n", m.spinner.View()))

		// autoscroll: show last N lines
		const maxLines = 12
		start := 0
		if len(m.streamMsgs) > maxLines {
			start = len(m.streamMsgs) - maxLines
		}
		lines := make([]string, 0, len(m.streamMsgs)-start)
		for _, s := range m.streamMsgs[start:] {
			lines = append(lines, fmt.Sprintf("[%s] %s", s.Source, s.Line))
		}

		panel := strings.Join(lines, "\n")
		if panel == "" {
			panel = "(waiting for output...)"
		}

		// simple styled box
		box := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1).Render(panel)
		out.WriteString(box)
	}

	return out.String()
}
