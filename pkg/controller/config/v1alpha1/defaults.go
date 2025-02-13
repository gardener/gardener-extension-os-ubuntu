package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ExtensionConfig(obj *ExtensionConfig) {
	if obj.NTP == nil {
		obj.NTP = &NTPConfig{}
	}
	obj.DisableUnattendedUpgrades = ptr.To(false)
}

func SetDefaults_NTPConfig(obj *NTPConfig) {
	if obj.Daemon == "" {
		obj.Daemon = SystemdTimesyncd
	}
}

func SetDefaults_NTPDConfig(obj *NTPDConfig) {
	if obj.Servers == nil {
		obj.Servers = make([]string, 0)
	}
}
