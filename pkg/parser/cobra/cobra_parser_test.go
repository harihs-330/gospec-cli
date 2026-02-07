package cobra

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCobraParser_Name(t *testing.T) {
	parser := NewCobraParser()
	if parser.Name() != "cobra" {
		t.Errorf("Expected parser name 'cobra', got '%s'", parser.Name())
	}
}

func TestCobraParser_Supports(t *testing.T) {
	parser := NewCobraParser()

	tests := []struct {
		name     string
		source   interface{}
		expected bool
	}{
		{
			name:     "Valid Cobra Command",
			source:   &cobra.Command{},
			expected: true,
		},
		{
			name:     "Invalid type - string",
			source:   "not a command",
			expected: false,
		},
		{
			name:     "Invalid type - nil",
			source:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Supports(tt.source)
			if result != tt.expected {
				t.Errorf("Expected Supports() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCobraParser_Parse(t *testing.T) {
	parser := NewCobraParser()

	// Create a sample Cobra command
	var verbose bool
	var config string

	rootCmd := &cobra.Command{
		Use:     "testapp",
		Short:   "A test application",
		Long:    "This is a test application for unit testing",
		Version: "1.0.0",
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Config file path")

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management",
	}

	var username string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a user",
	}
	createCmd.Flags().StringVarP(&username, "username", "u", "", "Username")

	userCmd.AddCommand(createCmd)
	rootCmd.AddCommand(userCmd)

	// Parse the command
	parsed, err := parser.Parse(rootCmd)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify parsed structure
	if parsed == nil {
		t.Fatal("Expected non-nil ParsedCLI")
	}

	if parsed.RootCommand == nil {
		t.Fatal("Expected non-nil RootCommand")
	}

	if parsed.RootCommand.Name != "testapp" {
		t.Errorf("Expected root command name 'testapp', got '%s'", parsed.RootCommand.Name)
	}

	if parsed.RootCommand.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", parsed.RootCommand.Version)
	}

	// Check persistent flags
	if len(parsed.RootCommand.PersistentFlags) != 2 {
		t.Errorf("Expected 2 persistent flags, got %d", len(parsed.RootCommand.PersistentFlags))
	}

	// Check subcommands
	if len(parsed.RootCommand.Subcommands) != 1 {
		t.Errorf("Expected 1 subcommand, got %d", len(parsed.RootCommand.Subcommands))
	}

	// Check commands map
	if len(parsed.Commands) < 3 {
		t.Errorf("Expected at least 3 commands in map, got %d", len(parsed.Commands))
	}
}

func TestCobraParser_ParseInvalidSource(t *testing.T) {
	parser := NewCobraParser()

	_, err := parser.Parse("invalid source")
	if err == nil {
		t.Error("Expected error for invalid source, got nil")
	}
}

func TestCobraParser_ParseFlags(t *testing.T) {
	parser := NewCobraParser()

	var (
		stringFlag   string
		intFlag      int
		boolFlag     bool
		sliceFlag    []string
		requiredFlag string
	)

	cmd := &cobra.Command{
		Use: "test",
	}

	cmd.Flags().StringVarP(&stringFlag, "string", "s", "default", "A string flag")
	cmd.Flags().IntVarP(&intFlag, "int", "i", 42, "An int flag")
	cmd.Flags().BoolVarP(&boolFlag, "bool", "b", false, "A bool flag")
	cmd.Flags().StringSliceVar(&sliceFlag, "slice", []string{}, "A slice flag")
	cmd.Flags().StringVar(&requiredFlag, "required", "", "A required flag")
	cmd.MarkFlagRequired("required")

	parsed, err := parser.Parse(cmd)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(parsed.RootCommand.Flags) != 5 {
		t.Errorf("Expected 5 flags, got %d", len(parsed.RootCommand.Flags))
	}

	// Verify flag types
	flagTypes := make(map[string]string)
	for _, flag := range parsed.RootCommand.Flags {
		flagTypes[flag.Name] = flag.Type
	}

	expectedTypes := map[string]string{
		"string": "string",
		"int":    "int",
		"bool":   "bool",
		"slice":  "stringSlice",
	}

	for name, expectedType := range expectedTypes {
		if flagTypes[name] != expectedType {
			t.Errorf("Flag '%s': expected type '%s', got '%s'", name, expectedType, flagTypes[name])
		}
	}
}

func TestCobraParser_ParseMetadata(t *testing.T) {
	parser := NewCobraParser()

	rootCmd := &cobra.Command{
		Use:     "myapp",
		Short:   "My application",
		Version: "2.0.0",
		Annotations: map[string]string{
			"author":     "Test Author",
			"license":    "MIT",
			"homepage":   "https://example.com",
			"repository": "https://github.com/user/repo",
		},
	}

	parsed, err := parser.Parse(rootCmd)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	metadata := parsed.Metadata
	if metadata.Name != "myapp" {
		t.Errorf("Expected name 'myapp', got '%s'", metadata.Name)
	}

	if metadata.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", metadata.Version)
	}

	if metadata.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", metadata.Author)
	}

	if metadata.License != "MIT" {
		t.Errorf("Expected license 'MIT', got '%s'", metadata.License)
	}
}

func TestCobraParser_ParseHiddenCommands(t *testing.T) {
	parser := NewCobraParser()

	rootCmd := &cobra.Command{
		Use: "root",
	}

	visibleCmd := &cobra.Command{
		Use:   "visible",
		Short: "A visible command",
	}

	hiddenCmd := &cobra.Command{
		Use:    "hidden",
		Short:  "A hidden command",
		Hidden: true,
	}

	rootCmd.AddCommand(visibleCmd)
	rootCmd.AddCommand(hiddenCmd)

	parsed, err := parser.Parse(rootCmd)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Both commands should be parsed
	if len(parsed.RootCommand.Subcommands) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(parsed.RootCommand.Subcommands))
	}

	// Check hidden flag
	for _, subcmd := range parsed.RootCommand.Subcommands {
		if subcmd.Name == "hidden" && !subcmd.Hidden {
			t.Error("Expected 'hidden' command to have Hidden=true")
		}
		if subcmd.Name == "visible" && subcmd.Hidden {
			t.Error("Expected 'visible' command to have Hidden=false")
		}
	}
}
