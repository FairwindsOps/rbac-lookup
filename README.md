# rbac-lookup

This is a simple project that allows you to easily find Kubernetes roles and cluster roles bound to any user, service account, or group name. Binaries are generated with goreleaser for each release for simple installation.

## Usage

```
rbac-lookup rob

SUBJECT                   SCOPE             ROLE
rob@example.com           cluster-wide      ClusterRole/view
rob@example.com           nginx-ingress     ClusterRole/edit
```

```
rbac-lookup rob -owide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
```

```
rbac-lookup ro -owide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/ross@example.com     cluster-wide      ClusterRole/admin   ClusterRoleBinding/ross-admin
User/ron@example.com      web               ClusterRole/edit    RoleBinding/ron-edit
ServiceAccount/rops       infra             ClusterRole/admin   RoleBinding/rops-admin
```

### Kubernetes Configuration
If a `KUBECONFIG` environment variable is specified, rbac-lookup will attempt to use the config at that path, otherwise it will default to `~/.kube/config`.

## License
Apache License 2.0
