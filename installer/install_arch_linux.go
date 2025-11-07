package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runCommands(commands [][]string) {
	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Command %q failed: %v\nOutput:\n%s", args, err, string(out))
		}
		fmt.Printf("✅ %v succeeded\n%s\n", args, string(out))
	}
}

func main() {

	// format and setup swap
	prepareCommands := [][]string{
		{"mkfs.fat", "-F32", efiDisk},
		{"mkfs.ext4", rootDisk},
		{"mkswap", swapDisk},
		{"swapon", swapDisk},
	}

  runCommands(prepareCommands)
}
