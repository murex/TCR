#!/usr/bin/env bash

# prepare-release.sh - Prepare TCR for a new release
# This script automates version updates across the project and creates a git tag

set -e  # Exit on error
set -u  # Exit on undefined variable

# Color codes for output
RED='\e[1;31m'
GREEN='\e[1;32m'
YELLOW='\e[1;33m'
BLUE='\e[1;34m'
CYAN='\e[1;36m'
NC='\e[1;0m' # No Color

# Function to print colored messages
print_line() {
    printf "$1\n" >&2
}

print_info() {
    print_line "${CYAN}[INFO]${NC} $1"
}

print_warning() {
    print_line "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    print_line "${RED}[ERROR]${NC} $1"
}

print_success() {
    print_line "${GREEN}[DONE]${NC} $1"
}

# Function to display usage
print_usage() {
    cat << EOF

Usage: $0 <version>

Prepare TCR for a new release by updating version information across the project.

Arguments:
  version    The new TCR version in the form vX.Y.Z (e.g., v1.5.0)

Example:
  $0 v1.5.0

This script will:
  1. Update tools/tcr/version.txt
  2. Update auto-install/version.txt
  3. Update examples/tcr/version.txt
  4. Update webapp/package.json (version without 'v' prefix)
  5. Update tar commands in README.md with new version numbers
  6. Commit all changes with message "Upgrade TCR version to X.Y.Z"
  7. Create a git tag vX.Y.Z
EOF
    exit 1
}

# Check if version argument is provided
if [ $# -ne 1 ]; then
    print_error "Missing version argument"
    print_usage
fi

NEW_VERSION="$1"

# Validate version format (vX.Y.Z)
if ! [[ "$NEW_VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format: $NEW_VERSION"
    print_line "Version must be in the form vX.Y.Z (e.g., v1.5.0)"
    exit 1
fi

# Extract version components
VERSION_WITH_V="$NEW_VERSION"
VERSION_WITHOUT_V="${NEW_VERSION#v}"  # Remove leading 'v'
VERSION_X=$(echo "$VERSION_WITHOUT_V" | cut -d. -f1)
VERSION_Y=$(echo "$VERSION_WITHOUT_V" | cut -d. -f2)
VERSION_Z=$(echo "$VERSION_WITHOUT_V" | cut -d. -f3)

print_info "Preparing release for version: $VERSION_WITH_V"
# print_info "Version components: X=$VERSION_X, Y=$VERSION_Y, Z=$VERSION_Z"

# Get the project root directory (assumes script is in tools/scripts/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

print_info "Project root: $PROJECT_ROOT"

# Navigate to project root
cd "$PROJECT_ROOT"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    print_warning "You have uncommitted changes in your working directory"
    read -p "Do you want to continue? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Aborted by user"
        exit 0
    fi
fi

# Check if tag already exists
if git rev-parse "$VERSION_WITH_V" >/dev/null 2>&1; then
    print_error "Tag $VERSION_WITH_V already exists"
    exit 1
fi

# Array to track modified files
declare -a MODIFIED_FILES=()

# Update tools/tcr/version.txt
print_info "Updating tools/tcr/version.txt"
echo "tcr $VERSION_WITHOUT_V" > tools/tcr/version.txt
MODIFIED_FILES+=("tools/tcr/version.txt")

# Update auto-install/version.txt
print_info "Updating auto-install/version.txt"
echo "tcr $VERSION_WITHOUT_V" > auto-install/version.txt
MODIFIED_FILES+=("auto-install/version.txt")

# Update examples/tcr/version.txt
print_info "Updating examples/tcr/version.txt"
echo "tcr $VERSION_WITHOUT_V" > examples/tcr/version.txt
MODIFIED_FILES+=("examples/tcr/version.txt")

# Update webapp/package.json
print_info "Updating webapp/package.json"
if [ -f "webapp/package.json" ]; then
    # Use sed to update the version field (compatible with both Linux and macOS)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS requires -i with an argument (can be empty string)
        sed -i '' "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_WITHOUT_V\"/" webapp/package.json
    else
        # Linux sed
        sed -i "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_WITHOUT_V\"/" webapp/package.json
    fi
    MODIFIED_FILES+=("webapp/package.json")
else
    print_warning "webapp/package.json not found, skipping"
fi

# Update README.md - replace tar commands with new version
print_info "Updating version in README.md"
if [ -f "README.md" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS sed
        sed -i '' "s/tcr_[0-9]\+\.[0-9]\+\.[0-9]\+_\(.*\.tar\.gz\)/tcr_${VERSION_WITHOUT_V}_\1/g" README.md
    else
        # Linux sed
        sed -i "s/tcr_[0-9]\+\.[0-9]\+\.[0-9]\+_\(.*\.tar\.gz\)/tcr_${VERSION_WITHOUT_V}_\1/g" README.md
    fi
    MODIFIED_FILES+=("README.md")
else
    print_warning "README.md not found, skipping"
fi

# Display changes
print_info "Modified files:"
for file in "${MODIFIED_FILES[@]}"; do
    print_line "  - $file"
done

# Git add all modified files
print_info "Staging changes"
git add "${MODIFIED_FILES[@]}"

# Commit changes
COMMIT_MESSAGE="Upgrade TCR version to $VERSION_WITHOUT_V"
print_info "Committing changes: $COMMIT_MESSAGE"
git commit -m "$COMMIT_MESSAGE"

# Create git tag
print_info "Creating git tag: $VERSION_WITH_V"
git tag -a "$VERSION_WITH_V" -m "Release $VERSION_WITH_V"

# Success message
print_line ""
print_info "✓ Successfully prepared release $VERSION_WITH_V"
print_line ""
print_line "Next steps:"
print_line "  1. Review the commit and tag:"
print_line "     git show $VERSION_WITH_V"
print_line ""
print_line "  2. Push the commit and tag:"
print_line "     git push origin main"
print_line "     git push origin $VERSION_WITH_V"
print_line ""
print_line "  3. The release is then triggered automatically through GitHub actions."
print_line ""
print_line "  4. Verify the release on GitHub:"
print_line "     https://github.com/murex/TCR/releases/tag/$VERSION_WITH_V"
print_line ""
