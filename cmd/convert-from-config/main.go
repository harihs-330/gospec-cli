package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/harihs-330/gospec-cli"
	"github.com/harihs-330/gospec-cli/pkg/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Info struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Version     string `yaml:"version"`
		Contact     struct {
			Name string `yaml:"name"`
			URL  string `yaml:"url"`
		} `yaml:"contact"`
		License struct {
			Name string `yaml:"name"`
		} `yaml:"license"`
	} `yaml:"info"`
	Source struct {
		Type            string `yaml:"type"`
		Path            string `yaml:"path"`
		LocalPath       string `yaml:"localPath"`
		Framework       string `yaml:"framework"`
		RootCommandFunc string `yaml:"rootCommandFunc"`
	} `yaml:"source"`
	Output struct {
		Directory string   `yaml:"directory"`
		Formats   []string `yaml:"formats"`
		Filename  string   `yaml:"filename"`
	} `yaml:"output"`
	Options struct {
		IncludeHidden        bool   `yaml:"includeHidden"`
		IncludeDeprecated    bool   `yaml:"includeDeprecated"`
		GenerateOperationIds bool   `yaml:"generateOperationIds"`
		InferResponses       bool   `yaml:"inferResponses"`
		TagStrategy          string `yaml:"tagStrategy"`
		ExtractComponents    bool   `yaml:"extractComponents"`
	} `yaml:"options"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: convert-from-config <configspec.yaml>")
		fmt.Println("\nExample:")
		fmt.Println("  convert-from-config /path/to/configspec.yaml")
		os.Exit(1)
	}

	configPath := os.Args[1]

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "🚀 Converting CLI to OpenCLI Specification...")
	fmt.Fprintf(os.Stderr, "📄 Config: %s\n", configPath)
	fmt.Fprintf(os.Stderr, "📦 CLI: %s\n", cfg.Info.Title)
	fmt.Fprintln(os.Stderr, "")

	// Get root command (this is the tricky part - needs to load the CLI dynamically)
	rootCmd, err := loadRootCommand(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error loading CLI: %v\n", err)
		fmt.Fprintln(os.Stderr, "\n💡 Tip: Make sure the CLI package is accessible and has a GetRootCmd() function")
		os.Exit(1)
	}

	// Create output directory
	configDir := filepath.Dir(configPath)
	outputDir := filepath.Join(configDir, cfg.Output.Directory)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Setup options
	options := gospec.DefaultOptions()
	options.CustomInfo = &spec.Info{
		Title:       cfg.Info.Title,
		Description: cfg.Info.Description,
		Version:     cfg.Info.Version,
		Contact: &spec.Contact{
			Name: cfg.Info.Contact.Name,
			URL:  cfg.Info.Contact.URL,
		},
		License: &spec.License{
			Name: cfg.Info.License.Name,
		},
	}
	options.IncludeHidden = cfg.Options.IncludeHidden
	options.IncludeDeprecated = cfg.Options.IncludeDeprecated
	options.GenerateOperationIDs = cfg.Options.GenerateOperationIds
	options.InferResponses = cfg.Options.InferResponses
	options.TagStrategy = cfg.Options.TagStrategy
	options.ExtractComponents = cfg.Options.ExtractComponents

	// Create converter
	gs := gospec.New()

	// Generate specs
	for _, format := range cfg.Output.Formats {
		outputPath := filepath.Join(outputDir, cfg.Output.Filename+"."+format)
		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		if format == "yaml" {
			if err := gs.ConvertToYAML(rootCmd, options, file); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error generating YAML: %v\n", err)
				os.Exit(1)
			}
		} else if format == "json" {
			if err := gs.ConvertToJSON(rootCmd, options, file); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error generating JSON: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Fprintf(os.Stderr, "✅ Generated: %s\n", outputPath)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "🎉 Conversion complete!")
	fmt.Fprintf(os.Stderr, "📁 Output: %s\n", outputDir)
}

func loadRootCommand(cfg Config) (*cobra.Command, error) {
	// For now, this requires the CLI to be compiled as a plugin or
	// we need to use the existing approach with a helper file

	// This is a placeholder - in reality, we'd need to:
	// 1. Load the Go package dynamically
	// 2. Call the GetRootCmd() function
	// 3. Return the root command

	// For the working solution, users should use the existing generate_spec.go approach
	// or we can use Go plugins (which have limitations)

	return nil, fmt.Errorf("dynamic loading not yet implemented - please use the generate_spec.go approach for now")
}
