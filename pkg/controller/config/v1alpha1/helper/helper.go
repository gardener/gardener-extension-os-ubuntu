// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

func IsDaemonConfigured(config *configv1alpha1.NTPConfig, daemon configv1alpha1.Daemon) bool {
	if config == nil {
		return false
	}
	if config.Daemon == daemon {
		return true
	}
	for _, override := range config.UbuntuVersionOverrides {
		if override.Daemon == daemon {
			return true
		}
	}
	return false
}
