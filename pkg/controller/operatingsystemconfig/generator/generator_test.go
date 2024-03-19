// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package generator_test

import (
	commongen "github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/generator"
	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/generator/test"
	"github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig/generator"
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig/generator/testfiles"
)

var logger = logr.Discard()

var _ = Describe("Ubuntu OS Generator Test", func() {
	var (
		osc *commongen.OperatingSystemConfig
	)

	BeforeEach(func() {
		osc = &commongen.OperatingSystemConfig{
			Object: &v1alpha1.OperatingSystemConfig{
				Spec: v1alpha1.OperatingSystemConfigSpec{
					Purpose: v1alpha1.OperatingSystemConfigPurposeProvision,
				},
			},

			CRI:       &v1alpha1.CRIConfig{Name: v1alpha1.CRINameContainerD},
			Bootstrap: true,
		}
	})

	Describe("Conformance Tests", func() {
		g := generator.CloudInitGenerator(false)
		test.DescribeTest(generator.CloudInitGenerator(false), testfiles.Files)()

		It("should render correctly with Containerd enabled during Bootstrap (osc.type = provision)", func() {
			expectedCloudInit, err := testfiles.Files.ReadFile("cloud-init-containerd-provision")
			Expect(err).NotTo(HaveOccurred())
			expected := string(expectedCloudInit)

			osc.Units = []*commongen.Unit{
				{
					Name: "cloud-config-downloader.service",
				},
			}

			cloudInit, _, err := g.Generate(logger, osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})

		It("should render correctly with Containerd enabled but not during Bootstrap (osc.type = reconcile)", func() {
			expectedCloudInit, err := testfiles.Files.ReadFile("cloud-init-containerd-reconcile")
			Expect(err).NotTo(HaveOccurred())
			expected := string(expectedCloudInit)

			osc.Files = []*commongen.File{
				{
					Path:    "/path",
					Content: []byte("content"),
				},
			}

			osc.Bootstrap = false
			osc.Object.Spec.Purpose = v1alpha1.OperatingSystemConfigPurposeReconcile
			cloudInit, _, err := g.Generate(logger, osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})

		It("should render correctly with drop-in units", func() {
			expectedCloudInit, err := testfiles.Files.ReadFile("cloud-init-with-drop-in")
			expected := string(expectedCloudInit)

			Expect(err).NotTo(HaveOccurred())

			content := []byte(`[Service]
ExecStartPre=/opt/bin/init-containerd`)

			osc.Bootstrap = false
			osc.CRI = nil
			osc.Units = []*commongen.Unit{
				{
					Name:    "abc.service",
					Content: nil,
					DropIns: []*commongen.DropIn{
						{
							Name:    "10-exec-start-pre-init-config.conf",
							Content: content,
						},
						{
							Name:    "12-exec-start-pre-init-config.conf",
							Content: content,
						},
					},
				},
				{
					Name:    "mtu-customizer.service",
					Content: content,
				},
				{
					Name:    "other.service",
					Content: content,
				},
				{
					Name: "cloud-config-downloader.service",
				},
			}

			cloudInit, _, err := g.Generate(logger, osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})

		It("should render correctly with unattended upgrades are disabled (osc.type = provision)", func() {
			g := generator.CloudInitGenerator(true)
			expectedCloudInit, err := testfiles.Files.ReadFile("cloud-init-disabled-unattended-upgrades")
			Expect(err).NotTo(HaveOccurred())
			expected := string(expectedCloudInit)

			osc.Units = []*commongen.Unit{
				{
					Name: "cloud-config-downloader.service",
				},
			}

			cloudInit, _, err := g.Generate(logger, osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})
	})
})
