package implementations

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	types "installer.malang/internal/types"
)

func CheckNetworkConnection() bool {
	cmd := exec.Command("ping", "-c", "1", "8.8.8.8")
	err := cmd.Run()
	return err == nil
}

func GetWiFiNetworks() ([]types.WiFiNetwork, error) {
	cmd := exec.Command("nmcli", "dev", "wifi", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w, %s", err, string(output))
	}

	return parseWiFiNetworks(string(output))
}

func ConnectWithWiFi(ssid string, password string) error {
	var cmd *exec.Cmd

	if password == "" {
		// Open network
		cmd = exec.Command("nmcli", "dev", "wifi", "connect", ssid)
	} else {
		// Secured network
		cmd = exec.Command("nmcli", "dev", "wifi", "connect", ssid, "password", password)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to connect to WiFi: %s, error: %w", string(output), err)
	}

	return nil
}

func parseWiFiNetworks(output string) ([]types.WiFiNetwork, error) {
	var networks []types.WiFiNetwork
	scanner := bufio.NewScanner(strings.NewReader(output))

	if !scanner.Scan() {
		return networks, nil
	}

	var ssids map[string]bool = make(map[string]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		network, err := parseWiFiLine(line)
		if err != nil {
			continue
		}

		_, exists := ssids[network.SSID]
		if exists || network.SSID == "--" || network.SSID == "" {
			continue
		}
		ssids[network.SSID] = true
		networks = append(networks, network)
	}

	return networks, nil
}

func parseWiFiLine(line string) (types.WiFiNetwork, error) {
	fields := strings.Fields(line)
	if len(fields) < 8 {
		return types.WiFiNetwork{}, fmt.Errorf("invalid wifi line format")
	}
	startIndex := 0
	inUse := strings.Contains(fields[0], "*")

	if inUse {
		startIndex = 1
	}
	bssid := fields[startIndex]
	ssid := fields[startIndex+1]

	chanStr := fields[startIndex+4]
	channel, err := strconv.Atoi(chanStr)
	if err != nil {
		channel = 0
	}

	signalStr := fields[startIndex+5]
	signal, err := strconv.Atoi(signalStr)
	if err != nil {
		signal = 0
	}

	security := ""
	if len(fields) > 8 {
		security = strings.Join(fields[8:], " ")
	}

	requiresPSK := strings.Contains(security, "WPA") || strings.Contains(security, "WEP")

	frequency := calculateFrequency(channel)

	return types.WiFiNetwork{
		SSID:        ssid,
		Security:    security,
		Signal:      signal,
		InUse:       inUse,
		Hidden:      false,
		BSSID:       bssid,
		Channel:     channel,
		Frequency:   frequency,
		RequiresPSK: requiresPSK,
	}, nil
}

func calculateFrequency(channel int) int {
	if channel >= 1 && channel <= 14 {
		// 2.4 GHz band
		return 2407 + (channel * 5)
	} else if channel >= 36 && channel <= 165 {
		// 5 GHz band (simplified)
		return 5000 + (channel * 5)
	}
	return 0
}
