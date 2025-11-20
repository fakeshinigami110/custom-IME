#!/bin/bash

# prerequisites.sh - Install system dependencies for custom-ime
set -e

echo "Installing system dependencies for custom-ime..."

# Check if running as root
if [ "$EUID" -eq 0 ]; then
    echo "❌ Please do not run as root. Run as normal user."
    exit 1
fi

# Detect distribution
if command -v pacman > /dev/null 2>&1; then
    # Arch Linux
    echo "Detected Arch Linux"
    echo ""
    echo "Required packages for Arch Linux:"
    echo "   fcitx5                    # Input method framework"
    echo "   cmake                     # Build system"
    echo "   base-devel                # C/C++ compiler and development tools"
    echo "   go                        # Go compiler"
    echo "   git                       # Version control"
    echo ""
    echo "Install with:"
    echo "   sudo pacman -S fcitx5 cmake base-devel go git"
    echo ""
    echo "Search for available fcitx5 packages:"
    echo "   pacman -Ss fcitx5"
    
elif command -v apt > /dev/null 2>&1; then
    # Debian/Ubuntu
    echo "Detected Debian/Ubuntu"
    echo "Please install these packages manually:"
    echo ""
    echo "sudo apt update"
    echo "sudo apt install fcitx5 cmake build-essential golang-go git"
    echo ""
    echo "🔍 Search for available fcitx5 packages:"
    echo "   apt-cache search fcitx5"
    
elif command -v dnf > /dev/null 2>&1; then
    # Fedora
    echo "Detected Fedora"
    echo "Please install these packages manually:"
    echo ""
    echo "sudo dnf install fcitx5 cmake gcc-c++ golang git"
    echo ""
    echo "🔍 Search for available fcitx5 packages:"
    echo "   dnf search fcitx5"
    
elif command -v yum > /dev/null 2>&1; then
    # RHEL/CentOS
    echo "Detected RHEL/CentOS"
    echo "Please install these packages manually:"
    echo ""
    echo "sudo yum install fcitx5 cmake gcc-c++ golang git"
    echo ""
    echo "🔍 Search for available fcitx5 packages:"
    echo "   yum search fcitx5"
    
elif command -v zypper > /dev/null 2>&1; then
    # openSUSE
    echo "Detected openSUSE"
    echo "Please install these packages manually:"
    echo ""
    echo "sudo zypper install fcitx5 cmake gcc-c++ golang git"
    echo ""
    echo "🔍 Search for available fcitx5 packages:"
    echo "   zypper search fcitx5"
    
else
    echo "Unsupported distribution."
    echo "Please install these packages manually:"
    echo "   - fcitx5"
    echo "   - cmake" 
    echo "   - C/C++ compiler (gcc, g++)"
    echo "   - build tools (make, etc.)"
    echo "   - go (golang)"
    echo "   - git"
fi

echo ""
echo "Next steps:"
echo "   1. Install the packages above if they havent"
echo "   2. Run: ./install.sh to install custom-ime"