#!/bin/bash

GITHUB_USERNAME=erokez0
APP_NAME=sweep

set -e

echo "installing $APP_NAME"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
linux*)
  OS="linux"
  ;;
darwin*)
  OS="darwin"
  ;;
msys* | cygwin* | mingw*)
  OS="windows"
  ;;
freebsd*)
  OS="freebsd"
  ;;
openbsd*)
  OS="openbsd"
  ;;
netbsd*)
  OS="netbsd"
  ;;
*)
  echo "Unsupported OS: $OS"
  exit 1
  ;;
esac

case $ARCH in
x86_64 | x64)
  ARCH="amd64"
  ;;
aarch64 | arm64)
  ARCH="arm64"
  ;;
  # armv7l|armv8l)
  # ARCH="arm"
  # ;;
  # i386|i686|x86)
  # ARCH="386"
  # ;;
*)
  echo "Unsupported architecture: $ARCH"
  exit 1
  ;;
esac

echo "Setting up directories..."

try_system_install() {
  local test_file="/usr/local/bin/.${APP_NAME}_test_$$"

  if touch "$test_file" 2>/dev/null; then
    rm -f "$test_file"
    return 0
  fi

  if sudo touch "$test_file" 2>/dev/null; then
    sudo rm -f "$test_file"
    return 0
  fi
  return 1
}

if [[ "$OS" == "windows" ]]; then
  BIN_DIR="${APPDATA:-$HOME/AppData/Roaming}/Programs/$APP_NAME"
  CONFIG_DIR="${APPDATA:-$HOME/AppData/Roaming}/$APP_NAME"
  EXT=".exe"
  USE_SUDO=""
else
  EXT=""

  if try_system_install; then
    BIN_DIR="/usr/local/bin"
    CONFIG_DIR="$HOME/.config/$APP_NAME"
    USE_SUDO="sudo"
  else
    BIN_DIR="$HOME/.local/bin"
    CONFIG_DIR="$HOME/.config/$APP_NAME"
    USE_SUDO=""
  fi
fi

BINARY_NAME="$APP_NAME$EXT"

echo "Binary directory: $BIN_DIR"
echo "Config directory: $CONFIG_DIR"

if [[ -n "$USE_SUDO" ]] && [[ "$BIN_DIR" == "/usr/local/bin" ]]; then
  mkdir -p "$CONFIG_DIR"
  chmod 755 "$CONFIG_DIR"
else
  mkdir -p "$BIN_DIR" "$CONFIG_DIR"
fi

REPO_URL="https://github.com/$GITHUB_USERNAME/$APP_NAME"
VERSION="${VERSION:-latest}"

