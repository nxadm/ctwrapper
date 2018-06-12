package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"os"
)

type Repo struct {
	URL, Branch, Commit, Dir string
	Depth                    int
	Auth                     transport.AuthMethod
}

func (repo *Repo) clone() error {

	// clone it
	referenceName :=
		plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", repo.Branch))
	repoObj, err := git.PlainClone(repo.Dir, false,
		&git.CloneOptions{
			URL:           repo.URL,
			ReferenceName: referenceName,
			SingleBranch:  true,
			Depth:         repo.Depth,
			Progress:      os.Stdout,
			Auth:          repo.Auth,
		})
	if err != nil {
		return err
	}

	// Get the working tree so we can change refs
	tree, err := repoObj.Worktree()
	if err != nil {
		return err
	}

	// Checkout the commit
	err = tree.Checkout(
		&git.CheckoutOptions{Hash: plumbing.NewHash(repo.Commit)})
	if err != nil {
		return err
	}

	return nil
}
