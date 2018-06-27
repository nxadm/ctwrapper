# ctwrapper 
[![Build Status](https://travis-ci.com/nxadm/ctwrapper.svg?token=3PQd6zsu83EBNA2LAEeq&branch=master)](https://travis-ci.com/nxadm/ctwrapper)

ctwrapper is a small git wrapper to interface with Hashicorp's 
[consul-template](https://github.com/hashicorp/consul-template).

The use case for this tool is providing remote configuration and secrets to
containers that require a more complex configuration. Many orchestrators and 
tools, like Hashicorp's [nomad](https://github.com/hashicorp/nomad), 
only provide mechanisms to provision containers with simple 1-file configuration 
requirements, e.g. by the 
[template stanza](https://www.nomadproject.io/docs/job-specification/template.html).

As an alternative, ctwrapper retrieves a git repository with static files and
templates used to create the configuration. The templates are passed as 
arguments to consul-template that will create regular files after 
injecting secrets. Options can be passed to consul-template, e.g. "-exec" to
run the actual application.

HTTP(s) Basic Authentication is supported for retrieving the git repo. The 
password can be retrieved from Hashicorp's 
[vault](https://github.com/hashicorp/vault) by providing the appropiate `VAULT_*`
[environment variables](https://www.vaultproject.io/docs/commands/index.html#environment-variables) vault address 
(at least VAULT_ADDR and VAULT_TOKEN) and a path to the secret (--vault-path).
If no authentication is provided, the repo will be retrieved anonymously.

ctwrapper, being a wrapper for consul-template, expects the latter to be in the
PATH or in the working directory.
  
```
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
  -b  | --branch    : Git branch [default: master]
  -c  | --commit    : Git commit [default: HEAD].
  -g  | --git-depth : Git depth  [default: unlimited].
  -u  | --user      : Git username.
  -p  | --password  : Git password.
  -s  | --vault-path: Vault path to the secret (including the backend).
  -e  | --ext       : Template extension [defaul: .tmpl].
  -h  | --help      : This help message.
  -v  | --version   : Version message.
  -- 				: Extra consul-template parameters, e.g. -exec.  

Examples:
  $ ctwrapper -g 10 -r https://github.com/nxadm/ctwrapper.git -d /project
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -d /project \ 
    -s "secret/production/third-party/repo-password"
  $ ctwrapper -r https://github.com/nxadm/ctwrapper.git -d /project \
    -- -o -vault-addr 'https://10.5.32.5:8200 -exec /sbin/my-server
```

You may want to set the depth to a low number (e.g.) in order not to 
unnecessarily retrieve the complete history of the repo. The `--commit` and 
`--git-depth` options were included in order to prevent a race condition 
between CI systems and git commits. If your setup ensures that the specified
commit is the last one, you can set `--git-depth` to 1.

Everything after `--` is directly passed as-is to consul-template.