package main

import (
	"installer.malang/internal/ui"
  "os"
  "fmt"
)

func main() {
  mode := os.Getenv("MODE")
  fmt.Println("Running in mode:", os.Getenv("IS_MOCKING"))

  if mode == "cli" {
    RunCliMode()
  } else {
    runAppMode()
  }
}

func runAppMode() {
	ui.App()
}
