// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package operatingsystemconfig_test

import (
	"context"
	_ "embed"

	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/utils/test"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
	. "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig"
)

//go:embed scripts/installNTP.sh
var ntpInstallScript string

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
			NTP: &v1alpha1.NTPConfig{
				Daemon: v1alpha1.SystemdTimesyncd,
			},
		}}
		actuator = NewActuator(mgr, false, extensionConfig)

		osc = &extensionsv1alpha1.OperatingSystemConfig{
			Spec: extensionsv1alpha1.OperatingSystemConfigSpec{
				Purpose: extensionsv1alpha1.OperatingSystemConfigPurposeProvision,
				Units:   []extensionsv1alpha1.Unit{{Name: "some-unit", Content: ptr.To("foo")}},
				Files:   []extensionsv1alpha1.File{{Path: "/some/file", Content: extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: "bar"}}}},
			},
		}
	})

	When("purpose is 'provision'", func() {
		expectedUserData := `#!/bin/bash
mkdir -p /etc/cloud/cloud.cfg.d/
cat <<EOF > /etc/cloud/cloud.cfg.d/custom-networking.cfg
network:
  config: disabled
EOF
chmod 0644 /etc/cloud/cloud.cfg.d/custom-networking.cfg

mkdir -p "/some"

cat << EOF | base64 -d > "/some/file"
YmFy
EOF


cat << EOF | base64 -d > "/etc/systemd/system/some-unit"
Zm9v
EOF
until apt-get update -qq && apt-get install --no-upgrade -qqy containerd runc docker.io socat nfs-common logrotate jq policykit-1; do sleep 1; done
ln -s /usr/bin/docker /bin/docker

if [ ! -s /etc/containerd/config.toml ]; then
  mkdir -p /etc/containerd/
  containerd config default > /etc/containerd/config.toml
  chmod 0644 /etc/containerd/config.toml
fi

mkdir -p /etc/systemd/system/containerd.service.d
cat <<EOF > /etc/systemd/system/containerd.service.d/11-exec_config.conf
[Service]
ExecStart=
ExecStart=/usr/bin/containerd --config=/etc/containerd/config.toml
EOF
chmod 0644 /etc/systemd/system/containerd.service.d/11-exec_config.conf

systemctl daemon-reload
systemctl enable containerd && systemctl restart containerd
systemctl enable docker && systemctl restart docker
systemctl enable 'some-unit' && systemctl restart --no-block 'some-unit'
`

		Describe("#Reconcile", func() {
			It("should not return an error", func() {
				userData, extensionUnits, extensionFiles, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(userData)).To(Equal(expectedUserData))
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
			})
		})

		Describe("#Reconcile with disabled unattended upgrades", func() {
			expectedUserData := `#!/bin/bash
mkdir -p /etc/cloud/cloud.cfg.d/
cat <<EOF > /etc/cloud/cloud.cfg.d/custom-networking.cfg
network:
  config: disabled
EOF
chmod 0644 /etc/cloud/cloud.cfg.d/custom-networking.cfg

mkdir -p "/some"

cat << EOF | base64 -d > "/some/file"
YmFy
EOF


cat << EOF | base64 -d > "/etc/systemd/system/some-unit"
Zm9v
EOF
until apt-get update -qq && apt-get install --no-upgrade -qqy containerd runc docker.io socat nfs-common logrotate jq policykit-1; do sleep 1; done
ln -s /usr/bin/docker /bin/docker

if [ ! -s /etc/containerd/config.toml ]; then
  mkdir -p /etc/containerd/
  containerd config default > /etc/containerd/config.toml
  chmod 0644 /etc/containerd/config.toml
fi

mkdir -p /etc/systemd/system/containerd.service.d
cat <<EOF > /etc/systemd/system/containerd.service.d/11-exec_config.conf
[Service]
ExecStart=
ExecStart=/usr/bin/containerd --config=/etc/containerd/config.toml
EOF
chmod 0644 /etc/systemd/system/containerd.service.d/11-exec_config.conf

mkdir -p /etc/apt/apt.conf.d
cat <<EOF > /etc/apt/apt.conf.d/99-auto-upgrades.conf
APT::Periodic::Unattended-Upgrade "0";
EOF
chmod 0644 /etc/apt/apt.conf.d/99-auto-upgrades.conf

systemctl daemon-reload
systemctl enable containerd && systemctl restart containerd
systemctl enable docker && systemctl restart docker
systemctl enable 'some-unit' && systemctl restart --no-block 'some-unit'
`
			It("should not return an error", func() {
				extensionConfig := Config{}
				actuator = NewActuator(mgr, true, extensionConfig)
				userData, extensionUnits, extensionFiles, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())

				Expect(string(userData)).To(Equal(expectedUserData))
				Expect(extensionUnits).To(BeEmpty())
				Expect(extensionFiles).To(BeEmpty())
			})
		})
	})

	When("purpose is 'reconcile'", func() {
		BeforeEach(func() {
			osc.Spec.Purpose = extensionsv1alpha1.OperatingSystemConfigPurposeReconcile
		})

		Describe("#Reconcile", func() {
			It("should not return an error", func() {
				userData, extensionUnits, extensionFiles, err := actuator.Reconcile(ctx, log, osc)
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
ExecStart=/bin/bash /opt/bin/install-ntp.sh systemd-timesyncd

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
					Path:        "/opt/bin/install-ntp.sh",
					Content:     extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: ntpInstallScript}},
					Permissions: ptr.To[uint32](0744),
				}))
			})
			It("should not return an error with ntp instead of systemd-timesyncd", func() {
				extensionConfig := Config{
					ExtensionConfig: &v1alpha1.ExtensionConfig{
						NTP: &v1alpha1.NTPConfig{
							Daemon: v1alpha1.NTPD,
							NTPD:   &v1alpha1.NTPDConfig{Servers: []string{"127.0.0.1"}},
						},
					},
				}
				actuator = NewActuator(mgr, true, extensionConfig)
				userData, extensionUnits, extensionFiles, err := actuator.Reconcile(ctx, log, osc)
				Expect(err).NotTo(HaveOccurred())
				Expect(userData).To(BeEmpty())

				Expect(extensionFiles).To(ContainElement(extensionsv1alpha1.File{
					Path:        "/opt/bin/install-ntp.sh",
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
ExecStart=/bin/bash /opt/bin/install-ntp.sh ntpd

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
