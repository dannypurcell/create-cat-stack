package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dannypurcell/create-cat-stack/internal/config"
)

// Generate walks the embedded template filesystem and renders templates
// into the output directory based on the provided config.
func Generate(cfg config.Config, outputDir string) error {
	return fs.WalkDir(TemplateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path from the templates/ root
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		// Skip the root "templates" entry itself
		if relPath == "." {
			return nil
		}

		// Skip .gitkeep files — they're only for preserving empty dirs in git
		if d.Name() == ".gitkeep" {
			return nil
		}

		// Determine which top-level template dir this belongs to
		topDir := strings.SplitN(relPath, string(filepath.Separator), 2)[0]

		// Skip backend template dirs that don't match the selected backend
		if shouldSkipDir(topDir, cfg) {
			return fs.SkipDir
		}

		// Skip etl directories if data processing is disabled
		if !cfg.DataProcessing && d.IsDir() && d.Name() == "etl" {
			return fs.SkipDir
		}

		// Skip cli directories if no CLI client is selected
		if cfg.CLIClient == "none" && d.IsDir() && d.Name() == "cli" {
			return fs.SkipDir
		}

		// Compute the output path
		outPath := filepath.Join(outputDir, relPath)

		// Strip the top-level template dir prefix so files land at project root.
		// Both the selected backend dir (e.g., "dotnet/") and "shared/" get
		// flattened into the output root.
		// e.g., "dotnet/api/Program.cs.tmpl" -> "api/Program.cs.tmpl"
		// e.g., "shared/compose.yaml.tmpl"   -> "compose.yaml.tmpl"
		if topDir == cfg.BackendDir() || topDir == "shared" {
			sub := strings.TrimPrefix(relPath, topDir+string(filepath.Separator))
			if sub == "" {
				// This is the top-level dir itself — map to output root
				outPath = outputDir
			} else {
				outPath = filepath.Join(outputDir, sub)
			}
		}

		// Render template expressions in path components (e.g., {{.ProjectName}})
		outPath = renderPathTemplates(outPath, cfg)

		if d.IsDir() {
			return os.MkdirAll(outPath, 0755)
		}

		// Read the file from the embedded FS
		data, err := fs.ReadFile(TemplateFS, path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", path, err)
		}

		// If it's a .tmpl file, render it as a Go template
		if strings.HasSuffix(outPath, ".tmpl") {
			outPath = strings.TrimSuffix(outPath, ".tmpl")

			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}

			tmpl, err := template.New(filepath.Base(path)).Parse(string(data))
			if err != nil {
				return fmt.Errorf("parsing template %s: %w", path, err)
			}

			f, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("creating %s: %w", outPath, err)
			}
			defer f.Close()

			if err := tmpl.Execute(f, cfg); err != nil {
				return fmt.Errorf("executing template %s: %w", path, err)
			}

			fmt.Printf("  rendered: %s\n", outPath)
		} else {
			// Copy non-template files as-is
			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}

			if err := os.WriteFile(outPath, data, 0644); err != nil {
				return fmt.Errorf("writing %s: %w", outPath, err)
			}

			fmt.Printf("  copied:   %s\n", outPath)
		}

		return nil
	})
}

// renderPathTemplates replaces Go template expressions in file/dir paths.
// For example, "{{.ProjectName}}.sln" becomes "test-claims.sln".
func renderPathTemplates(path string, cfg config.Config) string {
	if !strings.Contains(path, "{{") {
		return path
	}
	tmpl, err := template.New("path").Parse(path)
	if err != nil {
		return path
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return path
	}
	return buf.String()
}

// shouldSkipDir returns true if a top-level template directory should be
// skipped based on the current config. Backend dirs that don't match are skipped.
func shouldSkipDir(dir string, cfg config.Config) bool {
	backendDirs := map[string]bool{
		"python": true,
		"go":     true,
		"dotnet": true,
	}
	if backendDirs[dir] && dir != cfg.BackendDir() {
		return true
	}
	return false
}
