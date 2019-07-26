# rbac-lookup

[![Go Report Card](https://goreportcard.com/badge/github.com/FairwindsOps/rbac-lookup)](https://goreportcard.com/report/github.com/FairwindsOps/rbac-lookup) [![CircleCI](https://circleci.com/gh/FairwindsOps/rbac-lookup.svg?style=svg)](https://circleci.com/gh/FairwindsOps/rbac-lookup) [![codecov](https://codecov.io/gh/FairwindsOps/rbac-lookup/branch/master/graph/badge.svg)](https://codecov.io/gh/FairwindsOps/rbac-lookup)

This is a simple project that allows you to easily find Kubernetes roles and cluster roles bound to any user, service account, or group name. Binaries are generated with goreleaser for each release for simple installation.

**Want to learn more?** Fairwinds holds [office hours on Zoom](https://zoom.us/j/242508205) the first Friday of every month, at 12pm Eastern. You can also reach out via email at `opensource@fairwinds.com`

## Installation

### Homebrew
```
brew install FairwindsOps/tap/rbac-lookup
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
rbac-lookup rob --gke --output wide

SUBJECT                   SCOPE             ROLE                SOURCE
User/rob@example.com      cluster-wide      ClusterRole/view    ClusterRoleBinding/rob-cluster-view
User/rob@example.com      nginx-ingress     ClusterRole/edit    RoleBinding/rob-edit
User/rob@example.com      project-wide      IAM/gke-developer   IAMRole/container.developer
User/rob@example.com      project-wide      IAM/gcp-viewer      IAMRole/viewer
```

At this point this integration only supports standard IAM roles, and is not advanced enough to include any custom roles. For a full list of supported roles and how they are mapped, view [lookup/gke_roles.go](lookup/gke_roles.go).

## Flags Supported
```
      --context string   context to use for Kubernetes config
      --gke              enable GKE integration
  -h, --help             help for rbac-lookup
  -k, --kind string      filter by this RBAC subject kind (user, group, serviceaccount)
  -o, --output string    output format (normal, wide)
```

## RBAC Manager
While RBAC Lookup helps provide visibility into Kubernetes auth, RBAC Manager helps make auth simpler to manage. This is a Kubernetes operator that enables more concise RBAC configuration that is easier to scale and automate. For more information, see [RBAC Manager on GitHub](https://github.com/FairwindsOps/rbac-manager).

## Contributing
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Roadmap](ROADMAP.md)
- [Changelog](https://github.com/FairwindsOps/rbac-lookup/releases)

## License
Apache License 2.0
