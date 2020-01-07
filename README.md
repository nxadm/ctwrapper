# ctwrapper
[![Build Status](https://travis-ci.com/nxadm/ctwrapper.svg?branch=master)](https://travis-ci.com/nxadm/ctwrapper)

ctwrapper is a small git wrapper to interface with Hashicorp's
[consul-template](https://github.com/hashicorp/consul-template). The use case
for this tool is providing remote configuration and secrets to containers that
require a complex configuration. Many orchestrators and tools, like
Hashicorp's [nomad](https://github.com/hashicorp/nomad), only provide
mechanisms to provision containers with simple 1-file configuration
requirements, e.g. by the
[template stanza](https://www.nomadproject.io/docs/job-specification/template.html).

As an alternative, ctwrapper retrieves a git repository with static files and
templates. Templates are passed as arguments to consul-template in order to
let consul-template run them and, by example, inject secrets from Vault.
Options can be passed to consul-template after "--", e.g. "-exec" to run the
actual application. In order to disable Vault support (e.g. when you inject
secret by environment variables), pass the "-vault-renew-token=false" option
to consul-template.

In order to use the
[consul](https://github.com/hashicorp/consul) and
[vault](https://github.com/hashicorp/vault) backends you need to define the
necessary environment variables (like VAULT_ADDR, VAULT_TOKEN and/or
CONSUL_TOKEN) or pass the '-consul-addr' and/or '-vault-addr' options to
consul-template (as passthrough after the '--'). Consult the
[consul-template documentation](https://github.com/hashicorp/consul-template)
for the parameters for consul-template.

Anonymous and authenticated git cloning is supported through SSH and HTTP(s).
Next to SSH (where the authentication is done by an SSH agent), HTTP(S) Basic
Authentication can use the username/password combination supplied on the command
line or retrieve the password from Vault. If no authentication is provided,
the repo will be retrieved anonymously.

ctwrapper, being a wrapper for consul-template, expects the latter to be in the
PATH or in the working directory. When run from a Docker container, use
[the exec variant of ENTRYPOINT and not the shell variant](https://docs.docker.com/engine/reference/builder/#entrypoint).
This allows ctwrapper to to preserve the signals received by the container and
pass it to consul-template and your application.

## Usage

```
Usage:
  ctwrapper [-r <URL>] [-b <branch>] [-c <commit>] [-g <gitDepth>] [-d <dir>]
            [-e <extension>] [-u <git user>]
			[-p <git password> | -s <vault path for git password>]
            [-- <extra consul-template parameters>]
  ctwrapper [-h]
  ctwrapper [-v]


Parameters:
  -r  | --repo                : Git repo URL.
  -d  | --dir                 : Directory to download the repo [default: /project].
  -b  | --branch              : Git branch [default: master]
  -c  | --commit              : Git commit [default: HEAD].
  -g  | --git-depth           : Git depth  [default: 0 (unlimited)].
  -u  | --git-user            : Git HTTPS username (when not using SSH).
  -p  | --git-password        : Git HTTPS password (when not using SSH).
  -s  | --git-pass-vault-path : Retrieve the git HTTPS password at Vault path
                                (including the backend).
  -e  | --ext                 : Template extension [defaul: .tmpl].
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
  $ ctwrapper -- "echo lala"
  $ ctwrapper -r git@github.com:nxadm/ctwrapper.git
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -d /var/tmp/project \
    -s "secret/production/third-party/repo-password"
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -u foo -p bar \
    -d /project -- -vault-addr 'https://10.5.32.5:8200 -exec /sbin/my-server
```

You may want to set the depth to a low number (e.g.) in order not to
unnecessarily retrieve the complete history of the repo. The `--commit` and
`--git-depth` options were included in order to prevent a race condition
between CI systems and git commits. If your setup ensures that the specified
commit is the last one, you can set `--git-depth` to 1.

Everything after `--` is directly passed as-is to consul-template. In the most
cases you'll pass an `-exec` command to consul-template to start your
application this way (quote the command).

## Releases

The creation of binaries found on the
[releases tab](https://github.com/nxadm/ctwrapper/releases) is completely
automated by Travis CI and created from a version tag in the repo. The
sha512 checkums files can be verified with the output of the
[Travis build](https://travis-ci.com/nxadm/ctwrapper/branches).
