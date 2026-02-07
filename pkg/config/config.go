package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SpecConfig represents the configspec.yaml structure
type SpecConfig struct {
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
			URL  string `yaml:"url"`
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
	Platforms []struct {
		Name          string   `yaml:"name"`
		Architectures []string `yaml:"architectures"`
	} `yaml:"platforms"`
	Environment []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"environment"`
	Tags []struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"tags"`
}

// LoadConfig reads and parses a configspec.yaml file
func LoadConfig(path string) (*SpecConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config SpecConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Validate checks if the config is valid
func (c *SpecConfig) Validate() error {
	if c.Info.Title == "" {
		return fmt.Errorf("info.title is required")
	}
	if c.Output.Directory == "" {
		return fmt.Errorf("output.directory is required")
	}
	if len(c.Output.Formats) == 0 {
		return fmt.Errorf("output.formats must contain at least one format")
	}
	if c.Output.Filename == "" {
		return fmt.Errorf("output.filename is required")
	}
	return nil
}