if [[ "$VERSION" == "latest" ]]; then
  LATEST_TAG=$(curl -s "https://api.github.com/repos/$GITHUB_USERNAME/$APP_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "v1.0.0")
  VERSION="$LATEST_TAG"
fi

BINARY_URL="$REPO_URL/releases/download/$VERSION/${APP_NAME}-${OS}-${ARCH}${EXT}"

TEMP_DIR=$(mktemp -d)
TEMP_BINARY="$TEMP_DIR/$BINARY_NAME"

if command -v curl &>/dev/null; then
  curl -L -f "$BINARY_URL" -o "$TEMP_BINARY" || {
    echo "Failed to download binary"
    exit 1
  }
elif command -v wget &>/dev/null; then
  wget -q "$BINARY_URL" -O "$TEMP_BINARY" || {
    echo "Failed to download binary"
    exit 1
  }
else
  echo "Need curl or wget to download files"
  exit 1
fi

if [[ ! -f "$TEMP_BINARY" ]]; then
  echo "Download failed - file not found"
  exit 1
fi

FILE_SIZE=$(stat -f%z "$TEMP_BINARY" 2>/dev/null || stat -c%s "$TEMP_BINARY" 2>/dev/null)
if [[ $FILE_SIZE -lt 1000 ]]; then
  echo "Download failed - file too small ($FILE_SIZE bytes)"
  exit 1
fi

echo "Binary downloaded ($((FILE_SIZE / 1024)) KB)"

INSTALL_PATH="$BIN_DIR/$BINARY_NAME"

if [[ -n "$USE_SUDO" ]]; then
  sudo cp "$TEMP_BINARY" "$INSTALL_PATH"
  sudo chmod +x "$INSTALL_PATH" 2>/dev/null || true
else
  cp "$TEMP_BINARY" "$INSTALL_PATH"
  chmod +x "$INSTALL_PATH" 2>/dev/null || true
fi

rm -rf "$TEMP_DIR"

CONFIG_BASE_URL="$REPO_URL/raw/master/"

CONFIG_FILES=("config.schema.json" "config.default.json")

for file in "${CONFIG_FILES[@]}"; do
  CONFIG_URL="$CONFIG_BASE_URL/$file"
  CONFIG_PATH="$CONFIG_DIR/$file"

  if [[ ! -f "$CONFIG_PATH" ]]; then
    if command -v curl &>/dev/null; then
      curl -L -f "$CONFIG_URL" -o "$CONFIG_PATH" 2>/dev/null || {
        echo "Could not download $file, creating empty"
        touch "$CONFIG_PATH"
      }
    elif command -v wget &>/dev/null; then
      wget -q "$CONFIG_URL" -O "$CONFIG_PATH" 2>/dev/null || {
        echo "Could not download $file, creating empty"
        touch "$CONFIG_PATH"
      }
    fi

    if [[ "$OS" != "windows" ]]; then
      chmod 600 "$CONFIG_PATH"
    fi
  fi
done

DEFAULT_CONFIG="$CONFIG_DIR/config.default.json"
USER_CONFIG="$CONFIG_DIR/config.json"

if [[ -f "$DEFAULT_CONFIG" ]] && [[ ! -f "$USER_CONFIG" ]]; then
  cp "$DEFAULT_CONFIG" "$USER_CONFIG"

  if [[ "$OS" != "windows" ]]; then
    chmod 600 "$USER_CONFIG"
  fi
fi

echo "Setting up PATH"

add_to_path() {
  local bin_dir="$1"

  if [[ "$OS" == "windows" ]]; then
    if [[ "$PATH" != *"$bin_dir"* ]]; then
      echo "Add this to your PATH manually:"
      echo "  $bin_dir"
    fi
    return
  fi

  if [[ ":$PATH:" != *":$bin_dir:"* ]]; then
    local shell_rc

    case "$SHELL" in
    */bash)
      shell_rc="$HOME/.bashrc"
      ;;
    */zsh)
      shell_rc="$HOME/.zshrc"
      ;;
    */fish)
      shell_rc="$HOME/.config/fish/config.fish"
      ;;
    *)
      shell_rc="$HOME/.profile"
      ;;
    esac

    if [[ -f "$shell_rc" ]]; then
      if ! grep -q "$bin_dir" "$shell_rc" 2>/dev/null; then
        echo "" >>"$shell_rc"
        echo "# Added by $APP_NAME installer" >>"$shell_rc"
        echo "export PATH=\"\$PATH:$bin_dir\"" >>"$shell_rc"
        echo "Added $bin_dir to PATH in $shell_rc"
        echo "Run 'source $shell_rc' or restart your terminal"
      else
        echo "$bin_dir is already in your PATH"
      fi
    else
      echo "Add this to your shell configuration:"
      echo "export PATH=\"\$PATH:$bin_dir\""
    fi
  else
    echo "$bin_dir is already in your PATH"
  fi
}

if [[ "$BIN_DIR" != "/usr/local/bin" ]] && [[ "$BIN_DIR" != "/usr/bin" ]]; then
  add_to_path "$BIN_DIR"
fi

echo "Installation Complete"

echo "If you can't run it immediately, restart your terminal or run:"
if [[ "$SHELL" == *"bash" ]]; then
  echo "  source ~/.bashrc"
elif [[ "$SHELL" == *"zsh" ]]; then
  echo "  source ~/.zshrc"
fi

