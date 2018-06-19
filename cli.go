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
  vault-wrapper [-r <URL>] [-d <dir>] [-b <branch>] [-c <commit>] [-g <depth>]
                [-u <user> -p <password>]
				[-u <user> -a <vault address> -s <vault path/key>]
                [-e <extension>] 
                [-o <quoted options for consul-template>]  
  vault-wrapper [-h]
  vault-wrapper [-v]

Parameters:
  -r  | --repo      : Git repo URL.
  -d  | --dir       : Directory to download the repo.
  -b  | --branch    : Git branch [default: ` + defaultBranch + `]
  -c  | --commit    : Git commit [default: ` + defaultCommit + `].
  -g  | --git-depth : Git depth  [default: unlimited].
  -u  | --user      : Git username.
  -p  | --password  : Git password.
  -a  | --vault-addr: Vault address (will be passed to consul-template).
  -s  | --secret    : Vault path (backend and key to retrieve).
  -e  | --ext       : Template extension [defaul: ` + defaultExt + `].
  -o  | --ct-opt    : Extra (quoted) options to pass to consul-template.
  -h  | --help      : This help message.
  -v  | --version   : Version message.
`

// Define the flags
var help, progVersion bool
var address, branch, commit, ctOpt, dir, ext, password, repo, secret, user string
var depth int

type Config struct {
	Address, Branch, Commit, Dir, Ext, Repo, Password, User string
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
	flag.StringVar(&secret, "a", "", "")
	flag.StringVar(&secret, "address", "", "")
	flag.StringVar(&secret, "s", "", "")
	flag.StringVar(&secret, "secret", "", "")
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
	config.Address = address
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
	err := config.retrievePassword(user, password, secret)
	if err != nil {
		return err, false
	}

	// Verify the parameters
	return config.verifyParams(), false
}

func (config *Config) retrievePassword(user, password, secret string) error {
	var err error
	switch {
	// Anonymous
	case user == "":
	// CLI passoword
	case password != "":
		config.Password = password
	// Password from Vault
	case secret != "":
		if  address != "" {
			secretFromVault, err := retrieveVaultSecret(address, secret)
			if err != nil {
				return err
			}
			config.Password = secretFromVault
		} else {
			return errors.New("When providing a secret, vault-addr must be provided.")
		}
	}

	return err
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
