# kdump
Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form.

Quick and dirty hack. Don't expect pretty code :).

##### *My use case: poor man's etcd -> git sync*

Dumps all api-resources from all configured contexts.

* Calls `kubectl api-resources` to figure out what it has access to, then starts downloading all of it using `kubectl get <resource> -o yaml > <file>`.

NOTE: also dumps secrets, if you explicitly tell it to do so and provide an encryption key

### WARNING 
This is a test branch for rewriting the entire thing in go. You have been warned
