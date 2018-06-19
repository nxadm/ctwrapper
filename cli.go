package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const usage = `ctwrapper, ` + version + `.
A small git wrapper for Hashicorp's consul-template.
See ` + website + ` for more information.
Author: ` + author + `

Usage:
  ctwrapper [-r <URL>] [-b <branch>] [-c <commit>] [-g <depth>] [-d <dir>]
            [-u <user>] [-p <password> | -k <vault path>]
            [-e <extension>] 
            [-o <quoted options for consul-template>]  
  ctwrapper [-h]
  ctwrapper [-v]


Parameters:
  -r  | --repo      : Git repo URL.
  -d  | --dir       : Directory to download the repo.
  -b  | --branch    : Git branch [default: ` + defaultBranch + `]
  -c  | --commit    : Git commit [default: ` + defaultCommit + `].
  -g  | --git-depth : Git depth  [default: unlimited].
  -u  | --user      : Git username.
  -p  | --password  : Git password.
  -s  | --vault-path: Vault path to the secret (including the backend).
  -e  | --ext       : Template extension [defaul: ` + defaultExt + `].
  -o  | --ct-opt    : Extra (quoted) options to pass to consul-template.
  -h  | --help      : This help message.
  -v  | --version   : Version message.
`

// Define the flags
var help, progVersion bool
var branch, commit, ctOpt, dir, ext, password, path, repo, user string
var depth int

type Config struct {
	Address, Branch, Commit, Dir, Ext, Password, Repo, User string
	CTOptions                                               []string
	Depth                                                   int
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
	flag.IntVar(&depth, "g", defaultDepth, "")
	flag.IntVar(&depth, "git-depth", defaultDepth, "")
	flag.StringVar(&dir, "d", "", "")
	flag.StringVar(&dir, "dir", "", "")
	flag.StringVar(&ext, "e", defaultExt, "")
	flag.StringVar(&ext, "ext", defaultExt, "")
	flag.StringVar(&user, "u", "", "")
	flag.StringVar(&user, "user", "", "")
	flag.StringVar(&password, "p", "", "")
	flag.StringVar(&password, "password", "", "")
	flag.StringVar(&path, "s", "", "")
	flag.StringVar(&path, "vault-path", "", "")
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
	case flag.NFlag() == 0:
		flag.Usage()
		return errors.New("No parameters supplied."), false
	case help == true:
		flag.Usage()
		return nil, true
	case progVersion == true:
		fmt.Println(version)
		return nil, true
	}

	// importValues the values from CLI switches
	config.Branch = branch
	config.Commit = commit
	config.Dir = dir
	config.Ext = ext
	config.Repo = repo
	config.User = user

	// Convert ctOpt
	if ctOpt != "" {
		config.CTOptions = splitArg(ctOpt)
	}

	// Retrieve Password
	err := config.retrievePassword(user, password, path)
	if err != nil {
		return err, false
	}

	// Verify the parameters
	return config.verifyParams(), false
}

func (config *Config) retrievePassword(user, password, path string) error {
	switch {
	// Anonymous
	case user == "":
	// CLI passoword
	case password != "":
		config.Password = password
	// Password from Vault
	case path != "":
		secret, err := retrieveVaultSecret(path)
		config.Password = secret
		if err != nil {
			return err
		}
		config.Password = secret
	default:
		return errors.New("No password can be retrieved.")
	}
	return nil
}

func splitArg(arg string) []string {
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	return strings.FieldsFunc(arg, f)
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
