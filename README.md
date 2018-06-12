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
[vault](https://github.com/hashicorp/vault) by using the standard
VAULT_* environment values. If no authentication is provided, the repo will be 
retrieved anonymously.

ctwrapper, being a wrapper for consul-template, expects the latter to be in the
PATH or in the working directory.
  
```
Usage:
  vault-wrapper [-r <URL>] [-b <branch>] [-c <commit>] [-gd  <nr of commits>]
                [-u <user>] [-p <password> | -vp <path> -vk <key>]
                [-d <dir>] [-e <extension>] 
                [-o <options>]  
  vault-wrapper [-h]
  vault-wrapper [-v]

Parameters:
  -r  | --repo      : Git repo URL.
  -b  | --branch    : Git branch [default: master]
  -c  | --commit    : Git commit [default: HEAD].
  -gd | --git-depth : Git depth  [default: unlimited].
  -u  | --user      : Git username.
  -p  | --password  : Git password.
  -vp | --vault-path: Vault path (include backend).
  -vk | --vault-key : Vault key.
  -d  | --dir       : Directory with templates [default: . ].
  -e  | --ext       : Template extension [defaul: .tmpl].
  -o  | --ct-opt    : Extra options to pass to consul-template.
  -h  | --help      : This help message.
  -v  | --version   : Version message.

Examples:
  ctwrapper -r https://github.com/nxadm/ctwrapper.git
  ctwrapper -r https://github.com/nxadm/ctwrapper.git \ 
      -vp "secret/production/third-party" -kp "api-key"
  ctwrapper -r https://github.com/nxadm/ctwrapper.git \
        -o "-vault-addr 'https://10.5.32.5:8200' -exec '/sbin/my-server'
```