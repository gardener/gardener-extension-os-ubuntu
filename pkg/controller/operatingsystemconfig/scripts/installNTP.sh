#!/usr/bin/env bash

set -e

install_ntp() {
  echo "Installing ntpd client..."
  if package_installed ntp; then
    echo "Package is already installed!"
    return
  fi
  echo "apt update && apt install -y ntp"
  apt update && DEBIAN_FRONTEND=noninteractive apt install -o Dpkg::Options::="--force-confold" -y ntp
  echo "ntp installed successfully!"
}

ntpd() {
  install_ntp
}

systemd_timesyncd() {
  install_systemd_timesyncd
}

install_systemd_timesyncd() {
  echo "Installing systemd-timesyncd client..."
  if package_installed systemd-timesyncd; then
    echo "Package is already installed!"
    return
  fi
  echo "apt update && apt install -y systemd-timesyncd"
  apt update && DEBIAN_FRONTEND=noninteractive apt install -y systemd-timesyncd
  echo "systemd-timesyncd installed successfully!"
}

package_installed() {
  if ! dpkg-query -W -f='${Status}' "$1" 2>/dev/null | grep -c "ok installed"; then
    return 1
  fi
  return 0
}

resolve_daemon () {
  local daemon="$1"
  shift
  local ubuntu_version=""
  if [ -f "$OS_RELEASE_FILE" ]; then
    ubuntu_version=$(grep '^VERSION_ID=' "$OS_RELEASE_FILE" | cut -d= -f2 | tr -d '"')
  fi
  for override in "$@"; do
    if [ "${override%%=*}" = "$ubuntu_version" ]; then
      daemon="${override#*=}"
      break
    fi
  done
  echo "$daemon"
}

OS_RELEASE_FILE="${OS_RELEASE_FILE:-/etc/os-release}"

if [ -z "$1" ]; then
    echo "Usage: $0 <option> [<ubuntu-version>=<option>...]"
    echo "Options:"
    echo "  ntpd                Install ntp"
    echo "  systemd-timesyncd   Install systemd-timesyncd"
    echo "  none                Do not manage the NTP client, keep the default"
    exit 1
fi

DAEMON=$(resolve_daemon "$@")

# Process the argument
case "$DAEMON" in
    ntpd)
        ntpd
        ;;
    systemd-timesyncd)
        systemd_timesyncd
        ;;
    none)
        echo "NTP client management is disabled for this Ubuntu version, keeping the default."
        ;;
    *)
        echo "Invalid option: $DAEMON"
        echo "Please use 'ntpd' or 'systemd-timesyncd' or 'none'."
        exit 1
        ;;
esac

echo "Installation complete."