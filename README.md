# SSH-TUI

A fully interactive Terminal User Interface (TUI) for selecting and connecting to SSH hosts, built with Go and the Bubbletea framework.

## Features

- **Host Discovery**: Automatically parses SSH hosts from `~/.ssh/config` and `~/.ssh/known_hosts`
- **Interactive Host Selection**: Scrollable menu with search/filter functionality
- **SSH Options Entry**: Input custom SSH options and arguments (e.g., `-L 8080:localhost:80`, `-i ~/.ssh/id_rsa`)
- **Command Preview**: Shows the final SSH command before execution with confirm/cancel options
- **Cross-Platform**: Works on Linux, macOS, and Windows with OpenSSH

## Installation

### Prerequisites

- Go 1.19 or later
- SSH client (OpenSSH) installed and available in PATH

### Build from Source

```bash
git clone <repository-url>
cd ssh-tui
go build -o ssh-tui
```

### Install Globally

```bash
go install
```

## Usage

Simply run the command without any arguments:

```bash
./ssh-tui
```

### TUI Flow

1. **Host Selection**: Use arrow keys or `j`/`k` to navigate, `/` to search, `Enter` to select
2. **Options Entry**: Enter SSH options and arguments (optional), press `Enter` to continue
3. **Command Preview**: Review the final SSH command, press `Y` to confirm or `N` to go back
4. **Connection**: SSH connection is established

### Keyboard Shortcuts

#### Host Selection Screen
- `↑`/`↓` or `j`/`k`: Navigate hosts
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

#### Command Preview Screen
- `Y` or `Enter`: Confirm and connect
- `N` or `Esc`: Go back
- `Ctrl+C`: Quit

## Configuration

SSH-TUI reads host information from standard SSH configuration files:

### ~/.ssh/config
```
Host myserver
    HostName example.com
    User myuser
    Port 2222

Host production
    HostName prod.company.com
    User deploy
    IdentityFile ~/.ssh/prod_key
```

### ~/.ssh/known_hosts
SSH-TUI also reads hostnames from your known_hosts file, making previously connected hosts available for selection.

## Examples

### Basic Connection
1. Select a host from the list
2. Press Enter (no additional options)
3. Confirm the SSH command
4. Connect!

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
- **Invalid options**: Validation before command execution

## Project Structure

```
ssh-tui/
├── main.go                    # Main application entry point
├── internal/
│   ├── parser/               # SSH config and known_hosts parsing
│   │   ├── config.go         # SSH config file parser
│   │   ├── known_hosts.go    # known_hosts file parser
│   │   └── discovery.go      # Host discovery and merging
│   ├── tui/                  # Terminal User Interface components
│   │   ├── host_selector.go  # Host selection screen
│   │   ├── options_entry.go  # SSH options entry screen
│   │   └── command_preview.go # Command preview and confirmation
│   └── ssh/                  # SSH execution
│       └── executor.go       # Cross-platform SSH command execution
├── go.mod
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
