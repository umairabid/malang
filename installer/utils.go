package main

import (
	"fmt"
	"os/exec"
	"strings"
)

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
