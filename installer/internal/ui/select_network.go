package ui

import (
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/bubbles/textinput"

    tea "github.com/charmbracelet/bubbletea"
    services "installer.malang/internal/services"
    types "installer.malang/internal/types"
)

const WIRED = "Wired"
const WIFI = "Wi-Fi"

type SelectNetworkModel struct {
    networkTypes                []string
    selectedNetworkType string
    wifiNetworks                []types.WiFiNetwork
    selectedWifiNetwork int
    ssid                                string
    password                        textinput.Model
    connected                     bool
    spinner                         spinner.Model
    loadingMessage            string
    errorMessage                string
}

func initTextInput() textinput.Model {
    ti := textinput.New()
    ti.Placeholder = "Enter Wi-Fi Password"
    ti.CharLimit = 64
    ti.Width = 30
    ti.EchoMode = textinput.EchoPassword
    ti.EchoCharacter = '•'
    return ti
}

func InitNetworkStep() tea.Model {
    return SelectNetworkModel{
        networkTypes:                []string{WIRED, WIFI},
        selectedNetworkType: WIRED,
        wifiNetworks:                []types.WiFiNetwork{},
        selectedWifiNetwork: 0,
        ssid:                                "",
        password:                        initTextInput(),
        connected:                     false,
        spinner:                         spinner.New(),
        loadingMessage:            "",
        errorMessage:                "",
    }
}

func checkConnection() tea.Cmd {
    return func() tea.Msg {
        connected := services.CheckNetworkConnection()
        if connected {
            return types.NetworkConnectedMsg(true)
        } else {
            return types.NetworkFailureMsg(true)
        }
    }
}

func collectWifiConnections() tea.Cmd {
    return func() tea.Msg {
        networks, error := services.GetWiFiNetworks()
        if error != nil {
            return types.WifiNetworkError(error.Error())
        }
        return types.WiFiNetworksMsg(networks)
    }
}

func connectWithWifi(ssid string, password string) tea.Cmd {
    return func() tea.Msg {
        error := services.ConnectWithWiFi(ssid, password)
        if error != nil {
            return types.WifiNetworkError(error.Error())
        } else {
            return checkConnection()()
        }
    }
}

func handleChange(m SelectNetworkModel, key string) SelectNetworkModel {
    if len(m.wifiNetworks) > 0 {
        if key == "up" || key == "shift+tab" {
            m.selectedWifiNetwork = (m.selectedWifiNetwork - 1 + int(len(m.wifiNetworks))) % int(len(m.wifiNetworks))
        } else if key == "down" || key == "tab" {
            m.selectedWifiNetwork = (m.selectedWifiNetwork + 1) % int(len(m.wifiNetworks))
        }
        return m
    }

    if len(m.networkTypes) > 0 {
        if m.selectedNetworkType == WIRED {
            m.selectedNetworkType = WIFI
        } else {
            m.selectedNetworkType = WIRED
        }
    }
    return m
}

func handleSelection(m SelectNetworkModel) (SelectNetworkModel, tea.Cmd) {
    if m.ssid != "" {
        m.loadingMessage = "Checking network connection for " + m.ssid + "..."
        return m, connectWithWifi(m.ssid, m.password.Value())
    } else if len(m.wifiNetworks) > 0 {
        m.ssid = m.wifiNetworks[m.selectedWifiNetwork].SSID
        return m, m.password.Focus()
    } else {
        if m.selectedNetworkType == WIFI {
            m.loadingMessage = "Collecting Wi-Fi networks..."
            return m, collectWifiConnections()
        } else {
            m.loadingMessage = "Checking network connection..."
            return m, checkConnection()
        }
    }
}

func (m SelectNetworkModel) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m SelectNetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab", "shift+tab", "up", "down":
            m = handleChange(m, msg.String())
            return m, nil
        case "enter":
            m.connected = false
            m.errorMessage = ""
            return handleSelection(m)
        }
    case types.NetworkFailureMsg:
        m.loadingMessage = ""
        m.errorMessage = "No network connection detected. Please check your connection and try again."
        return m, nil
    case types.WifiNetworkError:
        m.loadingMessage = ""
        m.errorMessage = string(msg)
        return m, nil
    case types.WiFiNetworksMsg:
        m.loadingMessage = ""
        m.networkTypes = []string{}
        m.wifiNetworks = []types.WiFiNetwork(msg)
        return m, nil
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    m.password, cmd = m.password.Update(msg)
    return m, cmd
}

func (m SelectNetworkModel) View() string {
    var lines string
    if m.loadingMessage == "" {
        if m.ssid != "" {
            lines += renderWifiPassword(m.password, m.ssid)
        } else if len(m.wifiNetworks) > 0 {
            lines += renderWifiNetworks(m.wifiNetworks, m.selectedWifiNetwork)
        } else {
            lines += renderNetworkChoices(m.selectedNetworkType)
        }
    } else {
        lines += m.spinner.View() + "\t" + m.loadingMessage
    }

    if m.connected {
        lines += "\n\n" + focusedStyle().Render("Network connected successfully!")
    }

    if m.errorMessage != "" {
        lines += "\n\n" + errorStyle().Render(m.errorMessage)
    }

    help := "\n\n(↑/↓/tab: navigate • enter: confirm • q: quit)"
    return lines + help
}

func renderWifiPassword(password textinput.Model, ssid string) string {
    var lines string
    lines += "Selected Network: " + focusedStyle().Render(ssid)
    lines += "\nPlease enter the Wi-Fi password"
    lines += "\n\n" + password.View()
    return lines
}

func renderNetworkChoices(networkType string) string {
    types := [2]string{WIRED, WIFI}
    var lines string
    lines += "Select Network Type:\n\n"
    for i := range types {
        lines += applyStyle(types[i], networkType == types[i]) + "\n"
    }
    return lines
}

func renderWifiNetworks(networks []types.WiFiNetwork, selectedNetwork int) string {
    var lines string
    lines += "Select Wi-Fi Network:\n\n"
    for i := range networks {
        lines += applyStyle(networks[i].SSID, selectedNetwork == i) + "\n"
    }
    return lines
}
