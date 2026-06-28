#!/usr/bin/env sh
set -eu

BINARY_NAME="pkgui"

# ── OS detection ──────────────────────────────────────────────────────────
case "$(uname -s)" in
  Darwin)  GOOS="darwin"  ;;
  Linux)   GOOS="linux"   ;;
  CYGWIN*|MINGW*|MSYS*)  GOOS="windows"  ;;
  *)       echo "error: unsupported OS '$(uname -s)'" >&2; exit 1 ;;
esac

# ── Arch detection ────────────────────────────────────────────────────────
case "$(uname -m)" in
  x86_64|amd64) GOARCH="amd64" ;;
  aarch64|arm64) GOARCH="arm64" ;;
  armv7l|armv6l|arm) GOARCH="arm" ;;
  i686|i386)     GOARCH="386" ;;
  *) echo "error: unsupported arch '$(uname -m)'" >&2; exit 1 ;;
esac

# ── Determine install directory ───────────────────────────────────────────
if [ -n "${GOBIN:-}" ]; then
    INSTALL_DIR="$GOBIN"
elif [ -n "${GOPATH:-}" ]; then
    INSTALL_DIR="$GOPATH/bin"
else
    INSTALL_DIR="$HOME/go/bin"
fi

# ── Build ─────────────────────────────────────────────────────────────────
echo "building $BINARY_NAME for $GOOS/$GOARCH"
mkdir -p build
GOOS="$GOOS" GOARCH="$GOARCH" go build -o "build/$BINARY_NAME" .

# ── Install ───────────────────────────────────────────────────────────────
echo "installing to $INSTALL_DIR/$BINARY_NAME"
mkdir -p "$INSTALL_DIR"
cp "build/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

# ── PATH check ────────────────────────────────────────────────────────────
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "warning: $INSTALL_DIR is not on your PATH. add it to your shell rc file:" >&2
     echo "  export PATH=\"\$PATH:$INSTALL_DIR\"" >&2 ;;
esac

echo "done."
