# bucketeer-code-refs

Command line program for generating Bucketeer feature flag code references.

This repository provides solutions for configuring [Bucketeer code references](https://bucketeer.io) with various systems out-of-the-box, as well as the ability to automate code reference discovery on your own infrastructure using the provided command line interface.

### Documentation quick links

- [Execution via CLI](#execution-via-cli)
  - [Prerequisites](#prerequisites)
  - [Installing](#installing)
    - [MacOS](#macOS)
    - [Linux](#linux)
    - [Windows](#windows)
    - [Docker](#docker)
- [Configuration](#configuration)
  - [Required arguments](docs/CONFIGURATION.md#required-arguments)
  - [All arguments](docs/CONFIGURATION.md#command-line)
  - [Using environment variables](docs/CONFIGURATION.md#environment-variables)
  - [Using a YAML file](docs/CONFIGURATION.md#YAML)
  - [Delimiters](docs/CONFIGURATION.md#delimiters)
  - [Ignoring files and directories](docs/CONFIGURATION.md#ignoring-files-and-directories)

## Execution via CLI

The command line program may be run manually, and executed in an environment of your choosing. The program requires your `git` repo to be cloned locally, and the currently checked out branch will be scanned for code references.

We recommend incorporating `bucketeer-code-refs` into your CI/CD build process. `bucketeer-code-refs` should run whenever a commit is pushed to your repository.

### Prerequisites

If you are scanning a git repository, `bucketeer-code-refs` requires git (tested with version 2.21.0) to be installed on the system path.

### Installing

#### macOS

```bash
brew tap bucketeer-io/code-refs https://github.com/bucketeer-io/code-refs
brew install bucketeer-find-code-refs
```

You can now run `bucketeer-find-code-refs`.

#### Linux

Supports `x86_64`, `arm64`, and `i386`. The script auto-detects your architecture and installs via `dpkg`, `rpm`, or a binary depending on your system.

```bash
curl -sfL https://raw.githubusercontent.com/bucketeer-io/code-refs/main/scripts/install.sh | bash
```

You can now run `bucketeer-find-code-refs`.

#### Windows

A Windows executable of `bucketeer-find-code-refs` is available on the [releases page](https://github.com/bucketeer-io/code-refs/releases/latest).

#### Docker

`bucketeer-find-code-refs` is available as a Docker image on the GitHub Container Registry. The image provides an entrypoint for `bucketeer-find-code-refs`, to which command line arguments may be passed. Your repository to be scanned should be mounted as a volume.

```bash
docker pull ghcr.io/bucketeer-io/bucketeer-find-code-refs:latest
docker run \
  -v /path/to/your/repo:/repo \
  ghcr.io/bucketeer-io/bucketeer-find-code-refs:latest \
  --dir="/repo"
```

### Configuration

`bucketeer-code-refs` provides a number of configuration options to customize how code references are generated and surfaced in your Bucketeer dashboard.

Required configuration:
- `API_KEY`: Your Bucketeer API key
- `API_ENDPOINT`: Your Bucketeer API endpoint
- `ENVIRONMENT_ID`: The Bucketeer environment ID
- `REPO_NAME`: The name of your repository
- `REPO_OWNER`: The owner of your repository
- `REPO_TYPE`: The type of repository (GITHUB, GITLAB, BITBUCKET, or CUSTOM)

Optional configuration:
- `BRANCH`: The git branch to scan (defaults to current branch)
- `REVISION`: The git commit SHA to scan (defaults to current commit)
- `DEBUG`: Enable debug logging
- `DRY_RUN`: Run without sending data to Bucketeer
- `IGNORE_SERVICE_ERRORS`: Continue execution even if there are API errors
- `OUT_DIR`: Directory to write CSV output files
- `USER_AGENT`: Custom user agent string for API requests

For detailed configuration options, please refer to the [configuration documentation](docs/CONFIGURATION.md).

