#!/bin/bash

# install.sh - Install custom-ime tool only
set -e

echo "Installing custom-ime tool..."

echo "Note: For most users, installing from AUR is recommended:"
echo "   yay -S custom-ime  OR  paru -S custom-ime"
echo ""

read -p "Continue with manual installation? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Installation canceled."
    exit 0
fi

# Simple dependency check
echo "🔍 Checking basic dependencies..."

if ! command -v go > /dev/null 2>&1; then
    echo "Go is not installed."
    echo "   Run ./prerequisites.sh to see installation instructions"
    exit 1
fi

if ! command -v cmake > /dev/null 2>&1; then
    echo "CMake not found."
    echo "   Run ./prerequisites.sh to see installation instructions"
    exit 1
fi

echo "✅ Basic dependencies found"

# Build custom-ime
echo "Building custom-ime..."
if ! go build -o custom-ime; then
    echo "Build failed!"
    echo "   Make sure all dependencies are installed"
    exit 1
fi

# Install to system
echo "Installing to /usr/local/bin/"
sudo mv custom-ime /usr/local/bin/

echo "custom-ime installed successfully!"
echo ""
echo " Usage:"
echo "   custom-ime create -p myproject -n imename"
echo "   custom-ime install -p myproject"
echo "   custom-ime list"
echo ""
echo "💡 Tip: Run 'custom-ime create --help' for more options"