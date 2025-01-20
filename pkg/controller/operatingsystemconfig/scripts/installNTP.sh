#!/usr/bin/env bash

set -e

install_ntp() {
  echo "Installing ntpd client..."
  if package_installed ntp; then
    echo "Package is already installed!"
    return
  fi
  echo "apt update && apt install -y ntp"
  apt update && DEBIAN_FRONTEND=noninteractive apt install -y ntp
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
  apt update && apt install -y systemd-timesyncd
  echo "systemd-timesyncd installed successfully!"
}

package_installed() {
  if ! dpkg-query -W -f='${Status}' "$1" 2>/dev/null | grep -c "ok installed"; then
    return 1
  fi
  return 0
}

if [ -z "$1" ]; then
    echo "Usage: $0 <option>"
    echo "Options:"
    echo "  ntpd                Install ntp"
    echo "  systemd-timesyncd   Install systemd-timesyncd"
    exit 1
fi

# Process the argument
case "$1" in
    ntpd)
        ntpd
        ;;
    systemd-timesyncd)
        systemd_timesyncd
        ;;
    *)
        echo "Invalid option: $1"
        echo "Please use 'ntpd' or 'systemd-timesyncd'."
        exit 1
        ;;
esac

echo "Installation complete."