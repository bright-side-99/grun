#!/usr/bin/env bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Installation paths
BINARY_NAME="grun"
INSTALL_DIR_USER="$HOME/.local/bin"
INSTALL_DIR_SYSTEM="/usr/local/bin"

# Determine install directory (prefer user-local, fallback to system)
if [ -w "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
    INSTALL_DIR="$INSTALL_DIR_USER"
    USE_SUDO=false
else
    INSTALL_DIR="$INSTALL_DIR_SYSTEM"
    USE_SUDO=true
fi

INSTALL_PATH="$INSTALL_DIR/$BINARY_NAME"

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if binary is already installed
is_installed() {
    [ -f "$INSTALL_PATH" ]
}

# Print usage
usage() {
    echo "Usage: $0 {install|uninstall}"
    echo ""
    echo "  install   - Build and install $BINARY_NAME to $INSTALL_DIR"
    echo "  uninstall - Remove $BINARY_NAME from $INSTALL_DIR"
    exit 1
}

# Install function
install() {
    echo -e "${BLUE}Installing $BINARY_NAME...${NC}"
    
    # Check if Go is installed
    if ! command_exists go; then
        echo -e "${RED}Error: Go is not installed.${NC}"
        echo "Please install Go first: https://golang.org/doc/install"
        exit 1
    fi
    
    echo -e "${YELLOW}Go version: $(go version)${NC}"
    
    # Build the binary
    echo -e "${YELLOW}Building $BINARY_NAME...${NC}"
    if ! go build -o "$BINARY_NAME" ./main.go; then
        echo -e "${RED}Error: Build failed${NC}"
        exit 1
    fi
    
    # Ensure install directory exists
    if [ "$USE_SUDO" = true ]; then
        echo -e "${YELLOW}Creating directory $INSTALL_DIR (requires sudo)...${NC}"
        sudo mkdir -p "$INSTALL_DIR"
    else
        mkdir -p "$INSTALL_DIR"
    fi
    
    # Install the binary
    echo -e "${YELLOW}Installing to $INSTALL_PATH...${NC}"
    if [ "$USE_SUDO" = true ]; then
        sudo cp "$BINARY_NAME" "$INSTALL_PATH"
        sudo chmod +x "$INSTALL_PATH"
    else
        cp "$BINARY_NAME" "$INSTALL_PATH"
        chmod +x "$INSTALL_PATH"
    fi
    
    # Clean up local binary
    rm -f "$BINARY_NAME"
    
    # Verify installation
    if [ -f "$INSTALL_PATH" ]; then
        echo -e "${GREEN}✓ Successfully installed $BINARY_NAME${NC}"
        echo ""
        echo -e "${BLUE}Installation location: $INSTALL_PATH${NC}"
        
        # Check if it's in PATH and add if needed
        if echo "$PATH" | grep -q "$INSTALL_DIR"; then
            echo -e "${GREEN}✓ $INSTALL_DIR is already in your PATH${NC}"
        else
            # Detect shell and profile file
            SHELL_NAME=$(basename "$SHELL" 2>/dev/null || echo "bash")
            PROFILE_FILE=""
            
            case "$SHELL_NAME" in
                zsh)
                    PROFILE_FILE="$HOME/.zshrc"
                    ;;
                bash)
                    PROFILE_FILE="$HOME/.bashrc"
                    ;;
                *)
                    # Try common profile files
                    if [ -f "$HOME/.zshrc" ]; then
                        PROFILE_FILE="$HOME/.zshrc"
                    elif [ -f "$HOME/.bashrc" ]; then
                        PROFILE_FILE="$HOME/.bashrc"
                    elif [ -f "$HOME/.profile" ]; then
                        PROFILE_FILE="$HOME/.profile"
                    fi
                    ;;
            esac
            
            if [ -n "$PROFILE_FILE" ] && [ -w "$PROFILE_FILE" ]; then
                # Check if PATH export already exists
                PATH_EXPORT="export PATH=\"\$HOME/.local/bin:\$PATH\""
                if ! grep -qF ".local/bin" "$PROFILE_FILE" 2>/dev/null; then
                    echo -e "${YELLOW}Adding $INSTALL_DIR to PATH in $PROFILE_FILE...${NC}"
                    echo "" >> "$PROFILE_FILE"
                    echo "# Added by grun setup.sh" >> "$PROFILE_FILE"
                    echo "$PATH_EXPORT" >> "$PROFILE_FILE"
                    echo -e "${GREEN}✓ Added to $PROFILE_FILE${NC}"
                    echo -e "${YELLOW}Note: Run 'source $PROFILE_FILE' or start a new shell to use $BINARY_NAME${NC}"
                else
                    echo -e "${GREEN}✓ $INSTALL_DIR already configured in $PROFILE_FILE${NC}"
                fi
            else
                echo -e "${YELLOW}⚠ $INSTALL_DIR is not in your PATH${NC}"
                echo "Please add this to your shell profile:"
                echo -e "${BLUE}  export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
            fi
        fi
        
        echo ""
        echo -e "${GREEN}You can now use '$BINARY_NAME' from anywhere!${NC}"
        if echo "$PATH" | grep -q "$INSTALL_DIR"; then
            echo -e "${BLUE}Try running: $BINARY_NAME${NC}"
        else
            echo -e "${BLUE}Try running: source ~/.${SHELL_NAME}rc && $BINARY_NAME${NC}"
            echo -e "${BLUE}Or in a new shell: $BINARY_NAME${NC}"
        fi
    else
        echo -e "${RED}Error: Installation verification failed${NC}"
        echo -e "${RED}Binary not found at: $INSTALL_PATH${NC}"
        exit 1
    fi
}

# Uninstall function
uninstall() {
    echo -e "${BLUE}Uninstalling $BINARY_NAME...${NC}"
    
    if ! is_installed; then
        echo -e "${YELLOW}$BINARY_NAME is not installed at $INSTALL_PATH${NC}"
        exit 0
    fi
    
    # Remove the binary
    echo -e "${YELLOW}Removing $INSTALL_PATH...${NC}"
    if [ "$USE_SUDO" = true ]; then
        sudo rm -f "$INSTALL_PATH"
    else
        rm -f "$INSTALL_PATH"
    fi
    
    # Verify removal
    if ! is_installed; then
        echo -e "${GREEN}✓ Successfully uninstalled $BINARY_NAME${NC}"
    else
        echo -e "${RED}Error: Uninstallation verification failed${NC}"
        exit 1
    fi
}

# Main script logic
main() {
    if [ $# -eq 0 ]; then
        usage
    fi
    
    case "${1:-}" in
        install)
            if is_installed; then
                echo -e "${YELLOW}$BINARY_NAME is already installed at $INSTALL_PATH${NC}"
                read -p "Do you want to reinstall? (y/N) " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    echo "Installation cancelled."
                    exit 0
                fi
            fi
            install
            ;;
        uninstall)
            uninstall
            ;;
        *)
            echo -e "${RED}Error: Unknown command '$1'${NC}"
            usage
            ;;
    esac
}

main "$@"

