package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const usage = `ctwrapper, ` + version + `.
A small git wrapper for Hashicorp's consul-template.
See ` + website + ` for more information.
Author: ` + author + `

Usage:
  vault-wrapper [-r <URL>] [-b <branch>] [-c <commit>] [-gd  <nr of commits>]
                [-u <user>] [-p <password> | -vp <path> -vk <key>]
                [-d <dir>] [-e <extension>] 
                [-o <options>]  
  vault-wrapper [-h]
  vault-wrapper [-v]

Parameters:
  -r  | --repo      : Git repo URL.
  -b  | --branch    : Git branch [default: ` + defaultBranch + `]
  -c  | --commit    : Git commit [default: ` + defaultCommit + `].
  -gd | --git-depth : Git depth  [default: unlimited].
  -u  | --user      : Git username.
  -p  | --password  : Git password.
  -vp | --vault-path: Vault path (include backend).
  -vk | --vault-key : Vault key.
  -d  | --dir       : Directory with templates [default: . ].
  -e  | --ext       : Template extension [defaul: ` + defaultExt + `].
  -o  | --ct-opt    : Extra options to pass to consul-template.
  -h  | --help      : This help message.
  -v  | --version   : Version message.
`

// Define the flags
var help, progVersion bool
var repo, branch, commit, dir, ext, user, password, vaultPath, vaultKey, ctOpt string
var depth int

type Config struct {
	Repo, Branch, Commit, Dir, Ext, User, Password string
	CTOptions                                      []string
	Depth                                          int
}

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")
	flag.BoolVar(&progVersion, "v", false, "")
	flag.BoolVar(&progVersion, "version", false, "")
	flag.StringVar(&repo, "r", "", "")
	flag.StringVar(&repo, "repo", "", "")
	flag.StringVar(&branch, "b", defaultBranch, "")
	flag.StringVar(&branch, "branch", defaultBranch, "")
	flag.StringVar(&commit, "c", defaultCommit, "")
	flag.StringVar(&commit, "commit", defaultCommit, "")
	flag.IntVar(&depth, "gd", defaultDepth, "")
	flag.IntVar(&depth, "git-depth", defaultDepth, "")
	flag.StringVar(&dir, "d", "", "")
	flag.StringVar(&dir, "dir", "", "")
	flag.StringVar(&ext, "e", defaultExt, "")
	flag.StringVar(&ext, "ext", defaultExt, "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&user, "user", "", "")
	flag.StringVar(&password, "p", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.StringVar(&vaultPath, "vp", "", "")
	flag.StringVar(&vaultPath, "vault-path", "", "")
	flag.StringVar(&vaultKey, "vk", "", "")
	flag.StringVar(&vaultKey, "vault-key", "", "")
	flag.StringVar(&ctOpt, "o", "", "")
	flag.StringVar(&ctOpt, "ct-opt", "", "")

	// Set a custom usage message
	flag.Usage = func() { fmt.Println(usage) }

	// Parse it
	flag.Parse()
}

func (config *Config) importValues() error {
	// Read the CLI
	err, earlyExit := config.readCliParams()
	switch {
	case earlyExit == true:
		os.Exit(0)
	case err != nil:
		return err
	}

	return nil
}

func (config *Config) readCliParams() (error, bool) {
	// Handle early exits
	switch {
	case help == true:
		flag.Usage()
		return nil, true
	case progVersion == true:
		fmt.Println(version)
		return nil, true
	}

	// importValues the values from CLI switches
	config.Repo = repo
	config.Branch = branch
	config.Commit = commit
	config.Dir = dir
	config.Ext = ext
	config.User = user

	// Convert ctOpt
	if ctOpt != "" {
		config.CTOptions = strings.Split(ctOpt, " ")
	}

	// Retrieve Password
	err := config.retrievePassword(user, password, vaultPath, vaultKey)
	if err != nil {
		return err, false
	}

	// Verify the parameters
	return config.verifyParams(), false
}

func (config *Config) retrievePassword(user, password, vaultPath, vaultKey string) error {
	var err error
	switch {
	// Anonymous
	case user == "":
	// CLI passoword
	case password != "":
		config.Password = password
	// Password from Vault
	case vaultPath != "" && vaultKey != "":
		secret, err := retrieveVaultSecret(vaultPath, vaultKey)
		if err != nil {
			return err
		}
		config.Password = secret
	// Error cases
	case vaultPath != "":
		return errors.New("vault-path is required when vault-key is used.")
	case vaultKey != "":
		return errors.New("vault-key is required when vault-path is used.")
	}

	return err
}

func (config *Config) verifyParams() error {
	nonEmpty := map[string]string{
		"repo": config.Repo,
		"dir":  config.Dir,
	}
	for k, v := range nonEmpty {
		if v == "" {
			return errors.New(k + " is required.")
		}
	}

	return nil
}
