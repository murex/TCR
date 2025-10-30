#!/usr/bin/env bash
#
# TCR Auto-Install Script
# Usage: curl -sSL https://raw.githubusercontent.com/murex/TCR/main/auto-install/tcr-setup.sh | bash
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
GITHUB_API_URL="https://api.github.com/repos/${GITHUB_REPO}"
GITHUB_RELEASES_URL="${GITHUB_API_URL}/releases/latest"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="tcr"

GITHUB_RAW_URL="https://raw.githubusercontent.com"
GITHUB_BRANCH="experiment/one-liner-install-script"
GITHUB_AUTO_INSTALL_URL="${GITHUB_RAW_URL}/${GITHUB_REPO}/${GITHUB_BRANCH}/auto-install"

# Print colored messages
print_info() {
    echo -e "${CYAN}[INFO]${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" >&2
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" >&2
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Display TCR banner
show_banner() {
    echo "" >&2
    curl -sSL "${GITHUB_AUTO_INSTALL_URL}/banner.sh" | bash
    echo "" >&2
}

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Linux*)     os="Linux";;
        Darwin*)    os="Darwin";;
        MINGW*|MSYS*|CYGWIN*) os="Windows";;
        *)          print_error "Unsupported operating system: $(uname -s)"; exit 1;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)   arch="x86_64";;
        i386|i686)      arch="i386";;
        arm64|aarch64)  arch="arm64";;
        *)              print_error "Unsupported architecture: $(uname -m)"; exit 1;;
    esac

    echo "${os}_${arch}"
}

# Get latest release information
get_latest_release() {
    print_info "Fetching latest release information..."
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "${GITHUB_RELEASES_URL}"
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "${GITHUB_RELEASES_URL}"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Extract download URL from release JSON
get_download_url() {
    local platform="$1"
    local release_json="$2"

    echo "${release_json}" | grep -o "\"browser_download_url\":\"[^\"]*tcr_[^\"]*_${platform}\.tar\.gz\"" | cut -d'"' -f4 | head -n1
}

# Extract version from release JSON
get_version() {
    local release_json="$1"
    echo "${release_json}" | grep -o '"tag_name": *"[^"]*' | cut -d'"' -f4
}

# Download and extract TCR
download_and_install() {
    local download_url="$1"
    local version="$2"
    local platform="$3"

    print_info "Downloading TCR ${version} for ${platform}..."

    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    cd "${temp_dir}"

    # Download archive
    if command -v curl >/dev/null 2>&1; then
        curl -sSL "${download_url}" -o tcr.tar.gz
    else
        wget -q "${download_url}" -O tcr.tar.gz
    fi

    print_info "Extracting TCR..."
    tar -xzf tcr.tar.gz

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Move binary to install directory
    if [[ "${platform}" == *"Windows"* ]]; then
        mv tcr.exe "${INSTALL_DIR}/${BINARY_NAME}"
    else
        mv tcr "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Cleanup
    cd - >/dev/null
    rm -rf "${temp_dir}"

    print_success "TCR ${version} installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Check if install directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        print_warning "Install directory ${INSTALL_DIR} is not in your PATH."
        print_info "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
        print_info "Or run: echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.bashrc"
    fi
}

# Verify installation
verify_installation() {
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        local version
        version=$("${BINARY_NAME}" --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation verified: ${version}"
    elif [[ -x "${INSTALL_DIR}/${BINARY_NAME}" ]]; then
        local version
        version=$("${INSTALL_DIR}/${BINARY_NAME}" --version 2>/dev/null | head -n1 || echo "unknown")
        print_success "Installation verified: ${version}"
        print_info "Run with: ${INSTALL_DIR}/${BINARY_NAME}"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Show usage information
show_usage() {
    print_info "Usage examples:"
    echo "  ${BINARY_NAME} -b <base-directory> -w <work-directory> -l <language> -t <toolchain>"
    echo "  ${BINARY_NAME} --help"
    echo ""
    print_info "For more information, visit: https://github.com/${GITHUB_REPO}"
}

# Main installation function
main() {
    # Show banner
    show_banner

    print_info "Starting TCR setup for current directory..."

    # Detect platform
    local platform
    platform=$(detect_platform)
    print_info "Detected platform: ${platform}"

    # Get latest release
    local release_json
    release_json=$(get_latest_release)

    # Extract version and download URL
    local version download_url
    version=$(get_version "${release_json}")

    print_info "Latest version: ${version}"

    download_url=$(get_download_url "${platform}" "${release_json}")

    if [[ -z "${download_url}" ]]; then
        print_error "Could not find download URL for platform: ${platform}"
        print_error "Please check https://github.com/${GITHUB_REPO}/releases for manual installation"
        exit 1
    fi

    print_info "Latest version: ${version}"

    # Check if already installed
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        local current_version
        current_version=$("${BINARY_NAME}" --version 2>/dev/null | head -n1 || echo "unknown")
        if [[ "${current_version}" == *"${version}"* ]]; then
            print_success "TCR ${version} is already installed and up to date"
            show_usage
            exit 0
        else
            print_warning "TCR is already installed (${current_version}), updating to ${version}..."
        fi
    fi

    # Download and install
    download_and_install "${download_url}" "${version}" "${platform}"

    # Check PATH
    check_path

    # Verify installation
    verify_installation

    # Show usage
    show_usage

    print_success "TCR installation completed!"
}

# Run main function
main "$@"
