package options

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bucketeer-io/code-refs/internal/validation"
)

const (
	maxProjKeyLength = 20 // Maximum project key length
)

type RepoType string

func (repoType RepoType) isValid() error {
	switch repoType {
	case GITHUB, GITLAB, BITBUCKET, CUSTOM:
		return nil
	default:
		return fmt.Errorf(`invalid value %q for "repoType": must be %s, %s, %s, or %s`, repoType, GITHUB, GITLAB, BITBUCKET, CUSTOM)
	}
}

const (
	GITHUB    RepoType = "github"
	GITLAB    RepoType = "gitlab"
	BITBUCKET RepoType = "bitbucket"
	CUSTOM    RepoType = "custom"
)

type Project struct {
	Key     string  `mapstructure:"key"`
	Dir     string  `mapstructure:"dir"`
	Aliases []Alias `mapstructure:"aliases"`
}
type Options struct {
	ApiKey              string `mapstructure:"apiKey"`
	BaseUri             string `mapstructure:"baseUri"`
	Branch              string `mapstructure:"branch"`
	CommitUrlTemplate   string `mapstructure:"commitUrlTemplate"`
	DefaultBranch       string `mapstructure:"defaultBranch"`
	Dir                 string `mapstructure:"dir" yaml:"-"`
	OutDir              string `mapstructure:"outDir"`
	EnvironmentID       string `mapstructure:"environmentId"`
	RepoOwner           string `mapstructure:"repoOwner"`
	RepoName            string `mapstructure:"repoName"`
	RepoType            string `mapstructure:"repoType"`
	RepoUrl             string `mapstructure:"repoUrl"`
	Revision            string `mapstructure:"revision"`
	Subdirectory        string `mapstructure:"subdirectory"`
	UserAgent           string `mapstructure:"userAgent"`
	ContextLines        int    `mapstructure:"contextLines"`
	UpdateSequenceId    int    `mapstructure:"updateSequenceId"`
	AllowTags           bool   `mapstructure:"allowTags"`
	Debug               bool   `mapstructure:"debug"`
	DryRun              bool   `mapstructure:"dryRun"`
	IgnoreServiceErrors bool   `mapstructure:"ignoreServiceErrors"`
	Prune               bool   `mapstructure:"prune"`

	// The following options can only be configured via YAML configuration
	Aliases    []Alias    `mapstructure:"aliases"`
	Delimiters Delimiters `mapstructure:"delimiters"`
}

type Delimiters struct {
	// If set to `true`, the default delimiters (single-quote, double-qoute, and backtick) will not be used unless provided as `additional` delimiters
	DisableDefaults bool     `mapstructure:"disableDefaults"`
	Additional      []string `mapstructure:"additional"`
}

func Init(flagSet *pflag.FlagSet) error {
	for _, f := range flags {
		usage := strings.ReplaceAll(f.usage, "\n", " ")
		switch value := f.defaultValue.(type) {
		case string:
			flagSet.StringP(f.name, f.short, value, usage)
		case int:
			flagSet.IntP(f.name, f.short, value, usage)
		case bool:
			flagSet.BoolP(f.name, f.short, value, usage)
		}
	}

	flagSet.VisitAll(func(f *pflag.Flag) {
		viper.BindEnv(f.Name, "BUCKETEER_"+strcase.ToScreamingSnake(f.Name))
	})

	return viper.BindPFlags(flagSet)
}

func InitYAML() error {
	err := validateYAMLPreconditions()
	if err != nil {
		return err
	}
	absPath, err := validation.NormalizeAndValidatePath(viper.GetString("dir"))
	if err != nil {
		return err
	}
	subdirectoryPath := viper.GetString("subdirectory")
	viper.SetConfigName("coderefs")
	viper.SetConfigType("yaml")
	configPath := filepath.Join(absPath, subdirectoryPath, ".bucketeer")
	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return err
	}
	return nil
}

