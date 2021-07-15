package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/livegrep/livegrep/cmd/livegrep-index-builder/codehost"
)

var (
	flagIndexName           = flag.String("index-name", "livegrep-index", "name of the generated index")
	flagRepoBasePath        = flag.String("repo-base-path", "/tmp", "working directory on disk for storing repositories")
	flagOutIndexConfig      = flag.String("out-index-config", "index.json", "path on disk to save the livegrep index config")
	flagSSHPrivateKeyPath   = flag.String("ssh-private-key-path", "", "path to the private key for SSH authentication; defers to the SSH agent if empty")
	flagSSHSkipHostVerify   = flag.Bool("ssh-skip-host-verify", false, "skip server host identity verification for SSH authentication")
	flagRepoSyncConcurrency = flag.Int("repo-sync-concurrency", 5, "concurrency limit for repository synchronization")
	flagCodeHost            = newChoicesFlag([]string{codehost.Gitlab, codehost.Static}, "")
	flagCodeHostParams      = newStringMapFlag()

	codehostBackends = map[string]codehost.Factory{
		codehost.Gitlab: codehost.NewGitLabCodeHost,
		codehost.Static: codehost.NewStaticCodeHost,
	}
)

func init() {
	flag.Var(flagCodeHost, "codehost", "name of the codehost backend to use")
	flag.Var(flagCodeHostParams, "codehost-param", "key-value parameter used to configure behavior of the codehost backend")

	flag.Parse()
}

func main() {
	/* Initialization of the code host client */

	codeHostFactory, ok := codehostBackends[flagCodeHost.Choice()]
	if !ok {
		panic(fmt.Errorf("unknown codehost backend: codehost=%s", flagCodeHost.Choice()))
	}

	codeHost, err := codeHostFactory(flagCodeHostParams.Values())
	if err != nil {
		panic(err)
	}

	projects, err := codeHost.ListProjects()
	if err != nil {
		panic(err)
	}

	log.Printf("listed projects from codehost backend: codehost=%s projects=%d", flagCodeHost.Choice(), len(projects))

	/* Generate a livegrep index config and write it to disk */

	cfg := indexConfig(*flagIndexName, *flagRepoBasePath, projects)

	cfgJSON, err := json.Marshal(&cfg)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(*flagOutIndexConfig, cfgJSON, 0644); err != nil {
		panic(err)
	}

	log.Printf("wrote livegrep index configuration: path=%s", *flagOutIndexConfig)

	/* Synchronize the repositories requested for indexing */

	var syncTasks []taskClosure

	for _, project := range projects {
		syncTasks = append(syncTasks, func(project *codehost.Project) taskClosure {
			return func() error {
				log.Printf(
					"synchronizing repository: name=%s remote=%s revision=%s path=%s",
					project.Name,
					project.Remote,
					project.Revision,
					project.RepositoryPath(*flagRepoBasePath),
				)

				opts := &authenticationOptions{
					sshPrivateKeyPath: *flagSSHPrivateKeyPath,
					sshSkipHostVerify: *flagSSHSkipHostVerify,
				}

				err := syncRepository(project, *flagRepoBasePath, opts)
				if err != nil {
					log.Printf("encountered synchronization error: name=%s err=%v", project.Name, err)
				}

				return err
			}
		}(project))
	}

	if errs := executeTasks(syncTasks, *flagRepoSyncConcurrency); len(errs) > 0 {
		panic(fmt.Errorf("repository synchronization encountered nonzero errors: errs=%d", len(errs)))
	}

	log.Printf("repository synchronization complete")
}
