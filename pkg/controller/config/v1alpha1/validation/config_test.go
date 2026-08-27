// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

var _ = Describe("ExtensionConfig validation", func() {
	var (
		config *configv1alpha1.ExtensionConfig
	)

	BeforeEach(func() {
		config = &configv1alpha1.ExtensionConfig{
			NTP: &configv1alpha1.NTPConfig{
				Daemon: configv1alpha1.SystemdTimesyncd,
			},
		}
	})

	It("should allow valid config", func() {
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with incorrect daemon name", func() {
		config.NTP.Daemon = "foo"
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeNotSupported))
		Expect(errs[0].Field).To(Equal("daemon"))
	})

	It("should succeed with valid NTPd config", func() {
		config.NTP.Daemon = configv1alpha1.NTPD
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{}
		config.NTP.NTPD.Servers = []string{"ntp.ubuntu.com"}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(BeEmpty())
	})

	It("should fail with invalid NTPd config (no ntp servers provided)", func() {
		config.NTP.Daemon = configv1alpha1.NTPD
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{}
		config.NTP.NTPD.Servers = []string{}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeRequired))
		Expect(errs[0].Field).To(Equal("ntpd.servers"))
	})

	It("should fail with daemon systemd-timesyncd and ntpd config set", func() {
		config.NTP.Daemon = configv1alpha1.SystemdTimesyncd
		config.NTP.NTPD = &configv1alpha1.NTPDConfig{Servers: []string{"foo.bar"}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeForbidden))
		Expect(errs[0].Field).To(Equal("ntpd"))
	})

	It("should succeed with a valid architecture, URI and search for primary apt mirror", func() {
		config.APTConfig = &configv1alpha1.APTConfig{Primary: []configv1alpha1.APTArchive{{
			Arches:    []configv1alpha1.Architecture{configv1alpha1.ARM64},
			URI:       "http://packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
			Search:    []string{"http://archive.ubuntu.com/ubuntu/"},
			SearchDNS: false,
		}}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(BeEmpty())
	})

	It("should fail with an invalid URI and invalid search for primary apt mirror", func() {
		config.APTConfig = &configv1alpha1.APTConfig{Primary: []configv1alpha1.APTArchive{{
			Arches:    []configv1alpha1.Architecture{configv1alpha1.ARM64},
			URI:       "packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
			Search:    []string{"archive.ubuntu.com/ubuntu/"},
			SearchDNS: false,
		}}}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(2))
	})

	It("should succeed with valid apt repositories", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with invalid apt repository", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "", URI: "not-a-valid-url"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(2))
	})

	It("should succeed with valid dependencies", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd.io", Version: "1.7.29-1~ubuntu.22.04~jammy", Hold: true},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with dependency missing name", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Version: "1.7.29-1~ubuntu.22.04~jammy"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeRequired))
		Expect(errs[0].Field).To(Equal("dependencies[0].name"))
	})

	It("should fail with invalid package name containing shell metacharacters", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "pkg\"; rm -rf / #"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("dependencies[0].name"))
	})

	It("should fail with invalid package version containing shell metacharacters", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", Version: "1.0\"; malicious_cmd #"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("dependencies[0].version"))
	})

	It("should fail with invalid ubuntu version format", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", UbuntuVersion: "22; rm -rf /"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("dependencies[0].ubuntuVersion"))
	})

	It("should fail with invalid build serial format", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", UbuntuBuildSerial: "20250725; echo pwned"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("dependencies[0].ubuntuBuildSerial"))
	})

	It("should succeed with build serial containing dots", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", UbuntuBuildSerial: "20261201.1"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with duplicate unconstrained dependencies", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", Version: "1.0.0"},
			{Name: "containerd", Version: "2.0.0"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeDuplicate))
		Expect(errs[0].Field).To(Equal("dependencies[1].name"))
	})

	It("should succeed with same dependency name but different ubuntu versions", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", Version: "1.0.0", UbuntuVersion: "22.04"},
			{Name: "containerd", Version: "2.0.0", UbuntuVersion: "24.04"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should succeed with one unconstrained and one constrained dependency", func() {
		config.Dependencies = []configv1alpha1.DependencyConfig{
			{Name: "containerd", Version: "1.0.0"},
			{Name: "containerd", Version: "2.0.0", UbuntuVersion: "24.04"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with duplicate repository names", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
			{Name: "docker", URI: "https://mirror.example.com/docker"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeDuplicate))
		Expect(errs[0].Field).To(Equal("aptRepositories[1].name"))
	})

	It("should succeed with apt repository using key", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", Key: "-----BEGIN PGP PUBLIC KEY BLOCK-----"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should succeed with apt repository using key URL", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "https://download.docker.com/linux/ubuntu/gpg"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail when both Key and KeyURL are set", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", Key: "-----BEGIN PGP PUBLIC KEY BLOCK-----", KeyURL: "https://example.com/gpg"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeForbidden))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].keyUrl"))
	})

	It("should fail with invalid KeyURL", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "not-a-valid-url"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].keyUrl"))
	})

	It("should succeed with keyUrl and keyFormat asc", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "https://download.docker.com/linux/ubuntu/gpg", KeyFormat: configv1alpha1.KeyFormatASC},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should succeed with keyUrl and keyFormat gpg", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "https://download.docker.com/linux/ubuntu/gpg", KeyFormat: configv1alpha1.KeyFormatGPG},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should not fail with uppercase keyFormat", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "https://download.docker.com/linux/ubuntu/gpg", KeyFormat: "ASC"},
		}
		Expect(ValidateExtensionConfig(config)).To(BeEmpty())
	})

	It("should fail with invalid keyFormat", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyURL: "https://download.docker.com/linux/ubuntu/gpg", KeyFormat: "binary"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeNotSupported))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].keyFormat"))
	})

	It("should fail when keyFormat is set without keyUrl", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyFormat: configv1alpha1.KeyFormatASC},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeForbidden))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].keyFormat"))
	})

	It("should fail with both forbidden and not-supported errors when keyFormat is invalid and keyUrl is unset", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "docker", URI: "https://download.docker.com/linux/ubuntu", KeyFormat: "binary"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(2))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeForbidden))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].keyFormat"))
		Expect(errs[1].Type).To(Equal(field.ErrorTypeNotSupported))
		Expect(errs[1].Field).To(Equal("aptRepositories[0].keyFormat"))
	})

	It("should fail with path traversal in apt repository name", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "../../etc/cron.d/evil", URI: "https://download.docker.com/linux/ubuntu"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].name"))
	})

	It("should fail with invalid characters in apt repository name", func() {
		config.AptRepositories = []configv1alpha1.AptRepository{
			{Name: "my repo", URI: "https://download.docker.com/linux/ubuntu"},
		}
		errs := ValidateExtensionConfig(config)
		Expect(errs).To(HaveLen(1))
		Expect(errs[0].Type).To(Equal(field.ErrorTypeInvalid))
		Expect(errs[0].Field).To(Equal("aptRepositories[0].name"))
	})
})
