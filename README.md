# rbac-lookup

[![Go Report Card](https://goreportcard.com/badge/github.com/reactiveops/rbac-lookup)](https://goreportcard.com/report/github.com/reactiveops/rbac-lookup) [![CircleCI](https://circleci.com/gh/reactiveops/rbac-lookup.svg?style=svg)](https://circleci.com/gh/reactiveops/rbac-lookup)

This is a simple project that allows you to easily find Kubernetes roles and cluster roles bound to any user, service account, or group name. Binaries are generated with goreleaser for each release for simple installation.

## Installation

### Homebrew
```
brew install reactiveops/tap/rbac-lookup
```

### Krew
```
kubectl krew install rbac-lookup
```

## Usage

In the simplest use case, rbac-lookup will return any matching user, service account, or group along with the roles it has been given.
```
rbac-lookup rob

SUBJECT                   SCOPE             ROLE
rob@example.com           cluster-wide      ClusterRole/view
rob@example.com           nginx-ingress     ClusterRole/edit
```

The wide output option includes the kind of subject (user, service account, or group), along with the source role binding.

```
rbac-lookup rob -owide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
```

With a more generic query, we can see that a variety of users and service accounts can be returned, as long as they match the query.
```
rbac-lookup ro -owide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/ron@example.com      web               ClusterRole/edit    RoleBinding/ron-edit
ServiceAccount/rops       infra             ClusterRole/admin   RoleBinding/rops-admin
```

Of course a query is an optional parameter for rbac-lookup. You could simply run `rbac-lookup` to get a full picture of authorization in your cluster, and then pipe that output to something like grep for your own more advanced filtering.
```
rbac-lookup | grep rob

User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
```

### GKE IAM Integration

If you're connected to a GKE cluster, RBAC is only half the story here. Google Cloud IAM roles can grant cluster access. Cluster access is effectively determined by a union of IAM and RBAC roles. To see the relevant IAM roles along with RBAC roles, use the `--gke` flag.

```
rbac-lookup rob --gke

SUBJECT              SCOPE             ROLE
rob@example.com      cluster-wide      ClusterRole/view
rob@example.com      nginx-ingress     ClusterRole/edit
rob@example.com      project-wide      IAM/gke-developer
rob@example.com      project-wide      IAM/viewer
```

Of course this GKE integration also supports wide output, in this case referencing the specific IAM roles that are assigned to a user.

```
rbac-lookup rob --gke -owide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/rob@example.com      project-wide      IAM/gke-developer   IAMRole/container.developer
User/rob@example.com      project-wide      IAM/gcp-viewer      IAMRole/viewer
```

At this point this integration only supports standard IAM roles, and is not advanced enough to include any custom roles. For a full list of supported roles and how they are mapped, view [lookup/gke_roles.go](lookup/gke_roles.go).

### Kubernetes Configuration
If a `KUBECONFIG` environment variable is specified, rbac-lookup will attempt to use the config at that path, otherwise it will default to `~/.kube/config`.

## RBAC Manager
While RBAC Lookup helps provide visibility into Kubernetes auth, RBAC Manager helps make auth simpler to manage. This is a Kubernetes operator that enables more concise RBAC configuration that is easier to scale and automate. For more information, see [RBAC Manager on GitHub](https://github.com/reactiveops/rbac-lookup).

## License
Apache License 2.0
