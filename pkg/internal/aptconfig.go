package internal

type APTConfigSnake struct {
	PreserveSourcesList bool              `json:"preserve_sources_list"`
	Primary             []APTArchiveSnake `json:"primary,omitempty"`
	Security            []APTArchiveSnake `json:"security,omitempty"`
}

type APTArchiveSnake struct {
	Arches    []string `json:"arches,omitempty"`
	URI       string   `json:"uri,omitempty"`
	Search    []string `json:"search,omitempty"`
	SearchDNS bool     `json:"search_dns,omitempty"`
}

type APTCloudInit struct {
	APT APTConfigSnake `json:"apt,omitempty"`
}

type FilePart struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
