package codehost

const (
	// Gitlab is the identifier for the Gitlab codehost backend.
	Gitlab = "gitlab"
	// Static is the identifier for the static codehost backend.
	Static = "static"
)

// Factory types a function that creates a CodeHost instance given a map of key-value parameters.
type Factory func(params map[string]string) (CodeHost, error)

// CodeHost is any backend that can provide a list of Projects.
type CodeHost interface {
	// ListProjects retrieves a list of eligible projects for indexing.
	ListProjects() ([]*Project, error)
}
