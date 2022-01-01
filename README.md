# kdump
Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form.

This is implemented in several languages (just for the heck of it) 

* [Go](https://github.com/gigurra/kdump/tree/use-go)
* [Rust](https://github.com/gigurra/kdump/tree/use-rust)
* [Js (node, the original experiment)](https://github.com/gigurra/kdump/tree/use-node)
* [Scala (not yet complete, missing cli)](https://github.com/gigurra/kdump/tree/use-scala)

##### *My use case: poor man's etcd -> git sync*

Dumps all api-resources from all configured contexts.

* Calls `kubectl api-resources` to figure out what it has access to, then starts downloading all of it using `kubectl get <resource> -o yaml > <file>`.

NOTE: also dumps secrets, if you explicitly tell it to do so by providing an encryption key (aes gcm)
