// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package operatingsystemconfig

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig"
	oscommonactuator "github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/actuator"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig/generator"
)

type actuator struct {
	client                    client.Client
	useGardenerNodeAgent      bool
	disableUnattendedUpgrades bool
}

// NewActuator creates a new Actuator that updates the status of the handled OperatingSystemConfig resources.
func NewActuator(mgr manager.Manager, useGardenerNodeAgent bool, disableUnattendedUpgrades bool) operatingsystemconfig.Actuator {
	return &actuator{
		client:                    mgr.GetClient(),
		useGardenerNodeAgent:      useGardenerNodeAgent,
		disableUnattendedUpgrades: disableUnattendedUpgrades,
	}
}

func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, *string, []string, []string, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	cloudConfig, command, err := oscommonactuator.CloudConfigFromOperatingSystemConfig(ctx, log, a.client, osc, generator.CloudInitGenerator())
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("could not generate cloud config: %w", err)
	}

	switch purpose := osc.Spec.Purpose; purpose {
	case extensionsv1alpha1.OperatingSystemConfigPurposeProvision:
		if !a.useGardenerNodeAgent {
			return cloudConfig, command, oscommonactuator.OperatingSystemConfigUnitNames(osc), oscommonactuator.OperatingSystemConfigFilePaths(osc), nil, nil, nil
		}
		userData, err := a.handleProvisionOSC(ctx, osc)
		return []byte(userData), nil, nil, nil, nil, nil, err

	case extensionsv1alpha1.OperatingSystemConfigPurposeReconcile:
		extensionUnits, extensionFiles, err := a.handleReconcileOSC(osc)
		return cloudConfig, command, oscommonactuator.OperatingSystemConfigUnitNames(osc), oscommonactuator.OperatingSystemConfigFilePaths(osc), extensionUnits, extensionFiles, err

	default:
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("unknown purpose: %s", purpose)
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

func (a *actuator) Restore(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, *string, []string, []string, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
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

func (a *actuator) handleReconcileOSC(_ *extensionsv1alpha1.OperatingSystemConfig) ([]extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	var (
		extensionUnits []extensionsv1alpha1.Unit
		extensionFiles []extensionsv1alpha1.File
	)

	// add scripts and dropins for kubelet
	filePathKubeletConfigureResolvConfScript := filepath.Join("/", "opt", "gardener", "bin", "configure_kubelet_resolv_conf.sh")
	extensionFiles = append(extensionFiles, extensionsv1alpha1.File{
		Path:        filePathKubeletConfigureResolvConfScript,
		Permissions: ptr.To[int32](0755),
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
