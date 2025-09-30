package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/afnank19/homer-code/internal/agent"
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

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	appConfigDir := filepath.Join(configDir, "homer-code/plugins")

	err = os.MkdirAll(appConfigDir, 0o755)
	if err != nil {
		panic(err)
	}

	// plugin.LoadPlugins()
	agent.BuildSystemPrompt()
	// fmt.Println(agent.SYSTEM_PROMPT)

	tui.StartTUI()

	// agent.StartLoop(" you stage my files for me")

	log.Println("--- PROGRAM END ---")
}
