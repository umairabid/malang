package utils

import (
	"encoding/json"
	"os/exec"
)

type LSBLKOutput struct {
	Blockdevices []struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Children []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"children"`
	} `json:"blockdevices"`
}


func FetchPartitions(diskName string) []string {
	cmd := exec.Command("lsblk", "--json", "-o", "NAME,TYPE")
	out, _ := cmd.Output()

	var data LSBLKOutput
	json.Unmarshal(out, &data)

	var parts []string
	for _, dev := range data.Blockdevices {
		if dev.Name == diskName && len(dev.Children) > 0 {
			for _, c := range dev.Children {
				if c.Type == "part" {
					parts = append(parts, "/dev/"+c.Name)
				}
			}
			break
		}
	}

	return parts
}
