// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package operatingsystemconfig_test

import (
	"context"
	_ "embed"
	"path/filepath"
	"strings"

	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/utils/test"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
	. "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig"
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig/fixtures"
)

//go:embed scripts/installNTP.sh
var ntpInstallScript string

func expectUserDataToMatch(actual []byte, expected string) {
	GinkgoHelper()
	actualLines := strings.Split(string(actual), "\n")
	expectedLines := strings.Split(expected, "\n")
	if diff := cmp.Diff(expectedLines, actualLines); diff != "" {
		Fail("UserData mismatch (-want +got):\n" + diff)
	}
}

var _ = Describe("Actuator", func() {
	var (
		ctx        = context.TODO()
		log        = logr.Discard()
		fakeClient client.Client
		mgr        manager.Manager

		osc      *extensionsv1alpha1.OperatingSystemConfig
		actuator operatingsystemconfig.Actuator
	)

	BeforeEach(func() {
		fakeClient = fakeclient.NewClientBuilder().Build()
		mgr = test.FakeManager{Client: fakeClient}
		extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
			DisableUnattendedUpgrades: ptr.To(false),
			NTP: &v1alpha1.NTPConfig{
				Daemon: v1alpha1.SystemdTimesyncd,
			},
			Dependencies: []v1alpha1.DependencyConfig{
				{Name: "containerd"},
				{Name: "jq"},
				{Name: "logrotate"},
				{Name: "nfs-common"},
				{Name: "policykit-1"},
				{Name: "runc"},
				{Name: "socat"},
			},
		}}
		actuator = NewActuator(mgr, extensionConfig)

		osc = &extensionsv1alpha1.OperatingSystemConfig{
			Spec: extensionsv1alpha1.OperatingSystemConfigSpec{
				Purpose: extensionsv1alpha1.OperatingSystemConfigPurposeProvision,
				Units:   []extensionsv1alpha1.Unit{{Name: "some-unit", Content: ptr.To("foo")}},
				Files:   []extensionsv1alpha1.File{{Path: "/some/file", Content: extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: "bar"}}}},
			},
		}
	})

	When("purpose is 'provision'", func() {
		Describe("#Reconcile", func() {
			It("should not return an error", func() {
				userData, extensionUnits, extensionFiles, inplaceUpdateStatus, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataBasic)
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
				Expect(inplaceUpdateStatus).To(BeNil())
			})
		})

		Describe("#Reconcile with disabled unattended upgrades", func() {
			It("should not return an error", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(true),
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd"},
						{Name: "jq"},
						{Name: "logrotate"},
						{Name: "nfs-common"},
						{Name: "policykit-1"},
						{Name: "runc"},
						{Name: "socat"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, extensionUnits, extensionFiles, inplaceUpdateStatus, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataDisabledUnattendedUpgrades)
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
				Expect(inplaceUpdateStatus).To(BeNil())
			})
		})

		Describe("#Reconcile with custom apt config", func() {
			It("should not return an error", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd"},
						{Name: "jq"},
						{Name: "logrotate"},
						{Name: "nfs-common"},
						{Name: "policykit-1"},
						{Name: "runc"},
						{Name: "socat"},
					},
					APTConfig: &v1alpha1.APTConfig{
						PreserveSourcesList: false,
						Primary: []v1alpha1.APTArchive{
							v1alpha1.APTArchive{
								Arches: []v1alpha1.Architecture{v1alpha1.ArchDefault},
								URI:    "http://packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
							},
						},
						Security: []v1alpha1.APTArchive{
							v1alpha1.APTArchive{
								Arches: []v1alpha1.Architecture{v1alpha1.ArchDefault},
								URI:    "http://packages.ubuntu-mirror.example.com/apt-mirror/ubuntu",
							},
						},
					}}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, extensionUnits, extensionFiles, inplaceUpdateStatus, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataCustomAptConfig)
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
				Expect(inplaceUpdateStatus).To(BeNil())
			})
		})

		Describe("#Reconcile with docker apt repository and default dependencies", func() {
			It("should configure the docker apt source and install default packages", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(false),
					NTP: &v1alpha1.NTPConfig{
						Daemon: v1alpha1.SystemdTimesyncd,
					},
					AptRepositories: []v1alpha1.AptRepository{
						{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
						{Name: "runc"},
						{Name: "socat"},
						{Name: "nfs-common"},
						{Name: "logrotate"},
						{Name: "jq"},
						{Name: "policykit-1"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, extensionUnits, extensionFiles, inplaceUpdateStatus, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataDockerRepo)
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
				Expect(inplaceUpdateStatus).To(BeNil())
			})
		})

		Describe("#Reconcile with custom apt repository mirror", func() {
			It("should use the provided mirror URI", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					AptRepositories: []v1alpha1.AptRepository{
						{Name: "docker", URI: "http://mirror.example.com/linux/ubuntu"},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataCustomMirror)
			})
		})

		Describe("#Reconcile with apt repository and GPG key", func() {
			It("should include signed-by option when GPG key URL is provided", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					AptRepositories: []v1alpha1.AptRepository{
						{
							Name:   "docker",
							URI:    "https://download.docker.com/linux/ubuntu",
							KeyURL: "https://download.docker.com/linux/ubuntu/gpg",
						},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataGpgKey)
			})

			It("should include signed-by option with binary key format when GPG key URL is provided", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					AptRepositories: []v1alpha1.AptRepository{
						{
							Name:      "docker",
							URI:       "https://download.docker.com/linux/ubuntu",
							KeyURL:    "https://download.docker.com/linux/ubuntu/gpg",
							KeyFormat: v1alpha1.KeyFormatGPG,
						},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataGpgKeyBinary)
			})

			It("should include signed-by option with ascii-armored key format when GPG key URL is provided", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					AptRepositories: []v1alpha1.AptRepository{
						{
							Name:      "docker",
							URI:       "https://download.docker.com/linux/ubuntu",
							KeyURL:    "https://download.docker.com/linux/ubuntu/gpg",
							KeyFormat: v1alpha1.KeyFormatASC,
						},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataGpgKey)
			})

			It("should include signed-by option when inline GPG key is provided", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					AptRepositories: []v1alpha1.AptRepository{
						{
							Name: "docker",
							URI:  "https://download.docker.com/linux/ubuntu",
							Key: `-----BEGIN PGP PUBLIC KEY BLOCK-----
mQINBFit2ioBEADhWpZ8/wvZ6hUTiXOwQHXMAlaFHcPH9hAtr4F1y2+OYdbtMuth
lO49r1nayzQb8T14bZf2DdpOjXn8b5XrZQ9JZ5hZ1XyZ5XyZ5XyZ5XyZ5XyZ5XyZ
-----END PGP PUBLIC KEY BLOCK-----`,
						},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{Name: "containerd.io"},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataGpgKeyInline)
			})
		})

		Describe("#Reconcile with pinned dependencies", func() {
			It("should generate version and build serial specific install commands", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(false),
					NTP: &v1alpha1.NTPConfig{
						Daemon: v1alpha1.SystemdTimesyncd,
					},
					AptRepositories: []v1alpha1.AptRepository{
						{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{
							Name:              "containerd.io",
							Version:           "1.7.29-1~ubuntu.22.04~jammy",
							UbuntuVersion:     "22.04",
							UbuntuBuildSerial: "20250725",
							Hold:              true,
						},
						{
							Name:              "containerd.io",
							Version:           "2.2.4-1~ubuntu.26.04~resolute",
							UbuntuVersion:     "26.04",
							UbuntuBuildSerial: "20260520",
							Hold:              true,
						},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataPinnedDeps)
			})
		})

		Describe("#Reconcile with pinned dependencies and fallback", func() {
			It("should generate version and build serial specific install commands with a fallback", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(false),
					NTP: &v1alpha1.NTPConfig{
						Daemon: v1alpha1.SystemdTimesyncd,
					},
					AptRepositories: []v1alpha1.AptRepository{
						{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{
							Name:              "containerd.io",
							Version:           "1.7.29-1~ubuntu.22.04~jammy",
							UbuntuVersion:     "22.04",
							UbuntuBuildSerial: "20250725",
							Hold:              true,
						},
						{
							Name:              "containerd.io",
							Version:           "2.2.4-1~ubuntu.26.04~resolute",
							UbuntuVersion:     "26.04",
							UbuntuBuildSerial: "20260520",
							Hold:              true,
						},
						{
							Name: "containerd.io",
						},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataPinnedDepsWithFallback)
			})

			It("should handle fallback regardless of dependency order", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(false),
					NTP: &v1alpha1.NTPConfig{
						Daemon: v1alpha1.SystemdTimesyncd,
					},
					AptRepositories: []v1alpha1.AptRepository{
						{Name: "docker", URI: "https://download.docker.com/linux/ubuntu"},
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{
							Name: "containerd.io",
						},
						{
							Name:              "containerd.io",
							Version:           "1.7.29-1~ubuntu.22.04~jammy",
							UbuntuVersion:     "22.04",
							UbuntuBuildSerial: "20250725",
							Hold:              true,
						},
						{
							Name:              "containerd.io",
							Version:           "2.2.4-1~ubuntu.26.04~resolute",
							UbuntuVersion:     "26.04",
							UbuntuBuildSerial: "20260520",
							Hold:              true,
						},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataPinnedDepsWithFallback)
			})
		})

		Describe("#Reconcile with pinned dependencies using different names per OS version", func() {
			It("should not exit on unmatched constraints and continue with remaining packages", func() {
				extensionConfig := Config{ExtensionConfig: &v1alpha1.ExtensionConfig{
					DisableUnattendedUpgrades: ptr.To(false),
					NTP: &v1alpha1.NTPConfig{
						Daemon: v1alpha1.SystemdTimesyncd,
					},
					Dependencies: []v1alpha1.DependencyConfig{
						{
							Name:          "containerd.io",
							Version:       "1.7.29-1~ubuntu.22.04~jammy",
							UbuntuVersion: "22.04",
						},
						{
							Name:          "containerd",
							Version:       "2.2.4-1~ubuntu.26.04~resolute",
							UbuntuVersion: "26.04",
						},
					},
				}}
				actuator = NewActuator(mgr, extensionConfig)
				userData, _, _, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				expectUserDataToMatch(userData, fixtures.UserDataPinnedDepsDifferentNames)
			})
		})
	})

	When("purpose is 'reconcile'", func() {
		BeforeEach(func() {
			osc.Spec.Purpose = extensionsv1alpha1.OperatingSystemConfigPurposeReconcile
		})

		Describe("#Reconcile", func() {
			It("should not return an error", func() {
				userData, extensionUnits, extensionFiles, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				Expect(userData).To(BeEmpty())
				Expect(extensionUnits).To(ConsistOf(
					extensionsv1alpha1.Unit{
						Name: "kubelet.service",
						DropIns: []extensionsv1alpha1.DropIn{{
							Name: "10-configure-resolv-conf.conf",
							Content: `[Service]
ExecStartPre=/opt/gardener/bin/configure_kubelet_resolv_conf.sh
`,
						}},
						FilePaths: []string{"/opt/gardener/bin/configure_kubelet_resolv_conf.sh"},
					},
					extensionsv1alpha1.Unit{
						Name:    "install-ntp-client.service",
						Command: ptr.To(extensionsv1alpha1.CommandRestart),
						Content: ptr.To(`[Unit]
Description=Oneshot service to install requested ntp client

[Service]
Type=oneshot
ExecStart=/bin/bash /opt/gardener/bin/install-ntp.sh systemd-timesyncd

[Install]
WantedBy=multi-user.target
`),
					},
				),
				)
				Expect(extensionFiles).To(ContainElement(extensionsv1alpha1.File{
					Path:        "/opt/gardener/bin/configure_kubelet_resolv_conf.sh",
					Permissions: ptr.To[uint32](0755),
					Content: extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: `#!/bin/bash
if grep -q 'resolvConf: /etc/resolv.conf' /var/lib/kubelet/config/kubelet; then
  sed -i -e 's|resolvConf: /etc/resolv.conf|resolvConf: /run/systemd/resolve/resolv.conf|g' /var/lib/kubelet/config/kubelet;
fi
`}},
				}))
				Expect(extensionFiles).To(ContainElement(extensionsv1alpha1.File{
					Path:        "/opt/gardener/bin/install-ntp.sh",
					Content:     extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: ntpInstallScript}},
					Permissions: ptr.To[uint32](0744),
				}))
			})
			It("should enable ntpd service", func() {
				extensionConfig := Config{
					ExtensionConfig: &v1alpha1.ExtensionConfig{
						NTP: &v1alpha1.NTPConfig{
							Daemon: v1alpha1.NTPD,
							NTPD: &v1alpha1.NTPDConfig{
								Servers:    []string{"foo.bar", "bar.foo"},
								Interfaces: []string{"dev1", "dev2"},
							},
						},
					},
				}
				actuator = NewActuator(mgr, extensionConfig)
				userData, extensionUnits, extensionFiles, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())
				Expect(userData).To(BeEmpty())
				Expect(extensionUnits).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Name": Equal("install-ntp-client.service"),
				})))
				Expect(extensionFiles).To(ContainElement(extensionsv1alpha1.File{
					Path:        filepath.Join(string(filepath.Separator), "etc", "ntp.conf"),
					Permissions: ptr.To[uint32](0644),
					Content: extensionsv1alpha1.FileContent{
						Inline: &extensionsv1alpha1.FileContentInline{
							Data: `
server foo.bar iburst
server bar.foo iburst

driftfile /var/lib/ntp/ntp.drift
restrict default nomodify nopeer noquery notrap limited kod
restrict 127.0.0.1
restrict [::1]

interface ignore wildcard
interface listen 127.0.0.1
interface listen dev1
interface listen dev2
`,
						},
					},
				}))
			})
			It("should not return an error with ntp instead of systemd-timesyncd", func() {
				extensionConfig := Config{
					ExtensionConfig: &v1alpha1.ExtensionConfig{
						DisableUnattendedUpgrades: ptr.To(true),
						NTP: &v1alpha1.NTPConfig{
							Daemon: v1alpha1.NTPD,
							NTPD:   &v1alpha1.NTPDConfig{Servers: []string{"127.0.0.1"}},
						},
					},
				}
				actuator = NewActuator(mgr, extensionConfig)
				userData, extensionUnits, extensionFiles, _, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())
				Expect(userData).To(BeEmpty())

				Expect(extensionFiles).To(ContainElement(extensionsv1alpha1.File{
					Path:        "/opt/gardener/bin/install-ntp.sh",
					Content:     extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: ntpInstallScript}},
					Permissions: ptr.To[uint32](0744),
				}))

				Expect(extensionUnits).To(ConsistOf(
					extensionsv1alpha1.Unit{
						Name: "kubelet.service",
						DropIns: []extensionsv1alpha1.DropIn{{
							Name: "10-configure-resolv-conf.conf",
							Content: `[Service]
ExecStartPre=/opt/gardener/bin/configure_kubelet_resolv_conf.sh
`,
						}},
						FilePaths: []string{"/opt/gardener/bin/configure_kubelet_resolv_conf.sh"},
					},
					extensionsv1alpha1.Unit{
						Name:    "install-ntp-client.service",
						Command: ptr.To(extensionsv1alpha1.CommandRestart),
						Content: ptr.To(`[Unit]
Description=Oneshot service to install requested ntp client

[Service]
Type=oneshot
ExecStart=/bin/bash /opt/gardener/bin/install-ntp.sh ntpd

[Install]
WantedBy=multi-user.target
`),
					},
				),
				)
			})
		})
	})
})
