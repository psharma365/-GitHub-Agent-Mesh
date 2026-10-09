package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psharma365/-GitHub-Agent-Mesh/internal/config"
	gh "github.com/psharma365/-GitHub-Agent-Mesh/internal/github"
)

func TestDispatcherSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/users/octo/repos":
			_, _ = w.Write([]byte(`[{"name":"demo-repo","full_name":"octo/demo-repo","description":"demo","private":false,"default_branch":"main"}]`))
		case "/repos/octo/demo/actions/runs":
			_, _ = w.Write([]byte(`{"workflow_runs":[{"id":1,"name":"ci","status":"completed","conclusion":"success","html_url":"https://example.com/run"}]}`))
		case "/repos/octo/demo/issues":
			_, _ = w.Write([]byte(`[{"number":42,"title":"Investigate issue","state":"open","body":"example body"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := gh.NewClient(config.Settings{Host: server.URL})
	dispatcher := NewDispatcher(client)

	result, err := dispatcher.Execute(context.Background(), "summary", "octo", "demo")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result["action"] != "summary" {
		t.Fatalf("result[action] = %v, want summary", result["action"])
	}
	if _, ok := result["insight"]; !ok {
		t.Fatalf("summary payload missing insight field")
	}
}
