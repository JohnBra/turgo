package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type choice struct {
	text  string
	color string
}

type model struct {
	choices []choice
	cursor  int
}

func initialModel() model {
	return model{
		choices: []choice{
			{text: "[ADD] : feature commits, adding lines of code.", color: "#7FD85A"},
			{text: "[FIX] : bug fixing commits.", color: "#E3F9B4"},
			{text: "[REF] : small and big changes without new features.", color: "#F471F7"},
			{text: "[BRK] : breaking changes.", color: "#FA6A5A"},
		},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			fmt.Println(m.cursor)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select a tag for your commit\n"

	for i, choice := range m.choices {

		cursor := " " // no cursor
		if m.cursor == i {
			cursor = "❯" // cursor!
		}

		var style = gloss.NewStyle().Foreground(gloss.Color(choice.color))
		s += fmt.Sprintf("%s %s\n", cursor, style.Render(choice.text))
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
