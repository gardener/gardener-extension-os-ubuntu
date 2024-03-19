// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"embed"

	ostemplate "github.com/gardener/gardener/extensions/pkg/controller/operatingsystemconfig/oscommon/template"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"k8s.io/apimachinery/pkg/util/runtime"
)

var cmd = "/usr/bin/env bash %s"
var cloudInitGenerator *ostemplate.CloudInitGenerator
var cloudInitGeneratorWithDisabledUnattendedUpgrades *ostemplate.CloudInitGenerator

//go:embed templates/*
var templates embed.FS

func init() {
	cloudInitTemplateString, err := templates.ReadFile("templates/cloud-init-ubuntu.template")
	runtime.Must(err)

	cloudInitTemplate, err := ostemplate.NewTemplate("cloud-init").Parse(string(cloudInitTemplateString))
	runtime.Must(err)
	cloudInitGenerator = ostemplate.NewCloudInitGenerator(cloudInitTemplate, ostemplate.DefaultUnitsPath, cmd, nil)

	cloudInitGeneratorWithDisabledUnattendedUpgrades = ostemplate.NewCloudInitGenerator(cloudInitTemplate, ostemplate.DefaultUnitsPath, cmd,
		func(*extensionsv1alpha1.OperatingSystemConfig) (map[string]interface{}, error) {
			return map[string]interface{}{
				"DisableUnattendedUpgrades": true,
			}, nil
		})
}

// CloudInitGenerator is the generator which will generate the cloud init yaml.
func CloudInitGenerator(disableUnattendedUpgrades bool) *ostemplate.CloudInitGenerator {
	if disableUnattendedUpgrades {
		return cloudInitGeneratorWithDisabledUnattendedUpgrades
	}
	return cloudInitGenerator
}
