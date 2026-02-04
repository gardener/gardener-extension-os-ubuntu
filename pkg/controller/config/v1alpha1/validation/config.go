// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"net/url"
	"slices"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

func ValidateExtensionConfig(config *configv1alpha1.ExtensionConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	var rootPath *field.Path

	validDaemonNames := sets.New(configv1alpha1.SystemdTimesyncd, configv1alpha1.NTPD)

	if config.NTP != nil {
		// Make sure daemon name is valid
		if !validDaemonNames.Has(config.NTP.Daemon) {
			allErrs = append(allErrs, field.NotSupported(rootPath.Child("daemon"), config.NTP.Daemon, validDaemonNames.UnsortedList()))
		}

		// Check if user configured systemd-timesyncd daemon with ntpd config
		if config.NTP.Daemon == configv1alpha1.SystemdTimesyncd && config.NTP.NTPD != nil {
			allErrs = append(allErrs, field.Forbidden(rootPath.Child("ntpd"), "NTPD config is not allowed if systemd-timesyncd is selected"))
		}

		if config.NTP.NTPD != nil {
			allErrs = append(allErrs, validateNTPDConfig(config.NTP.NTPD, rootPath.Child("ntpd"))...)
		}
	}

	if config.APTConfig != nil {
		allErrs = append(allErrs, validateAPTConfig(config.APTConfig, rootPath)...)
	}

	return allErrs
}

func validateNTPDConfig(config *configv1alpha1.NTPDConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(config.Servers) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("servers"), "a list of NTP servers is required"))
	}
	return allErrs
}

func validateAPTConfig(config *configv1alpha1.APTConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateAPTArchive(config.Primary, fldPath, "primary")...)
	allErrs = append(allErrs, validateAPTArchive(config.Security, fldPath, "security")...)
	return allErrs
}

func validateAPTArchive(config []configv1alpha1.APTArchive, fldPath *field.Path, archiveName string) field.ErrorList {
	validArchitectureNames := sets.New(configv1alpha1.Default, configv1alpha1.AMD64, configv1alpha1.ARM64)
	allErrs := field.ErrorList{}
	for _, configArchive := range config {
		for _, arch := range configArchive.Arches {
			if !slices.Contains(validArchitectureNames.UnsortedList(), arch) {
				allErrs = append(allErrs, field.NotSupported(fldPath.Child("apt").Child(archiveName).Child("arches"), configArchive.Arches, validArchitectureNames.UnsortedList()))
			}
		}
		if !isValidURL(configArchive.URI) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("apt").Child(archiveName).Child("uri"), configArchive.URI, "invalid URL"))
		}
		for _, search := range configArchive.Search {
			if !isValidURL(search) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("apt").Child(archiveName).Child("search"), search, "invalid URL"))
			}
		}
	}
	return allErrs
}

func isValidURL(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}
