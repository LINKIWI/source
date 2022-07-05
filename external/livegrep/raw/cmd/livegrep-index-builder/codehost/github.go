package codehost

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubCodeHost is an implementation of CodeHost for Github.
type GithubCodeHost struct {
	parentType string
	parentName string
	timeout    time.Duration

	client *github.Client
}

// NewGithubCodeHost creates a new GithubCodeHost given CLI-supplied key-value parameters.
// An access token, alongside the repository parent type, must be supplied.
func NewGithubCodeHost(params map[string]string) (CodeHost, error) {
	var err error

	parentType := params["parent-type"]
	if parentType != "user" && parentType != "organization" {
		return nil, fmt.Errorf("github: repository parent must either be a user or organization")
	}

	parentName := params["parent-name"]
	if parentType == "organization" && parentName == "" {
		return nil, fmt.Errorf("github: parent name must be non-empty for an organization parent")
	}

	accessToken := params["access-token"]
	if accessToken == "" {
		accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("github: no access token specified")
	}

	timeout := 0 * time.Second
	if t := params["timeout"]; t != "" {
		timeout, err = time.ParseDuration(t)
		if err != nil {
			return nil, fmt.Errorf("github: unable to parse timeout: timeout=%s err=%v", t, err)
		}
	}

	oauthCtx := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}),
	)

	return &GithubCodeHost{
		client:     github.NewClient(oauthCtx),
		parentType: parentType,
		parentName: parentName,
		timeout:    timeout,
	}, nil
}

// ListProjects uses the Github API to query a list of repositories visible to the authentication
// context.
func (gh *GithubCodeHost) ListProjects() ([]*Project, error) {
	var cancel context.CancelFunc
	var fetchPaginatedRepositories func(page int) ([]*github.Repository, error)

	ctx := context.Background()
	if gh.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, gh.timeout)
		defer cancel()
	}

	fetchPaginatedRepositories = func(page int) ([]*github.Repository, error) {
		opts := github.ListOptions{
			Page:    page,
			PerPage: 100,
		}

		switch gh.parentType {
		case "user":
			repos, resp, err := gh.client.Repositories.List(
				ctx,
				gh.parentName,
				&github.RepositoryListOptions{ListOptions: opts},
			)
			if err != nil {
				return nil, fmt.Errorf("github: error querying user repositories: err=%v", err)
			}

			if resp.NextPage > 0 {
				additionalRepos, err := fetchPaginatedRepositories(resp.NextPage)
				if err != nil {
					return nil, err
				}

				repos = append(repos, additionalRepos...)
			}

			return repos, nil

		case "organization":
			repos, resp, err := gh.client.Repositories.ListByOrg(
				ctx,
				gh.parentName,
				&github.RepositoryListByOrgOptions{ListOptions: opts},
			)
			if err != nil {
				return nil, fmt.Errorf("github: error querying organization repositories: err=%v", err)
			}

			if resp.NextPage > 0 {
				additionalRepos, err := fetchPaginatedRepositories(resp.NextPage)
				if err != nil {
					return nil, err
				}

				repos = append(repos, additionalRepos...)
			}

			return repos, nil

		default:
			return nil, fmt.Errorf("github: unknown parent type: type=%s", gh.parentType)
		}
	}

	repos, err := fetchPaginatedRepositories(0)
	if err != nil {
		return nil, err
	}

	projects := make([]*Project, len(repos))
	for i, repo := range repos {
		projects[i] = &Project{
			Name:       *repo.FullName,
			Revision:   *repo.DefaultBranch,
			Remote:     *repo.SSHURL,
			URLPattern: fmt.Sprintf("%s/blob/{version}/{path}#L{lno}", *repo.HTMLURL),
		}
	}

	return projects, nil
}
