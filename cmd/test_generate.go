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
		ProjectName:    "claims-flow",
		Backend:        "dotnet",
		Frontend:       "nextjs",
		Auth:           "jumpcloud",
		DataProcessing: true,
		CLIClient:      "none",
		Deployment:     "local",
		CICD:           "bitbucket",
	}

	outputDir := "/tmp/claims-flow"
	os.RemoveAll(outputDir)

	if err := generator.Generate(cfg, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Generated successfully to", outputDir)
}
