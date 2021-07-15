package codehost

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// config is the expected YAML config file format.
type config struct {
	Projects []struct {
		Name       string   `yaml:"name"`
		Revision   string   `yaml:"revision"`
		Remote     string   `yaml:"remote"`
		URLPattern string   `yaml:"url_pattern"`
		Labels     []string `yaml:"labels"`
	} `yaml:"projects"`
}

// StaticCodeHost describes a code host whose projects are defined statically in a config file.
type StaticCodeHost struct {
	cfg *config
}

// NewStaticCodeHost creates a CodeHost implementation based on a config file on disk.
func NewStaticCodeHost(params map[string]string) (CodeHost, error) {
	path := params["path"]
	if path == "" {
		return nil, fmt.Errorf("static: must specify path to the config file on disk as codehost param `path`")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("static: error reading config: err=%v", err)
	}

	var cfg *config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("static: error parsing config: err=%v", err)
	}

	return &StaticCodeHost{cfg}, nil
}

// ListProjects reshapes the information parsed from the YAML config file into Project structs.
func (s *StaticCodeHost) ListProjects() ([]*Project, error) {
	var projects []*Project

	for _, project := range s.cfg.Projects {
		projects = append(projects, &Project{
			Name:       project.Name,
			Revision:   project.Revision,
			Remote:     project.Remote,
			URLPattern: project.URLPattern,
			Labels:     project.Labels,
		})
	}

	return projects, nil
}
