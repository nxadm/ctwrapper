# ctwrapper 

ctwrapper is a small git wrapper to interface with Hashicorp's 
[consul-template](https://github.com/hashicorp/consul-template).

ctwrapper retrieves a git repository and passes all the templates it finds 
as arguments to consul-template. Addionally other options can be passed for 
consul-template (e.g. -exec).

HTTP(s) Basic Authentication is supported. The password can be retrieved from
Hashicorp's [vault](https://github.com/hashicorp/vault) by using the standard
VAULT_* environment values. If no authentication is provided, the repo will be 
retrieved anonymously.

consul-template 
must be in the PATH or in the same directory as ctwrapper.
  
```
Usage:
  vault-wrapper [-r <URL>] [-b <branch>] [-c <commit>]
				[-u <user>] [-p <password>] 
				[-vp <path> -vk <key>]
				[-d <dir>] [-e <extension>] 
                [-o <options>]  
  vault-wrapper [-h]
  vault-wrapper [-v]

Parameters:
  -r  | --url        : Git repo URL.
  -b  | --branch     : Git branch [default: master].
  -c  | --commit     : Git commit [default: HEAD].
  -u  | --user       : Git username.
  -p  | --password   : Git password.
  -vp | --vault-path : Vault path (include backend).
  -vk | --vault-key  : Vault key.
  -d  | --dir        : directory with templates [default: . ].
  -e  | --ext        : template extension [default: templ].
  -o  | --ct-opt     : extra options to pass to consul-template.
  -h  | --help       : this help message.
  -v  | --version    : version message.

Examples:
  ctwrapper -r https://github.com/nxadm/ctwrapper.git
    ctwrapper -r https://github.com/nxadm/ctwrapper.git \ 
      -vp "secret/production/third-party" -kp "api-key"
  ctwrapper -r https://github.com/nxadm/ctwrapper.git \
        -o "-vault-addr 'https://10.5.32.5:8200' -exec '/sbin/my-server'
```