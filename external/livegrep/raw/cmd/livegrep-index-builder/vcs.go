package main

import (
	"os"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/livegrep/livegrep/cmd/livegrep-index-builder/codehost"
)

// authenticationOptions formalizes optional repository authentication options
// to use while synchronizing the remote repository.
type authenticationOptions struct {
	sshPrivateKeyPath string
	sshSkipHostVerify bool
}

// syncRepository clones a new repository or fetches an existing repository on disk.
func syncRepository(project *codehost.Project, basePath string, opts *authenticationOptions) error {
	var auth transport.AuthMethod

	repoPath := project.RepositoryPath(basePath)

	// Use an explicit private key for authentication when supplied; otherwise, use the
	// session SSH agent
	if opts != nil && opts.sshPrivateKeyPath != "" {
		sshKey, err := os.ReadFile(opts.sshPrivateKeyPath)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(sshKey)
		if err != nil {
			return err
		}

		hostVerificationCallback := gitssh.HostKeyCallbackHelper{}
		if opts.sshSkipHostVerify {
			hostVerificationCallback.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		auth = &gitssh.PublicKeys{
			User:                  gitssh.DefaultUsername,
			Signer:                signer,
			HostKeyCallbackHelper: hostVerificationCallback,
		}
	} else {
		agentAuth, err := gitssh.NewSSHAgentAuth(gitssh.DefaultUsername)
		if err != nil {
			return err
		}

		if opts != nil && opts.sshSkipHostVerify {
			agentAuth.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		auth = agentAuth
	}

	repo, err := git.PlainOpen(repoPath)

	// Clone the repository if it doesn't already exist
	if err == git.ErrRepositoryNotExists {
		cloneOpts := &git.CloneOptions{
			RemoteName: "origin",
			URL:        project.Remote,
			Auth:       auth,
			NoCheckout: true,
		}
		_, err = git.PlainClone(repoPath, false, cloneOpts)
		return err
	} else if err != nil {
		return err
	}

	// Otherwise, simply fetch against the existing remote
	fetchOpts := &git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
	}
	if err := repo.Fetch(fetchOpts); err != nil {
		if err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	return nil
}
