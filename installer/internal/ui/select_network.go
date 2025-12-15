package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
  tea "github.com/charmbracelet/bubbletea"
  services "installer.malang/internal/services"
  types "installer.malang/internal/types"
)

const WIRED = "Wired"
const WIFI = "Wi-Fi"

type SelectNetworkModel struct {
  networkTypes []string
  selectedNetworkType string
  wifiNetworks []types.WiFiNetwork
  selectedWifiNetwork int
  ssid string
  password string
  collectWifiCredentials bool
  connected bool
  spinner spinner.Model
  loadingMessage string
  errorMessage string
}

func InitNetworkStep() tea.Model {
  return SelectNetworkModel{
    networkTypes: []string{WIRED, WIFI},
    selectedNetworkType: WIRED,
    wifiNetworks: []types.WiFiNetwork{},
    selectedWifiNetwork: 0,
    ssid: "",
    password: "",
    connected: false,
    collectWifiCredentials: false,
    spinner: spinner.New(),
    loadingMessage: "",
    errorMessage: "",
  }
}

func checkConnection() tea.Cmd {
  return func() tea.Msg {
    connected := services.CheckNetworkConnection()
    return types.NetworkStatusMsg(connected)
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

func handleSelectionChange(m SelectNetworkModel, key string) SelectNetworkModel {
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

func (m SelectNetworkModel) Init() tea.Cmd {
  return m.spinner.Tick
}

func (m SelectNetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "tab", "shift+tab", "up", "down":
      m = handleSelectionChange(m, msg.String())
      return m, nil
    case "enter":
      m.connected = false
      m.errorMessage = ""
      if m.selectedNetworkType == WIFI {
        m.collectWifiCredentials = true
        m.loadingMessage = "Collecting Wi-Fi networks..."
        return m, collectWifiConnections()
      } else {
        m.loadingMessage = "Checking network connection..." 
        return m, checkConnection()
      }
    }
  case types.NetworkStatusMsg:
    m.connected = bool(msg)
    m.loadingMessage = ""
    if !m.connected {
      m.errorMessage = "No network connection detected. Please check your connection and try again."
    } else {
      m.errorMessage = ""
    }
    return m, nil
  case types.WifiNetworkError:
    m.loadingMessage = ""
    m.collectWifiCredentials = false
    m.errorMessage = string(msg)
    return m, nil
  case types.WiFiNetworksMsg:
    m.loadingMessage = ""
    m.collectWifiCredentials = false
    m.networkTypes = []string{}
    m.wifiNetworks = []types.WiFiNetwork(msg)
    return m, nil
  case spinner.TickMsg:
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
  } 
  return m, nil
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

func (m SelectNetworkModel) View() string {
  var lines string 
  if m.loadingMessage == "" { 
    if len(m.wifiNetworks) > 0 {
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
