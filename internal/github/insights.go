package github

import "context"

// RepositoryInsight provides a compact status snapshot for an engineering workspace.
type RepositoryInsight struct {
	Owner         string `json:"owner"`
	Repo          string `json:"repo"`
	RepositoryCount int `json:"repository_count"`
	OpenIssueCount  int `json:"open_issue_count"`
	WorkflowRunCount int `json:"workflow_run_count"`
	LatestStatus    string `json:"latest_status,omitempty"`
}

// Insight builds a compact repository health summary for a specific project context.
func (c *Client) Insight(ctx context.Context, owner, repo string) (RepositoryInsight, error) {
	repositories, err := c.Repositories(ctx, owner)
	if err != nil {
		return RepositoryInsight{}, err
	}

	issues, err := c.Issues(ctx, owner, repo)
	if err != nil {
		return RepositoryInsight{}, err
	}

	workflowRuns, err := c.WorkflowRuns(ctx, owner, repo)
	if err != nil {
		return RepositoryInsight{}, err
	}

	insight := RepositoryInsight{
		Owner:           owner,
		Repo:            repo,
		RepositoryCount: len(repositories),
		OpenIssueCount:  len(issues),
		WorkflowRunCount: len(workflowRuns),
	}

	if len(workflowRuns) > 0 {
		insight.LatestStatus = workflowRuns[0].Status
	}

	return insight, nil
}
