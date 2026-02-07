package converter

import (
	"fmt"
	"strings"

	"github.com/harihs-330/gospec-cli/pkg/parser"
	"github.com/harihs-330/gospec-cli/pkg/spec"
)

// DefaultConverter implements the Converter interface
type DefaultConverter struct{}

// NewDefaultConverter creates a new default converter
func NewDefaultConverter() *DefaultConverter {
	return &DefaultConverter{}
}

// Convert transforms ParsedCLI into OpenCLI Specification
func (c *DefaultConverter) Convert(parsed *parser.ParsedCLI, options *parser.ConvertOptions) (*spec.OpenCLISpec, error) {
	if parsed == nil {
		return nil, fmt.Errorf("parsed CLI is nil")
	}

	if options == nil {
		options = DefaultConvertOptions()
	}

	openCLI := &spec.OpenCLISpec{
		OpenCLI:  options.SpecVersion,
		Commands: make(map[string]spec.Command),
	}

	// Set info
	openCLI.Info = c.convertInfo(parsed.Metadata, options)

	// Convert tags
	if len(parsed.Metadata.Tags) > 0 {
		openCLI.Tags = c.convertTags(parsed.Metadata.Tags)
	}

	// Convert environment variables
	if len(parsed.Metadata.EnvVars) > 0 {
		openCLI.Environment = c.convertEnvironment(parsed.Metadata.EnvVars)
	}

	// Convert platforms
	if len(parsed.Metadata.Platforms) > 0 {
		openCLI.Platforms = c.convertPlatforms(parsed.Metadata.Platforms)
	}

	// Convert commands
	for path, cmdInfo := range parsed.Commands {
		if !options.IncludeHidden && cmdInfo.Hidden {
			continue
		}
		if !options.IncludeDeprecated && cmdInfo.Deprecated != "" {
			continue
		}

		command := c.convertCommand(cmdInfo, options)

		// Use path as key, ensuring it starts with /
		key := path
		if !strings.HasPrefix(key, "/") && key != cmdInfo.Name {
			key = "/" + key
		}
		// For root command, use the command name directly
		if path == cmdInfo.Name {
			key = cmdInfo.Name
		}

		openCLI.Commands[key] = command
	}

	// Extract components if requested
	if options.ExtractComponents {
		openCLI.Components = c.extractComponents(parsed, options)
	}

	return openCLI, nil
}

// convertInfo creates Info from metadata
func (c *DefaultConverter) convertInfo(metadata *parser.CLIMetadata, options *parser.ConvertOptions) spec.Info {
	if options.CustomInfo != nil {
		return *options.CustomInfo
	}

	info := spec.Info{
		Title:       metadata.Name,
		Description: metadata.Description,
		Version:     metadata.Version,
	}

	if metadata.Author != "" || metadata.Homepage != "" {
		info.Contact = &spec.Contact{
			Name: metadata.Author,
			URL:  metadata.Homepage,
		}
	}

	if metadata.License != "" {
		info.License = &spec.License{
			Name: metadata.License,
		}
	}

	return info
}

// convertTags converts tag info to spec tags
func (c *DefaultConverter) convertTags(tags []parser.TagInfo) []spec.Tag {
	result := make([]spec.Tag, len(tags))
	for i, tag := range tags {
		result[i] = spec.Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}
	}
	return result
}

// convertEnvironment converts environment variables
func (c *DefaultConverter) convertEnvironment(envVars []parser.EnvVarInfo) []spec.EnvironmentVariable {
	result := make([]spec.EnvironmentVariable, len(envVars))
	for i, env := range envVars {
		result[i] = spec.EnvironmentVariable{
			Name:        env.Name,
			Description: env.Description,
			Required:    env.Required,
			Default:     env.Default,
		}
	}
	return result
}

// convertPlatforms converts platform info
func (c *DefaultConverter) convertPlatforms(platforms []parser.PlatformInfo) []spec.Platform {
	result := make([]spec.Platform, len(platforms))
	for i, platform := range platforms {
		result[i] = spec.Platform{
			Name:          platform.OS,
			Architectures: platform.Architectures,
		}
	}
	return result
}

// convertCommand converts CommandInfo to Command
func (c *DefaultConverter) convertCommand(cmdInfo *parser.CommandInfo, options *parser.ConvertOptions) spec.Command {
	command := spec.Command{
		Summary:     cmdInfo.Short,
		Description: cmdInfo.Long,
		Aliases:     cmdInfo.Aliases,
		Tags:        cmdInfo.Tags,
		Deprecated:  cmdInfo.Deprecated != "",
		Hidden:      cmdInfo.Hidden,
		Parameters:  make([]spec.Parameter, 0),
		Responses:   make(map[string]spec.Response),
		Extensions:  make(map[string]interface{}),
	}

	// Generate operation ID if requested
	if options.GenerateOperationIDs {
		command.OperationID = generateOperationID(cmdInfo)
	}

	// Convert flags to parameters
	for _, flag := range cmdInfo.Flags {
		if !options.IncludeHidden && flag.Hidden {
			continue
		}
		param := c.convertFlag(flag, "local")
		command.Parameters = append(command.Parameters, param)
	}

	// Convert persistent flags
	for _, flag := range cmdInfo.PersistentFlags {
		if !options.IncludeHidden && flag.Hidden {
			continue
		}
		param := c.convertFlag(flag, "inherited")
		command.Parameters = append(command.Parameters, param)
	}

	// Convert arguments
	for _, arg := range cmdInfo.Args {
		param := c.convertArgument(arg)
		command.Parameters = append(command.Parameters, param)
	}

	// Add default responses if requested
	if options.InferResponses {
		command.Responses = c.generateDefaultResponses()
	}

	// Copy extensions
	for key, value := range cmdInfo.Extensions {
		command.Extensions[key] = value
	}

	return command
}

