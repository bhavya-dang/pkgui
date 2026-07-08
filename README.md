# pkgui

A terminal UI for managing packages across multiple package managers.

> Written in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea)

<br/>

<<<<<<< Updated upstream
![Screenshot 1](./preview/2.png)
=======
<video src="https://github.com/bhavya-dang/pkgui/blob/master/preview/demo.mp4?raw=true" width="100%" controls></video>
>>>>>>> Stashed changes

## Features

- List installed packages
- Fuzzy search installed packages (`/` to search)
- View package details: version (installed/latest), description, homepage, license, dependencies, installation path
- Scrollable package list with keyboard navigation

## Roadmap

- [x] installed formulae
- [ ] installed casks/taps
- [x] installed npm packages
- [ ] installed pip packages
- [ ] upgrade/remove packages
- [ ] search packages

## Currently Supported PMs

- **Homebrew** (formulae)
- **npm**

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
curl -sSL https://raw.githubusercontent.com/bhavya-dang/pkgui/refs/heads/master/install.sh | sh
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

| Key            | Action                         |
| -------------- | ------------------------------ |
| `↑` / `↓`      | Navigate package list          |
| `/`            | Toggle search (type to filter) |
| `Esc`          | Exit search                    |
| `q` / `Ctrl+C` | Quit                           |

## Support

- Homebrew
  - installed formulae (with detail view from the Homebrew API)
- npm
  - installed packages

## License

MIT

## Contributions

I am actively working on this project. Feel free to raise any issues you find.
If you want to contribute something, let me know or raise an issue, fork the repo, and start contributing!
