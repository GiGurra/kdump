# kdump
Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form.

Quick and dirty hack. Don't expect pretty code :).

##### *My use case: poor man's etcd -> git sync*

Dumps all api-resources from all configured contexts.

* Calls `kubectl api-resources` to figure out what it has access to, then starts downloading all of it using `kubectl get <resource> -o yaml > <file>`.

NOTE: also dumps secrets, if you explicitly tell it to do so by providing an encryption key (aes gcm)

#### Usage


```

╰─>$ ./kdump --help

NAME:
   kdump - Dump all kubernetes resources as yaml files to a dir

USAGE:
   kdump [global options] [arguments...]

VERSION:
   2.0-go-beta

GLOBAL OPTIONS:
   --output-dir value, -o value    output directory to create
   --delete-previous-dir           if to delete previous output directory (default: false)
   --secrets-encryption-key value  symmetric secrets encryption hex key for aes GCM (lower case 64 chars)
   --help, -h                      show help (default: false)
   --version, -v                   print the version (default: false)
   ```

#### Code quality

code quality level = "hacky hack". Go is a new language to me.
tests = none.
