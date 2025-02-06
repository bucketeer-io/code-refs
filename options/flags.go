package options

type flag struct {
	name         string
	short        string
	defaultValue interface{}
	usage        string
}

const (
	defaultContextLines = 2 // Default number of context lines to include
)

// Options that are available as command line flags
var flags = []flag{
	{
		name:         "apiKey",
		short:        "k",
		defaultValue: "",
		usage:        "Bucketeer API key with write-level access.",
	},
	{
		name:         "allowTags",
		defaultValue: false,
		usage:        "Enables storing references for tags. The tag will be listed as a branch.",
	},
	{
		name:         "baseUri",
		short:        "U",
		defaultValue: "",
		usage:        "Bucketeer base URI.",
	},
	{
		name:         "branch",
		short:        "b",
		defaultValue: "",
		usage: `The currently checked out branch. If not provided, branch
name will be auto-detected. Provide this option when using CI systems that
leave the repository in a detached HEAD state.`,
	},
	{
		name:         "commitUrlTemplate",
		defaultValue: "",
		usage: `If provided, Bucketeer will attempt to generate links to
your VCS service provider per commit.
Example: https://github.com/bucketeer/bucketeer-find-code-refs/commit/${sha}.
Allowed template variables: 'branchName', 'sha'. If "commitUrlTemplate" is not provided, but "repoUrl" is provided and "repoType" is not custom, Bucketeer will attempt to automatically generate source code links for the given "repoType".`,
	},
	{
		name:         "contextLines",
		short:        "C",
		defaultValue: defaultContextLines,
		usage: `The number of context lines to send to Bucketeer. If < 0, no
source code will be sent to Bucketeer. If 0, only the lines containing
flag references will be sent. If > 0, will send that number of context
lines above and below the flag reference. A maximum of 5 context lines
may be provided.`,
	},
	{
		name:         "debug",
		defaultValue: false,
		usage:        "Enables verbose debug logging",
	},
	{
		name:         "defaultBranch",
		short:        "B",
		defaultValue: "main",
		usage: `The default branch. The Bucketeer UI will default to this branch.
If not provided, will fallback to 'main'.`,
	},
	{
		name:         "dir",
		short:        "d",
		defaultValue: "",
		usage:        "Path to existing checkout of the repository.",
	},
	{
		name:         "dryRun",
		defaultValue: false,
		usage: `If enabled, the scanner will run without sending code references to
Bucketeer. Combine with the outDir option to output code references to a CSV.`,
	},
	{
		name:         "environmentId",
		short:        "e",
		defaultValue: "",
		usage:        `Bucketeer environment ID. The API key must have access to this environment.`,
	},
	{
		name:         "ignoreServiceErrors",
		short:        "i",
		defaultValue: false,
		usage: `If enabled, the scanner will terminate with exit code 0 when the
Bucketeer API is unreachable or returns an unexpected response.`,
	},
	{
		name:         "outDir",
		short:        "o",
		defaultValue: "",
		usage: `If provided, will output a csv file containing all code references for
the project to this directory.`,
	},
	{
		name:         "prune",
		defaultValue: true,
		usage:        `If enabled, branches that are not found in the remote repository will be deleted from Bucketeer.`,
	},
	{
		name:         "repoName",
		short:        "r",
		defaultValue: "",
		usage: `Repository name. Will be displayed in Bucketeer. Case insensitive.
Repository names must only contain letters, numbers, '.', '_' or '-'."`,
	},
	{
		name:         "repoType",
		short:        "T",
		defaultValue: "custom",
		usage: `The repo service provider. Used to correctly categorize repositories in the
Bucketeer UI. Acceptable values: bitbucket|custom|github|gitlab.`,
	},
	{
		name:         "repoUrl",
		short:        "u",
		defaultValue: "",
		usage:        `The URL for the repository. If provided and "repoType" is not custom, Bucketeer will attempt to automatically generate source code links for the given "repoType".`,
	},
	{
		name:         "revision",
		short:        "R",
		defaultValue: "",
		usage:        `Use this option to scan non-git codebases. The current revision of the repository to be scanned. If set, the version string for the scanned repository will not be inferred, and branch garbage collection will be disabled. The "branch" option is required when "revision" is set.`,
	},
	{
		name:         "subdirectory",
		defaultValue: "",
		usage: `If the .bucketeer/coderefs.yaml file is not in the root of
the repository, provide the path to the subdirectory containing the configuration,
relative to the root. Code references will only run on this provided subdirectory.
This allows a monorepo to have multiple configuration files, one per subdirectory.`,
	},
	{
		name:         "updateSequenceId",
		short:        "s",
		defaultValue: -1,
		usage: `An integer representing the order number of code reference updates.
Used to version updates across concurrent executions of the flag finder.
If not provided, data will always be updated. If provided, data will
only be updated if the existing "updateSequenceId" is less than the new
"updateSequenceId". Examples: the time a "git push" was initiated, CI
build number, the current unix timestamp.`,
	},
	{
		name:         "userAgent",
		defaultValue: "",
		usage:        `(Internal) Platform where code references is run.`,
	},
}
