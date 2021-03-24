# GKE IAM Integration

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
