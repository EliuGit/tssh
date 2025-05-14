# tssh - Terminal SSH Connection Manager

tssh is a terminal-based SSH connection manager written in Go, featuring an interactive TUI for managing and connecting to SSH servers.

## Features

- 🖥️ Interactive terminal user interface (TUI)
- 🔒 Secure storage of connection details
  - Passwords are encrypted before storage
- 🔑 Supports both password and SSH key authentication
- 📋 Easy management of SSH connections
  - Add, edit, delete connections
  - Quick connect to saved servers
- 🔍 Filter connections by name/host
- 🛠️ Simple configuration in `~/.tssh/`

## Installation

### Prerequisites

- Go 1.16 or later
- `sshpass` (for password authentication)
  - On Ubuntu/Debian: `sudo apt install sshpass`
  - On macOS: `brew install hudochenkov/sshpass/sshpass`

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/tssh.git
   cd tssh
   ```

2. Build and install:
   ```bash
   go build -o tssh
   ```

3. Run the program:
   ```bash
   ./tssh
   ```

## Usage

After starting tssh, you'll see the main interface with your saved connections.

### Key Bindings

| Key       | Action               |
|-----------|----------------------|
| `Enter`   | Connect to server    |
| `p`       | Connect to sftpserver|
| `a`       | Add new connection   |
| `e`       | Edit connection      |
| `d`       | Delete connection    |
| `/`       | Filter connections   |
| `Esc`     | Cancel filter        |
| `q`       | Quit program        |

### Adding a Connection

1. Press `a` to add a new connection
2. Fill in the connection details:
   - Name: Descriptive name for the connection
   - Host: Server hostname or IP
   - Port: SSH port (default: 22)
   - Username: Login username
   - Authentication method: Password or SSH key
3. Press `Enter` to save

## Configuration

tssh stores its configuration and database in `~/.tssh/`:
- `connections.db` - SQLite database with connection info
- Configuration files (if any) will be stored here

## Security

- Passwords are encrypted before being stored in the database
- SSH keys are not stored - only paths to key files are saved
- Database file has restricted permissions (0600)

## License

MIT
