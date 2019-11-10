# kdump
Dumps all kubernetes api resources (pods, deployments, namespaces, etc..) to files in yaml form.

Quick and dirty hack. Don't expect pretty code :).

##### *My use case: poor man's etcd -> git sync*

Dumps all api-resources from all configured contexts.

* Calls `kubectl api-resources` to figure out what it has access to, then starts downloading all of it using `kubectl get <resource> -o yaml > <file>`.

*WARNING: also dumps secrets. If you use this in for example a backup script and then commit to git (my use case), be sure to put `*secret*` in your .gitignore (or encrypt them)*


##### Options

```
╰─>$ kdump --help
Options:
  --version                    Show version number                     [boolean]
  --context, -c                Specify contexts. If omitted - use all available
                                                                         [array]
  --exclude-context, --ec      Exclude contexts                          [array]
  --namespace, -n              Specify namespaces. If omitted - use all
                               available                                 [array]
  --exclude-namespace, --en    Exclude namespaces                        [array]
  --exclude-global, --eg       Exclude global (non-namespaced) data
                                                      [boolean] [default: false]
  --namespaced-resource, --nr  Specify namespaced resources. If omitted - use
                               all available                             [array]
  --global-resource, --gr      Specify global resources. If omitted - use all
                               available                                 [array]
  --exclude-resource, --er     Exclude resource                          [array]
  --include-secrets            If to include secrets, default false. you will
                               need to include encrypt-password or set
                               encrypt-secrets false  [boolean] [default: false]
  --encrypt-secrets            If to encrypt the secrets resource. Default and
                               recommended. To decrypt:
                               openssl enc -d -aes-256-cbc -iv hexIV -K hexKey
                                                       [boolean] [default: true]
  --encrypt-password, -p       Password for aes-256-cbc encryption of secrets
                               resource. This must be 32 bytes hex (64
                               characters).You can generate one using:
                               openssl rand -hex 32                     [string]
  --encrypt-algorithm          Encryption algorithm to use for secrets
                                               [string] [default: "aes-256-cbc"]
  --prev-dump-dir              Directory with contents of previous dump. Useful
                               to compare encrypted secrets to only replace file
                               if something actually changed. (otherwise you
                               will get a git diff every time because encryption
                               IV changes)                              [string]
  --encrypt-prev-password      encrypt-password used for prev-dump-dir, if
                               different than current                   [string]
  --encrypt-prev-algorithm     encrypt-algorithm used for prev-dump-dir, if
                               different than current                   [string]
  --output-dir, -o             Output directory                         [string]
  --include-empty-resources    If to write yaml files for resources with no
                               entries                [boolean] [default: false]
  --help                       Show help                               [boolean]

```

### Output: Downloaded directory/file structure

```
.
├── context1
│   ├── <non-namespaced-resources..>.yml
│   ├── context1_namespace1
│   |   ├── <namespaced-resources..>.yml
│   ├── context1_namespace2
│       ├── <namespaced-resources..>.yml
└── context2
    ├── <non-namespaced-resources..>.yml
    ├── context2_namespace1
    |   ├── <namespaced-resources..>.yml
    ├── context2_namespace2
        ├── <namespaced-resources..>.yml
```

## Warning

Written in Node.js, which I am a beginner in. All advice on style and best practices very welcome!
Expect things to break :S.


### Example

A rancher2 (kubernetes management platform) setup with 2 clusters (test and prod), each with a default namespace:

`cd somewhere`

`kdump`

`tree` -->


