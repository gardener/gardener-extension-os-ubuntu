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

	Context("with nil dependencies", func() {
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

	Context("with explicitly empty dependencies", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{}
		})

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(BeEmpty())
		})
	})

	Context("with pinned dependency for containerd", func() {
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

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(HaveLen(1))
		})
	})

	Context("with unconstrained dependency for containerd", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "containerd"},
			}
		})

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(HaveLen(1))
		})
	})

	Context("with custom unconstrained package", func() {
		BeforeEach(func() {
			config.Dependencies = []DependencyConfig{
				{Name: "custom-package"},
			}
		})

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(HaveLen(1))
			Expect(config.Dependencies[0].Name).To(Equal("custom-package"))
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

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(HaveLen(1))
			Expect(config.Dependencies[0].UbuntuVersion).NotTo(BeEmpty())
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

		It("should not add any defaults", func() {
			Expect(config.Dependencies).To(HaveLen(1))
			Expect(config.Dependencies[0].Version).To(Equal("1.0.0"))
		})
	})
})
