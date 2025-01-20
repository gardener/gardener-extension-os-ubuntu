// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package operatingsystemconfig

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
)

//go:embed templates/ntp-config.conf.tpl
var ntpConfigTemplateContent string

//go:embed scripts/installNTP.sh
var ntpInstallScript string

var ntpConfigTemplate *template.Template

type actuator struct {
	client                    client.Client
	disableUnattendedUpgrades bool
	extensionConfig           Config
}

// Config contains configuration for the extension service.
type Config struct {
	// Embed the entire Extension config here for direct access in the controller.
	*configv1alpha1.ExtensionConfig
}

func init() {
	var err error
	ntpConfigTemplate, err = template.New("ntp-config").Parse(ntpConfigTemplateContent)
	if err != nil {
		panic(fmt.Errorf("failed to parse NTP config template: %w", err))
	}
}

// NewActuator creates a new Actuator that updates the status of the handled OperatingSystemConfig resources.
func NewActuator(mgr manager.Manager, disableUnattendedUpgrades bool, extensionConfig Config) operatingsystemconfig.Actuator {
	return &actuator{
		client:                    mgr.GetClient(),
		disableUnattendedUpgrades: disableUnattendedUpgrades,
		extensionConfig:           extensionConfig,
	}
}

func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	switch purpose := osc.Spec.Purpose; purpose {
	case extensionsv1alpha1.OperatingSystemConfigPurposeProvision:
		userData, err := a.handleProvisionOSC(ctx, osc)
		return []byte(userData), nil, nil, err

	case extensionsv1alpha1.OperatingSystemConfigPurposeReconcile:
		extensionUnits, extensionFiles, err := a.handleReconcileOSC(osc)
		return nil, extensionUnits, extensionFiles, err

	default:
		return nil, nil, nil, fmt.Errorf("unknown purpose: %s", purpose)
	}
}

func (a *actuator) Delete(_ context.Context, _ logr.Logger, _ *extensionsv1alpha1.OperatingSystemConfig) error {
	return nil
}

func (a *actuator) Migrate(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) error {
	return a.Delete(ctx, log, osc)
}

func (a *actuator) ForceDelete(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) error {
	return a.Delete(ctx, log, osc)
}

func (a *actuator) Restore(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	return a.Reconcile(ctx, log, osc)
}

func (a *actuator) handleProvisionOSC(ctx context.Context, osc *extensionsv1alpha1.OperatingSystemConfig) (string, error) {
	writeFilesToDiskScript, err := operatingsystemconfig.FilesToDiskScript(ctx, a.client, osc.Namespace, osc.Spec.Files)
	if err != nil {
		return "", err
	}
	writeUnitsToDiskScript := operatingsystemconfig.UnitsToDiskScript(osc.Spec.Units)

	script := `#!/bin/bash
mkdir -p /etc/cloud/cloud.cfg.d/
cat <<EOF > /etc/cloud/cloud.cfg.d/custom-networking.cfg
network:
  config: disabled
EOF
chmod 0644 /etc/cloud/cloud.cfg.d/custom-networking.cfg
` + writeFilesToDiskScript + `
` + writeUnitsToDiskScript + `
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
` + disableUnattendedUpgradesScript(a.disableUnattendedUpgrades) + `
systemctl daemon-reload
systemctl enable containerd && systemctl restart containerd
systemctl enable docker && systemctl restart docker
`

	for _, unit := range osc.Spec.Units {
		script += fmt.Sprintf(`systemctl enable '%s' && systemctl restart --no-block '%s'
`, unit.Name, unit.Name)
	}

	return script, nil
}

func (a *actuator) generateNTPConfig() (string, error) {
	templateData := a.extensionConfig.NTP.NTPD
	var templateOutput strings.Builder

	err := ntpConfigTemplate.Execute(&templateOutput, templateData)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return templateOutput.String(), nil
}

