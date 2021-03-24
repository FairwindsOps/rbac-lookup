# Usage

In the simplest use case, rbac-lookup will return any matching user, service account, or group along with the roles it has been given.
```
rbac-lookup rob

SUBJECT                   SCOPE             ROLE
rob@example.com           cluster-wide      ClusterRole/view
rob@example.com           nginx-ingress     ClusterRole/edit
```

The wide output option includes the kind of subject along with the source role binding.

```
rbac-lookup rob --output wide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/ron@example.com      web               ClusterRole/edit    RoleBinding/ron-edit
ServiceAccount/rops       infra             ClusterRole/admin   RoleBinding/rops-admin
```

It's also possible to filter output by the kind of RBAC Subject. The `--kind` or `-k` parameter accepts `user`, `group`, and `serviceaccount` as values.

```
rbac-lookup ro --output wide --kind user

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/ron@example.com      web               ClusterRole/edit    RoleBinding/ron-edit
```

## Flags Supported
```
      --context string      context to use for Kubernetes config
      --gke                 enable GKE integration
  -h, --help                help for rbac-lookup
  -k, --kind string         filter by this RBAC subject kind (user, group, serviceaccount)
      --kubeconfig string   config file location
  -o, --output string       output format (normal, wide)
```
