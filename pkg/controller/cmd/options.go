package cmd

import (
	"fmt"
	"os"

	configv1alpha1 "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1/validation"
	"github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/operatingsystemconfig"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/ptr"
)

var (
	// DisableUnattendedUpgrades is the name of the command line flag to disable unattended upgrades in ubuntu.
	DisableUnattendedUpgrades = "disable-unattended-upgrades"
	configDecoder             runtime.Decoder
	// config is the parsed configFile
	extensionConfig *configv1alpha1.ExtensionConfig
)

func init() {
	configScheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(
		configv1alpha1.AddToScheme,
	)
	utilruntime.Must(schemeBuilder.AddToScheme(configScheme))
	configDecoder = serializer.NewCodecFactory(configScheme).UniversalDecoder()
}

// UbuntuOptions are command line options that can be set for ubuntu configuration.
type UbuntuOptions struct {
	// DisableUnattendedUpgrades is the flag to disable unattended upgrades in ubuntu.
	DisableUnattendedUpgrades bool
	// configFile path of the extension config
	configFile string
}

// AddFlags implements cmd.Option.
func (u *UbuntuOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&u.DisableUnattendedUpgrades, DisableUnattendedUpgrades, u.DisableUnattendedUpgrades, "The flag to disable unattended upgrades in ubuntu.")
	fs.StringVar(&u.configFile, "config", u.configFile, "Path to configuration file.")
}

// Complete implements cmd.Option.
func (u *UbuntuOptions) Complete() error {
	extensionConfig = &configv1alpha1.ExtensionConfig{}
	if u.configFile != "" {
		data, err := os.ReadFile(u.configFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
		if err = runtime.DecodeInto(configDecoder, data, extensionConfig); err != nil {
			return fmt.Errorf("error decoding config: %w", err)
		}
	} else {
		configv1alpha1.SetObjectDefaults_ExtensionConfig(extensionConfig)
	}

	// If disable-unattended-upgrades is true then set DisableUnattendedUpgrades in ExtensionConfig to true
	if u.DisableUnattendedUpgrades {
		extensionConfig.DisableUnattendedUpgrades = ptr.To(u.DisableUnattendedUpgrades)
	}

	return nil
}

// Completed implements cmd.Option.
func (u *UbuntuOptions) Completed() *UbuntuOptions {
	return u
}

func (u *UbuntuOptions) Validate() error {
	if errs := validation.ValidateExtensionConfig(extensionConfig); len(errs) > 0 {
		return fmt.Errorf("invalid extension config: %w", errs.ToAggregate())
	}
	return nil
}

func (u *UbuntuOptions) Apply(config *operatingsystemconfig.Config) {
	config.ExtensionConfig = extensionConfig
}
