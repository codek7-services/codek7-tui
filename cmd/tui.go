package main

import (
	"log"

	"github.com/codek7-services/codek7-tui/internal/tui"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	app := tui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("TUI failed: %v", err)
	}
}
