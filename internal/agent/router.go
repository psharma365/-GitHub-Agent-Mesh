package agent

import (
	"context"
	"fmt"

	gh "github.com/psharma365/-GitHub-Agent-Mesh/internal/github"
)

// Dispatcher routes agent actions to the appropriate GitHub tooling.
type Dispatcher struct {
	client *gh.Client
}

func NewDispatcher(client *gh.Client) *Dispatcher {
	return &Dispatcher{client: client}
}

// Execute runs a requested action and returns a structured payload.
func (d *Dispatcher) Execute(ctx context.Context, action, owner, repo string) (map[string]any, error) {
	switch action {
	case "repositories":
		repos, err := d.client.Repositories(ctx, owner)
		if err != nil {
			return nil, err
		}
		return map[string]any{"action": action, "repositories": repos}, nil
	case "issues":
		issues, err := d.client.Issues(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		return map[string]any{"action": action, "issues": issues}, nil
	case "workflows":
		runs, err := d.client.WorkflowRuns(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		return map[string]any{"action": action, "workflow_runs": runs}, nil
	case "summary", "":
		repos, err := d.client.Repositories(ctx, owner)
		if err != nil {
			return nil, err
		}
		issues, err := d.client.Issues(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		runs, err := d.client.WorkflowRuns(ctx, owner, repo)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"action": "summary",
			"owner": owner,
			"repo": repo,
			"repositories": repos,
			"issues": issues,
			"workflow_runs": runs,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}
