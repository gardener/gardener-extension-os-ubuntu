#!/usr/bin/env bash

set -euo pipefail

# Suppress interactive debconf prompts globally for all apt-get commands in this script
export DEBIAN_FRONTEND=noninteractive

OS_RELEASE_FILE="${OS_RELEASE_FILE:-/etc/os-release}"

log_info() {
  printf '[INFO] %s\n' "$*"
}

log_err() {
  printf '[ERROR] %s\n' "$*" >&2
}

# Ensure the script is run with root privileges before running package management actions
require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    log_err "This script must be run as root."
    exit 1
  fi
}

package_installed() {
  local pkg="$1"
  dpkg-query -W -f='${Status}' "${pkg}" 2>/dev/null | grep -q "ok installed"
}

install_ntp() {
  log_info "Installing ntpd client..."
  if package_installed ntp; then
    log_info "Package 'ntp' is already installed!"
    return 0
  fi

  require_root
  log_info "Updating package cache and installing ntp..."
  apt-get update -qq
  apt-get install -y -o Dpkg::Options::="--force-confold" ntp
  log_info "ntp installed successfully!"
}

install_systemd_timesyncd() {
  log_info "Installing systemd-timesyncd client..."
  if package_installed systemd-timesyncd; then
    log_info "Package 'systemd-timesyncd' is already installed!"
    return 0
  fi

  require_root
  log_info "Updating package cache and installing systemd-timesyncd..."
  apt-get update -qq
  apt-get install -y systemd-timesyncd
  log_info "systemd-timesyncd installed successfully!"
}

resolve_daemon() {
  local daemon="$1"
  shift
  local ubuntu_version=""

  if [[ -f "${OS_RELEASE_FILE}" ]]; then
    # Safely extract VERSION_ID from /etc/os-release without polluting shell environment
    ubuntu_version=$(awk -F'=' '/^VERSION_ID=/ {gsub(/"/, "", $2); print $2}' "${OS_RELEASE_FILE}")
  fi

  for override in "$@"; do
    if [[ "${override%%=*}" == "${ubuntu_version}" ]]; then
      daemon="${override#*=}"
      break
    fi
  done

  echo "${daemon}"
}

usage() {
  cat <<EOF
Usage: $0 <option> [<ubuntu-version>=<option>...]

Options:
  ntpd               Install ntp
  systemd-timesyncd  Install systemd-timesyncd
  none               Do not manage the NTP client, keep default
EOF
  exit 1
}

main() {
  if [[ $# -lt 1 ]]; then
    usage
  fi

  local daemon
  daemon=$(resolve_daemon "$@")

  case "${daemon}" in
    ntpd)
      install_ntp
      ;;
    systemd-timesyncd)
      install_systemd_timesyncd
      ;;
    none)
      log_info "NTP client management is disabled for this Ubuntu version, keeping the default."
      ;;
    *)
      log_err "Invalid option: '${daemon}'"
      log_err "Please use 'ntpd', 'systemd-timesyncd', or 'none'."
      exit 1
      ;;
  esac

  log_info "Installation complete."
}

main "$@"