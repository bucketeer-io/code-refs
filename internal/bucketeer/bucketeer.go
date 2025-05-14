package bucketeer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bucketeer-io/code-refs/internal/log"
	"github.com/bucketeer-io/code-refs/options"
	"github.com/cenkalti/backoff/v4"
)

const (
	DefaultPageSize    = 1000 // Maximum number of items to return per page
	serverErrorMinCode = 500  // HTTP 5xx status codes indicate server errors
)

type ReferenceHunksRep struct {
	Path  string
	Hunks []HunkRep
}

type HunkRep struct {
	FlagKey            string
	StartingLineNumber int
	Lines              string
	ContentHash        string
	Aliases            []string
	FileExt            string
}

// NumLines returns the number of lines in the hunk
func (h HunkRep) NumLines() int {
	return strings.Count(h.Lines, "\n") + 1
}

// Overlap returns the number of overlapping lines between the receiver (h) and the parameter (hr) hunkreps
// The return value will be negative if the hunks do not overlap
func (h HunkRep) Overlap(hr HunkRep) int {
	aLines := strings.Split(h.Lines, "\n")
	bLines := strings.Split(hr.Lines, "\n")

	aStart := h.StartingLineNumber
	aEnd := aStart + len(aLines)
	bStart := hr.StartingLineNumber
	bEnd := bStart + len(bLines)

	if bStart > aEnd || aStart > bEnd {
		return -1
	}

	if bStart >= aStart {
		return len(aLines) - (bStart - aStart)
	}
	return len(bLines) - (aStart - bStart)
}

type ApiClient interface {
	GetFlagKeyList(ctx context.Context, opts options.Options) ([]string, error)
	CreateCodeReference(ctx context.Context, opts options.Options, ref CodeReference) error
	UpdateCodeReference(ctx context.Context, opts options.Options, id string, ref CodeReference) error
	DeleteCodeReference(ctx context.Context, opts options.Options, id string) error
	ListCodeReferences(ctx context.Context, opts options.Options, featureId string, pageSize int64) (codeRefs []CodeReference, cursor string, totalCount string, err error)
}

type ApiOptions struct {
	ApiKey      string
	ApiEndpoint string
	UserAgent   string
	RetryMax    *int
}

type RepoParams struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Url  string `json:"url,omitempty"`
}

type CodeReference struct {
	ID               string   `json:"id,omitempty"`
	FeatureID        string   `json:"featureId"`
	FilePath         string   `json:"filePath"`
	FileExtension    string   `json:"fileExtension"`
	LineNumber       int      `json:"lineNumber"`
	CodeSnippet      string   `json:"codeSnippet"`
	ContentHash      string   `json:"contentHash"`
	Aliases          []string `json:"aliases"`
	RepositoryName   string   `json:"repositoryName"`
	RepositoryOwner  string   `json:"repositoryOwner"`
	RepositoryType   string   `json:"repositoryType"`
	RepositoryBranch string   `json:"repositoryBranch"`
	CommitHash       string   `json:"commitHash"`
	EnvironmentID    string   `json:"environmentId"`
	CreatedAt        string   `json:"createdAt,omitempty"`
	UpdatedAt        string   `json:"updatedAt,omitempty"`
}

type apiClient struct {
	apiKey      string
	apiEndpoint string
	userAgent   string
	retryMax    int
	client      *http.Client
}

//nolint:ireturn
func InitApiClient(opts ApiOptions) ApiClient {
	retryMax := 3
	if opts.RetryMax != nil {
		retryMax = *opts.RetryMax
	}
	return &apiClient{
		apiKey:      opts.ApiKey,
		apiEndpoint: opts.ApiEndpoint,
		userAgent:   opts.UserAgent,
		retryMax:    retryMax,
		client:      &http.Client{},
	}
}

func (c *apiClient) do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	op := func() error {
		var err error
		//nolint:bodyclose
		resp, err = c.client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode >= serverErrorMinCode {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("server error: %s", body)
		}
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Duration(c.retryMax) * time.Second

	err := backoff.Retry(op, b)
	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return nil, err
	}

	return resp, nil
}

func (c *apiClient) GetFlagKeyList(ctx context.Context, opts options.Options) ([]string, error) {
	url := c.apiEndpoint + "/v1/features?pageSize=" + strconv.Itoa(DefaultPageSize)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if opts.Debug {
		body, _ := io.ReadAll(resp.Body)
		log.Debug.Printf("[GetFlagKeyList] Response Status: %d, Body: %s", resp.StatusCode, string(body))
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var response struct {
		Features []struct {
			ID string `json:"id"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	flags := make([]string, 0, len(response.Features))
	for _, feature := range response.Features {
		flags = append(flags, feature.ID)
	}
	return flags, nil
}

func (c *apiClient) CreateCodeReference(ctx context.Context, opts options.Options, ref CodeReference) error {
	url := c.apiEndpoint + "/v1/code_reference"
	body, err := json.Marshal(ref)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if opts.Debug {
		body, _ := io.ReadAll(resp.Body)
		log.Debug.Printf("[CreateCodeReference] Response Status: %d, Body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *apiClient) UpdateCodeReference(ctx context.Context, opts options.Options, id string, ref CodeReference) error {
	ref.ID = id
	url := c.apiEndpoint + "/v1/code_reference"
	body, err := json.Marshal(ref)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if opts.Debug {
		body, _ := io.ReadAll(resp.Body)
		log.Debug.Printf("[UpdateCodeReference] Response Status: %d, Body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *apiClient) DeleteCodeReference(ctx context.Context, opts options.Options, id string) error {
	url := c.apiEndpoint + "/v1/code_reference?id=" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if opts.Debug {
		body, _ := io.ReadAll(resp.Body)
		log.Debug.Printf("[DeleteCodeReference] Response Status: %d, Body: %s", resp.StatusCode, string(body))
	}
	return nil
}

type ListCodeReferencesResponse struct {
	CodeReferences []CodeReference `json:"codeReferences"`
	Cursor         string          `json:"cursor"`
	TotalCount     string          `json:"totalCount"`
}

func (c *apiClient) ListCodeReferences(ctx context.Context, opts options.Options, featureId string, pageSize int64) (codeRefs []CodeReference, cursor string, totalCount string, err error) {
	url := c.apiEndpoint + "/v1/code_references?featureId=" + featureId
	if pageSize > 0 {
		url += "&pageSize=" + strconv.FormatInt(pageSize, 10)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", "", err
	}
	req.Header.Add("Authorization", c.apiKey)
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.do(req)
	if err != nil {
		return nil, "", "", err
	}
	defer resp.Body.Close()

	if opts.Debug {
		body, _ := io.ReadAll(resp.Body)
		log.Debug.Printf("[ListCodeReferences] Response Status: %d, Body: %s", resp.StatusCode, string(body))
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var response ListCodeReferencesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, "", "", err
	}

	return response.CodeReferences, response.Cursor, response.TotalCount, nil
}
