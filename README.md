# GoGraph: A Graph Theory Golang Package

&emsp;&emsp;![Go Version](https://img.shields.io/badge/Go-v1.15.1-blue)
&emsp;&emsp;[![Build Status](https://cloud.drone.io/api/badges/jrcasso/gograph/status.svg?ref=refs/heads/develop)](https://cloud.drone.io/jrcasso/gograph)
&emsp;&emsp;![](https://img.shields.io/github/issues/jrcasso/gograph)
&emsp;&emsp;![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/jrcasso/gograph?include_prereleases)

## Summary

While taking some time off work, I wanted to learn some Golang for purposes of professional growth. This is the first of several packages I took a solid swing at. Golang has become a very pervasive language in recent years, often found under the hood of a lot of paradigm-shifting tooling in the software industry. Software such as Docker, Kubernetes, Terraform, and several recent continuous integration systems (e.g. Drone) have used it to great effect. There is good evidence to suggest that learning Golang would be valuable for one's continuous improvement as an engineer.

This project endeavors to implement the following:
 - a working toolbox of graph theory primitives in the form of structs and functions
 - a complete set of functional unit tests using Golang's native `go test` functionality
 - a demonstration of the functionalities it provides via a main file
 - a modest, project-size appropriate CI build process
 - a semantically versioned progression of package improvements

# Development Setup

Ensure you have the following prerequisites satisfied:
 - VS Code Extensions: Remote Containers
   - Download and install Microsoft's VS Code extension for developing in [Remote Containers](vscode:extension/ms-vscode-remote.remote-containers)


>Note: This is a VS Code Remote Containers development project: all development is done within a container to reduce initial time-to-develop. Getting this project up and running on your machine can be as simple as pulling down the repository, opening the project in VS Code, and clicking twice.


## Directions

- Clone the repository

```sh
git clone git@github.com:jrcasso/gograph
```

- Open the repository in VS Code
```sh
code gograph
```

- In the bottom-left corner of the VS Code window, click, the highlighted "><" button (or navigate to the Remote Containers extension).
- From the dropdown, select "Remote Containers: Reopen in Container"

That's it.

## Development Details

VS Code will begin to build an image that is specified in `.devcontainer/`; it will be the container image that you develop in. When it's done, it'll automatically throw your entire VS Code interface/environment inside that container where you may begin deveopment. Utilitarian tools like git and all the things needed to run a Go program are in that environment. It's still a container, so all of the idempotency and innate destructivity of containers are in fact *features* of this development strategy. If everyone develops in the same way, the time-to-develop becomes incredibly small. This model, of course, fails for large projects - but this is not a large project. In any case, one can simply add Docker composition orchestration when that time comes.

Additional tooling that might be needed can be done so during container runtime; however, if it is something that should stick around for every other developer too (i.e. they might also run into this same issue), please modify the `.devcontainer/Dockerfile` and open a pull request.

Launch configurations have been created to facilitate debugging with VS Code's native debugger for both normal script execution as well as tests.
