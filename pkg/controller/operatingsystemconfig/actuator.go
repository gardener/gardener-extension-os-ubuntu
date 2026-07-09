// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package operatingsystemconfig

import (
	"context"
	_ "embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/yaml"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/internal"
)

//go:embed templates/ntp-config.conf.tpl
var ntpConfigTemplateContent string

//go:embed scripts/installNTP.sh
var ntpInstallScript string

var ntpConfigTemplate *template.Template

//go:embed templates/install-dependencies.sh.tpl
var installDependenciesTemplateContent string
var installDependenciesTemplate *template.Template

type packageInstruction struct {
	Name     string
	Blocks   []conditionalBlock
	Fallback *installCommand
}

type conditionalBlock struct {
	Condition string
	Command   installCommand
}

type installCommand struct {
	Version string
	Hold    bool
}

type actuator struct {
	client          client.Client
	extensionConfig Config
}

// Config contains configuration for the extension service.
type Config struct {
	// Embed the entire Extension config here for direct access in the controller.
	*configv1alpha1.ExtensionConfig
}

func init() {
	var err error
	ntpConfigTemplate, err = template.New("ntp-config").Funcs(sprig.TxtFuncMap()).Parse(ntpConfigTemplateContent)
	if err != nil {
		panic(fmt.Errorf("failed to parse NTP config template: %w", err))
	}

	installDependenciesTemplate, err = template.New("install-deps").Parse(installDependenciesTemplateContent)
	if err != nil {
		panic(fmt.Errorf("failed to parse install dependencies template: %w", err))
	}
}

// NewActuator creates a new Actuator that updates the status of the handled OperatingSystemConfig resources.
func NewActuator(mgr manager.Manager, extensionConfig Config) operatingsystemconfig.Actuator {
	return &actuator{
		client:          mgr.GetClient(),
		extensionConfig: extensionConfig,
	}
}

