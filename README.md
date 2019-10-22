# kdump
Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form

Dumps everything from all configured contexts, and I mean everything.

* Calls `kubectl api-resources` to figure out what it has access to, then starts downloading all of it :).

WARNING: also dumps secrets. If you use this in for example a backup script and then commit to git (my use case), be sure to put secret in your .gitignore (or encrypt them)
