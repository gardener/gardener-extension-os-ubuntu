// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"net/url"
	"regexp"
	"slices"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

var (
	validPackageName       = regexp.MustCompile(`^[a-z0-9][a-z0-9.+-]+$`)
	validPackageVersion    = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9.+:~\-]*$`)
	validUbuntuVersion     = regexp.MustCompile(`^[0-9]+\.[0-9]+$`)
	validBuildSerial       = regexp.MustCompile(`^[0-9.]+$`)
	validAptRepositoryName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
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
		allErrs = append(allErrs, validateAPTConfig(config.APTConfig, rootPath.Child("apt"))...)
	}

	allErrs = append(allErrs, validateAptRepositories(config.AptRepositories, rootPath.Child("aptRepositories"))...)
	allErrs = append(allErrs, validateDependencies(config.Dependencies, rootPath.Child("dependencies"))...)

	return allErrs
}

func validateAptRepositories(config []configv1alpha1.AptRepository, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	seenNames := sets.New[string]()
	for i, repo := range config {
		repoPath := fldPath.Index(i)
		if repo.Name == "" {
			allErrs = append(allErrs, field.Required(repoPath.Child("name"), "name is required"))
		} else if !validAptRepositoryName.MatchString(repo.Name) {
			allErrs = append(allErrs, field.Invalid(repoPath.Child("name"), repo.Name, "must be a valid name (alphanumeric, dots, underscores, hyphens)"))
		} else if seenNames.Has(repo.Name) {
			allErrs = append(allErrs, field.Duplicate(repoPath.Child("name"), repo.Name))
		} else {
			seenNames.Insert(repo.Name)
		}
		if !isValidURL(repo.URI) {
			allErrs = append(allErrs, field.Invalid(repoPath.Child("uri"), repo.URI, "invalid URL"))
		}
		if repo.Key != "" && repo.KeyURL != "" {
			allErrs = append(allErrs, field.Forbidden(repoPath.Child("keyUrl"), "key and keyUrl are mutually exclusive"))
		}
		if repo.KeyURL != "" && !isValidURL(repo.KeyURL) {
			allErrs = append(allErrs, field.Invalid(repoPath.Child("keyUrl"), repo.KeyURL, "invalid URL"))
		}
		if repo.KeyFormat != "" {
			if repo.KeyURL == "" {
				allErrs = append(allErrs, field.Forbidden(repoPath.Child("keyFormat"), "keyFormat is only valid when keyUrl is set"))
			}
			if strings.ToLower(string(repo.KeyFormat)) != string(configv1alpha1.KeyFormatASC) && strings.ToLower(string(repo.KeyFormat)) != string(configv1alpha1.KeyFormatGPG) {
				allErrs = append(allErrs, field.NotSupported(repoPath.Child("keyFormat"), repo.KeyFormat, []string{string(configv1alpha1.KeyFormatASC), string(configv1alpha1.KeyFormatGPG)}))
			}
		}
	}
	return allErrs
}

func validateDependencies(config []configv1alpha1.DependencyConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	unconstrainedDeps := sets.New[string]()

	for i, dep := range config {
		depPath := fldPath.Index(i)
		if dep.Name == "" {
			allErrs = append(allErrs, field.Required(depPath.Child("name"), "name is required"))
		} else if !validPackageName.MatchString(dep.Name) {
			allErrs = append(allErrs, field.Invalid(depPath.Child("name"), dep.Name, "must be a valid apt package name (lowercase letters, digits, dots, hyphens, plus)"))
		}
		if dep.Version != "" && !validPackageVersion.MatchString(dep.Version) {
			allErrs = append(allErrs, field.Invalid(depPath.Child("version"), dep.Version, "must be a valid apt package version"))
		}
		if dep.UbuntuVersion != "" && !validUbuntuVersion.MatchString(dep.UbuntuVersion) {
			allErrs = append(allErrs, field.Invalid(depPath.Child("ubuntuVersion"), dep.UbuntuVersion, "must be a valid Ubuntu version (e.g. \"22.04\")"))
		}
		if dep.UbuntuBuildSerial != "" && !validBuildSerial.MatchString(dep.UbuntuBuildSerial) {
			allErrs = append(allErrs, field.Invalid(depPath.Child("ubuntuBuildSerial"), dep.UbuntuBuildSerial, "must be a numeric build serial (digits and dots)"))
		}

		if dep.Name != "" && dep.UbuntuVersion == "" && dep.UbuntuBuildSerial == "" {
			if unconstrainedDeps.Has(dep.Name) {
				allErrs = append(allErrs, field.Duplicate(depPath.Child("name"), dep.Name))
			} else {
				unconstrainedDeps.Insert(dep.Name)
			}
		}
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
	allErrs = append(allErrs, validateAPTArchive(config.Primary, fldPath.Child("primary"))...)
	allErrs = append(allErrs, validateAPTArchive(config.Security, fldPath.Child("security"))...)
	return allErrs
}

func validateAPTArchive(config []configv1alpha1.APTArchive, fldPath *field.Path) field.ErrorList {
	validArchitectureNames := sets.New(configv1alpha1.ArchDefault, configv1alpha1.AMD64, configv1alpha1.ARM64)
	allErrs := field.ErrorList{}
	for _, configArchive := range config {
		for _, arch := range configArchive.Arches {
			if !slices.Contains(validArchitectureNames.UnsortedList(), arch) {
				allErrs = append(allErrs, field.NotSupported(fldPath.Child("arches"), configArchive.Arches, validArchitectureNames.UnsortedList()))
			}
		}
		if !isValidURL(configArchive.URI) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("uri"), configArchive.URI, "invalid URL"))
		}
		for _, search := range configArchive.Search {
			if !isValidURL(search) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("search"), search, "invalid URL"))
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