// convertFlag converts FlagInfo to Parameter
func (c *DefaultConverter) convertFlag(flag *parser.FlagInfo, scope string) spec.Parameter {
	param := spec.Parameter{
		Name:        flag.Name,
		In:          "flag",
		Description: flag.Usage,
		Required:    flag.Required,
		Scope:       scope,
		Deprecated:  flag.Deprecated != "",
		Hidden:      flag.Hidden,
		Schema:      c.createSchema(flag.Type, flag.DefaultValue, flag.ValidValues),
	}

	// Add aliases
	if flag.Shorthand != "" {
		param.Alias = []string{flag.Shorthand}
	}

	return param
}

// convertArgument converts ArgumentInfo to Parameter
func (c *DefaultConverter) convertArgument(arg *parser.ArgumentInfo) spec.Parameter {
	param := spec.Parameter{
		Name:        arg.Name,
		In:          "argument",
		Description: arg.Description,
		Required:    arg.Required,
		Scope:       "local",
		Position:    arg.Position,
		Schema:      c.createSchema(arg.Type, nil, arg.ValidValues),
	}

	// Set arity if specified
	if arg.MinArgs > 0 || arg.MaxArgs != 0 {
		param.Arity = &spec.Arity{
			Min: arg.MinArgs,
		}
		if arg.MaxArgs > 0 {
			param.Arity.Max = &arg.MaxArgs
		}
	}

	return param
}

// createSchema creates a Schema from type information
func (c *DefaultConverter) createSchema(typeName string, defaultValue interface{}, validValues []string) *spec.Schema {
	schema := &spec.Schema{
		Type:    mapTypeToSchemaType(typeName),
		Default: defaultValue,
	}

	// Add enum if valid values are specified
	if len(validValues) > 0 {
		schema.Enum = make([]interface{}, len(validValues))
		for i, v := range validValues {
			schema.Enum[i] = v
		}
	}

	return schema
}

// generateDefaultResponses creates default response definitions
func (c *DefaultConverter) generateDefaultResponses() map[string]spec.Response {
	return map[string]spec.Response{
		"0": {
			Description: "Command executed successfully",
			Content: map[string]spec.MediaType{
				"text/plain": {
					Example: "Operation completed successfully",
				},
			},
		},
		"1": {
			Description: "Command execution failed",
			Content: map[string]spec.MediaType{
				"text/plain": {
					Example: "Error: operation failed",
				},
			},
		},
	}
}

// extractComponents extracts reusable components
func (c *DefaultConverter) extractComponents(parsed *parser.ParsedCLI, options *parser.ConvertOptions) *spec.Components {
	components := &spec.Components{
		Schemas:    make(map[string]*spec.Schema),
		Parameters: make(map[string]*spec.Parameter),
		Responses:  make(map[string]*spec.Response),
	}

	// Extract common parameters
	paramCounts := make(map[string]int)
	paramMap := make(map[string]*parser.FlagInfo)

	for _, cmdInfo := range parsed.Commands {
		for _, flag := range cmdInfo.Flags {
			key := flag.Name
			paramCounts[key]++
			if paramCounts[key] == 1 {
				paramMap[key] = flag
			}
		}
	}

	// Add parameters that appear in multiple commands
	for name, count := range paramCounts {
		if count > 1 {
			flag := paramMap[name]
			param := c.convertFlag(flag, "local")
			components.Parameters[name] = &param
		}
	}

	// Add common responses
	components.Responses["Success"] = &spec.Response{
		Description: "Operation completed successfully",
		Content: map[string]spec.MediaType{
			"text/plain": {
				Example: "Success",
			},
		},
	}

	components.Responses["Error"] = &spec.Response{
		Description: "Operation failed",
		Content: map[string]spec.MediaType{
			"text/plain": {
				Example: "Error: operation failed",
			},
		},
	}

	return components
}

// Helper functions

func generateOperationID(cmdInfo *parser.CommandInfo) string {
	// Convert path to camelCase operation ID
	parts := strings.Split(strings.Trim(cmdInfo.Path, "/"), "/")
	if len(parts) == 0 {
		return "rootCommand"
	}

	operationID := parts[0]
	for i := 1; i < len(parts); i++ {
		operationID += strings.Title(parts[i])
	}
	operationID += "Command"

	return operationID
}

func mapTypeToSchemaType(typeName string) string {
	switch strings.ToLower(typeName) {
	case "bool", "boolean":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "float", "float32", "float64":
		return "number"
	case "string":
		return "string"
	case "duration":
		return "string"
	case "stringslice", "[]string":
		return "array"
	default:
		return "string"
	}
}

// DefaultConvertOptions returns default conversion options
func DefaultConvertOptions() *parser.ConvertOptions {
	return &parser.ConvertOptions{
		SpecVersion:          "1.0.0",
		IncludeHidden:        false,
		IncludeDeprecated:    true,
		GenerateOperationIDs: true,
		InferResponses:       true,
		TagStrategy:          "auto",
		ExtractComponents:    true,
	}
}
