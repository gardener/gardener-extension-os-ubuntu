// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// +k8s:deepcopy-gen=package
// +k8s:defaulter-gen=TypeMeta

//go:generate gen-crd-api-reference-docs -api-dir . -config ../../../../hack/api-reference/config.json -template-dir $GARDENER_HACK_DIR/api-reference/template -out-file ../../../../hack/api-reference/config.md

// Package v1alpha1 contains the API for configuring the os-ubuntu extension.
// +groupName=config.ubuntu.os.extensions.gardener.cloud
package v1alpha1 // import "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"
