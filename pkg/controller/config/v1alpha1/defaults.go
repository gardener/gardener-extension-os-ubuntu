// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ExtensionConfig(obj *ExtensionConfig) {
	if obj.NTP == nil {
		obj.NTP = &NTPConfig{}
	}
	if obj.DisableUnattendedUpgrades == nil {
		obj.DisableUnattendedUpgrades = ptr.To(false)
	}

	requiredPackages := []string{
		"containerd",
		"runc",
		"socat",
		"nfs-common",
		"logrotate",
		"jq",
		"policykit-1",
	}

	disabledPackages := sets.New[string]()
	hasUnconstrained := sets.New[string]()
	for _, dep := range obj.Dependencies {
		if dep.Disabled {
			disabledPackages.Insert(dep.Name)
		} else if dep.UbuntuVersion == "" && dep.UbuntuBuildSerial == "" {
			hasUnconstrained.Insert(dep.Name)
		}
	}

	for _, pkg := range requiredPackages {
		if !disabledPackages.Has(pkg) && !hasUnconstrained.Has(pkg) {
			obj.Dependencies = append(obj.Dependencies, DependencyConfig{Name: pkg})
		}
	}
}

func SetDefaults_NTPConfig(obj *NTPConfig) {
	if obj.Daemon == "" {
		obj.Daemon = SystemdTimesyncd
	}
}

func SetDefaults_NTPDConfig(obj *NTPDConfig) {
	if obj.Servers == nil {
		obj.Servers = make([]string, 0)
	}
}
