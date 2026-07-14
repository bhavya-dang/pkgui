<p align="center">
  <h1 align="center">pkgui</h1>
  <p align="center">A terminal UI for browsing packages across multiple package managers.</p>

  <p align="center">
    <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go version"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=for-the-badge" alt="MIT license"></a>
    <a href="https://github.com/bhavya-dang/pkgui/releases"><img src="https://img.shields.io/github/v/release/bhavya-dang/pkgui?style=for-the-badge&logo=github&logoColor=white&label=release" alt="Release"></a>
    <a href="https://github.com/bhavya-dang/pkgui/stargazers"><img src="https://img.shields.io/github/stars/bhavya-dang/pkgui?style=for-the-badge&logo=github&logoColor=white&label=stars" alt="Stars"></a>
  </p>

  <p align="center">
    <i>Built with <a href="https://github.com/charmbracelet/bubbletea">Bubble Tea</a></i>
  </p>

  <p align="center">
    <img src="./preview/demo2.gif" alt="demo">
  </p>
</p>

## Features

- List installed packages
- Fuzzy search installed packages (`/` to search)
- View package details:
  - version (installed/latest)
  - description
  - homepage
  - license
  - dependencies
  - installation path
  - binary size
- Scrollable package list with keyboard navigation
- Switch themes (5 theme palettes to select from)
  - Default theme is [Solace](https://github.com/bhavya-dang/Solace)

## Roadmap

- [x] installed formulae
- [ ] installed casks/taps ([#29](https://github.com/bhavya-dang/pkgui/issues/29))
- [x] installed npm packages
- [/] installed pip packages ([#20](https://github.com/bhavya-dang/pkgui/issues/20))
  - [x] globally installed packages
- [ ] upgrade/remove packages ([#3](https://github.com/bhavya-dang/pkgui/issues/3))
- [ ] search packages ([#19](https://github.com/bhavya-dang/pkgui/issues/19))
- [x] multi-theme support ([#9](https://github.com/bhavya-dang/pkgui/issues/9))
- [ ] persist user configuration ([#27](https://github.com/bhavya-dang/pkgui/issues/27))
- [x] list packages across all PMs in one list

## Currently Supported PMs

- **Homebrew** (formulae)
- **npm**
- **pip** (global packages)

## Prerequisites

- Go 1.25+ (if building from source)

## Installation

### Using Go

```bash
go install github.com/bhavya-dang/pkgui@latest
```

### Using install.sh

```bash
curl -sSL https://raw.githubusercontent.com/bhavya-dang/pkgui/refs/heads/master/install.sh | sh
```

### Using Makefile

```bash
git clone https://github.com/bhavya-dang/pkgui.git
cd pkgui
make install
```

### Manual

```bash
git clone https://github.com/bhavya-dang/pkgui.git
cd pkgui
go build -o build/pkgui .
cp build/pkgui "$GOPATH/bin/pkgui"
```

## Usage

```bash
pkgui
```

### Keybindings

| Key            | Action                                 |
| -------------- | -------------------------------------- |
| `↑` / `↓`      | Navigate package list                  |
| `←` / `→`      | Switch between package managers (tabs) |
| `/`            | Toggle search (type to filter)         |
| `t`            | Open theme selector                    |
| `Esc`          | Exit search / close overlay            |
| `q` / `Ctrl+C` | Quit                                   |

## Support

- Homebrew
  - installed formulae (with detail view from the Homebrew API)
- npm
  - installed packages
- pip
  - globally installed packages

## License

MIT

## Contributions

This is an open-source project. Feel free to raise any issues you find or contribute something.
