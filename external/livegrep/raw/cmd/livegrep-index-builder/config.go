package main

import (
	"fmt"

	pb "github.com/livegrep/livegrep/src/proto/config"

	"github.com/livegrep/livegrep/cmd/livegrep-index-builder/codehost"
)

// indexConfig generates a Livegrep index config struct given some parameters and a list of GitLab
// projects.
func indexConfig(name string, repoBasePath string, projects []*codehost.Project) *pb.IndexSpec {
	var repositories []*pb.RepoSpec

	for _, project := range projects {
		repo := &pb.RepoSpec{
			Path:      project.RepositoryPath(repoBasePath),
			Name:      project.Name,
			Revisions: []string{fmt.Sprintf("origin/%s", project.Revision)},
			Metadata: &pb.Metadata{
				UrlPattern: project.URLPattern,
				Remote:     project.Remote,
				Labels:     project.Labels,
			},
		}

		repositories = append(repositories, repo)
	}

	return &pb.IndexSpec{
		Name:         name,
		Repositories: repositories,
	}
}
