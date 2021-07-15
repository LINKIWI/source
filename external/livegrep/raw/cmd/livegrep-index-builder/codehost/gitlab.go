package codehost

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

// GitLabCodeHost is an implementation of CodeHost for a GitLab instance.
type GitLabCodeHost struct {
	namespaces []string

	client *gitlab.Client
}

// NewGitLabCodeHost creates a new GitLabCodeHost given CLI-supplied key-value parameters.
// It allows specification of a custom GitLab base URL, and a whitelist of project namespaces.
func NewGitLabCodeHost(params map[string]string) (CodeHost, error) {
	// GitLab instance base URL
	baseURL := params["base-url"]
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	apiBaseURL := strings.TrimSuffix(baseURL, "/") + "/api/v4"

	// API access token, supplied either as a parameter or stored in the environment
	accessToken := params["access-token"]
	if accessToken == "" {
		accessToken = os.Getenv("GITLAB_ACCESS_TOKEN")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("gitlab: no access token specified")
	}

	// Optionally specify repository namespaces to list for indexing, comma-delimited
	var namespaces []string
	if joinedNamespaces := params["namespaces"]; joinedNamespaces != "" {
		namespaces = strings.Split(joinedNamespaces, ",")
	}

	client := gitlab.NewClient(nil, accessToken)

	if err := client.SetBaseURL(apiBaseURL); err != nil {
		return nil, fmt.Errorf("gitlab: error setting base URL: err=%v", err)
	}

	return &GitLabCodeHost{namespaces, client}, nil
}

// ListProjects paginates through the GitLab API to retrieve all projects. This method also filters
// the retrieved projects so that only whitelisted namespaces are returned.
func (gl *GitLabCodeHost) ListProjects() ([]*Project, error) {
	var fetchPaginatedProjects func(page int) ([]*Project, error)

	// Reshape the GitLab project object into the generic, interface-compliant Project struct.
	mapProject := func(project *gitlab.Project) *Project {
		baseURL := gl.client.BaseURL()
		urlPattern := fmt.Sprintf(
			"%s://%s/{name}/blob/{version}/{path}#L{lno}",
			baseURL.Scheme,
			baseURL.Host,
		)

		return &Project{
			Name:       project.PathWithNamespace,
			Revision:   project.DefaultBranch,
			Remote:     project.SSHURLToRepo,
			URLPattern: urlPattern,
		}
	}

	// Apply the project namespace whitelist filter predicate to a list of Projects.
	applyNamespaceFilter := func(projects []*Project) []*Project {
		var filtered []*Project

		predicate := func(project *Project) bool {
			for _, namespace := range gl.namespaces {
				if strings.HasPrefix(project.Name, fmt.Sprintf("%s/", namespace)) {
					return true
				}
			}

			// Only consider the filter if a nonzero number of whitelist entries are present
			return len(gl.namespaces) == 0
		}

		for _, project := range projects {
			if predicate(project) {
				filtered = append(filtered, project)
			}
		}

		return filtered
	}

	// Use the GitLab API to fetch all projects, page by page.
	fetchPaginatedProjects = func(page int) ([]*Project, error) {
		var projects []*Project
		var empty []*Project

		opts := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}

		gitlabProjects, resp, err := gl.client.Projects.ListProjects(opts)
		if err != nil {
			return empty, err
		}

		for _, gitlabProject := range gitlabProjects {
			projects = append(projects, mapProject(gitlabProject))
		}

		if resp.NextPage == 0 {
			return applyNamespaceFilter(projects), nil
		}

		additionalProjects, err := fetchPaginatedProjects(resp.NextPage)
		if err != nil {
			return empty, nil
		}

		return applyNamespaceFilter(append(projects, additionalProjects...)), nil
	}

	return fetchPaginatedProjects(1)
}