func (a *actuator) Reconcile(ctx context.Context, _ logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, *extensionsv1alpha1.InPlaceUpdatesStatus, error) {
	switch purpose := osc.Spec.Purpose; purpose {
	case extensionsv1alpha1.OperatingSystemConfigPurposeProvision:
		userData, err := a.handleProvisionOSC(ctx, osc)
		return []byte(userData), nil, nil, nil, err

	case extensionsv1alpha1.OperatingSystemConfigPurposeReconcile:
		extensionUnits, extensionFiles, err := a.handleReconcileOSC(osc)
		return nil, extensionUnits, extensionFiles, nil, err

	default:
		return nil, nil, nil, nil, fmt.Errorf("unknown purpose: %s", purpose)
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

func (a *actuator) Restore(ctx context.Context, log logr.Logger, osc *extensionsv1alpha1.OperatingSystemConfig) ([]byte, []extensionsv1alpha1.Unit, []extensionsv1alpha1.File, *extensionsv1alpha1.InPlaceUpdatesStatus, error) {
	return a.Reconcile(ctx, log, osc)
}

func (a *actuator) handleProvisionOSC(ctx context.Context, osc *extensionsv1alpha1.OperatingSystemConfig) (string, error) {
	writeFilesToDiskScript, err := operatingsystemconfig.FilesToDiskScript(ctx, a.client, osc.Namespace, osc.Spec.Files)
	if err != nil {
		return "", err
	}
	writeUnitsToDiskScript := operatingsystemconfig.UnitsToDiskScript(osc.Spec.Units)

	installScript := a.generateInstallDependenciesScript()

	script := `#!/bin/bash
mkdir -p /etc/cloud/cloud.cfg.d/
cat <<EOF > /etc/cloud/cloud.cfg.d/custom-networking.cfg
network:
  config: disabled
EOF
chmod 0644 /etc/cloud/cloud.cfg.d/custom-networking.cfg
` + writeFilesToDiskScript + `
` + writeUnitsToDiskScript + `
` + installScript + `

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
` + disableUnattendedUpgradesScript(a.extensionConfig.DisableUnattendedUpgrades) + `
systemctl daemon-reload
systemctl enable containerd && systemctl restart containerd
`

	for _, unit := range osc.Spec.Units {
		script += fmt.Sprintf(`systemctl enable '%s' && systemctl restart --no-block '%s'
`, unit.Name, unit.Name)
	}

	script = operatingsystemconfig.WrapProvisionOSCIntoOneshotScript(script)

	parts := []internal.FilePart{
		{
			Type:    "text/x-shellscript",
			Content: script,
		},
	}

	if a.extensionConfig.APTConfig != nil || len(a.extensionConfig.AptRepositories) > 0 {
		aptConfig, err := a.createAPTCloudConfig()
		if err != nil {
			return "", err
		}
		parts = append(parts, aptConfig)
	}
	header := "#cloud-config-archive\n"
	yamlBody, err := yaml.Marshal(parts)
	if err != nil {
		return "", fmt.Errorf("error marshalling archives: %v", err)
	}
	archive := header + string(yamlBody)

	return archive, nil
}

func (a *actuator) createAPTCloudConfig() (internal.FilePart, error) {
	aptConfig := internal.APTConfig{}
	aptCloudConfig := internal.FilePart{
		Type: "text/cloud-config",
	}
	if a.extensionConfig.APTConfig != nil {
		aptConfig.PreserveSourcesList = a.extensionConfig.APTConfig.PreserveSourcesList
		aptConfig.Primary = make([]internal.APTArchive, 0, len(a.extensionConfig.APTConfig.Primary))
		for _, primary := range a.extensionConfig.APTConfig.Primary {
			aptConfig.Primary = append(aptConfig.Primary, internal.APTArchive{
				Arches:    primary.Arches,
				URI:       primary.URI,
				Search:    primary.Search,
				SearchDNS: primary.SearchDNS,
			})
		}
		aptConfig.Security = make([]internal.APTArchive, 0, len(a.extensionConfig.APTConfig.Security))
		for _, security := range a.extensionConfig.APTConfig.Security {
			aptConfig.Security = append(aptConfig.Security, internal.APTArchive{
				Arches:    security.Arches,
				URI:       security.URI,
				Search:    security.Search,
				SearchDNS: security.SearchDNS,
			})
		}
	}

	if len(a.extensionConfig.AptRepositories) > 0 {
		aptConfig.Sources = make(map[string]internal.APTSource, len(a.extensionConfig.AptRepositories))
		for _, repo := range a.extensionConfig.AptRepositories {
			key := repo.Key
			suite := repo.Suite
			if suite == "" {
				suite = "$RELEASE"
			}
			components := repo.Components
			if len(components) == 0 {
				components = []string{"stable"}
			}
			source := fmt.Sprintf("deb [signed-by=$KEY_FILE] %s %s %s", repo.URI, suite, strings.Join(components, " "))
			aptConfig.Sources[repo.Name] = internal.APTSource{
				Source: source,
				Key:    key,
			}
		}
	}

	cloudInitApt := internal.APTCloudInit{APT: aptConfig}
	cloudInitAptYaml, err := yaml.Marshal(cloudInitApt)
	if err != nil {
		return aptCloudConfig, fmt.Errorf("failed to marshal cloud-init apt config to yaml: %w", err)
	}
	aptCloudConfig.Content = "#cloud-config\n" + string(cloudInitAptYaml)

	return aptCloudConfig, nil
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

func disableUnattendedUpgradesScript(disableAutoUpgrades *bool) string {
	if disableAutoUpgrades != nil && *disableAutoUpgrades {
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

func (a *actuator) generateInstallDependenciesScript() string {
	// Group dependencies by package name
	byName := make(map[string][]configv1alpha1.DependencyConfig)
	for _, dep := range a.extensionConfig.Dependencies {
		byName[dep.Name] = append(byName[dep.Name], dep)
	}

	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	var instructions []packageInstruction

	// Build data model for the template
	for _, name := range names {
		instruction := packageInstruction{Name: name}

		for _, dep := range byName[name] {
			isUnconstrained := dep.UbuntuVersion == "" && dep.UbuntuBuildSerial == ""

			if isUnconstrained {
				instruction.Fallback = &installCommand{
					Version: dep.Version,
					Hold:    dep.Hold,
				}
				break
			}

			var constraints []string
			if dep.UbuntuVersion != "" {
				constraints = append(constraints, fmt.Sprintf(`(-z "$UBUNTU_VERSION" || "%s" == "$UBUNTU_VERSION")`, dep.UbuntuVersion))
			}
			if dep.UbuntuBuildSerial != "" {
				constraints = append(constraints, fmt.Sprintf(`(-z "$BUILD_SERIAL" || "%s" == "$BUILD_SERIAL")`, dep.UbuntuBuildSerial))
			}

			instruction.Blocks = append(instruction.Blocks, conditionalBlock{
				Condition: strings.Join(constraints, " && "),
				Command: installCommand{
					Version: dep.Version,
					Hold:    dep.Hold,
				},
			})
		}
		instructions = append(instructions, instruction)
	}

	var sb strings.Builder
	err := installDependenciesTemplate.Execute(&sb, instructions)
	if err != nil {
		return fmt.Sprintf("echo 'Template generation failed: %v'", err)
	}

	return strings.TrimSpace(sb.String())
}

// configureNTPDaemon configures the VM either with systemd-timesyncd or ntpd as the time syncing client
func (a *actuator) configureNTPDaemon(extensionUnits []extensionsv1alpha1.Unit, extensionFiles []extensionsv1alpha1.File) ([]extensionsv1alpha1.Unit, []extensionsv1alpha1.File, error) {
	filePathNTPScript := filepath.Join(string(filepath.Separator), "opt", "gardener", "bin", "install-ntp.sh")
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
