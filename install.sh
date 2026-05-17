#!/usr/bin/env bash

set -euo pipefail

REPO_URL="https://github.com/acacMAX/accil.git"
INSTALL_DIR="$HOME/.accil/bin"
CACHE_DIR="${TMPDIR:-/tmp}/accil-cache"
BINARY_NAME="accil"
DEFAULT_VERSION="1.4.6"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMP_DIR=""
SOURCE_DIR=""

info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}[OK]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

fail() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
    exit 1
}

cleanup() {
    if [[ -n "$TEMP_DIR" && -d "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}

trap cleanup EXIT

check_go() {
    command -v go >/dev/null 2>&1 || fail "Go is not installed or not in PATH."
    success "Go detected: $(go version)"
}

resolve_source_dir() {
    if [[ -f "$SCRIPT_DIR/go.mod" ]]; then
        SOURCE_DIR="$SCRIPT_DIR"
        info "Using local source tree: $SOURCE_DIR"
        return
    fi

    command -v git >/dev/null 2>&1 || fail "Git is required when installing without a local source tree."
    TEMP_DIR="$(mktemp -d)"
    info "Downloading ACCIL source package..."
    git clone --depth 1 "$REPO_URL" "$TEMP_DIR" >/dev/null 2>&1 || fail "Failed to download source package."
    SOURCE_DIR="$TEMP_DIR"
    success "Download completed"
}

read_version() {
    if [[ -f "$SOURCE_DIR/VERSION" ]]; then
        tr -d '\r\n' < "$SOURCE_DIR/VERSION"
    else
        printf "%s" "$DEFAULT_VERSION"
    fi
}

build_and_install() {
    local version
    version="$(read_version)"
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$CACHE_DIR/go-build" "$CACHE_DIR/gomod"
    export GOCACHE="$CACHE_DIR/go-build"
    export GOMODCACHE="$CACHE_DIR/gomod"

    info "Building ACCIL v$version..."
    (
        cd "$SOURCE_DIR"
        go build -buildvcs=false -ldflags="-X github.com/accil/accil/cmd.Version=$version" -o "$INSTALL_DIR/$BINARY_NAME" .
    ) || fail "Build failed."

    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    success "Installed to $INSTALL_DIR/$BINARY_NAME"
}

add_to_path() {
    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            success "Install directory already present in PATH"
            return
            ;;
    esac

    local shell_rc=""
    if [[ -n "${ZSH_VERSION:-}" ]]; then
        shell_rc="$HOME/.zshrc"
    elif [[ -n "${BASH_VERSION:-}" ]]; then
        shell_rc="$HOME/.bashrc"
    elif [[ -f "$HOME/.profile" ]]; then
        shell_rc="$HOME/.profile"
    fi

    if [[ -n "$shell_rc" ]]; then
        if ! grep -Fq "$INSTALL_DIR" "$shell_rc" 2>/dev/null; then
            {
                printf "\n# ACCIL\n"
                printf "export PATH=\"\$PATH:%s\"\n" "$INSTALL_DIR"
            } >> "$shell_rc"
        fi
        success "Updated PATH in $shell_rc"
    else
        warn "Could not determine shell rc file automatically."
        warn "Add this directory to PATH manually: $INSTALL_DIR"
    fi

    export PATH="$PATH:$INSTALL_DIR"
}

verify_install() {
    "$INSTALL_DIR/$BINARY_NAME" version >/dev/null 2>&1 || fail "Installed binary verification failed."
    success "Installation verified"
}

main() {
    info "ACCIL installer starting"
    check_go
    resolve_source_dir
    build_and_install
    add_to_path
    verify_install
    printf "\n"
    success "Installation complete"
    printf "Run: %s\n" "$BINARY_NAME"
}

main "$@"
