# pkgui

A terminal UI for managing packages across multiple package managers.

> Written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea)

<br/>

![Screenshot 1](./preview/2.png)

## Features

- List installed packages
- Fuzzy search installed packages (`/` to search)
- View package details: version (installed/latest), description, homepage, license, dependencies, installation path
- Scrollable package list with keyboard navigation

## Roadmap

- Cask support
- npm / yarn / pnpm
- pip
- upgrade/remove packages

## Currently Supported PMs

- **Homebrew** (formulae)

## Prerequisites

- [Homebrew](https://brew.sh)
- Go 1.25+ (if building from source)

## Installation

### Using Homebrew

```bash
brew install bhavya-dang/pkgui/pkgui
```

### Using Go

```bash
go install github.com/bhavyadang/pkgui@latest
```

### Using install.sh

```bash
curl -sSL https://raw.githubusercontent.com/bhavyadang/pkgui/main/install.sh | sh
```

### Using Makefile

```bash
git clone https://github.com/bhavyadang/pkgui.git
cd pkgui
make install
```

### Manual

```bash
git clone https://github.com/bhavyadang/pkgui.git
cd pkgui
go build -o build/pkgui .
cp build/pkgui "$GOPATH/bin/pkgui"
```

## Usage

```bash
pkgui
```

### Keybindings

| Key            | Action                        |
| -------------- | ----------------------------- |
| `↑` / `↓`      | Navigate package list         |
| `/`            | Start search (type to filter) |
| `Esc`          | Exit search                   |
| `Enter`        | Confirm search                |
| `q` / `Ctrl+C` | Quit                          |

## Support

- Homebrew
  - formulae (with detail view from the Homebrew API)

## License

MIT
