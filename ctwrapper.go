package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

const author = "Claudio Ramirez <pub.claudio@gmail.com>"
const version = "v0.3.4"
const website = "https://github.com/nxadm/ctwrapper"
const defaultExt = ".tmpl"
const defaultBranch = "master"
const defaultCommit = "HEAD"
const defaultDepth = 0 // disable

func main() {

	/* Read the CLI */
	config := Config{}
	exitIfErr(config.importValues())

	/* Retrieve git project */
	// Create authentication method only if provided at the CLI
	var authMethod transport.AuthMethod
	if config.User != "" {
		authMethod = transport.AuthMethod(&http.BasicAuth{
			Username: config.User, Password: config.Password})
	}

	// clone the repo and go to the specified commit
	repo := Repo{
		URL:    config.Repo,
		Branch: config.Branch,
		Commit: config.Commit,
		Depth:  config.Depth,
		Dir:    config.Dir,
		Auth:   authMethod,
	}
	exitIfErr(repo.clone())

	/* Look for templates */
	files, err := findTemplates(config.Dir, config.Ext)
	exitIfErr(err)

	/* Interface with consul-template */
	err = runCt(config.Ext, files, config.CtParams)
	exitIfErr(err)
}

func exitIfErr(err error) {
	switch err {
	case nil:
	default:
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
