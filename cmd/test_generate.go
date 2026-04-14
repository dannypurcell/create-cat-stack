//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/dannypurcell/create-cat-stack/internal/config"
	"github.com/dannypurcell/create-cat-stack/internal/generator"
)

func main() {
	cfg := config.Config{
		ProjectName:    "test-claims",
		Backend:        "dotnet",
		Frontend:       "nextjs",
		Auth:           "pocket-id",
		DataProcessing: true,
		CLIClient:      "tui",
		Deployment:     "local",
		CICD:           "github-actions",
	}

	outputDir := "/tmp/test-claims"
	os.RemoveAll(outputDir)

	if err := generator.Generate(cfg, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated successfully to", outputDir)
}
