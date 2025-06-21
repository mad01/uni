# uni

A super minimal task management CLI written in Go using the spf13/cobra framework.

## Features

- **Simple task management**: Create, list, and manage tasks with different statuses
- **Flexible storage**: Stores tasks in `~/.uni` by default, or `.uni` directory in git repositories
- **Multiple output formats**: Support for normal, text, json, and yaml output formats
- **Status tracking**: Tasks can be open, working, blocked, done, or cancelled
- **Auto-incrementing IDs**: Each task gets a unique incrementing ID
- **Short command aliases**: All commands have short aliases (a, l, b, d, c, w, e, h)
- **Advanced filtering**: Filter tasks by status (--left for active, --closed for completed)
- **Interactive editing**: Edit tasks using your preferred editor via EDITOR environment variable
- **Comprehensive testing**: Full unit test coverage for all core functionality

## Installation

### Option 1: Go Install (Recommended)
```bash
go install github.com/mad01/uni@latest
```

### Option 2: Download Binary
Download the latest release for your platform from the [releases page](https://github.com/mad01/uni/releases).

**Available platforms:**
- macOS (Intel): `uni_Darwin_x86_64.tar.gz`
- macOS (Apple Silicon): `uni_Darwin_arm64.tar.gz`
- Linux (Intel): `uni_Linux_x86_64.tar.gz`
- Linux (ARM64): `uni_Linux_arm64.tar.gz`

**Install downloaded binary:**
```bash
# Extract the binary
tar -xzf uni_*.tar.gz

# Move to PATH
sudo mv uni /usr/local/bin/

# Make executable (if needed)
chmod +x /usr/local/bin/uni
```

### Option 3: Build from Source
```bash
git clone https://github.com/mad01/uni.git
cd uni

# Simple build
go build -o uni

# Build with version info (recommended for development)
make build

# See all build options
make help
```

## Usage

### Basic Commands

```bash
# Add a new task (using flags)
uni add --name "Fix database connection" --description "Connect to production database"
# Or with short flags and alias
uni a -n "Fix database connection" -d "Connect to production database"

# List all tasks
uni list
# Or with short alias
uni l

# Get a specific task
uni get 1

# Change task status
uni working 1    # Mark task as working (or: uni w 1)
uni blocked 1    # Mark task as blocked (or: uni b 1)
uni done 1       # Mark task as done (or: uni d 1)
uni cancel 1     # Mark task as cancelled (or: uni c 1)

# Edit a task interactively
uni edit 1       # Opens task in your $EDITOR (or: uni e 1)
```

### Output Formats & Filtering

```bash
# Default format (colorized)
uni list

# Text format (tabular)
uni list -o text

# JSON format
uni list -o json

# YAML format
uni list -o yaml

# Filter by status
uni list --left      # Show only active tasks (open, working, blocked)
uni list --closed    # Show only completed tasks (done, cancelled)
uni l --left -o json  # Combine filtering with output format
```

## Task Storage

- **Default**: Tasks are stored in `~/.uni/tasks.json`
- **Git repositories**: If you're in a git repository and a `.uni` directory exists, tasks are stored locally in `.uni/tasks.json`

This allows you to have project-specific tasks that can be checked into version control if desired.

## Task Structure

Each task has the following fields:
- `id`: Auto-incrementing unique identifier
- `name`: Task name
- `description`: Optional task description
- `status`: One of `open`, `working`, `blocked`, `done`, `cancel`
- `created_at`: Task creation timestamp
- `updated_at`: Last update timestamp

## Commands

### Task Management
- `uni add` (`a`) - Add a new task using `--name/-n` and `--description/-d` flags
- `uni list` (`l`) - List all tasks
- `uni get <id>` - Get a specific task
- `uni edit <id>` (`e`) - Edit a task using your default editor

### Status Changes
- `uni working <id>` (`w`) - Mark task as working
- `uni blocked <id>` (`b`) - Mark task as blocked
- `uni done <id>` (`d`) - Mark task as done
- `uni cancel <id>` (`c`) - Mark task as cancelled

## Global Flags

- `-o, --output`: Output format (normal, text, json, yaml)
- `--left`: Show only active tasks (open, working, blocked)
- `--closed`: Show only completed tasks (done, cancelled)

## Testing

Run the test suite:
```bash
go test ./...
```

The project includes comprehensive unit tests covering:
- Task creation and management
- Status updates
- Filtering functionality
- Data persistence
- Edge cases and error handling

## CI/CD

The project uses GitHub Actions for continuous integration and automated releases:

- **Continuous Integration**: Runs tests, vet, and format checks on every push
- **Automated Releases**: Creates releases with cross-platform binaries when tags are pushed
- **Snapshot Builds**: Creates development builds on pushes to main branch

### Creating a Release

The project includes a built-in version bump command:

```bash
# Bump patch version (v1.0.0 -> v1.0.1)
uni bump patch

# Bump minor version (v1.0.1 -> v1.1.0)  
uni bump minor

# Bump major version (v1.1.0 -> v2.0.0)
uni bump major

# Preview what would happen
uni bump patch --dry-run
```

Or manually create a release:

```bash
# Tag a new version
git tag v1.0.0

# Push the tag
git push origin v1.0.0
```

This will automatically trigger the release workflow and create:
- GitHub release with changelog
- Cross-platform binaries (macOS/Linux, AMD64/ARM64)
- Checksums for verification

### Development Builds

For local development with proper version info:

```bash
# Build with git hash and version info
make build

# Check version info before building
make version

# Development workflow
make dev  # lint, test, and build
``` 