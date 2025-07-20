package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Cythonic1/bubleTea/pkg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // pinkish
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")). // bright yellow
			Background(lipgloss.Color("57")).  // dark blue bg
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")) // light gray

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")). // blue
			Bold(true).
			Underline(true)
)

type Module struct {
	choice   []string
	cursor   int
	selected map[int]struct{}
}

func InitialModule() Module {
	TasksParser := pkg.Init()
	TasksParser.ParseFile("/home/groot/Documents/Obsidian Vault/Weekly Tasks/week14.md")

	return Module{
		choice:   TasksParser.Tasks,
		selected: make(map[int]struct{}),
	}

}

func (m Module) Init() tea.Cmd {
	return nil
}

func (m Module) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choice)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m Module) View() string {
	var body string

	_, month, day := time.Now().Date()
	weekday := time.Now().Weekday()
	header := headerStyle.Render(fmt.Sprintf("📅 What would you like to do today: %d/%d (%s)\n", month, day, weekday.String()))

	for i, choice := range m.choice {
		cursor := " "
		if m.cursor == i {
			cursor = "➤  "
		}
		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✔"
		}
		line := fmt.Sprintf("%s [%s] %s", cursor, checked, choice)

		if m.cursor == i {
			body += selectedStyle.Render(line) + "\n"
		} else {
			body += normalStyle.Render(line) + "\n"
		}
	}

	body += "\n Press q or Ctrl+c to quit\n"
	return header + "\n\n" + body
}

func main() {
	p := tea.NewProgram(InitialModule())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bee, there's been an error: %v", err)
		os.Exit(1)
	}

}
