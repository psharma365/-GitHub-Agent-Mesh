package github

import "context"

// SecuritySignal summarizes repository health around delivery stability and issue load.
type SecuritySignal struct {
	Owner             string `json:"owner"`
	Repo              string `json:"repo"`
	OpenIssueCount    int    `json:"open_issue_count"`
	WorkflowRunCount  int    `json:"workflow_run_count"`
	LatestWorkflowState string `json:"latest_workflow_state,omitempty"`
}

// SecuritySignals produces a lightweight operational risk signal for agentic workflows.
func (c *Client) SecuritySignals(ctx context.Context, owner, repo string) (SecuritySignal, error) {
	issues, err := c.Issues(ctx, owner, repo)
	if err != nil {
		return SecuritySignal{}, err
	}

	workflowRuns, err := c.WorkflowRuns(ctx, owner, repo)
	if err != nil {
		return SecuritySignal{}, err
	}

	signal := SecuritySignal{
		Owner:            owner,
		Repo:             repo,
		OpenIssueCount:   len(issues),
		WorkflowRunCount: len(workflowRuns),
	}

	if len(workflowRuns) > 0 {
		signal.LatestWorkflowState = workflowRuns[0].Status
	}

	return signal, nil
}
