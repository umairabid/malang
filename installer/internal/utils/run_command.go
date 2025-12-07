package utils

import (
	"os/exec"
	"strings"
)

type Command struct {
	Args  []string
	Stdin *string
}

func RunCommands(commands []Command) error {
	for _, command := range commands {
		args := command.Args
		cmd := exec.Command(args[0], args[1:]...)
		if command.Stdin != nil {
			input := *command.Stdin
			cmd.Stdin = strings.NewReader(input)
		}
		_, err := cmd.CombinedOutput()
		if err != nil {
      return err
		}
	}
  return nil
}
