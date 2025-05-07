# Tox - Simple CLI Todo Manager

Tox is a command-line todo management tool built with Go, Cobra, and SQLite.

## Features

Tox provides a simple command-line interface for managing your todos. Below are the available commands and their descriptions:

| Commands               | Description                                      |
|------------------------|--------------------------------------------------|
| `add`                  | Add a new todo item to the list                  |
| `list`                 | Display all pending or completed todos           |
| `done`                 | Mark a specific todo item as completed           |
| `delete`               | Remove a todo item from the list                 |
| `reindex`              | Reset and reorder all todo IDs to be sequential  |

| Flags                  | Description                                      |
|------------------------|--------------------------------------------------|
| `--all`, `-a`          | Show all todos, including completed ones         |
| `--help`, `-h`         | Show help information for a command              |
| `--version`, `-v`      | Show the version of the application              |

## Installation

### Prerequisites

- Go 1.18 or later
- SQLite

### Building from source

1. Clone the repository:

  ```bash
  git clone https://github.com/yourusername/tox.git
  cd tox
  ```

2. Build the application:

  ```bash
  go build -o tox
  ```

3. Add the binary to your PATH or move it to a location in your PATH:

  ```bash
  # Option 1: Move to /usr/local/bin (may require sudo)
  sudo mv tox /usr/local/bin/

  # Option 2: Move to ~/go/bin (if in your PATH)
  mv tox ~/go/bin/
  ```

## Usage

### Add a todo

```bash
tox add Buy groceries
tox add "Call mom at 5pm"
```

### List todos

```bash
# List only pending todos (default)
tox list

# List all todos including completed ones
tox list --all
tox list -a
```

### Mark a todo as done

```bash
tox done 1
```

### Delete a todo

```bash
tox delete 1
# Or use aliases
tox del 1
tox rm 1
```

### Reindex todos

```bash
tox reindex
```

This will reset and reorder all todo IDs to be sequential starting from 1. This happens automatically when you delete a todo, but you can also run it manually if needed.

## Data Storage

Your todos are stored in an SQLite database located at `~/.tox/todos.db`. You can back up or move this file to transfer your todos to another machine. The database is created automatically when you first run the application.

## To contribute

This project is open for contributions! If you have suggestions, bug reports, or feature requests, please open an issue or submit a pull request.

### Steps to contribute

1. Fork the repository
2. Create a new branch for your feature or bug fix:

   ```bash
   git checkout -b feature/my-feature
   ```

3. Make your changes and commit them:

   ```bash
   git commit -m "Add my feature"
   ```

4. Push to your forked repository:

   ```bash
   git push origin feature/my-feature
   ```

5. Open a pull request against the main repository

## License

Go through the license [MIT](LICENSE)
