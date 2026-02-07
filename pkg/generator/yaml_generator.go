package generator

import (
	"io"

	"github.com/harihs-330/gospec-cli/pkg/spec"
	"gopkg.in/yaml.v3"
)

// YAMLGenerator generates YAML output for OpenCLI specifications
type YAMLGenerator struct {
	indent int
}

// NewYAMLGenerator creates a new YAML generator
func NewYAMLGenerator() *YAMLGenerator {
	return &YAMLGenerator{
		indent: 2,
	}
}

// Generate writes the OpenCLI spec as YAML to the writer
func (g *YAMLGenerator) Generate(spec *spec.OpenCLISpec, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(g.indent)
	defer encoder.Close()

	return encoder.Encode(spec)
}

// GenerateToString returns the OpenCLI spec as a YAML string
func (g *YAMLGenerator) GenerateToString(spec *spec.OpenCLISpec) (string, error) {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetIndent sets the indentation level for YAML output
func (g *YAMLGenerator) SetIndent(indent int) {
	g.indent = indent
}
