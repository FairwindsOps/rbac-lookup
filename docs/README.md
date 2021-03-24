<div align="center">
    <img src="img/logo.svg" height="200" alt="RBAC Lookup" style="padding-bottom: 20px" />
    <br>
    <a href="https://github.com/FairwindsOps/rbac-lookup/releases">
        <img src="https://img.shields.io/github/v/release/FairwindsOps/rbac-lookup">
    </a>
    <a href="https://goreportcard.com/report/github.com/FairwindsOps/rbac-lookup">
        <img src="https://goreportcard.com/badge/github.com/FairwindsOps/rbac-lookup">
    </a>
    <a href="https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g">
      <img src="https://img.shields.io/static/v1?label=Slack&message=Join+our+Community&color=4a154b&logo=slack">
    </a>
</div>

RBAC Lookup is a CLI that allows you to easily find Kubernetes roles and cluster roles bound to any user, service account, or group name. Binaries are generated with goreleaser for each release for simple installation.

**Want to learn more?** Reach out on [the Slack channel](https://fairwindscommunity.slack.com/messages/rbac-projects) ([request invite](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)), send an email to `opensource@fairwinds.com`, or join us for [office hours on Zoom](https://fairwindscommunity.slack.com/messages/office-hours)


## Installation

### Homebrew
```
brew install FairwindsOps/tap/rbac-lookup
```

### ASDF
```
asdf plugin add rbac-lookup
asdf install rbac-lookup latest
asdf global rbac-lookup latest
```

## RBAC Manager
While RBAC Lookup helps provide visibility into Kubernetes auth, RBAC Manager helps make auth simpler to manage. This is a Kubernetes operator that enables more concise RBAC configuration that is easier to scale and automate. For more information, see [RBAC Manager on GitHub](https://github.com/FairwindsOps/rbac-manager).

