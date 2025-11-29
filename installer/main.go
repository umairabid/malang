package main

import (
	"installer.malang/internal/ui/app"
)

const MODE = "cli"

func main() {
  if MODE == "app" {
    runAppMode()
  } else {
    RunCliMode()
  }
}

func runAppMode() {
	app.App()
}