func (a *actuator) handleReconcileOSC(_ *extensionsv1alpha1.OperatingSystemConfig) ([]extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	var (
		extensionUnits []extensionsv1alpha1.Unit
		extensionFiles []extensionsv1alpha1.File
	)

	var err error
	if extensionUnits, extensionFiles, err = a.configureNTPDaemon(extensionUnits, extensionFiles); err != nil {
		return nil, nil, fmt.Errorf("error configuring NTP Daemon: %v", err)
	}

	// add scripts and dropins for kubelet
	filePathKubeletConfigureResolvConfScript := filepath.Join("/", "opt", "gardener", "bin", "configure_kubelet_resolv_conf.sh")
	extensionFiles = append(extensionFiles, extensionsv1alpha1.File{
		Path:        filePathKubeletConfigureResolvConfScript,
		Permissions: ptr.To[uint32](0755),
		Content: extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: `#!/bin/bash
if grep -q 'resolvConf: /etc/resolv.conf' /var/lib/kubelet/config/kubelet; then
  sed -i -e 's|resolvConf: /etc/resolv.conf|resolvConf: /run/systemd/resolve/resolv.conf|g' /var/lib/kubelet/config/kubelet;
fi
`}},
	})
	extensionUnits = append(extensionUnits, extensionsv1alpha1.Unit{
		Name: "kubelet.service",
		DropIns: []extensionsv1alpha1.DropIn{{
			Name: "10-configure-resolv-conf.conf",
			Content: `[Service]
ExecStartPre=` + filePathKubeletConfigureResolvConfScript + `
`,
		}},
		FilePaths: []string{filePathKubeletConfigureResolvConfScript},
	})

	return extensionUnits, extensionFiles, nil
}

func disableUnattendedUpgradesScript(disableAutoUpgrades bool) string {
	if disableAutoUpgrades {
		return `
mkdir -p /etc/apt/apt.conf.d
cat <<EOF > /etc/apt/apt.conf.d/99-auto-upgrades.conf
APT::Periodic::Unattended-Upgrade "0";
EOF
chmod 0644 /etc/apt/apt.conf.d/99-auto-upgrades.conf
`
	}
	return ""
}

// configureNTPDaemon configures the VM either with systemd-timesyncd or ntpd as the time syncing client
func (a *actuator) configureNTPDaemon(extensionUnits []extensionsv1alpha1.Unit, extensionFiles []extensionsv1alpha1.File) ([]extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	filePathNTPScript := filepath.Join(string(filepath.Separator), "opt", "bin", "install-ntp.sh")
	extensionFiles = append(extensionFiles, extensionsv1alpha1.File{
		Path:        filePathNTPScript,
		Content:     extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: ntpInstallScript}},
		Permissions: ptr.To[uint32](0744),
	})

	switch a.extensionConfig.NTP.Daemon {
	case configv1alpha1.SystemdTimesyncd:
	case configv1alpha1.NTPD:
		templateData, err := a.generateNTPConfig()
		if err != nil {
			return nil, nil, fmt.Errorf("error generating NTP config: %v", err)
		}
		extensionFiles = append(extensionFiles, extensionsv1alpha1.File{
			Path:        filepath.Join(string(filepath.Separator), "etc", "ntp.conf"),
			Content:     extensionsv1alpha1.FileContent{Inline: &extensionsv1alpha1.FileContentInline{Data: templateData}},
			Permissions: ptr.To[uint32](0644),
		})
	default:
		return nil, nil, fmt.Errorf("unsupported NTP daemon: %s", a.extensionConfig.NTP.Daemon)
	}

	extensionUnits = append(extensionUnits, extensionsv1alpha1.Unit{
		Name: "install-ntp-client.service",
		Content: ptr.To(`[Unit]
Description=Oneshot service to install requested ntp client

[Service]
Type=oneshot
ExecStart=` + fmt.Sprintf("/bin/bash %s %s", filePathNTPScript, a.extensionConfig.NTP.Daemon) + `

[Install]
WantedBy=multi-user.target
`),
		Command: ptr.To(extensionsv1alpha1.CommandRestart),
	})

	return extensionUnits, extensionFiles, nil
}
