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

## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap, and network with fellow Kubernetes users. [Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g) or [join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

# Documentation
Check out the [documentation at docs.fairwinds.com](https://rbac-lookup.docs.fairwinds.com/)

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



## Other Projects from Fairwinds

Enjoying rbac-lookup? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Pluto](https://github.com/FairwindsOps/Pluto) - Detect Kubernetes resources that have been deprecated or removed in future versions
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters

Or [check out the full list](https://www.fairwinds.com/open-source-software?utm_source=rbac-lookup&utm_medium=rbac-lookup&utm_campaign=rbac-lookup)
