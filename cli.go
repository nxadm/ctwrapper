package main

import (
	"errors"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strconv"
)

const usage = `ctwrapper, ` + version + `.
A small git wrapper for Hashicorp's consul-template.
See ` + website + ` for more information.
Author: ` + author + `

Usage:
  ctwrapper [-r <URL>] [-b <branch>] [-c <commit>] [-g <gitDepth>] [-d <dir>]
            [-e <extension>] [-u <git user>]  
			[-p <git password> | -s <vault path for git password>]
            [-- <extra consul-template parameters>] 
  ctwrapper [-h]
  ctwrapper [-v]


Parameters:
  -r  | --repo                : Git repo URL.
  -d  | --dir                 : Directory to download the repo [default: ` + defaultCloneDir + `].
  -b  | --branch              : Git branch [default: ` + defaultBranch + `]
  -c  | --commit              : Git commit [default: ` + defaultCommit + `].
  -g  | --git-depth           : Git depth  [default: 0 (unlimited)].
  -u  | --git-user            : Git HTTPS username (when not using SSH).
  -p  | --git-password        : Git HTTPS password (when not using SSH).
  -s  | --git-pass-vault-path : Retrieve the git HTTPS password at Vault path
                                (including the backend).
  -e  | --ext                 : Template extension [defaul: ` + defaultExt + `].
  -h  | --help                : This help message.
  -v  | --version             : Version message.
  --                          : Extra consul-template parameters, e.g. -exec.

Besides the default values when applicable, all the parameters can be 
passed as environment variables by using the full parameter name in capitals
without '-':
REPO, DIR, BRANCH, COMMIT, GITDEPTH, GITUSER, GITPASSWORD, VAULTPATH, EXT.

When both command line parameters and environment variables are defined,
the first type take precedence.

For the Vault parameters used in templates, these are retrieved from
environment values like VAULT_ADDR, VAULT_TOKEN and other VAULT_* variables).

Examples:                                                                       
  $ ctwrapper -r git@github.com:nxadm/ctwrapper.git       
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -d /var/tmp/project \           
    -s "secret/production/third-party/repo-password"                            
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -u foo -p bar \           
    -d /project	-- -vault-addr 'https://10.5.32.5:8200 -exec /sbin/my-server          

`

/* Flags */
var help, progVersion bool
var branch, commit, dir, ext, gitPassword, vaultPath, repo, gitUser string
var gitDepth int
var envValues = make(map[string]*string)

/* Object to hold the parameters */
type Config struct {
	Branch, Commit, Dir, Ext, GitPassword, GitUser, Repo string
	CtParams                                             []string
	GitDepth                                             int
}

/* Initialize the flags */
func init() {
	// Look up environment variables: string
	branch = defaultBranch
	commit = defaultCommit
	ext = defaultExt
	dir = defaultCloneDir
	envValues["REPO"] = &repo
	envValues["DIR"] = &dir
	envValues["BRANCH"] = &branch
	envValues["COMMIT"] = &commit
	envValues["GITUSER"] = &gitUser
	envValues["GITPASSWORD"] = &gitPassword
	envValues["VAULTPATH"] = &vaultPath
	envValues["EXT"] = &ext

	for envName, container := range envValues {
		v, ok := os.LookupEnv(envName)
		if ok {
			*container = v
		}
	}

	// Look up environment variables: int
	condDefGitDepth := defaultGitDepth
	v, ok := os.LookupEnv("GITDEPTH")
	if ok {
		i, err := strconv.Atoi(v)
		if err == nil {
			condDefGitDepth = i
		}
	}

	// Flags
	// Note: double check the map keys (nil pointers!)
	flag.BoolVarP(&help, "help", "h", false, "")
	flag.BoolVarP(&progVersion, "version", "v", false, "")

	flag.StringVarP(&repo, "repo", "r", *envValues["REPO"], "")
	flag.StringVarP(&branch, "branch", "b", *envValues["BRANCH"], "")
	flag.StringVarP(&commit, "commit", "c", *envValues["COMMIT"], "")
	flag.IntVarP(&gitDepth, "git-depth", "g", condDefGitDepth, "")
	flag.StringVarP(&dir, "dir", "d", *envValues["DIR"], "")
	flag.StringVarP(&ext, "ext", "e", *envValues["EXT"], "")
	flag.StringVarP(&gitUser, "git-user", "u", *envValues["GITUSER"], "")
	flag.StringVarP(&gitPassword, "git-password", "p", *envValues["GITPASSWORD"], "")
	flag.StringVarP(&vaultPath, "vault-path", "s", *envValues["VAULTPATH"], "")

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
		return errors.New("no parameters supplied"), false
	case help == true:
		flag.Usage()
		return nil, true
	case progVersion == true:
		fmt.Println(version)
		return nil, true
	}

	// importValues the values from CLI switches
	config.GitDepth = gitDepth
	config.Branch = branch
	config.Commit = commit
	config.Dir = dir
	config.Ext = ext
	config.Repo = repo
	config.GitUser = gitUser
	config.CtParams = flag.Args()

	// Retrieve GitPassword
	err := config.retrievePassword(gitUser, gitPassword, vaultPath)
	if err != nil {
		return err, false
	}

	// Verify the parameters
	fmt.Printf("%#v\n", config)
	return config.verifyParams(), false
}

func (config *Config) retrievePassword(user, password, vaultPath string) error {
	switch {
	// Anonymous or SSH
	case user == "":
	// CLI gitPassword
	case password != "":
		config.GitPassword = password
	// GitPassword from Vault
	case vaultPath != "":
		secret, err := retrieveVaultSecret(vaultPath)
		config.GitPassword = secret
		if err != nil {
			return err
		}
		config.GitPassword = secret
	default:
		return errors.New("git password can not be retrieved")
	}
	return nil
}

func (config *Config) verifyParams() error {
	nonEmpty := map[string]string{
		"repo": config.Repo,
	}
	for k, v := range nonEmpty {
		if v == "" {
			return errors.New(k + " is required")
		}
	}

	return nil
}
