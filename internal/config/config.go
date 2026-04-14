package config

import (
	"strings"
	"unicode"
)

// Config holds all user choices from the interactive prompts.
type Config struct {
	ProjectName    string
	Backend        string // python-fastapi, go-echo, dotnet
	Frontend       string // none, nextjs, flutter
	Auth           string // cognito, pocket-id, auth0, clerk, keycloak
	DataProcessing bool
	CLIClient      string // none, tui, git-like
	Deployment     string // local, instance, scalable, robust
	CICD           string // codepipeline, bitbucket, github-actions
}

// BackendDir returns the template directory name for the selected backend.
func (c Config) BackendDir() string {
	switch c.Backend {
	case "python-fastapi":
		return "python"
	case "go-echo":
		return "go"
	case "dotnet":
		return "dotnet"
	default:
		return ""
	}
}

// DotnetNamespace returns a C#-safe PascalCase identifier from ProjectName.
// e.g., "test-claims" -> "TestClaims", "my_app" -> "MyApp"
func (c Config) DotnetNamespace() string {
	var b strings.Builder
	upper := true
	for _, r := range c.ProjectName {
		if r == '-' || r == '_' {
			upper = true
			continue
		}
		if upper {
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
