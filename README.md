# Vista

**VFS Interactive Shell Terminal Application**

A beautiful, interactive terminal user interface for managing and exploring virtual filesystems powered by the [VFS library](https://github.com/mwantia/vfs).

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) 🧋

## Features

### 🗂️ Filesystem Browser
- Visual tree/list view of mounted filesystems
- Navigate directories with intuitive keyboard shortcuts
- View file metadata (size, permissions, timestamps)
- Syntax-highlighted file preview with [Chroma](https://github.com/alecthomas/chroma)
- Image preview support with [PixTerm](https://github.com/eliukblau/pixterm)

### 🔌 Mount Management
- Interactive mount creation (select driver, configure URI)
- Support for multiple storage backends:
  - Ephemeral (in-memory)
  - SQLite (local database)
  - Consul (distributed key-value store)
  - PostgreSQL (relational database)
  - S3 (object storage)
  - File (local filesystem)
- Unmount with safety checks
- Display mount hierarchy and relationships
- Show driver capabilities and constraints

### 🎨 Customization
- Tokyo Night color scheme (built-in)
- Extensible theme system
- Custom key bindings
- Responsive layout

## Installation

### Prerequisites
- Go 1.25.3 or later
- [VFS library](https://github.com/mwantia/vfs) (automatically installed via go.mod)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mwantia/vista.git
cd vista

# Build the binary
go build -o vista ./cmd/vista

# Run
./vista
```

### Using Task (recommended)

```bash
# Install Task: https://taskfile.dev/

# Build
task build

# Run
task run
```

## Usage

### Basic Commands

```bash
# Start Vista
vista

# Show version
vista version

# Show help
vista --help
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `←/h` | Go to parent directory |
| `→/l` or `Enter` | Open directory/file |
| `/` | Search/filter |
| `m` | Mount management |
| `?` | Show help |
| `q` or `Ctrl+C` | Quit |

## Architecture

Vista follows a clean scene-based architecture:

```
vista/
├── cmd/vista/           # CLI entry point
│   ├── main.go
│   └── cli/            # Cobra command definitions
├── internal/
│   ├── app/            # Main Bubble Tea application model
│   ├── scenes/         # UI scenes (browser, mounts, settings)
│   ├── vfs/            # VFS integration layer
│   ├── preview/        # File preview with syntax highlighting
│   └── resources/      # Theme, icons, key bindings
└── taskfile.yml        # Task automation
```

### Integration with VFS

Vista is a **client** of the VFS library and uses its public API:

```go
// Initialize VFS
vfs := pkg.NewVirtualFileSystem()

// Mount root
vfs.Mount(ctx, "/", mount.WithObjectStorage("ephemeral://"))

// Perform operations
entries, _ := vfs.ReadDirectory(ctx, "/data")
```

For VFS architecture details, see the [VFS repository](https://github.com/mwantia/vfs).

## Development

### Project Structure

- **`cmd/vista/`** - CLI application entry point with Cobra
- **`internal/app/`** - Main Bubble Tea model and message handling
- **`internal/scenes/`** - Different UI views (browser, mounts, settings)
- **`internal/vfs/`** - VFS manager and operation wrappers
- **`internal/preview/`** - File content preview with syntax highlighting
- **`internal/resources/`** - Reusable UI components (theme, icons, keymaps)

### Development Workflow

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Lint
go vet ./...

# Build
go build -o vista ./cmd/vista
```

### Adding New Scenes

1. Create scene package in `internal/scenes/`
2. Implement Bubble Tea's `tea.Model` interface
3. Register scene in main app model
4. Add navigation transitions

## Dependencies

### Primary Dependencies
- **[VFS](https://github.com/mwantia/vfs)** - Virtual filesystem library
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Style definitions
- **[Bubbles](https://github.com/charmbracelet/bubbles)** - TUI components
- **[Chroma](https://github.com/alecthomas/chroma)** - Syntax highlighting
- **[Cobra](https://github.com/spf13/cobra)** - CLI framework

See `go.mod` for complete dependency list.

## Related Projects

- **[VFS](https://github.com/mwantia/vfs)** - The underlying virtual filesystem library

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea)
