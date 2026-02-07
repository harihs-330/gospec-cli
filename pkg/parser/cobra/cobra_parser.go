package cobra

import (
	"reflect"
	"strings"

	"github.com/harihs-330/gospec-cli/pkg/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CobraParser implements the Parser interface for Cobra CLI framework
type CobraParser struct{}

// NewCobraParser creates a new Cobra parser
func NewCobraParser() *CobraParser {
	return &CobraParser{}
}

// Name returns the parser name
func (p *CobraParser) Name() string {
	return "cobra"
}

// Supports checks if the source is a Cobra command
func (p *CobraParser) Supports(source interface{}) bool {
	_, ok := source.(*cobra.Command)
	return ok
}

// Parse extracts CLI structure from a Cobra command
func (p *CobraParser) Parse(source interface{}) (*parser.ParsedCLI, error) {
	cmd, ok := source.(*cobra.Command)
	if !ok {
		return nil, &parser.ParserError{
			Message: "source is not a *cobra.Command",
			Cause:   parser.ErrInvalidSource,
		}
	}

	parsed := &parser.ParsedCLI{
		Commands:      make(map[string]*parser.CommandInfo),
		FrameworkData: make(map[string]interface{}),
	}

	// Parse root command
	rootInfo := p.parseCommand(cmd, nil, "")
	parsed.RootCommand = rootInfo
	parsed.Commands[rootInfo.Path] = rootInfo

	// Parse all subcommands recursively
	p.parseSubcommands(cmd, rootInfo, parsed.Commands)

	// Extract metadata
	parsed.Metadata = p.extractMetadata(cmd)

	// Store framework-specific data
	parsed.FrameworkData["framework"] = "cobra"
	parsed.FrameworkData["version"] = getCobraVersion()

	return parsed, nil
}

// parseCommand converts a Cobra command to CommandInfo
func (p *CobraParser) parseCommand(cmd *cobra.Command, parent *parser.CommandInfo, parentPath string) *parser.CommandInfo {
	// Build command path
	path := parentPath
	if cmd.Name() != "" {
		if path == "" {
			path = cmd.Name()
		} else {
			path = parentPath + "/" + cmd.Name()
		}
	}

	info := &parser.CommandInfo{
		Name:            cmd.Name(),
		Path:            path,
		Use:             cmd.Use,
		Short:           cmd.Short,
		Long:            cmd.Long,
		Example:         cmd.Example,
		Aliases:         cmd.Aliases,
		Version:         cmd.Version,
		Parent:          parent,
		Subcommands:     make([]*parser.CommandInfo, 0),
		Flags:           make([]*parser.FlagInfo, 0),
		Args:            make([]*parser.ArgumentInfo, 0),
		PersistentFlags: make([]*parser.FlagInfo, 0),
		Hidden:          cmd.Hidden,
		Deprecated:      cmd.Deprecated,
		RunFunc:         cmd.Run != nil || cmd.RunE != nil,
		Annotations:     cmd.Annotations,
		Tags:            extractTags(cmd),
		Extensions:      make(map[string]interface{}),
	}

	// Parse local flags
	if cmd.Flags() != nil {
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if !flag.Changed && cmd.PersistentFlags() != nil {
				// Skip if it's a persistent flag (will be handled separately)
				if cmd.PersistentFlags().Lookup(flag.Name) != nil {
					return
				}
			}
			info.Flags = append(info.Flags, p.parseFlag(flag, false))
		})
	}

	// Parse persistent flags
	if cmd.PersistentFlags() != nil {
		cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			info.PersistentFlags = append(info.PersistentFlags, p.parseFlag(flag, true))
		})
	}

	// Parse arguments from ValidArgs and Args
	info.Args = p.parseArguments(cmd)

	// Store Cobra-specific extensions
	info.Extensions["cobra_use"] = cmd.Use
	info.Extensions["cobra_suggest_for"] = cmd.SuggestFor
	info.Extensions["cobra_args_validator"] = cmd.Args != nil
	if cmd.ValidArgs != nil {
		info.Extensions["cobra_valid_args"] = cmd.ValidArgs
	}

	return info
}

// parseFlag converts a pflag.Flag to FlagInfo
func (p *CobraParser) parseFlag(flag *pflag.Flag, persistent bool) *parser.FlagInfo {
	// Convert annotations from map[string][]string to map[string]string
	annotations := make(map[string]string)
	if flag.Annotations != nil {
		for key, values := range flag.Annotations {
			if len(values) > 0 {
				annotations[key] = values[0]
			}
		}
	}

	flagInfo := &parser.FlagInfo{
		Name:         flag.Name,
		Shorthand:    flag.Shorthand,
		Usage:        flag.Usage,
		Type:         flag.Value.Type(),
		DefaultValue: flag.DefValue,
		Required:     isRequiredFlag(flag),
		Hidden:       flag.Hidden,
		Deprecated:   flag.Deprecated,
		Persistent:   persistent,
		Annotations:  annotations,
	}

	// Extract valid values for enum-like flags
	if validValues := extractValidValues(flag); len(validValues) > 0 {
		flagInfo.ValidValues = validValues
	}

	return flagInfo
}

