#!/usr/bin/env sh
set -eu

REPO="bhavyadang/pkgui"
BINARY_NAME="pkgui"
VERSION="${VERSION:-latest}"

# Parse args
BUILD=false
for arg in "$@"; do
  case "$arg" in
    --build) BUILD=true ;;
  esac
done

# OS detection
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

# Download from GitHub Releases

download() {
  if [ "$VERSION" = "latest" ]; then
    TAG=$(curl -sfL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
    if [ -z "$TAG" ]; then
      echo "error: could not determine latest release" >&2
      return 1
    fi
    VERSION="$TAG"
  fi

  TARBALL="${BINARY_NAME}_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
  URL="https://github.com/$REPO/releases/download/$VERSION/$TARBALL"

  echo "downloading $URL"
  mkdir -p build
  curl -sfL "$URL" -o "build/$TARBALL" || return 1

  echo "extracting $TARBALL"
  tar -xzf "build/$TARBALL" -C build/

  if [ ! -f "build/$BINARY_NAME" ]; then
    echo "error: binary not found in archive" >&2
    return 1
  fi

  chmod +x "build/$BINARY_NAME"
  echo "installing to $INSTALL_DIR/$BINARY_NAME"
  mkdir -p "$INSTALL_DIR"
  cp "build/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
  echo "done."
}

# Build from source

build_from_source() {
  echo "building $BINARY_NAME for $GOOS/$GOARCH"
  mkdir -p build
  GOOS="$GOOS" GOARCH="$GOARCH" go build -o "build/$BINARY_NAME" .

  echo "installing to $INSTALL_DIR/$BINARY_NAME"
  mkdir -p "$INSTALL_DIR"
  cp "build/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
  echo "done."
}

# Main

if [ "$BUILD" = false ] && command -v curl >/dev/null 2>&1; then
  if ! download; then
    echo "download failed, falling back to building from source" >&2
    if command -v go >/dev/null 2>&1; then
      build_from_source
    else
      echo "error: Go is not installed. install Go or use --build" >&2
      exit 1
    fi
  fi
else
  if ! command -v go >/dev/null 2>&1; then
    echo "error: Go is required to build from source" >&2
    exit 1
  fi
  build_from_source
fi

# check PATH

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) echo "warning: $INSTALL_DIR is not on your PATH. add it to your shell rc file:" >&2
     echo "  export PATH=\"\$PATH:$INSTALL_DIR\"" >&2 ;;
esac
