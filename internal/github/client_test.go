package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psharma365/-GitHub-Agent-Mesh/internal/config"
)

func TestClientRepositories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/octo/repos" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"demo-repo","full_name":"octo/demo-repo","description":"demo","private":false,"default_branch":"main"}]`))
	}))
	defer server.Close()

	client := NewClient(config.Settings{Host: server.URL})
	repos, err := client.Repositories(context.Background(), "octo")
	if err != nil {
		t.Fatalf("Repositories() error = %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("Repositories() len = %d, want 1", len(repos))
	}
	if repos[0].Name != "demo-repo" {
		t.Fatalf("Repositories()[0].Name = %q, want %q", repos[0].Name, "demo-repo")
	}
}
