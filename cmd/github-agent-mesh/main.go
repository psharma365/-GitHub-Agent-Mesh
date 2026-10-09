package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/psharma365/-GitHub-Agent-Mesh/internal/agent"
	"github.com/psharma365/-GitHub-Agent-Mesh/internal/config"
	"github.com/psharma365/-GitHub-Agent-Mesh/internal/github"
)

func main() {
	var (
		owner = flag.String("owner", "", "GitHub repository owner")
		repo  = flag.String("repo", "", "GitHub repository name")
		action = flag.String("action", "summary", "tool action to run (summary, repositories, issues, workflows)")
		token = flag.String("token", "", "GitHub personal access token")
		host  = flag.String("host", "https://api.github.com", "GitHub API host")
	)
	flag.Parse()

	settings := config.Load()
	if *owner != "" {
		settings.Owner = *owner
	}
	if *repo != "" {
		settings.Repo = *repo
	}
	if *token != "" {
		settings.Token = *token
	}
	if *host != "" {
		settings.Host = *host
	}

	if settings.Owner == "" || settings.Repo == "" {
		log.Fatal("both --owner and --repo are required")
	}

	ctx := context.Background()
	client := github.NewClient(settings)
	dispatcher := agent.NewDispatcher(client)

	result, err := dispatcher.Execute(ctx, *action, settings.Owner, settings.Repo)
	if err != nil {
		log.Fatalf("action %q failed: %v", *action, err)
	}

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("marshal response: %v", err)
	}

	fmt.Fprintln(os.Stdout, string(payload))
}
