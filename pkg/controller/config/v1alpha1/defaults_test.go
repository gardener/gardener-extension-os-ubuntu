// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetDefaults_ExtensionConfig", func() {
	var config *ExtensionConfig

	BeforeEach(func() {
		config = &ExtensionConfig{}
	})

	JustBeforeEach(func() {
		SetDefaults_ExtensionConfig(config)
	})

	Context("with empty dependencies", func() {
		It("should add all required packages as unconstrained", func() {
			Expect(config.Dependencies).To(HaveLen(7))

			requiredPackages := []string{"containerd", "runc", "socat", "nfs-common", "logrotate", "jq", "policykit-1"}
			for _, pkg := range requiredPackages {
				Expect(config.Dependencies).To(ContainElement(And(
					WithTransform(func(d DependencyConfig) string { return d.Name }, Equal(pkg)),
					WithTransform(func(d DependencyConfig) string { return d.UbuntuVersion }, BeEmpty()),
					WithTransform(func(d DependencyConfig) string { return d.UbuntuBuildSerial }, BeEmpty()),
				)))
			}
		})
	})

	Context("with pinned dependency for required package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{
					Name:              "containerd",
					Version:           "1.7.29-1~ubuntu.22.04~jammy",
					UbuntuVersion:     "22.04",
					UbuntuBuildSerial: "20250725",
				},
			}
		})

		It("should add unconstrained fallback and other required packages", func() {
			Expect(config.Dependencies).To(HaveLen(8))

			containerdEntries := filterByName(config.Dependencies, "containerd")
			Expect(containerdEntries).To(HaveLen(2))
			Expect(containerdEntries).To(ContainElement(SatisfyAll(
				WithTransform(func(d DependencyConfig) string { return d.UbuntuVersion }, BeEmpty()),
				WithTransform(func(d DependencyConfig) string { return d.UbuntuBuildSerial }, BeEmpty()),
			)))
		})
	})

	Context("with unconstrained dependency for required package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "containerd"},
			}
		})

		It("should not add duplicate and should add other required packages", func() {
			Expect(config.Dependencies).To(HaveLen(7))

			containerdEntries := filterByName(config.Dependencies, "containerd")
			Expect(containerdEntries).To(HaveLen(1))
		})
	})

	Context("with custom unconstrained package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "custom-package"},
			}
		})

		It("should preserve custom package and add required packages", func() {
			Expect(config.Dependencies).To(HaveLen(8))
			Expect(filterByName(config.Dependencies, "custom-package")).To(HaveLen(1))
		})
	})

	Context("with custom pinned package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{
					Name:              "custom-package",
					Version:           "1.0.0",
					UbuntuVersion:     "22.04",
					UbuntuBuildSerial: "20250725",
				},
			}
		})

		It("should not add unconstrained fallback for custom package", func() {
			Expect(config.Dependencies).To(HaveLen(8))

			customEntries := filterByName(config.Dependencies, "custom-package")
			Expect(customEntries).To(HaveLen(1))
			Expect(customEntries[0].UbuntuVersion).NotTo(BeEmpty())
		})
	})

	Context("with custom versioned package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{
					Name:    "custom-package",
					Version: "1.0.0",
				},
			}
		})

		It("should not add additional unconstrained fallback", func() {
			Expect(config.Dependencies).To(HaveLen(8))

			customEntries := filterByName(config.Dependencies, "custom-package")
			Expect(customEntries).To(HaveLen(1))
			Expect(customEntries[0].Version).To(Equal("1.0.0"))
		})
	})

	Context("with disabled required package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "socat", Disabled: true},
			}
		})

		It("should not add fallback for disabled required package", func() {
			Expect(config.Dependencies).To(HaveLen(7))

			socatEntries := filterByName(config.Dependencies, "socat")
			Expect(socatEntries).To(HaveLen(1))
			Expect(socatEntries[0].Disabled).To(BeTrue())
		})
	})

	Context("with version-specific disabled required package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "policykit-1", UbuntuVersion: "26.04", Disabled: true},
			}
		})

		It("should not add unconstrained fallback for version-specific disabled required package", func() {
			Expect(config.Dependencies).To(HaveLen(7))

			policykitEntries := filterByName(config.Dependencies, "policykit-1")
			Expect(policykitEntries).To(HaveLen(1))
			Expect(policykitEntries[0].Disabled).To(BeTrue())
			Expect(policykitEntries[0].UbuntuVersion).To(Equal("26.04"))
		})
	})

	Context("with disabled required package and pinned entry", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{
					Name:              "containerd",
					Version:           "1.7.29-1",
					UbuntuVersion:     "22.04",
					UbuntuBuildSerial: "20250725",
				},
				{Name: "containerd", Disabled: true},
			}
		})

		It("should not add unconstrained fallback for disabled required package", func() {
			containerdEntries := filterByName(config.Dependencies, "containerd")
			Expect(containerdEntries).To(HaveLen(2))
			Expect(containerdEntries).NotTo(ContainElement(SatisfyAll(
				WithTransform(func(d DependencyConfig) string { return d.UbuntuVersion }, BeEmpty()),
				WithTransform(func(d DependencyConfig) string { return d.UbuntuBuildSerial }, BeEmpty()),
				WithTransform(func(d DependencyConfig) bool { return !d.Disabled }, BeTrue()),
			)))
		})
	})
})

func filterByName(deps []DependencyConfig, name string) []DependencyConfig {
	var result []DependencyConfig
	for _, dep := range deps {
		if dep.Name == name {
			result = append(result, dep)
		}
	}
	return result
}
