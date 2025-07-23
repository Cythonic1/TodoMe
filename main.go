package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Cythonic1/bubleTea/pkg"
	"github.com/charmbracelet/bubbles/textinput"
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
	choice       []string
	cursor       int
	selected     map[int]struct{}
	tasks        pkg.TodayTasks
	textBox      textinput.Model
	textBoxState bool
}

func InitialModule() Module {
	ti := textinput.New()
	ti.CharLimit = 255
	ti.Width = 100
	TasksParser := pkg.Init("/home/groot/Desktop/go/todo_bubbleTea/testfile")
	TasksParser.ParseFile()

	return Module{
		choice:       TasksParser.Tasks,
		selected:     make(map[int]struct{}),
		tasks:        *TasksParser,
		textBox:      ti,
		textBoxState: false,
	}

}

func (m Module) Init() tea.Cmd {
	return textinput.Blink
}

func (m Module) updateTasks() {
	for index, _ := range m.choice {
		_, ok := m.selected[index]
		if ok {
			parts := strings.Split(m.tasks.Tasks[index], "]")
			m.tasks.Tasks[index] = "- [x]" + parts[1]
		} else {
			parts := strings.Split(m.tasks.Tasks[index], "]")
			m.tasks.Tasks[index] = "- [ ]" + parts[1]
		}
	}
	m.tasks.ReplaceTodos()
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
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "u":
			m.updateTasks()

		case "m":
			m.textBoxState = !m.textBoxState
			if m.textBoxState {
				m.textBox.Focus()
			} else {
				m.textBox.Blur()
			}
		case tea.KeyEnter.String():
			userInput := m.textBox.Value()
			if userInput != "" {
				m.choice = append(m.choice, userInput)
				if !strings.Contains(userInput, "- [ ]") {
					m.tasks.Tasks = append(m.tasks.Tasks, "- [ ] "+userInput)
				} else {
					m.tasks.Tasks = append(m.tasks.Tasks, userInput)
				}

			}
			m.textBox.SetValue("")
			m.textBox.Blur()
			m.textBoxState = false

		}

	}

	if m.textBoxState {
		var cmd tea.Cmd
		m.textBox, cmd = m.textBox.Update(msg)
		return m, cmd
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
	view := header + "\n\n" + body
	if m.textBoxState {
		view += "\n" + m.textBox.View()
	}
	return view
}

func main() {
	p := tea.NewProgram(InitialModule())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bee, there's been an error: %v", err)
		os.Exit(1)
	}
}
