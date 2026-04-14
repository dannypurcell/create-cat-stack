package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "create-cat-stack",
	Short: "Generate full-stack monorepo projects",
	Long: `create-cat-stack is a project generator that scaffolds full-stack
monorepo projects with your choice of backend, frontend, auth,
data processing, CLI client, deployment, and CI/CD configuration.

Run 'create-cat-stack create <project-name>' to get started.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
