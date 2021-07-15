package codehost

import (
	"crypto/sha256"
	"encoding/base64"
	"path"
)

// Project describes a single project returned from the code host backend. It is assumed that there
// exists a one-to-one mapping between a project and indexed repository.
type Project struct {
	Name       string
	Revision   string
	Remote     string
	URLPattern string
	Labels     []string
}

// Checksum returns a string that can be used to uniquely identify the repository storage associated
// with this project. It hashes the repository name and remote to generate a filesystem-safe path.
func (p *Project) Checksum() string {
	hasher := sha256.New()
	hasher.Write([]byte(p.Name))
	hasher.Write([]byte(p.Remote))
	hash := hasher.Sum(nil)

	return base64.URLEncoding.EncodeToString(hash)
}

// RepositoryPath generates a path on disk from a base root path that would be suitable for storage
// of this project's repository.
func (p *Project) RepositoryPath(basePath string) string {
	return path.Join(basePath, p.Checksum())
}