// parseArguments extracts argument information from Cobra command
func (p *CobraParser) parseArguments(cmd *cobra.Command) []*parser.ArgumentInfo {
	args := make([]*parser.ArgumentInfo, 0)

	// If ValidArgs is set, create arguments from it
	if len(cmd.ValidArgs) > 0 {
		for i, validArg := range cmd.ValidArgs {
			args = append(args, &parser.ArgumentInfo{
				Name:        validArg,
				Position:    i + 1,
				Required:    true,
				Type:        "string",
				ValidValues: []string{validArg},
			})
		}
		return args
	}

	// Infer from Args validator
	if cmd.Args != nil {
		minArgs, maxArgs := inferArgsArity(cmd)
		if minArgs > 0 || maxArgs > 0 {
			// Create generic argument info
			argName := "args"
			if cmd.Use != "" {
				// Try to extract argument name from Use string
				parts := strings.Fields(cmd.Use)
				if len(parts) > 1 {
					argName = parts[1]
				}
			}

			args = append(args, &parser.ArgumentInfo{
				Name:     argName,
				Position: 1,
				Required: minArgs > 0,
				Type:     "string",
				MinArgs:  minArgs,
				MaxArgs:  maxArgs,
			})
		}
	}

	return args
}

// parseSubcommands recursively parses all subcommands
func (p *CobraParser) parseSubcommands(cmd *cobra.Command, parentInfo *parser.CommandInfo, commandMap map[string]*parser.CommandInfo) {
	for _, subCmd := range cmd.Commands() {
		// Skip hidden commands if needed (can be controlled by options)
		subInfo := p.parseCommand(subCmd, parentInfo, parentInfo.Path)
		parentInfo.Subcommands = append(parentInfo.Subcommands, subInfo)
		commandMap[subInfo.Path] = subInfo

		// Recursively parse nested subcommands
		p.parseSubcommands(subCmd, subInfo, commandMap)
	}
}

// extractMetadata extracts global CLI metadata
func (p *CobraParser) extractMetadata(cmd *cobra.Command) *parser.CLIMetadata {
	// Find root command
	root := cmd
	for root.Parent() != nil {
		root = root.Parent()
	}

	metadata := &parser.CLIMetadata{
		Name:        root.Name(),
		Version:     root.Version,
		Description: root.Long,
		Tags:        make([]parser.TagInfo, 0),
		EnvVars:     make([]parser.EnvVarInfo, 0),
		Platforms:   make([]parser.PlatformInfo, 0),
	}

	// Extract from annotations if available
	if root.Annotations != nil {
		if author, ok := root.Annotations["author"]; ok {
			metadata.Author = author
		}
		if license, ok := root.Annotations["license"]; ok {
			metadata.License = license
		}
		if homepage, ok := root.Annotations["homepage"]; ok {
			metadata.Homepage = homepage
		}
		if repo, ok := root.Annotations["repository"]; ok {
			metadata.Repository = repo
		}
	}

	return metadata
}

// Helper functions

func extractTags(cmd *cobra.Command) []string {
	tags := make([]string, 0)

	if cmd.Annotations != nil {
		if tagStr, ok := cmd.Annotations["tags"]; ok {
			tags = strings.Split(tagStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}
	}

	return tags
}

func isRequiredFlag(flag *pflag.Flag) bool {
	// Check if flag has required annotation
	if flag.Annotations != nil {
		if _, ok := flag.Annotations["required"]; ok {
			return true
		}
	}
	return false
}

func extractValidValues(flag *pflag.Flag) []string {
	// Check annotations for valid values
	if flag.Annotations != nil {
		if values, ok := flag.Annotations["validValues"]; ok {
			if len(values) > 0 {
				return strings.Split(values[0], ",")
			}
		}
	}
	return nil
}

func inferArgsArity(cmd *cobra.Command) (min, max int) {
	// This is a best-effort inference based on common Cobra validators
	// In practice, you might need to use reflection or custom annotations

	if cmd.Args == nil {
		return 0, -1 // unlimited
	}

	// Try to infer from common validators
	argsType := reflect.TypeOf(cmd.Args)
	if argsType != nil {
		argsName := argsType.String()
		switch {
		case strings.Contains(argsName, "NoArgs"):
			return 0, 0
		case strings.Contains(argsName, "ExactArgs"):
			return 1, 1 // Default, would need more context
		case strings.Contains(argsName, "MinimumNArgs"):
			return 1, -1
		case strings.Contains(argsName, "MaximumNArgs"):
			return 0, 1
		case strings.Contains(argsName, "RangeArgs"):
			return 1, -1 // Default range
		}
	}

	return 0, -1 // Default: unlimited
}

func getCobraVersion() string {
	// This would ideally get the actual Cobra version
	// For now, return a placeholder
	return "1.8.0"
}
