package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Cythonic1/bubleTea/pkg"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	CursorStyle  lipgloss.Color
	Boarder      lipgloss.Color
	SelectedItem lipgloss.Style
	InputField   lipgloss.Style
	Header       lipgloss.Style
}

func InitStyle() *Styles {
	s := new(Styles)
	s.Boarder = lipgloss.Color(pkg.Catppuccin_base)
	s.CursorStyle = lipgloss.Color(pkg.Catppuccin_lavender)
	s.SelectedItem = lipgloss.NewStyle().BorderBackground(s.Boarder).BorderStyle(lipgloss.NormalBorder())
	s.InputField = lipgloss.NewStyle().BorderBackground(lipgloss.Color(pkg.Catppuccin_crust)).BorderStyle(lipgloss.RoundedBorder())
	s.Header = lipgloss.NewStyle().BorderForeground(lipgloss.Color(pkg.Catppuccin_blue)).BorderStyle(lipgloss.RoundedBorder())
	return s
}

type Module struct {
	choice       []string
	cursor       int
	selected     map[int]struct{}
	tasks        pkg.TodayTasks
	textBox      textinput.Model
	textBoxState bool
	width        int
	height       int
	styles       *Styles
}

func InitialModule() Module {
	// Setting text input
	ti := textinput.New()
	ti.CharLimit = 255
	ti.Width = 100
	ti.Placeholder = "What would you like todo ?"

	// Parsing file
	TasksParser := pkg.Init("/home/groot/Desktop/go/todo_bubbleTea/testfile")
	TasksParser.ParseFile()

	// To take full screent
	tea.WithAltScreen()

	// Setting styles
	style := InitStyle()

	return Module{
		choice:       TasksParser.Tasks,
		selected:     make(map[int]struct{}),
		tasks:        *TasksParser,
		textBox:      ti,
		textBoxState: false,
		styles:       style,
	}

}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal("Error executing os")
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	}

	if m.textBoxState {
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "m" {
			return m, nil
		}
		var cmd tea.Cmd
		m.textBox, cmd = m.textBox.Update(msg)
		return m, cmd
	}
	return m, nil
}

// for i, choice := range m.choice {
// 	cursor := " "
// 	if m.cursor == i {
// 		cursor = "➤  "
// 	}
// 	checked := " "
// 	if _, ok := m.selected[i]; ok {
// 		checked = "✔"
// 	}
// 	line := fmt.Sprintf("%s [%s] %s", cursor, checked, choice)
//
// 	if m.cursor == i {
// 		body += selectedStyle.Render(line) + "\n"
// 	} else {
// 		body += normalStyle.Render(line) + "\n"
// 	}
// }

func (m Module) View() string {

	// Header
	_, month, day := time.Now().Date()
	weekday := time.Now().Weekday()

	header := m.styles.Header.Render(fmt.Sprintf("📅 What would you like to do today: %d/%d (%s)\n", month, day, weekday.String()))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, header))

}

func main() {
	p := tea.NewProgram(InitialModule())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bee, there's been an error: %v", err)
		os.Exit(1)
	}
	clear()

}
