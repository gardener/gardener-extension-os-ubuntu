// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package internal

import "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"

// APTConfig with snake case is needed for cloud-init
type APTConfig struct {
	PreserveSourcesList bool                 `json:"preserve_sources_list"`
	Primary             []APTArchive         `json:"primary,omitempty"`
	Security            []APTArchive         `json:"security,omitempty"`
	Sources             map[string]APTSource `json:"sources,omitempty"`
}

// APTArchive with snake case is needed for cloud-init
type APTArchive struct {
	Arches    []v1alpha1.Architecture `json:"arches,omitempty"`
	URI       string                  `json:"uri,omitempty"`
	Search    []string                `json:"search,omitempty"`
	SearchDNS bool                    `json:"search_dns,omitempty"`
}

// APTSource describes a single cloud-init apt source entry.
type APTSource struct {
	Source string `json:"source"`
	Key    string `json:"key,omitempty"`
}

type APTCloudInit struct {
	APT        APTConfig   `json:"apt,omitempty"`
	WriteFiles []WriteFile `json:"write_files,omitempty"`
}

type WriteFile struct {
	Path        string           `json:"path"`
	Source      *WriteFileSource `json:"source,omitempty"`
	Permissions string           `json:"permissions"`
	Owner       string           `json:"owner"`
}

type WriteFileSource struct {
	URI string `json:"uri"`
}

type FilePart struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
