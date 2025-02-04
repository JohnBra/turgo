/*
Copyright © 2024 Jonathan Braat <me@jebraat.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type (
	errMsg error
)

type choice struct {
	tag         string
	description string
	color       string
}

type model struct {
	choices   []choice
	cursor    int
	tag       string
	step      int
	textInput input.Model
	res       string
	err       error
}

func initialModel() model {
	ti := input.New()
	ti.Placeholder = "..."
	ti.Focus()
	ti.CharLimit = 72
	ti.Width = 72

	return model{
		choices: []choice{
			{tag: "[ADD]", description: "feature commits, adding lines of code.", color: "#7FD85A"},
			{tag: "[FIX]", description: "bug fixing commits.", color: "#E3F9B4"},
			{tag: "[REF]", description: "small and big changes without new features.", color: "#F471F7"},
			{tag: "[BRK]", description: "breaking changes.", color: "#FA6A5A"},
		},
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return input.Blink
}

// moves cursor up and down and sets tag in model on selection
func selectTag(msg tea.KeyMsg, m model) (model, tea.Cmd) {
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
		m.tag = m.choices[m.cursor].tag
		m.step++
	}
	return m, nil
}

// sets commit message in model and execs commit on enter
func getCommitTitle(msg tea.KeyMsg, m model) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Type {

	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyEnter:
		title := fmt.Sprintf("%s %s\n", m.tag, m.textInput.Value())
		gc := exec.Command("git", "commit", "-m", title)
		out, _ := gc.Output()
		m.res = string(out)
		return m, tea.Quit
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		if m.step == 0 {
			m, cmd = selectTag(msg, m)
		} else { // for now we just have step 1 here
			m, cmd = getCommitTitle(msg, m)
		}
	case errMsg:
		m.err, cmd = msg, nil
	}

	return m, cmd
}

func (m model) View() string {
	if m.res != "" {
		return m.res
	}

	var s string
	if m.step == 0 {
		var style = gloss.NewStyle().Bold(true)
		s += fmt.Sprintf("%s\n", style.Render("Select a tag for your commit"))
		for i, choice := range m.choices {

			cursor := " " // no cursor
			if m.cursor == i {
				cursor = "❯" // cursor
			}

			var style = gloss.NewStyle().Foreground(gloss.Color(choice.color))
			s += fmt.Sprintf("%s %s\n", cursor, style.Render(choice.tag+" : "+choice.description))
		}
	} else {
		s += fmt.Sprintf("Enter a title for your %s commit:\n%s", m.tag, m.textInput.View())
	}

	return s
}

// commitCmd represents the commit command
var (
	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Commit staged files to git",
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(initialModel())

			if _, err := p.Run(); err != nil {
				fmt.Printf("Error executing commit: %v", err)
				os.Exit(1)
			}
		},
	}

	logCmd = &cobra.Command{
		Use:   "log",
		Short: "Show the git log, but with colors",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("turgo log here")
		},
	}
)

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
