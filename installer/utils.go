package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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

type Command struct {
	Args  []string
	Stdin *string
}

func runCommands(commands []Command) {
	for _, command := range commands {
		args := command.Args
		fmt.Printf("✅ %v running\n", args)
		cmd := exec.Command(args[0], args[1:]...)
		if command.Stdin != nil {
			input := *command.Stdin
			fmt.Printf("---- stdin ----\n%s\n---------------\n", input)
			cmd.Stdin = strings.NewReader(input)
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Command %q failed: %v\nOutput:\n%s", args, err, string(out))
			panic(err)
		}
		fmt.Printf("✅ %v succeeded\n%s\n", args, string(out))
	}
}

func getPartitions(diskName string) []string {
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
