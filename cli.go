package main

import (
    "errors"
    "fmt"
    flag "github.com/spf13/pflag"
    "os"
)

const usage = `ctwrapper, ` + version + `.
A small git wrapper for Hashicorp's consul-template.
See ` + website + ` for more information.
Author: ` + author + `

Usage:
  ctwrapper [-r <URL>] [-b <branch>] [-c <commit>] [-g <depth>] [-d <dir>]
            [-u <user>] [-p <password> | -s <vault path>]
            [-e <extension>] 
            [-- <extra consul-template parameters>] 
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
  -o  | --ct-opt    : Quoted paramters to pass to consul-template.
  -h  | --help      : This help message.
  -v  | --version   : Version message.
  --                : Extra consul-template parameters, e.g. -exec.  
`

/* Flags */
var help, progVersion bool
var branch, commit, dir, ext, password, path, repo, user string
var depth int

/* Object to hold the parameters */
type Config struct {
    Address, Branch, Commit, Dir, Ext, Password, Repo, User string
    CtParams                                                []string
    Depth                                                   int
}

/* Initialize the flags */
func init() {
    flag.BoolVarP(&help, "help", "h", false, "")
    flag.BoolVarP(&progVersion, "version", "v", false, "")
    flag.StringVarP(&repo, "repo", "r", "", "")
    flag.StringVarP(&branch, "branch", "b", defaultBranch, "")
    flag.StringVarP(&commit, "commit", "c", defaultCommit, "")
    flag.IntVarP(&depth, "git-depth", "g", defaultDepth, "")
    flag.StringVarP(&dir, "dir", "d", "", "")
    flag.StringVarP(&ext, "ext", "e", defaultExt, "")
    flag.StringVarP(&user, "user", "u", "", "")
    flag.StringVarP(&password, "password", "p", "", "")
    flag.StringVarP(&path, "vault-path", "s", "", "")

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
    config.CtParams = flag.Args()

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
        return errors.New("Password can not be retrieved.")
    }
    return nil
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
