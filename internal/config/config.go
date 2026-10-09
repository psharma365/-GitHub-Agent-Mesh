package config

import "os"

// Settings holds the runtime configuration for the GitHub Agent Mesh runtime.
type Settings struct {
	Host  string
	Token string
	Owner string
	Repo  string
}

// Load reads configuration from environment variables and returns a validated settings object.
func Load() Settings {
	return Settings{
		Host:  firstNonEmpty(os.Getenv("GITHUB_API_URL"), "https://api.github.com"),
		Token: os.Getenv("GITHUB_TOKEN"),
		Owner: os.Getenv("GITHUB_OWNER"),
		Repo:  os.Getenv("GITHUB_REPO"),
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
