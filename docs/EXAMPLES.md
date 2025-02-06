# Examples

The section provides examples of various `bash` commands to execute `bucketeer-find-code-refs` (when installed in the system PATH) with various configurations. We recommend reading through the following examples to gain an understanding of common configurations. For more information on advanced configuration, see [CONFIGURATION.md](CONFIGURATION.md)

## Basic configuration

```bash
bucketeer-find-code-refs \
  --apiKey=$YOUR_BUCKETEER_API_KEY \ # example: api-xxxx
  --baseUri=$YOUR_BUCKETEER_BASE_URI \ # example: https://api.bucketeer.io
  --environmentId=$YOUR_BUCKETEER_ENVIRONMENT_ID \ # example: development
  --repoName=$YOUR_REPOSITORY_NAME \ # example: my-repo
  --dir="/path/to/git/repo"
```

## Configuration with context lines

```bash
bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY" \
  --baseUri="$YOUR_BUCKETEER_BASE_URI" \
  --environmentId="$YOUR_BUCKETEER_ENVIRONMENT_ID" \
  --repoName="$YOUR_REPOSITORY_NAME" \
  --dir="/path/to/git/repo" \
  --contextLines=3 # can be up to 5. If < 0, no source code will be sent to Bucketeer
```

## Configuration with repository metadata

A configuration with the the `repoType` set to GitHub, and the `repoUrl` set to a GitHub URL. We recommend configuring these parameters so Bucketeer is able to generate reference links to your source code:

```bash
bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY" \
  --baseUri="$YOUR_BUCKETEER_BASE_URI" \
  --environmentId="$YOUR_BUCKETEER_ENVIRONMENT_ID" \
  --repoName="$YOUR_REPOSITORY_NAME" \
  --dir="/path/to/git/repo" \
  --contextLines=3 \
  --repoType="github" \
  --repoUrl="$YOUR_REPOSITORY_URL" # example: https://github.com/bucketeer/bucketeer-find-code-refs
```

## Scanning non-git repositories

By default, `bucketeer-find-code-refs` will attempt to infer repository metadata from a git configuration. If you are scanning a codebase with a version control system other than git, you must use the `--revision` and `--branch` options to manually provide information about your codebase.

```bash
bucketeer-find-code-refs \
  --apiKey=$YOUR_BUCKETEER_API_KEY \
  --baseUri=$YOUR_BUCKETEER_BASE_URI \
  --environmentId=$YOUR_BUCKETEER_ENVIRONMENT_ID \
  --repoName=$YOUR_REPOSITORY_NAME \ # example: my-repo
  --dir="/path/to/git/repo" \
  --revision="REPO_REVISION_STRING" \ # e.g. a version hash
  --branch="dev"
```

### Branch garbage collection for non-git repositories

When scanning a non-git repository, automatic [branch garbage collection](../README.md#branch-garbage-collection) is disabled. The `prune` sub-command may be used to manually delete code references for stale branches.

The following example instructs the `prune` command to delete code references for the branches named "branch1" and "branch2":

```bash
bucketeer-find-code-refs prune \
  --apiKey=$YOUR_BUCKETEER_API_KEY \
  --baseUri=$YOUR_BUCKETEER_BASE_URI \
  --environmentId=$YOUR_BUCKETEER_ENVIRONMENT_ID \
  --repoName=$YOUR_REPOSITORY_NAME \ # example: my-repo
  --dir="/path/to/git/repo" \
  "branch1" "branch2"
```