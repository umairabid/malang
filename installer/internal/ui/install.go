package ui

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/lipgloss"

    tea "github.com/charmbracelet/bubbletea"
    services "installer.malang/internal/services"
    types "installer.malang/internal/types"
)

type InstallModel struct {
    drives             types.PartitionConfigMsg
    progressChan chan types.ProgressUpdate
    streamChan     chan types.InstallPackageStream
    progressMsg    types.ProgressUpdate
    streamMsgs     []types.InstallPackageStream
    spinner            spinner.Model
}

type progressUpdateMsg types.ProgressUpdate
type streamMsg types.InstallPackageStream

func listenForInstallProgress(progressChan chan types.ProgressUpdate) tea.Cmd {
    return func() tea.Msg {
        update := <-progressChan
        return progressUpdateMsg(update)
    }
}

func listenForInstallStream(streamChan <-chan types.InstallPackageStream) tea.Cmd {
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
        mountPoints, err := services.Install([3]string(drives), progressChan, streamChan)
        progressChan <- types.ProgressUpdate{
            Message: "Installation complete.",
            Step:        4,
            Success: true,
        }
        if err != nil {
            progressChan <- types.ProgressUpdate{
                Message: fmt.Sprintf("Installation failed: %v", err),
                Step:        4,
                Success: false,
            }
            return nil
        } else {
            return types.InstallCompleteMsg{mountPoints[0], mountPoints[1]}
        }
    }
}

func InitInstallStep(drives types.PartitionConfigMsg) tea.Model {
    progressChan := make(chan types.ProgressUpdate, 10)
    streamChan := make(chan types.InstallPackageStream, 10)

    sp := spinner.New()
    sp.Style = lipgloss.NewStyle()

    return InstallModel{
        drives:             drives,
        progressChan: progressChan,
        streamChan:     streamChan,
        spinner:            sp,
        progressMsg:    types.ProgressUpdate{Message: "Starting installation...", Step: 0, Success: true},
        streamMsgs:     []types.InstallPackageStream{},
    }
}

func (m InstallModel) Init() tea.Cmd {
    return tea.Batch(
        m.spinner.Tick,
        listenForInstallProgress(m.progressChan),
        listenForInstallStream(m.streamChan),
        performInstall(m.drives, m.progressChan, m.streamChan),
    )
}

func (m InstallModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case progressUpdateMsg:
        m.progressMsg = types.ProgressUpdate(msg)
        return m, listenForInstallProgress(m.progressChan)

    case streamMsg:
        stream := types.InstallPackageStream(msg)
        m.streamMsgs = append(m.streamMsgs, stream)
        return m, listenForInstallStream(m.streamChan)

    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }

    return m, nil
}

func (m InstallModel) View() string {
    progressLine := m.progressMsg.Message
    out := strings.Builder{}
    out.WriteString("stating installer archlnux" + "\n")
    out.WriteString(m.spinner.View() + " ")
    out.WriteString(progressLine + "\n")

    // show stream window when packages are being installed or when any stream messages exist
    if m.progressMsg.Step == 2 || len(m.streamMsgs) > 0 {
        // spinner + label
        out.WriteString("\n")

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
