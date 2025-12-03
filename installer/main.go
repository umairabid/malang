package main

import (
	"installer.malang/internal/ui/app"
  "os"
)

func main() {
  mode := os.Getenv("MODE")

  if mode == "cli" {
    RunCliMode()
  } else {
    runAppMode()
  }
}

func runAppMode() {
	app.App()
}

