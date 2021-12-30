# kdump

Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form.

Quick and dirty hack. Don't expect pretty code :).

#### My use case: poor man's etcd -> git sync
Dumps all api-resources from all configured contexts.

Calls kubectl api-resources to figure out what it has access to, then starts downloading all of it using kubectl get <resource> -o yaml > <file>.
NOTE: also dumps secrets, if you explicitly tell it to do so by providing an encryption key (aes gcm)


Now written in rust.

###Usage:
```
kdump --help                                
kdump 0.1.0
Dump all kubernetes resources as yaml files to a dir

USAGE:
    kdump <SUBCOMMAND>

OPTIONS:
    -h, --help       Print help information
    -V, --version    Print version information

SUBCOMMANDS:
    download                  Normal usage. Download all resources
    cluster-resource-types    List resource types available for download in the cluster
    default-excluded-types    Don't download resources - instead show default excluded types
```
```
kdump download --help
kdump-download 
Normal usage. Download all resources

USAGE:
    kdump download [OPTIONS] --output-dir <OUTPUT_DIR>

OPTIONS:
    -o, --output-dir <OUTPUT_DIR>
            REQUIRED: output directory to create

        --delete-previous-dir
            if to delete previous output directory (default: false)

        --secrets-encryption-key <SECRETS_ENCRYPTION_KEY>
            symmetric secrets encryption hex key for aes GCM (lower case 64 chars)

        --no-default-excluded-types
            disable default excluded types

        --excluded-types <EXCLUDED_TYPES>
            add additional excluded types

    -h, --help
            Print help information
```
###Examples:
```
kdump download --output-dir test --delete-previous-dir --excluded-types deployments.apps --excluded-
types services
2021-12-30 22:10:36,500 INFO [kdump] Checking what k8s types to download...
2021-12-30 22:10:37,679 INFO [kdump] Downloading all objects...
2021-12-30 22:10:44,863 INFO [kdump] Deserializing yaml...
2021-12-30 22:10:47,302 INFO [kdump] Writing yaml files...
2021-12-30 22:10:47,782 INFO [kdump] DONE!
```
```
kdump cluster-resource-types
Cluster types:
 - configmaps
 - challenges.acme.cert-manager.io
 - orders.acme.cert-manager.io
 ...
```
```
kdump default-excluded-types
Default excluded types:
 - limitranges
 - podtemplates
 - replicationcontrollers
 - resourcequotas
 - events
 - jobs
 ...
```
