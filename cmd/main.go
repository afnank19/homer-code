package main

import (
	"fmt"
	"log"
	"os"

	"github.com/afnank19/homer-code/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

// new ubuntu ssh test

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("--- DEBUG BEGIN ---")

	tui.StartTUI()

	// agent.StartLoop(" you stage my files for me")

	log.Println("--- PROGRAM END ---")
}
