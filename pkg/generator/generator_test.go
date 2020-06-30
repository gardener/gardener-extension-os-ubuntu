// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator_test

import (
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/generator"

	commongen "github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/generator"
	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/generator/test"
	"github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gobuffalo/packr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		var box = packr.NewBox("./testfiles")
		g := generator.CloudInitGenerator()
		test.DescribeTest(generator.CloudInitGenerator(), box)()

		It("should render correctly with Containerd enabled during Bootstrap (osc.type = provision)", func() {
			expectedCloudInit, err := box.Find("cloud-init-containerd-provision")
			Expect(err).NotTo(HaveOccurred())
			expected := string(expectedCloudInit)

			cloudInit, _, err := g.Generate(osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})

		It("should render correctly with Containerd enabled but not during Bootstrap (osc.type = reconcile)", func() {
			expectedCloudInit, err := box.Find("cloud-init-containerd-reconcile")
			Expect(err).NotTo(HaveOccurred())
			expected := string(expectedCloudInit)
			osc.Bootstrap = false
			osc.Object.Spec.Purpose = v1alpha1.OperatingSystemConfigPurposeReconcile
			cloudInit, _, err := g.Generate(osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})

		It("should render correctly with drop-in units", func() {
			expectedCloudInit, err := box.Find("cloud-init-with-drop-in")
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
			}

			cloudInit, _, err := g.Generate(osc)

			Expect(err).NotTo(HaveOccurred())
			Expect(string(cloudInit)).To(Equal(expected))
		})
	})
})
