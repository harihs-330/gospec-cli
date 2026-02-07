package generator

import (
	"encoding/json"
	"io"

	"github.com/harihs-330/gospec-cli/pkg/spec"
)

// JSONGenerator generates JSON output for OpenCLI specifications
type JSONGenerator struct {
	indent string
	pretty bool
}

// NewJSONGenerator creates a new JSON generator
func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{
		indent: "  ",
		pretty: true,
	}
}

// Generate writes the OpenCLI spec as JSON to the writer
func (g *JSONGenerator) Generate(spec *spec.OpenCLISpec, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	if g.pretty {
		encoder.SetIndent("", g.indent)
	}
	return encoder.Encode(spec)
}

// GenerateToString returns the OpenCLI spec as a JSON string
func (g *JSONGenerator) GenerateToString(spec *spec.OpenCLISpec) (string, error) {
	var data []byte
	var err error

	if g.pretty {
		data, err = json.MarshalIndent(spec, "", g.indent)
	} else {
		data, err = json.Marshal(spec)
	}

	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetIndent sets the indentation string for JSON output
func (g *JSONGenerator) SetIndent(indent string) {
	g.indent = indent
}

// SetPretty enables or disables pretty printing
func (g *JSONGenerator) SetPretty(pretty bool) {
	g.pretty = pretty
}
