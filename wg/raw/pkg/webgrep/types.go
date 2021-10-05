package webgrep

const (
	// EndpointSearch is the path to the code search query endpoint.
	EndpointSearch = "/api/search"
	// EndpointSource is the path to the source code query endpoint.
	EndpointSource = "/api/source"
	// EndpointMetadata is the path to the metadata endpoint.
	EndpointMetadata = "/api/meta/info"
)

// CodeSearchResult formalizes fields for a single code search result.
type CodeSearchResult struct {
	Repository string `json:"repo"`
	Version    string `json:"version"`
	Path       string `json:"path"`
	Lines      []struct {
		Line   string `json:"line"`
		Number int    `json:"number"`
		Bounds []int  `json:"bounds"`
	} `json:"lines"`
}

// FileSearchResult formalizes fields for a single file path result.
type FileSearchResult struct {
	Repository string `json:"repo"`
	Version    string `json:"version"`
	Path       string `json:"path"`
	Bounds     []int  `json:"bounds"`
}

// SearchStats formalizes fields in server-side search statistics.
type SearchStats struct {
	ExitReason int `json:"exitReason"`
	Time       int `json:"time"`
}

// Repository describes an indexed source repository.
type Repository struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	URL     string   `json:"url"`
	Remote  string   `json:"remote"`
	Labels  []string `json:"labels"`
}

// SearchQueryRequest describes the top-level request for a search query.
type SearchQueryRequest struct {
	Query         string   `json:"query"`
	File          string   `json:"file"`
	Repositories  []string `json:"repos"`
	Regex         bool     `json:"regex"`
	CaseSensitive bool     `json:"caseSensitive"`
	MaxMatches    int      `json:"maxMatches"`
	Context       int      `json:"context"`
}

// SearchQueryResponse describes the top-level response for a search query.
type SearchQueryResponse struct {
	Stats SearchStats        `json:"stats"`
	Code  []CodeSearchResult `json:"code"`
	Files []FileSearchResult `json:"files"`
}

// SourceQueryRequest describes the top-level request for a source code query.
type SourceQueryRequest struct {
	Repository string `json:"repo"`
	Version    string `json:"version"`
	Path       string `json:"path"`
}

// SourceQueryResponse is a wrapper over the contents of a source code query response payload.
type SourceQueryResponse struct {
	Repository string   `json:"repo"`
	Version    string   `json:"version"`
	Path       string   `json:"path"`
	Lines      []string `json:"source"`
}

// MetadataResponse describes the top-level response for a metadata request.
type MetadataResponse struct {
	Name         string       `json:"name"`
	Timestamp    int          `json:"timestamp"`
	Version      string       `json:"version"`
	Repositories []Repository `json:"repositories"`
}
