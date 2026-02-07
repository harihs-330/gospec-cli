package gospec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/harihs-330/gospec-cli/pkg/config"
	"github.com/harihs-330/gospec-cli/pkg/converter"
	"github.com/harihs-330/gospec-cli/pkg/generator"
	"github.com/harihs-330/gospec-cli/pkg/parser"
	"github.com/harihs-330/gospec-cli/pkg/parser/cobra"
	"github.com/harihs-330/gospec-cli/pkg/spec"
)

// GoSpec is the main entry point for CLI to OpenCLI Spec conversion
type GoSpec struct {
	registry  *parser.ParserRegistry
	converter parser.Converter
}

// New creates a new GoSpec instance with default parsers
func New() *GoSpec {
	registry := parser.NewParserRegistry()

	// Register default parsers
	registry.Register(cobra.NewCobraParser())
	// Add more parsers here as they are implemented

	return &GoSpec{
		registry:  registry,
		converter: converter.NewDefaultConverter(),
	}
}

// NewWithRegistry creates a GoSpec instance with a custom parser registry
func NewWithRegistry(registry *parser.ParserRegistry) *GoSpec {
	return &GoSpec{
		registry:  registry,
		converter: converter.NewDefaultConverter(),
	}
}

// SetConverter sets a custom converter
func (g *GoSpec) SetConverter(conv parser.Converter) {
	g.converter = conv
}

// RegisterParser registers a new parser
func (g *GoSpec) RegisterParser(p parser.Parser) {
	g.registry.Register(p)
}

// Convert converts a CLI application to OpenCLI Specification
func (g *GoSpec) Convert(source interface{}, options *parser.ConvertOptions) (*spec.OpenCLISpec, error) {
	// Find suitable parser
	p, err := g.registry.FindParser(source)
	if err != nil {
		return nil, fmt.Errorf("failed to find parser: %w", err)
	}

	// Parse the CLI
	parsed, err := p.Parse(source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CLI: %w", err)
	}

	// Convert to OpenCLI spec
	openCLI, err := g.converter.Convert(parsed, options)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to OpenCLI spec: %w", err)
	}

	return openCLI, nil
}

// ConvertToYAML converts a CLI application to OpenCLI Specification in YAML format
func (g *GoSpec) ConvertToYAML(source interface{}, options *parser.ConvertOptions, writer io.Writer) error {
	openCLI, err := g.Convert(source, options)
	if err != nil {
		return err
	}

	gen := generator.NewYAMLGenerator()
	return gen.Generate(openCLI, writer)
}

// ConvertToJSON converts a CLI application to OpenCLI Specification in JSON format
func (g *GoSpec) ConvertToJSON(source interface{}, options *parser.ConvertOptions, writer io.Writer) error {
	openCLI, err := g.Convert(source, options)
	if err != nil {
		return err
	}

	gen := generator.NewJSONGenerator()
	return gen.Generate(openCLI, writer)
}

// ConvertToYAMLString converts a CLI application to OpenCLI Specification YAML string
func (g *GoSpec) ConvertToYAMLString(source interface{}, options *parser.ConvertOptions) (string, error) {
	openCLI, err := g.Convert(source, options)
	if err != nil {
		return "", err
	}

	gen := generator.NewYAMLGenerator()
	return gen.GenerateToString(openCLI)
}

// ConvertToJSONString converts a CLI application to OpenCLI Specification JSON string
func (g *GoSpec) ConvertToJSONString(source interface{}, options *parser.ConvertOptions) (string, error) {
	openCLI, err := g.Convert(source, options)
	if err != nil {
		return "", err
	}

	gen := generator.NewJSONGenerator()
	return gen.GenerateToString(openCLI)
}

// ListParsers returns a list of registered parser names
func (g *GoSpec) ListParsers() []string {
	return g.registry.List()
}

// GetParser retrieves a parser by name
func (g *GoSpec) GetParser(name string) (parser.Parser, bool) {
	return g.registry.Get(name)
}

// ParseWith parses a CLI using a specific parser by name
func (g *GoSpec) ParseWith(parserName string, source interface{}) (*parser.ParsedCLI, error) {
	p, ok := g.registry.Get(parserName)
	if !ok {
		return nil, fmt.Errorf("parser '%s' not found", parserName)
	}

	return p.Parse(source)
}

// DefaultOptions returns default conversion options
func DefaultOptions() *parser.ConvertOptions {
	return converter.DefaultConvertOptions()
}

// ConvertFromConfig reads a configspec.yaml file and generates OpenCLI specifications
// This is the simplest way for users - just provide a config file!
func (g *GoSpec) ConvertFromConfig(configPath string, source interface{}) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Get config directory for relative paths
	configDir := filepath.Dir(configPath)
	outputDir := filepath.Join(configDir, cfg.Output.Directory)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Setup conversion options from config
	options := &parser.ConvertOptions{
		SpecVersion:          "1.0.0",
		IncludeHidden:        cfg.Options.IncludeHidden,
		IncludeDeprecated:    cfg.Options.IncludeDeprecated,
		GenerateOperationIDs: cfg.Options.GenerateOperationIds,
		InferResponses:       cfg.Options.InferResponses,
		TagStrategy:          cfg.Options.TagStrategy,
		ExtractComponents:    cfg.Options.ExtractComponents,
		CustomInfo: &spec.Info{
			Title:       cfg.Info.Title,
			Description: cfg.Info.Description,
			Version:     cfg.Info.Version,
			Contact: &spec.Contact{
				Name: cfg.Info.Contact.Name,
				URL:  cfg.Info.Contact.URL,
			},
			License: &spec.License{
				Name: cfg.Info.License.Name,
				URL:  cfg.Info.License.URL,
			},
		},
	}

	// Generate specs in requested formats
	for _, format := range cfg.Output.Formats {
		outputPath := filepath.Join(outputDir, cfg.Output.Filename+"."+format)
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer file.Close()

		switch format {
		case "yaml":
			if err := g.ConvertToYAML(source, options, file); err != nil {
				return fmt.Errorf("failed to generate YAML: %w", err)
			}
		case "json":
			if err := g.ConvertToJSON(source, options, file); err != nil {
				return fmt.Errorf("failed to generate JSON: %w", err)
			}
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}

		fmt.Fprintf(os.Stderr, "✅ Generated: %s\n", outputPath)
	}

	return nil
}
