#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print messages
print_message() {
  echo -e "${BLUE}==>${NC} $1"
}

print_success() {
  echo -e "${GREEN}==>${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}==>${NC} $1"
}

print_error() {
  echo -e "${RED}==>${NC} $1"
}

# Detect operating system
detect_os() {
  OS="$(uname -s)"
  case "${OS}" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    *)          print_error "Unsupported operating system: ${OS}"; exit 1;;
  esac
  
  print_message "Operating system: ${OS}"
  return 0
}

# Detect architecture
detect_arch() {
  ARCH="$(uname -m)"
  case "${ARCH}" in
    x86_64*)    ARCH="amd64";;
    arm64*)     ARCH="arm64";;
    aarch64*)   ARCH="arm64";;
    *)          print_error "Unsupported architecture: ${ARCH}"; exit 1;;
  esac
  
  print_message "Architecture: ${ARCH}"
  return 0
}

# Get latest version
get_latest_version() {
  print_message "Checking latest version..."
  
  LATEST_VERSION=$(curl -s https://api.github.com/repos/Firdavs9512/aurora-agent/releases/latest | grep "tag_name" | cut -d '"' -f 4)
  
  if [ -z "$LATEST_VERSION" ]; then
    print_error "Error getting version"
    exit 1
  fi
  
  print_message "Latest version: ${LATEST_VERSION}"
  return 0
}

# Download binary
download_binary() {
  DOWNLOAD_URL="https://github.com/Firdavs9512/aurora-agent/releases/download/${LATEST_VERSION}/aurora-agent-${OS}-${ARCH}.tar.gz"
  TEMP_DIR=$(mktemp -d)
  TEMP_FILE="${TEMP_DIR}/aurora-agent.tar.gz"
  
  print_message "Downloading: ${DOWNLOAD_URL}"
  
  if ! curl -L -o "${TEMP_FILE}" "${DOWNLOAD_URL}"; then
    print_error "Error downloading"
    rm -rf "${TEMP_DIR}"
    exit 1
  fi
  
  print_message "Extracting archive..."
  tar -xzf "${TEMP_FILE}" -C "${TEMP_DIR}"
  
  BINARY_PATH="${TEMP_DIR}/aurora-agent-${OS}-${ARCH}"
  if [ ! -f "${BINARY_PATH}" ]; then
    print_error "Binary file not found"
    rm -rf "${TEMP_DIR}"
    exit 1
  fi
  
  chmod +x "${BINARY_PATH}"
  return 0
}

# Install binary
install_binary() {
  INSTALL_DIR="/usr/local/bin"
  BINARY_NAME="aurora"
  TARGET_PATH="${INSTALL_DIR}/${BINARY_NAME}"
  
  # Check install directory
  if [ ! -d "${INSTALL_DIR}" ]; then
    print_message "Creating install directory: ${INSTALL_DIR}"
    mkdir -p "${INSTALL_DIR}"
  fi
  
  # Check old version
  if [ -f "${TARGET_PATH}" ]; then
    CURRENT_VERSION=$(${TARGET_PATH} version 2>/dev/null | grep -o "v[0-9]*\.[0-9]*\.[0-9]*" || echo "unknown")
    
    if [ "${CURRENT_VERSION}" = "${LATEST_VERSION}" ]; then
      print_success "Aurora Agent already on latest version (${LATEST_VERSION})"
      rm -rf "${TEMP_DIR}"
      exit 0
    fi
    
    print_message "Current version: ${CURRENT_VERSION}, updating to: ${LATEST_VERSION}"
  else
    print_message "Installing: ${LATEST_VERSION}"
  fi
  
  # Install binary
  print_message "Installing: ${TARGET_PATH}"
  
  if [ -w "${INSTALL_DIR}" ]; then
    mv "${BINARY_PATH}" "${TARGET_PATH}"
  else
    print_warning "Install directory is not writable, sudo is required"
    sudo mv "${BINARY_PATH}" "${TARGET_PATH}"
  fi
  
  rm -rf "${TEMP_DIR}"
  
  # Installation successful
  print_success "Aurora Agent ${LATEST_VERSION} installed successfully!"
  print_message "To start, run the following command: ${BINARY_NAME}"
  return 0
}

# Main program
main() {
  print_message "Aurora Agent installer"
  
  detect_os
  detect_arch
  get_latest_version
  download_binary
  install_binary
  
  return 0
}

main "$@" 