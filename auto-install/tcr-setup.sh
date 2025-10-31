#!/usr/bin/env bash
#
# Simplified TCR Setup Script
# Purpose: Download tcrw and version.txt, create tcr directory, and show usage instructions
#

set -euo pipefail

# Colors for output
RED='\e[1;31m'
GREEN='\e[1;32m'
YELLOW='\e[1;33m'
BLUE='\e[1;34m'
CYAN='\e[1;36m'
NC='\e[1;0m' # No Color

# Configuration
GITHUB_REPO="murex/TCR"
GITHUB_RAW_URL="https://raw.githubusercontent.com/${GITHUB_REPO}/main/auto-install"

# Print messages
print_info() {
    echo -e "${CYAN}[INFO]${NC} $1" >&2
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" >&2
}

# Download tcrw to current directory
download_tcrw() {
    print_info "Downloading tcrw..."
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "${GITHUB_RAW_URL}/tcrw" -o tcrw
    elif command -v wget >/dev/null 2>&1; then
        wget -q "${GITHUB_RAW_URL}/tcrw" -O tcrw
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    # Make tcrw executable
    chmod +x tcrw
    print_success "Downloaded tcrw and made it executable"
}

# Create tcr directory
create_tcr_directory() {
    print_info "Creating tcr directory..."
    mkdir -p tcr
    print_success "Created tcr directory"
}

# Download version.txt to tcr directory
download_version_file() {
    print_info "Downloading version.txt..."
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "${GITHUB_RAW_URL}/version.txt" -o tcr/version.txt
    else
        wget -q "${GITHUB_RAW_URL}/version.txt" -O tcr/version.txt
    fi
    print_success "Downloaded version.txt to tcr/version.txt"
}

# Show usage instructions
show_usage() {
    echo ""
    print_success "TCR setup completed!"
    echo ""
    echo "To launch TCR, run:"
    echo "  ./tcrw"
    echo ""
    echo "For help and available options:"
    echo "  ./tcrw --help"
    echo ""
    echo "For more information, visit: https://github.com/${GITHUB_REPO}"
    echo ""
}

# Display TCR banner
show_banner() {
    echo ""
    curl -sSL "${GITHUB_RAW_URL}/banner.sh" | bash
    echo ""
}

# Main function
main() {
    show_banner
    print_info "Starting TCR setup..."

    download_tcrw
    create_tcr_directory
    download_version_file
    show_usage
}

# Run main function
main "$@"
