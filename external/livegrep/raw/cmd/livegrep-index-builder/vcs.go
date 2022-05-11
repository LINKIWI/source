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
	sshPrivateKeyPath  string
	sshCertificatePath string
	sshSkipHostVerify  bool
}

// syncRepository clones a new repository or fetches an existing repository on disk.
func syncRepository(project *codehost.Project, basePath string, opts *authenticationOptions) error {
	var auth transport.AuthMethod

	repoPath := project.RepositoryPath(basePath)

	endpoint, err := transport.NewEndpoint(project.Remote)
	if err != nil {
		return err
	}

	username := gitssh.DefaultUsername
	if endpoint.User != "" {
		username = endpoint.User
	}

	// Use an explicit private key for authentication when supplied; otherwise, use the
	// session SSH agent
	if opts != nil && opts.sshPrivateKeyPath != "" {
		var signers []ssh.Signer

		sshKey, err := os.ReadFile(opts.sshPrivateKeyPath)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(sshKey)
		if err != nil {
			return err
		}

		if opts.sshCertificatePath != "" {
			sshCert, err := os.ReadFile(opts.sshCertificatePath)
			if err != nil {
				return err
			}

			cert, _, _, _, err := ssh.ParseAuthorizedKey(sshCert)
			if err != nil {
				return err
			}

			certSigner, err := ssh.NewCertSigner(cert.(*ssh.Certificate), signer)
			if err != nil {
				return err
			}

			signers = append(signers, certSigner)
		}

		signers = append(signers, signer)

		publicKeysCallback := func() ([]ssh.Signer, error) {
			// Prefer the certificate signer when available, but also allow fallback to
			// the plain signer. This allows compatibility with servers that reject
			// certificate-based SSH due to lack of capabilities.
			return signers, nil
		}

		hostVerificationCallback := gitssh.HostKeyCallbackHelper{}
		if opts.sshSkipHostVerify {
			hostVerificationCallback.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		auth = &gitssh.PublicKeysCallback{
			User:                  username,
			Callback:              publicKeysCallback,
			HostKeyCallbackHelper: hostVerificationCallback,
		}
	} else {
		agentAuth, err := gitssh.NewSSHAgentAuth(username)
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
