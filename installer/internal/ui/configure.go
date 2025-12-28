package ui

import (
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"

    tea "github.com/charmbracelet/bubbletea"
    services "installer.malang/internal/services"
    types "installer.malang/internal/types"
)

type ConfigureModel struct {
    progressChan   chan types.ConfigureStream
    spinner        spinner.Model
    mountPoints    [2]string
    progressMsg    types.ConfigureStream
    userConfig     types.UserConfig
    username       textinput.Model
    password       textinput.Model
    focusedField   int
    formSubmitted  bool
}

func configureSystem(m ConfigureModel) tea.Cmd {
    return func() tea.Msg {
        services.ConfigureSystem(m.mountPoints, m.progressChan)
        return nil
    }
}

func createUser(m ConfigureModel) tea.Cmd {
    return func() tea.Msg {
        services.CreateUser(m.userConfig, m.progressChan)
        return nil
    }
}

func initUsernameInput() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "Enter username"
    ti.CharLimit = 32
    ti.Width = 30
    ti.Focus()
    return ti
}

func initPasswordInput() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "Enter password"
    ti.CharLimit = 128
    ti.Width = 30
    ti.EchoMode = textinput.EchoPassword
    ti.EchoCharacter = '•'
    return ti
}

func InitConfigureStep(mounts [2]string) tea.Model {
    return ConfigureModel{
        progressChan:  make(chan types.ConfigureStream, 10),
        spinner:       spinner.New(),
        mountPoints:   mounts,
        username:      initUsernameInput(),
        password:      initPasswordInput(),
        focusedField:  0,
        formSubmitted: false,
    }
}

func listenForConfigureProgress(progressChan chan types.ConfigureStream) tea.Cmd {
    return func() tea.Msg {
        update := <-progressChan
        return types.ConfigureStream(update)
    }
}

func (m ConfigureModel) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m ConfigureModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    // Handle form input if not submitted
    if !m.formSubmitted {
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "enter":
                if m.focusedField == 0 {
                    // Move to password field
                    m.focusedField = 1
                    m.username.Blur()
                    m.password.Focus()
                    return m, nil
                } else {
                    // Submit form and start configuration
                    if m.username.Value() == "" || m.password.Value() == "" {
                        return m, nil
                    }
                    m.userConfig = types.UserConfig{
                        Username: m.username.Value(),
                        Password: m.password.Value(),
                    }
                    m.formSubmitted = true
                    return m, tea.Batch(
                        configureSystem(m),
                        listenForConfigureProgress(m.progressChan),
                        m.spinner.Tick,
                    )
                }
            case "tab", "shift+tab":
                if m.focusedField == 0 {
                    m.focusedField = 1
                    m.username.Blur()
                    m.password.Focus()
                } else {
                    m.focusedField = 0
                    m.password.Blur()
                    m.username.Focus()
                }
                return m, nil
            }
        }

        // Update the focused input
        if m.focusedField == 0 {
            m.username, cmd = m.username.Update(msg)
        } else {
            m.password, cmd = m.password.Update(msg)
        }
        return m, cmd
    }

    // Handle configuration progress
    switch msg := msg.(type) {
    case spinner.TickMsg:
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    case types.ConfigureStream:
        m.progressMsg = msg
        // Check if configuration is complete and start user creation
        if msg.Line == "Configuration complete" {
            return m, tea.Batch(
                createUser(m),
                listenForConfigureProgress(m.progressChan),
            )
        }
        return m, listenForConfigureProgress(m.progressChan)
    }
    return m, nil
}

func (m ConfigureModel) View() string {
    if !m.formSubmitted {
        var s string
        s += "Create User Account\n\n"
        s += "Username:\n"
        s += m.username.View() + "\n\n"
        s += "Password:\n"
        s += m.password.View() + "\n\n"
        s += "(tab: switch fields • enter: submit • q: quit)"
        return s
    }
    return m.spinner.View() + " " + m.progressMsg.Line
}
