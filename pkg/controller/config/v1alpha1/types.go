// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Daemon string

const (
	SystemdTimesyncd Daemon = "systemd-timesyncd"
	NTPD             Daemon = "ntpd"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExtensionConfig is the configuration for the os-ubuntu extension.
type ExtensionConfig struct {
	metav1.TypeMeta `json:",inline"`

	// NTP to configure either systemd-timesyncd or ntpd
	// +optional
	NTP *NTPConfig `json:"ntp,omitempty"`
	// DisableUnattendedUpgrades to disable unattended upgrades in ubuntu
	// +optional
	DisableUnattendedUpgrades *bool `json:"disableUnattendedUpgrades,omitempty"`
	// AptRepositories is the list of additional apt repositories to configure
	// via cloud-init.
	// +optional
	AptRepositories []AptRepository `json:"aptRepositories,omitempty"`
	// Dependencies is the list of apt packages to install on the node. If empty,
	// a default set of unpinned packages is installed. Each dependency may
	// optionally target a specific Ubuntu version and/or build serial.
	// +optional
	Dependencies []DependencyConfig `json:"dependencies,omitempty"`
	// Mirror to set custom Ubuntu mirror
	// +optional
	APTConfig *APTConfig `json:"apt,omitempty"`
}

// AptRepository describes an additional apt repository to configure via
// cloud-init.
type AptRepository struct {
	// Name is a unique name for the apt source.
	Name string `json:"name"`
	// URI is the base URI of the apt repository. For the docker repository this
	// is typically https://download.docker.com/linux/ubuntu.
	URI string `json:"uri"`
	// Key is the ASCII-armored GPG key used to sign the repository. If empty,
	// the repository is configured without GPG signature verification.
	// +optional
	Key string `json:"key,omitempty"`
	// KeyURL is the URL to download the GPG key from. The key is downloaded
	// via cloud-init write_files to /etc/apt/keyrings/<name>.gpg and
	// referenced with signed-by. Mutually exclusive with Key.
	// +optional
	KeyURL string `json:"keyUrl,omitempty"`
	// Suite is the apt suite to use. Defaults to "$RELEASE" which cloud-init
	// substitutes with the release codename.
	// +optional
	Suite string `json:"suite,omitempty"`
	// Components is the list of apt components to use. Defaults to ["stable"].
	// +optional
	Components []string `json:"components,omitempty"`
}

// DependencyConfig describes an apt package to install.
type DependencyConfig struct {
	// Name is the name of the apt package.
	Name string `json:"name"`
	// Version is the exact apt package version to install. If empty, the latest
	// available version is installed.
	// +optional
	Version string `json:"version,omitempty"`
	// UbuntuVersion optionally restricts this dependency to a specific Ubuntu
	// version (e.g. "22.04"). The value is matched against VERSION_ID from
	// /etc/os-release.
	// +optional
	UbuntuVersion string `json:"ubuntuVersion,omitempty"`
	// UbuntuBuildSerial optionally restricts this dependency to a specific
	// Ubuntu build serial (e.g. "20250725"). The value is matched against the
	// serial field in /etc/cloud/build.info.
	// +optional
	UbuntuBuildSerial string `json:"ubuntuBuildSerial,omitempty"`
	// Hold, if true, marks the package on apt hold after installation so that
	// apt upgrade does not update it.
	// +optional
	Hold bool `json:"hold,omitempty"`
}

// NTPConfig General NTP Config for either systemd-timesyncd or ntpd
type NTPConfig struct {
	// Daemon One of either systemd-timesyncd or ntp
	Daemon Daemon `json:"daemon"`
	// NTPD to configure the ntpd client
	// +optional
	NTPD *NTPDConfig `json:"ntpd,omitempty"`
}

// NTPDConfig is the struct used in the ntp-config.conf.tpl template file
type NTPDConfig struct {
	// Servers List of ntp servers
	Servers []string `json:"servers"`
	// Interfaces for ntpd to bind to. Can be more than one.
	Interfaces []string `json:"interfaces,omitempty"`
}

type APTConfig struct {
	PreserveSourcesList bool         `json:"preserveSourcesList,omitempty"`
	Primary             []APTArchive `json:"primary,omitempty"`
	Security            []APTArchive `json:"security,omitempty"`
}

type APTArchive struct {
	Arches    []Architecture `json:"arches,omitempty"`
	URI       string         `json:"uri,omitempty"`
	Search    []string       `json:"search,omitempty"`
	SearchDNS bool           `json:"searchDNS,omitempty"`
}

type Architecture string

const (
	AMD64       Architecture = constants.ArchitectureAMD64
	ARM64       Architecture = constants.ArchitectureARM64
	ArchDefault Architecture = "default"
)
