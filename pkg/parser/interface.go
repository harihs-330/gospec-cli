package parser

import (
	"github.com/harihs-330/gospec-cli/pkg/spec"
)

// Parser is the interface that all CLI framework parsers must implement
type Parser interface {
	// Name returns the name of the parser (e.g., "cobra", "urfave-cli", "flag")
	Name() string

	// Parse analyzes the CLI application and extracts its structure
	Parse(source interface{}) (*ParsedCLI, error)

	// Supports checks if this parser can handle the given source
	Supports(source interface{}) bool
}

// ParsedCLI represents the extracted CLI structure before conversion to OpenCLI spec
type ParsedCLI struct {
	// Root command information
	RootCommand *CommandInfo

	// All commands in the CLI (including root)
	Commands map[string]*CommandInfo

	// Global metadata
	Metadata *CLIMetadata

	// Framework-specific data
	FrameworkData map[string]interface{}
}

// CommandInfo contains all information about a command
type CommandInfo struct {
	// Basic information
	Name    string
	Path    string // Full path like "/user/create"
	Use     string // Usage string
	Short   string
	Long    string
	Example string
	Aliases []string
	Version string

	// Command hierarchy
	Parent      *CommandInfo
	Subcommands []*CommandInfo

	// Parameters
	Flags           []*FlagInfo
	Args            []*ArgumentInfo
	PersistentFlags []*FlagInfo // Flags inherited by subcommands

	// Behavior
	Hidden     bool
	Deprecated string
	RunFunc    bool // Whether command has a run function

	// Annotations and metadata
	Annotations map[string]string
	Tags        []string

	// Framework-specific data
	Extensions map[string]interface{}
}

// FlagInfo represents a command flag/option
type FlagInfo struct {
	Name         string
	Shorthand    string
	Usage        string
	Type         string // string, bool, int, float, duration, etc.
	DefaultValue interface{}
	Required     bool
	Hidden       bool
	Deprecated   string
	Persistent   bool // Whether flag is inherited by subcommands

	// Validation
	ValidValues []string // For enum-like flags

	// Extensions
	Annotations map[string]string
}

// ArgumentInfo represents a positional argument
type ArgumentInfo struct {
	Name        string
	Description string
	Position    int
	Required    bool
	Type        string

	// Arity
	MinArgs int
	MaxArgs int // -1 for unlimited

	// Validation
	ValidValues []string
}

// CLIMetadata contains global CLI information
type CLIMetadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	License     string
	Homepage    string
	Repository  string

	// Environment variables
	EnvVars []EnvVarInfo

	// Supported platforms
	Platforms []PlatformInfo

	// Tags for categorization
	Tags []TagInfo
}

// EnvVarInfo represents an environment variable
type EnvVarInfo struct {
	Name        string
	Description string
	Required    bool
	Default     string
}

// PlatformInfo represents a supported platform
type PlatformInfo struct {
	OS            string
	Architectures []string
}

// TagInfo represents a tag for categorization
type TagInfo struct {
	Name        string
	Description string
}

// ParserRegistry manages available parsers
type ParserRegistry struct {
	parsers map[string]Parser
}

// NewParserRegistry creates a new parser registry
func NewParserRegistry() *ParserRegistry {
	return &ParserRegistry{
		parsers: make(map[string]Parser),
	}
}

// Register adds a parser to the registry
func (r *ParserRegistry) Register(parser Parser) {
	r.parsers[parser.Name()] = parser
}

// Get retrieves a parser by name
func (r *ParserRegistry) Get(name string) (Parser, bool) {
	parser, ok := r.parsers[name]
	return parser, ok
}

// FindParser finds a suitable parser for the given source
func (r *ParserRegistry) FindParser(source interface{}) (Parser, error) {
	for _, parser := range r.parsers {
		if parser.Supports(source) {
			return parser, nil
		}
	}
	return nil, ErrNoSuitableParser
}

// List returns all registered parser names
func (r *ParserRegistry) List() []string {
	names := make([]string, 0, len(r.parsers))
	for name := range r.parsers {
		names = append(names, name)
	}
	return names
}

// Converter converts ParsedCLI to OpenCLI Specification
type Converter interface {
	Convert(parsed *ParsedCLI, options *ConvertOptions) (*spec.OpenCLISpec, error)
}

// ConvertOptions provides options for conversion
type ConvertOptions struct {
	// Spec version
	SpecVersion string

	// Include hidden commands/flags
	IncludeHidden bool

	// Include deprecated items
	IncludeDeprecated bool

	// Generate operation IDs
	GenerateOperationIDs bool

	// Infer response schemas
	InferResponses bool

	// Custom metadata
	CustomInfo *spec.Info

	// Tag generation strategy
	TagStrategy string // "auto", "manual", "none"

	// Component extraction
	ExtractComponents bool
}

// Error types
var (
	ErrNoSuitableParser = &ParserError{Message: "no suitable parser found for source"}
	ErrInvalidSource    = &ParserError{Message: "invalid source provided"}
	ErrParsingFailed    = &ParserError{Message: "parsing failed"}
)

// ParserError represents a parser error
type ParserError struct {
	Message string
	Cause   error
}

func (e *ParserError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *ParserError) Unwrap() error {
	return e.Cause
}
