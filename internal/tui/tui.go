package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/afnank19/homer-code/internal/agent"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	hello        string
	viewport     viewport.Model
	messages     []string
	textarea     textarea.Model
	senderStyle  lipgloss.Style
	agentRunning bool
	agentSteps   int
}

type TestProcessMsg string

const gap = "\n\n"

func StartTUI() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "* "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return model{
		hello:        "blah",
		viewport:     vp,
		textarea:     ta,
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		agentRunning: false,
		agentSteps:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width - 2) // leaving space for the border
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()

	case TestProcessMsg:
		m.messages = append(m.messages, m.senderStyle.Render(string(msg)))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		// m.textarea.Reset()
		if m.agentSteps == 2 {
			m.agentRunning = false
			m.agentSteps = 0
			return m, nil
		}
		m.viewport.GotoBottom()
		if m.agentSteps < 2 {
			m.agentSteps++
			return m, processMessage()
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.textarea.Value() == "" {
				return m, tea.Batch(tiCmd, vpCmd)
			}

			// Some stuff to do, maybe even ponder about:
			// what will happen here is, take input send to LLM,
			// run Agent loop
			// PROFIT???

			if m.agentRunning {
				return m, tea.Batch(tiCmd, vpCmd)
			}

			m.agentRunning = true
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			m.messages = append(m.messages, m.senderStyle.Render("HOMER: ")+"PROCESSING")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, processMessage()
		}

		// We handle errors just like any other message
		// case errMsg:
		// 	m.err = msg
		// 	return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	// m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))

	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(m.textarea.View()),
	)
}

func processMessage() tea.Cmd {
	return func() tea.Msg {
		var ac agent.AgentContext = agent.AgentContext{
			Goal:    "can you stage my files for me",
			History: []agent.ToolHistory{},
		}

		ac = agent.RunAgentIteration(ac)

		h := ac.History[len(ac.History)-1] // latest item

		var userContent string
		switch h.PrevToolCalled {
		case "run_terminal_command":
			userContent += fmt.Sprintf(
				"Agent ran tool '%s' with command '%s' which output: '%s'\n",
				h.PrevToolCalled, h.PrevToolCmd, h.PrevToolOutput,
			)

		case "clarify_query":
			userContent += fmt.Sprintf("Agent asked the user: %s\n", h.PrevToolOutput)

		case "talk_to_user":
			userContent += fmt.Sprintf(
				"Agent ran tool %s said to the user: %s\n",
				h.PrevToolCalled, h.PrevToolOutput,
			)
		}

		// agent.TestFunc("lmao")
		return TestProcessMsg(userContent)
	}
}
