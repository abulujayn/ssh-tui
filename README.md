# SSH-TUI

A fully interactive Terminal User Interface (TUI) for selecting and connecting to SSH hosts, built with Go and the Bubbletea framework.

**Latest Version:** 1.0.0-beta

## Features

- **Host Discovery**: Automatically parses SSH hosts from `~/.ssh/config` and `~/.ssh/known_hosts`
- **Interactive Host Selection**: Scrollable menu with search/filter functionality
- **SSH Options Entry**: Input custom SSH options and arguments (e.g., `-L 8080:localhost:80`, `-i ~/.ssh/id_rsa`)
- **Command Preview**: Shows the final SSH command before execution when entering custom options, with confirm/cancel options
- **Cross-Platform**: Works on Linux, macOS, and Windows with OpenSSH

## Installation

### Prerequisites

- Go 1.25.1 or later
- SSH client (OpenSSH) installed and available in PATH

## Dependencies

This project uses the following Go modules:

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal app framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output

### Build from Source

```bash
git clone https://github.com/abulujayn/ssh-tui.git
cd ssh-tui
go build ./cmd/ssh-tui
```

### Install Globally

```bash
go install ./cmd/ssh-tui
```

## Usage

Simply run the command without any arguments to start the interactive TUI:

```bash
./ssh-tui
```

### Command-Line Options

- `--help`: Display SSH help information
- `--version`: Show the application version

### Direct SSH Execution

You can also pass SSH arguments directly to bypass the TUI and execute SSH immediately:

```bash
./ssh-tui user@host -p 2222
```

### TUI Flow

1. **Host Selection**: Use arrow keys to navigate, `/` to search, `Enter` to select
2. **Options Entry** (optional): Enter SSH options and arguments, press `Enter` to continue or `Esc` to skip
3. **Command Preview** (if options entered): Review the final SSH command, press `Enter` to confirm or `Esc` to go back
4. **Connection**: SSH connection is established

### Keyboard Shortcuts

#### Host Selection Screen
- `↑`/`↓`: Navigate hosts
- `/`: Toggle search mode
- `Enter`: Select host
- `Esc`: Exit search or quit
- `q`: Quit

#### Options Entry Screen
- `Ctrl+A`: Move to beginning
- `Ctrl+E`: Move to end
- `Ctrl+U`: Clear to beginning
- `Ctrl+K`: Clear to end
- `Ctrl+W`: Delete word backwards
- `Enter`: Continue
- `Esc`: Go back
- `Ctrl+C`: Quit

## Configuration

SSH-TUI reads host information from standard SSH configuration files:
- `~/.ssh/config`
- `~/.ssh/known_hosts`

## Examples

### Basic Connection
1. Select a host from the list
2. Press Enter (no additional options)
3. Connect directly!

### Port Forwarding
1. Select a host
2. Enter options: `-L 8080:localhost:80`
3. Review command: `ssh -L 8080:localhost:80 user@example.com`
4. Confirm and connect

### Custom Identity File
1. Select a host  
2. Enter options: `-i ~/.ssh/custom_key`
3. Review and confirm

## Error Handling

SSH-TUI provides helpful error messages for common issues:

- **No hosts found**: Guidance on setting up SSH config files
- **SSH not available**: Instructions for installing OpenSSH
- **Connection failures**: Display SSH error messages
- **Invalid options**: Validation before command execution (prevents potentially dangerous characters like `;` or `|`)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

The project includes comprehensive unit tests covering key components. Run `go test ./...` to execute all tests.

## License

MIT License - see LICENSE.md file for details.