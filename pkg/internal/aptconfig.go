package internal

import "github.com/gardener/gardener-extension-os-ubuntu/pkg/controller/config/v1alpha1"

// APTConfig with snake case is needed for cloud-init
type APTConfig struct {
	PreserveSourcesList bool         `json:"preserve_sources_list"`
	Primary             []APTArchive `json:"primary,omitempty"`
	Security            []APTArchive `json:"security,omitempty"`
}

// APTArchive with snake case is needed for cloud-init
type APTArchive struct {
	Arches    []v1alpha1.Architecture `json:"arches,omitempty"`
	URI       string                  `json:"uri,omitempty"`
	Search    []string                `json:"search,omitempty"`
	SearchDNS bool                    `json:"search_dns,omitempty"`
}

type APTCloudInit struct {
	APT APTConfig `json:"apt,omitempty"`
}

type FilePart struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
