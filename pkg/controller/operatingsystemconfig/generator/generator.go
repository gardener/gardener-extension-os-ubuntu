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

//go:embed templates/*
var templates embed.FS

func init() {
	cloudInitTemplateString, err := templates.ReadFile("templates/cloud-init-ubuntu.template")
	runtime.Must(err)

	cloudInitTemplate, err := ostemplate.NewTemplate("cloud-init").Parse(string(cloudInitTemplateString))
	runtime.Must(err)
	cloudInitGenerator = ostemplate.NewCloudInitGenerator(cloudInitTemplate, ostemplate.DefaultUnitsPath, cmd, additionalValues)
}

// CloudInitGenerator is the generator which will generate the cloud init yaml.
func CloudInitGenerator() *ostemplate.CloudInitGenerator {
	return cloudInitGenerator
}

func additionalValues(*extensionsv1alpha1.OperatingSystemConfig) (map[string]interface{}, error) {
	return nil, nil
}
