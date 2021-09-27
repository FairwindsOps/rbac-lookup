# Contributing

Issues, whether bugs, tasks, or feature requests are essential for keeping rbac-lookup great. We believe it should be as easy as possible to contribute changes that get things working in your environment. There are a few guidelines that we need contributors to follow so that we can keep on top of things.

## Code of Conduct

This project adheres to a [code of conduct](/contributing/code-of-conduct.md). Please review this document before contributing to this project.

## Sign the CLA
Before you can contribute, you will need to sign the [Contributor License Agreement](https://cla-assistant.io/fairwindsops/rbac-lookup).

## Project Structure

rbac-lookup is a relatively simple cobra cli tool that looks up information about rbac in a cluster. The [/cmd](https://github.com/FairwindsOps/rbac-lookup/tree/master/cmd) folder contains the flags and other cobra config, while the [/lookup](https://github.com/FairwindsOps/rbac-lookup/tree/master/lookup) folder has the code for looking up rbac information a cluster. There is additinal code that allows the user to see GKE IAM information as well, since GKE IAM is so closely tied to rbac.

## Getting Started

We label issues with the ["good first issue" tag](https://github.com/FairwindsOps/rbac-lookup/labels/good%20first%20issue) if we believe they'll be a good starting point for new contributors. If you're interested in working on an issue, please start a conversation on that issue, and we can help answer any questions as they come up.

## Setting Up Your Development Environment

### Prerequisites

* A properly configured Golang environment with Go 1.17 or higher
* Access to a cluster via a properly configured KUBECONFIG

### Installation

* Install the project with `go get github.com/fairwindsops/rbac-lookup`
* Change into the rbac-lookup directory which is installed at `$GOPATH/src/github.com/fairwindsops/rbac-lookup`
* Use `make build` to build the binary locally.
* Use `make test` to run the tests and generate a coverage report.
* Use `make lint` to run [golangci-lint](https://github.com/golangci/golangci-lint) (requires golangci-lint installed)

## Creating a New Issue

If you've encountered an issue that is not already reported, please create an issue that contains the following:

- Clear description of the issue
- Steps to reproduce it
- Appropriate labels

## Creating a Pull Request

Each new pull request should:

- Reference any related issues
- Add tests that show the issues have been solved
- Pass existing tests and linting
- Contain a clear indication of if they're ready for review or a work in progress
- Be up to date and/or rebased on the master branch

## Creating a new release

Push a new annotated tag.  This tag should contain a changelog of pertinent changes. Goreleaser will take care of the rest.