```.
├── prod
│   ├── apiservices.yml
│   ├── bgpconfigurations.yml
│   ├── cattle-system
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── certificatesigningrequests.yml
│   ├── clusterinformations.yml
│   ├── clusterrolebindings.yml
│   ├── clusterroles.yml
│   ├── componentstatuses.yml
│   ├── csidrivers.yml
│   ├── csinodes.yml
│   ├── customresourcedefinitions.yml
│   ├── default
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── felixconfigurations.yml
│   ├── globalnetworkpolicies.yml
│   ├── globalnetworksets.yml
│   ├── hostendpoints.yml
│   ├── ingress-nginx
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── ippools.yml
│   ├── kube-node-lease
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── kube-public
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── kube-system
│   │   ├── alertmanagers.yml
│   │   ├── clusterauthtokens.yml
│   │   ├── clusteruserattributes.yml
│   │   ├── configmaps.yml
│   │   ├── controllerrevisions.yml
│   │   ├── cronjobs.yml
│   │   ├── daemonsets.yml
│   │   ├── deployments.yml
│   │   ├── endpoints.yml
│   │   ├── events.yml
│   │   ├── horizontalpodautoscalers.yml
│   │   ├── ingresses.yml
│   │   ├── jobs.yml
│   │   ├── leases.yml
│   │   ├── limitranges.yml
│   │   ├── networkpolicies.yml
│   │   ├── persistentvolumeclaims.yml
│   │   ├── poddisruptionbudgets.yml
│   │   ├── pods.yml
│   │   ├── podtemplates.yml
│   │   ├── prometheuses.yml
│   │   ├── prometheusrules.yml
│   │   ├── replicasets.yml
│   │   ├── replicationcontrollers.yml
│   │   ├── resourcequotas.yml
│   │   ├── rolebindings.yml
│   │   ├── roles.yml
│   │   ├── secrets.yml
│   │   ├── serviceaccounts.yml
│   │   ├── servicemonitors.yml
│   │   ├── services.yml
│   │   └── statefulsets.yml
│   ├── mutatingwebhookconfigurations.yml
│   ├── namespaces.yml
│   ├── nodes.yml
│   ├── persistentvolumes.yml
│   ├── podsecuritypolicies.yml
│   ├── priorityclasses.yml
│   ├── runtimeclasses.yml
│   ├── storageclasses.yml
│   ├── validatingwebhookconfigurations.yml
│   └── volumeattachments.yml
└── test
    ├── apiservices.yml
    ├── bgpconfigurations.yml
    ├── cattle-system
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── certificatesigningrequests.yml
    ├── clusterinformations.yml
    ├── clusterrolebindings.yml
    ├── clusterroles.yml
    ├── componentstatuses.yml
    ├── csidrivers.yml
    ├── csinodes.yml
    ├── customresourcedefinitions.yml
    ├── default
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── felixconfigurations.yml
    ├── globalnetworkpolicies.yml
    ├── globalnetworksets.yml
    ├── hostendpoints.yml
    ├── ingress-nginx
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── ippools.yml
    ├── kube-node-lease
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── kube-public
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── kube-system
    │   ├── alertmanagers.yml
    │   ├── clusterauthtokens.yml
    │   ├── clusteruserattributes.yml
    │   ├── configmaps.yml
    │   ├── controllerrevisions.yml
    │   ├── cronjobs.yml
    │   ├── daemonsets.yml
    │   ├── deployments.yml
    │   ├── endpoints.yml
    │   ├── events.yml
    │   ├── horizontalpodautoscalers.yml
    │   ├── ingresses.yml
    │   ├── jobs.yml
    │   ├── leases.yml
    │   ├── limitranges.yml
    │   ├── networkpolicies.yml
    │   ├── persistentvolumeclaims.yml
    │   ├── poddisruptionbudgets.yml
    │   ├── pods.yml
    │   ├── podtemplates.yml
    │   ├── prometheuses.yml
    │   ├── prometheusrules.yml
    │   ├── replicasets.yml
    │   ├── replicationcontrollers.yml
    │   ├── resourcequotas.yml
    │   ├── rolebindings.yml
    │   ├── roles.yml
    │   ├── secrets.yml
    │   ├── serviceaccounts.yml
    │   ├── servicemonitors.yml
    │   ├── services.yml
    │   └── statefulsets.yml
    ├── mutatingwebhookconfigurations.yml
    ├── namespaces.yml
    ├── nodes.yml
    ├── persistentvolumes.yml
    ├── podsecuritypolicies.yml
    ├── priorityclasses.yml
    ├── runtimeclasses.yml
    ├── storageclasses.yml
    ├── validatingwebhookconfigurations.yml
    └── volumeattachments.yml

14 directories, 435 files```

