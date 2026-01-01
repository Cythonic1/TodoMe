# 📝 BubbleTea Todo App

A beautiful, terminal-based todo list manager built with Go and Bubble Tea. Features vim-like keybindings and a stunning Catppuccin color scheme.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## ✨ Features

- **Vim-like Navigation**: Use `j/k` or arrow keys for intuitive movement
- **Todo Management**: Add, edit, toggle, and update todos seamlessly
- **Beautiful UI**: Powered by Bubble Tea with Catppuccin theme
- **File Persistence**: Todos are automatically saved to your notes directory
- **Scrollable View**: Efficiently handles large todo lists
- **Modal Editing**: Dedicated modes for different operations

## 🚀 Installation

### Prerequisites

- Go 1.25 
- Git

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/bubleTea.git
cd bubleTea

# Build the application
go build -ldflags="-s -w" -o todo-app main.go

# Move to your PATH (optional)
sudo mv todo-app /usr/local/bin/
```

### Quick Install

```bash
go install github.com/yourusername/bubleTea@latest
```

## 📖 Usage

### Starting the App

```bash
./todo-app
```

### Keybindings

#### Normal Mode
| Key | Action |
|-----|--------|
| `j` or `↓` | Move cursor down |
| `k` or `↑` | Move cursor up |
| `Space` | Toggle todo (check/uncheck) |
| `a` | Add new todo |
| `e` | Edit current todo |
| `u` | Update/save todos to file |
| `q` or `Ctrl+C` | Quit application |

#### Edit/Add Mode
| Key | Action |
|-----|--------|
| `Enter` | Save changes |
| `Esc` | Cancel and return to normal mode |

## 🛠️ Configuration

The app reads todos from a file in your notes directory. Update the path in `main.go`:

```go
TasksParser := pkg.Init("/home/<YOURHOME>/notes/")
```

Change this to your preferred location:

```go
TasksParser := pkg.Init("/path/to/your/notes/")
```

### Todo File Format

Todos should be in Markdown checkbox format:

```markdown
- [ ] Buy groceries
- [x] Finish project
- [ ] Call dentist
```

## 🏗️ Building for Production

### Standard Build

```bash
go build -ldflags="-s -w" -o todo-app main.go
```

### Optimized Build

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o todo-app main.go
```

### Cross-Platform Builds

Tested only on linux so not sure about the rest but you can try.
```bash

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o todo-app-linux main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o todo-app-macos main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o todo-app-macos-arm main.go

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o todo-app.exe main.go
```

## 📦 Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## 🎨 Customization

The app uses the Catppuccin color scheme. To customize colors, modify the `InitStyle()` function in `main.go`:

```go
func InitStyle() *Styles {
    s := new(Styles)
    s.Boarder = lipgloss.Color(pkg.Catppuccin_base)
    s.CursorStyle = lipgloss.Color(pkg.Catppuccin_lavender)
    // Customize other colors here
    return s
}
```

## 🐛 Troubleshooting

### App doesn't start
- Ensure the notes directory exists and is readable
- Check that the file path in `pkg.Init()` is correct

### Todos not saving
- Verify write permissions for the notes directory
- Make sure to press `u` to update/save changes

### Display issues
- Ensure your terminal supports Unicode and colors
- Try resizing your terminal window


## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [Catppuccin](https://github.com/catppuccin/catppuccin) for the beautiful color scheme
- The Go community for excellent tooling and support

---

Made with ❤️ and Go


# todos
- [ ] Make a config file parse.
- [ ] use Cobra for cmd 
