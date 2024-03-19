package cmd

import "github.com/spf13/pflag"

var (
	// DisableUnattendedUpgrades is the name of the command line flag to disable unattended upgrades in ubuntu.
	DisableUnattendedUpgrades = "disable-unattended-upgrades"
)

// UbuntuOptions are command line options that can be set for ubuntu configuration.
type UbuntuOptions struct {
	// DisableUnattendedUpgrades is the flat to disable unattended upgrades in ubuntu.
	DisableUnattendedUpgrades bool
}

// AddFlags implements cmd.Option.
func (u *UbuntuOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&u.DisableUnattendedUpgrades, DisableUnattendedUpgrades, u.DisableUnattendedUpgrades, "The flag to disable unattended upgrades in ubuntu.")
}

// Complete implements cmd.Option.
func (u *UbuntuOptions) Complete() error {
	return nil
}

// Completed implements cmd.Option.
func (u *UbuntuOptions) Completed() *UbuntuOptions {
	return u
}
