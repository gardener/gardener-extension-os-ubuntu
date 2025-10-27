package v1alpha1

import (
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
	// Mirror to set custom Ubuntu mirror
	// +optional
	ATPConfig APTConfig `json:"apt,omitempty"`
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
}

type APTConfig struct {
	PreserveSourcesList bool         `json:"preserveSourcesList,omitempty"`
	Primary             []APTArchive `json:"primary,omitempty"`
	Security            []APTArchive `json:"security,omitempty"`
}

type APTArchive struct {
	Arches    []string `json:"arches,omitempty"`
	URI       string   `json:"uri,omitempty"`
	Search    []string `json:"search,omitempty"`
	SearchDNS bool     `json:"searchDNS,omitempty"`
}
