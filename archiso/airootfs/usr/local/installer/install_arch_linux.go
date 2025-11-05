package main

import (
	"flag"
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
	var efiDisk string
	var swapDisk string
	var rootDisk string

	flag.StringVar(&efiDisk, "efiDisk", "", "Name of efi disk partition")
	flag.StringVar(&swapDisk, "swapDisk", "", "Name of swap disk partition")
	flag.StringVar(&rootDisk, "rootDisk", "", "Name of root disk partition")

	// format and setup swap
	prepareCommands := [][]string{
		{"mkfs.fat", "-F32", efiDisk},
		{"mkfs.ext4", rootDisk},
		{"mkswap", swapDisk},
		{"swapon", swapDisk},
	}

  runCommands(prepareCommands)
}
