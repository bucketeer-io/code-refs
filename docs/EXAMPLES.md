# Examples

The section provides examples of various `bash` commands to execute `bucketeer-find-code-refs` (when installed in the system PATH) with various configurations. We recommend reading through the following examples to gain an understanding of common configurations. For more information on advanced configuration, see [CONFIGURATION.md](CONFIGURATION.md)

## Basic configuration

```bash
bucketeer-find-code-refs \
  --apiKey=$YOUR_BUCKETEER_API_KEY \
  --apiEndpoint=$YOUR_BUCKETEER_API_ENDPOINT \
  --repoName=$YOUR_REPOSITORY_NAME \
  --dir="/path/to/git/repo"
```

## Configuration with multiple API keys

You can set multiple API keys using commas or by passing the `--apiKey` multiple times.

```bash
bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY1,$YOUR_BUCKETEER_API_KEY2" \
  --apiEndpoint=$YOUR_BUCKETEER_API_ENDPOINT \
  --repoName=$YOUR_REPOSITORY_NAME \
  --dir="/path/to/git/repo"

# Or

bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY1" \
  --apiKey="$YOUR_BUCKETEER_API_KEY2" \
  --apiEndpoint=$YOUR_BUCKETEER_API_ENDPOINT \
  --repoName=$YOUR_REPOSITORY_NAME \
  --dir="/path/to/git/repo"
```

## Using environment variables with multiple API keys

```bash
export BUCKETEER_APIKEY="api_key1,api_key2,api_key3"
bucketeer-find-code-refs \
  --apiEndpoint=$YOUR_BUCKETEER_API_ENDPOINT \
  --repoName=$YOUR_REPOSITORY_NAME \
  --dir="/path/to/git/repo"
```

## Configuration with context lines

```bash
bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY" \
  --apiEndpoint="$YOUR_BUCKETEER_API_ENDPOINT" \
  --repoName="$YOUR_REPOSITORY_NAME" \
  --dir="/path/to/git/repo" \
  --contextLines=3 # can be up to 5. If < 0, no source code will be sent to Bucketeer
```

## Configuration with repository metadata

A configuration with the the `repoType` set to GitHub, and the `repoUrl` set to a GitHub URL. We recommend configuring these parameters so Bucketeer is able to generate reference links to your source code:

```bash
bucketeer-find-code-refs \
  --apiKey="$YOUR_BUCKETEER_API_KEY" \
  --apiEndpoint="$YOUR_BUCKETEER_API_ENDPOINT" \
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
  --apiEndpoint=$YOUR_BUCKETEER_API_ENDPOINT \
  --repoName=$YOUR_REPOSITORY_NAME \
  --dir="/path/to/git/repo" \
  --revision="REPO_REVISION_STRING" \
  --branch="dev"
```
