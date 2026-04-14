package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dannypurcell/create-cat-stack/internal/config"
	"github.com/dannypurcell/create-cat-stack/internal/generator"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new full-stack monorepo project",
	Long:  `Interactively configure and generate a full-stack monorepo project.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	cfg := config.Config{}

	// Pre-fill project name from arg if provided
	if len(args) > 0 {
		cfg.ProjectName = args[0]
	}

	// Build interactive form
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("The name for your new project (lowercase, hyphens ok)").
				Value(&cfg.ProjectName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("project name is required")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Backend Framework").
				Options(
					huh.NewOption("Python / FastAPI", "python-fastapi"),
					huh.NewOption("Go / Echo", "go-echo"),
					huh.NewOption("C# / .NET", "dotnet"),
				).
				Value(&cfg.Backend),

			huh.NewSelect[string]().
				Title("Frontend Framework").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("Next.js (web)", "nextjs"),
					huh.NewOption("Flutter (web + mobile)", "flutter"),
				).
				Value(&cfg.Frontend),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Authentication Provider").
				Options(
					huh.NewOption("JumpCloud (external OIDC)", "jumpcloud"),
					huh.NewOption("Cloud-native (Cognito)", "cognito"),
					huh.NewOption("Pocket-ID (self-hosted)", "pocket-id"),
					huh.NewOption("Auth0", "auth0"),
					huh.NewOption("Clerk", "clerk"),
					huh.NewOption("Keycloak", "keycloak"),
				).
				Value(&cfg.Auth),

			huh.NewConfirm().
				Title("Include data processing (ETL pipeline)?").
				Value(&cfg.DataProcessing),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("CLI Client").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("TUI (Bubble Tea)", "tui"),
					huh.NewOption("Git-like CLI", "git-like"),
				).
				Value(&cfg.CLIClient),

			huh.NewSelect[string]().
				Title("Deployment Target").
				Description("Local only = docker-compose; Instance = single EC2/VM; Scalable = ECS/K8s; Robust = multi-region").
				Options(
					huh.NewOption("Local only", "local"),
					huh.NewOption("Instance", "instance"),
					huh.NewOption("Scalable", "scalable"),
					huh.NewOption("Robust", "robust"),
				).
				Value(&cfg.Deployment),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("CI/CD Platform").
				Options(
					huh.NewOption("AWS CodePipeline", "codepipeline"),
					huh.NewOption("Bitbucket Pipelines", "bitbucket"),
					huh.NewOption("GitHub Actions", "github-actions"),
				).
				Value(&cfg.CICD),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("prompt cancelled: %w", err)
	}

	// Print summary
	fmt.Println()
	fmt.Printf("Creating project: %s\n", cfg.ProjectName)
	fmt.Printf("  Backend:         %s\n", cfg.Backend)
	fmt.Printf("  Frontend:        %s\n", cfg.Frontend)
	fmt.Printf("  Auth:            %s\n", cfg.Auth)
	fmt.Printf("  Data Processing: %v\n", cfg.DataProcessing)
	fmt.Printf("  CLI Client:      %s\n", cfg.CLIClient)
	fmt.Printf("  Deployment:      %s\n", cfg.Deployment)
	fmt.Printf("  CI/CD:           %s\n", cfg.CICD)
	fmt.Println()

	// Generate project
	outputDir := cfg.ProjectName
	if err := generator.Generate(cfg, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %q created successfully in ./%s/\n", cfg.ProjectName, outputDir)
	return nil
}