// validatePreconditions ensures required flags have been set
func validateYAMLPreconditions() error {
	baseUri := viper.GetString("baseUri")
	apiKey := viper.GetString("apiKey")
	dir := viper.GetString("dir")
	missingRequiredOptions := []string{}
	if baseUri == "" {
		missingRequiredOptions = append(missingRequiredOptions, "baseUri")
	}
	if apiKey == "" {
		missingRequiredOptions = append(missingRequiredOptions, "apiKey")
	}
	if dir == "" {
		missingRequiredOptions = append(missingRequiredOptions, "dir")
	}
	if len(missingRequiredOptions) > 0 {
		return fmt.Errorf("missing required option(s): %v", missingRequiredOptions)
	}
	return nil
}

func GetOptions() (Options, error) {
	var opts Options
	err := viper.Unmarshal(&opts)
	return opts, err
}

func GetWrapperOptions(dir string, merge func(Options) (Options, error)) (Options, error) {
	flags := pflag.CommandLine

	err := Init(flags)
	if err != nil {
		return Options{}, err
	}

	// Set precondition flags
	err = flags.Set("apiKey", os.Getenv("BUCKETEER_API_KEY"))
	if err != nil {
		return Options{}, err
	}
	err = flags.Set("dir", dir)
	if err != nil {
		return Options{}, err
	}

	err = InitYAML()
	if err != nil {
		return Options{}, err
	}

	opts, err := GetOptions()
	if err != nil {
		return opts, err
	}

	return merge(opts)
}

func (o Options) ValidateRequired() error {
	missingRequiredOptions := []string{}
	if o.ApiKey == "" {
		missingRequiredOptions = append(missingRequiredOptions, "apiKey")
	}
	if o.BaseUri == "" {
		missingRequiredOptions = append(missingRequiredOptions, "baseUri")
	}
	if o.Dir == "" {
		missingRequiredOptions = append(missingRequiredOptions, "dir")
	}
	if o.EnvironmentID == "" {
		missingRequiredOptions = append(missingRequiredOptions, "environmentId")
	}
	if o.RepoOwner == "" {
		missingRequiredOptions = append(missingRequiredOptions, "repoOwner")
	}
	if o.RepoName == "" {
		missingRequiredOptions = append(missingRequiredOptions, "repoName")
	}
	if len(missingRequiredOptions) > 0 {
		return fmt.Errorf("missing required option(s): %v", missingRequiredOptions)
	}

	return nil
}

// Validate ensures all options have been set to a valid value
func (o Options) Validate() error {
	if err := o.ValidateRequired(); err != nil {
		return err
	}

	maxContextLines := 5
	if o.ContextLines > maxContextLines {
		return fmt.Errorf(`invalid value %q for "contextLines": must be <= %d`, o.ContextLines, maxContextLines)
	}

	repoType := RepoType(strings.ToLower(o.RepoType))
	if err := repoType.isValid(); err != nil {
		return err
	}

	if o.RepoUrl != "" {
		if _, err := url.ParseRequestURI(o.RepoUrl); err != nil {
			return fmt.Errorf(`invalid value %q for "repoUrl": %+v`, o.RepoUrl, err)
		}
	}

	// match all non-control ASCII characters
	validDelims := regexp.MustCompile("^[\x20-\x7E]$")
	for i, d := range o.Delimiters.Additional {
		if !validDelims.MatchString(d) {
			return fmt.Errorf(`invalid value %q for "delimiters.additional[%d]": each delimiter must be a valid non-control ASCII character`, d, i)
		}
	}

	if _, err := validation.NormalizeAndValidatePath(o.Dir); err != nil {
		return fmt.Errorf(`invalid value for "dir": %+v`, err)
	}

	if o.OutDir != "" {
		if _, err := validation.NormalizeAndValidatePath(o.OutDir); err != nil {
			return fmt.Errorf(`invalid valid for "outDir": %+v`, err)
		}
	}

	for _, a := range o.Aliases {
		if err := a.IsValid(); err != nil {
			return err
		}
	}

	if o.Revision != "" && o.Branch == "" {
		return errors.New(`"branch" option is required when "revision" option is set`)
	}

	return nil
}

func projKeyValidation(projKey string) error {
	if strings.HasPrefix(projKey, "sdk-") {
		return fmt.Errorf("provided project key (%s) appears to be a Bucketeer SDK key", "sdk-xxxx")
	} else if strings.HasPrefix(projKey, "api-") {
		return fmt.Errorf("provided project key (%s) appears to be a Bucketeer API access token", "api-xxxx")
	}

	return nil
}
