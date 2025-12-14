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
  networkType string
  ssid string
  password string
  connected bool
  checking_connection bool
  connection_failed bool
  spinner spinner.Model
}

func InitNetworkStep() tea.Model {
  return SelectNetworkModel{
    networkType: WIRED,
    ssid: "",
    password: "",
    connected: false,
    checking_connection: false,
    connection_failed: false,
    spinner: spinner.New(),
  }
}

func checkConnection() tea.Cmd {
  return func() tea.Msg {
    connected := services.CheckNetworkConnection()
    return types.NetworkStatusMsg(connected)
  }
}


func (m SelectNetworkModel) Init() tea.Cmd {
  return m.spinner.Tick
}

func (m SelectNetworkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "tab", "shift+tab", "up", "down":
      if m.networkType == WIRED {
        m.networkType = WIFI
      } else {
        m.networkType = WIRED
      }
      return m, nil
    case "enter":
      m.connected = false
      m.connection_failed = false
      m.checking_connection = true
      return m, checkConnection()
    }
  case types.NetworkStatusMsg:
    m.connected = bool(msg)
    m.connection_failed = !m.connected
    m.checking_connection = false
    return m, nil
  case spinner.TickMsg:
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
  } 
  return m, nil
}

func (m SelectNetworkModel) View() string {
  types := [2]string{WIRED, WIFI}
  var lines string 
  if !m.checking_connection { 
    lines += "Select Network Type:\n\n"
    for i := range types {
      lines += applyStyle(types[i], m.networkType == types[i]) + "\n"
    }
  } else {
    lines += m.spinner.View() + "\tChecking network connection..."
  }

  if m.connected {
    lines += "\n\n" + focusedStyle().Render("Network connected successfully!")
  }

  if m.connection_failed {
    lines += "\n\n" + errorStyle().Render("Failed to connect to the network. Please try again.")
  }

  help := "\n\n(↑/↓/tab: navigate • enter: confirm • q: quit)"
  return lines + help
}
