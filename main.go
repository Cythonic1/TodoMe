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

type Mode int

const (
	NormalMode Mode = iota
	UpdateMode
)

type Styles struct {
	CursorStyle    lipgloss.Color
	Boarder        lipgloss.Color
	SelectedItem   lipgloss.Style
	NormalItem     lipgloss.Style
	InputField     lipgloss.Style
	Header         lipgloss.Style
	ContentBoarder lipgloss.Style
}

func InitStyle() *Styles {
	s := new(Styles)
	s.Boarder = lipgloss.Color(pkg.Catppuccin_base)
	s.CursorStyle = lipgloss.Color(pkg.Catppuccin_lavender)
	s.SelectedItem = lipgloss.NewStyle().BorderForeground(lipgloss.Color(pkg.Catppuccin_pink)).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(120)
	s.NormalItem = lipgloss.NewStyle().BorderForeground(s.Boarder).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(120)
	s.InputField = lipgloss.NewStyle().BorderForeground(lipgloss.Color(pkg.Catppuccin_crust)).BorderStyle(lipgloss.RoundedBorder())
	s.Header = lipgloss.NewStyle().Foreground(lipgloss.Color(pkg.Catppuccin_green))
	s.ContentBoarder = lipgloss.NewStyle().BorderForeground(lipgloss.Color(pkg.Catppuccin_peach)).Padding(5).BorderStyle(lipgloss.RoundedBorder())
	return s
}

type Module struct {
	appMode      Mode
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
	TasksParser := pkg.Init("/home/groot/projects/go/TodoMe/testfile")
	TasksParser.ParseFile()

	// To take full screent
	tea.WithAltScreen()

	// Setting styles
	style := InitStyle()

	return Module{
		appMode:      NormalMode,
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
	switch m.appMode {
	case NormalMode:
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

				m.HandleSelectedItem(m.cursor)
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}

			case "u":
				m.updateTasks()

			case "e":
				m.textBoxState = !m.textBoxState
				m.appMode = UpdateMode
				if m.textBoxState {
					m.textBox.SetValue(m.choice[m.cursor])
					return m, m.textBox.Focus()
				}
			case "a":
				// Todo add
				m.textBoxState = !m.textBoxState
				m.textBox.SetValue("")
				m.appMode = UpdateMode
				return m, m.textBox.Focus()

			}

		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height

		}

		return m, nil

	case UpdateMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "e":
				text := m.textBox.Value()
				m.choice[m.cursor] = text
				m.textBoxState = false
				m.appMode = NormalMode
				return m, nil

			case "esc":
				// Cancel editing
				m.textBoxState = false
				m.appMode = NormalMode
				return m, nil

			case "a":
				text := m.textBox.Value()
				if text == "" {
					m.textBoxState = false
					m.appMode = NormalMode
					return m, nil
				}
				m.choice = append(m.choice, text)
				m.textBoxState = false
				m.appMode = NormalMode
				return m, nil
			}
		}

		if m.textBoxState {
			var cmd tea.Cmd
			m.textBox, cmd = m.textBox.Update(msg)
			return m, cmd
		}

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

func (m Module) HandleSelectedItem(i int) {

	if i < 0 || i > len(m.choice) {
		return
	}
	if strings.HasPrefix(m.choice[i], "- [ ]") {
		parts := strings.Split(m.choice[i], "]")
		m.choice[i] = "- [x]" + parts[1]
	} else {
		parts := strings.Split(m.choice[i], "]")
		m.choice[i] = "- [ ]" + parts[1]

	}

}

func (m Module) View() string {

	var todos []string
	// Header
	_, month, day := time.Now().Date()
	weekday := time.Now().Weekday()

	header := m.styles.Header.Render(fmt.Sprintf("📅 What would you like to do today: %d/%d (%s)\n", month, day, weekday.String()))
	textBox := m.styles.InputField.Render(m.textBox.View())

	for index, todo := range m.choice {
		var line string
		if index == m.cursor {
			line = m.styles.SelectedItem.Render(todo)
		} else {
			line = m.styles.NormalItem.Render(todo)
		}

		todos = append(todos, line)
	}

	todoList := lipgloss.JoinVertical(lipgloss.Left, todos...)

	if m.textBoxState {
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			todoList,
			textBox,
		)

		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.styles.ContentBoarder.Render(content))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		todoList,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.styles.ContentBoarder.Render(content))
}

func main() {
	p := tea.NewProgram(InitialModule())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bee, there's been an error: %v", err)
		os.Exit(1)
	}
	clear()

}
