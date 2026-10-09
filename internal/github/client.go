package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/psharma365/-GitHub-Agent-Mesh/internal/config"
)

const defaultPerPage = 10

// Client is a lightweight GitHub REST client focused on repository and workflow insight.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

// Repository represents a summary of a GitHub repository.
type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
}

// Issue represents a GitHub issue payload returned from the REST API.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	Body   string `json:"body"`
}

// WorkflowRun describes a GitHub Actions workflow run.
type WorkflowRun struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL string `json:"html_url"`
}

// Summary encapsulates the data returned to an agent for a repository.
type Summary struct {
	Owner     string       `json:"owner"`
	Repo      string       `json:"repo"`
	Repos     []Repository `json:"repositories,omitempty"`
	Issues    []Issue      `json:"issues,omitempty"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs,omitempty"`
}

func NewClient(settings config.Settings) *Client {
	return &Client{
		BaseURL: strings.TrimRight(settings.Host, "/"),
		Token:   settings.Token,
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) request(ctx context.Context, endpoint string, out any) error {
	url := c.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
		req.Header.Set("Accept", "application/vnd.github+json")
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("request %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		var payload bytes.Buffer
		_, _ = payload.ReadFrom(resp.Body)
		return fmt.Errorf("github request %s failed with %d: %s", url, resp.StatusCode, payload.String())
	}

	if out == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response from %s: %w", url, err)
	}
	return nil
}

// Repositories fetches a small set of repositories for the current owner.
func (c *Client) Repositories(ctx context.Context, owner string) ([]Repository, error) {
	endpoint := fmt.Sprintf("/users/%s/repos?per_page=%d", owner, defaultPerPage)
	var repos []Repository
	if err := c.request(ctx, endpoint, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

// Issues fetches open issues for a given repository.
func (c *Client) Issues(ctx context.Context, owner, repo string) ([]Issue, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/issues?state=open&per_page=%d", owner, repo, defaultPerPage)
	var issues []Issue
	if err := c.request(ctx, endpoint, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// WorkflowRuns fetches the latest workflow runs for a repository.
func (c *Client) WorkflowRuns(ctx context.Context, owner, repo string) ([]WorkflowRun, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/actions/runs?per_page=%d", owner, repo, defaultPerPage)
	var response struct {
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}
	if err := c.request(ctx, endpoint, &response); err != nil {
		return nil, err
	}
	return response.WorkflowRuns, nil
}
