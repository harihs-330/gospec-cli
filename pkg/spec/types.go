package spec

// OpenCLISpec represents the complete OpenCLI Specification structure
type OpenCLISpec struct {
	OpenCLI      string                `yaml:"opencli" json:"opencli"`
	Info         Info                  `yaml:"info" json:"info"`
	ExternalDocs *ExternalDocs         `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Platforms    []Platform            `yaml:"platforms,omitempty" json:"platforms,omitempty"`
	Environment  []EnvironmentVariable `yaml:"environment,omitempty" json:"environment,omitempty"`
	Tags         []Tag                 `yaml:"tags,omitempty" json:"tags,omitempty"`
	Commands     map[string]Command    `yaml:"commands" json:"commands"`
	Components   *Components           `yaml:"components,omitempty" json:"components,omitempty"`
}

// Info contains metadata about the CLI application
type Info struct {
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Version     string   `yaml:"version" json:"version"`
	Contact     *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License     *License `yaml:"license,omitempty" json:"license,omitempty"`
}

// Contact information for the CLI
type Contact struct {
	Name  string `yaml:"name,omitempty" json:"name,omitempty"`
	URL   string `yaml:"url,omitempty" json:"url,omitempty"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// License information
type License struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

// ExternalDocs points to external documentation
type ExternalDocs struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	URL         string `yaml:"url" json:"url"`
}

// Platform represents a supported platform
type Platform struct {
	Name          string   `yaml:"name" json:"name"`
	Architectures []string `yaml:"architectures,omitempty" json:"architectures,omitempty"`
}

// EnvironmentVariable represents an environment variable
type EnvironmentVariable struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty" json:"required,omitempty"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty"`
}

// Tag for grouping commands
type Tag struct {
	Name         string        `yaml:"name" json:"name"`
	Description  string        `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalDocs *ExternalDocs `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
}

// Command represents a CLI command
type Command struct {
	Summary     string                 `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	OperationID string                 `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Aliases     []string               `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Tags        []string               `yaml:"tags,omitempty" json:"tags,omitempty"`
	Parameters  []Parameter            `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses   map[string]Response    `yaml:"responses,omitempty" json:"responses,omitempty"`
	Deprecated  bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Hidden      bool                   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Extensions  map[string]interface{} `yaml:",inline" json:"-"`
}

// Parameter represents a command parameter/flag
type Parameter struct {
	Name        string                 `yaml:"name" json:"name"`
	In          string                 `yaml:"in,omitempty" json:"in,omitempty"` // argument, flag, option
	Alias       []string               `yaml:"alias,omitempty" json:"alias,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool                   `yaml:"required,omitempty" json:"required,omitempty"`
	Scope       string                 `yaml:"scope,omitempty" json:"scope,omitempty"` // local, inherited, global
	Position    int                    `yaml:"position,omitempty" json:"position,omitempty"`
	Schema      *Schema                `yaml:"schema,omitempty" json:"schema,omitempty"`
	Arity       *Arity                 `yaml:"arity,omitempty" json:"arity,omitempty"`
	Deprecated  bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Hidden      bool                   `yaml:"hidden,omitempty" json:"hidden,omitempty"`
	Extensions  map[string]interface{} `yaml:",inline" json:"-"`
}

// Arity defines the number of values a parameter can accept
type Arity struct {
	Min int  `yaml:"min,omitempty" json:"min,omitempty"`
	Max *int `yaml:"max,omitempty" json:"max,omitempty"` // nil means unlimited
}

// Schema defines the data type and validation rules
type Schema struct {
	Type       string             `yaml:"type,omitempty" json:"type,omitempty"`
	Format     string             `yaml:"format,omitempty" json:"format,omitempty"`
	Enum       []interface{}      `yaml:"enum,omitempty" json:"enum,omitempty"`
	Default    interface{}        `yaml:"default,omitempty" json:"default,omitempty"`
	Example    interface{}        `yaml:"example,omitempty" json:"example,omitempty"`
	Pattern    string             `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	MinLength  *int               `yaml:"minLength,omitempty" json:"minLength,omitempty"`
	MaxLength  *int               `yaml:"maxLength,omitempty" json:"maxLength,omitempty"`
	Minimum    *float64           `yaml:"minimum,omitempty" json:"minimum,omitempty"`
	Maximum    *float64           `yaml:"maximum,omitempty" json:"maximum,omitempty"`
	Items      *Schema            `yaml:"items,omitempty" json:"items,omitempty"`
	Properties map[string]*Schema `yaml:"properties,omitempty" json:"properties,omitempty"`
}

// Response represents a command response
type Response struct {
	Description string                 `yaml:"description" json:"description"`
	Content     map[string]MediaType   `yaml:"content,omitempty" json:"content,omitempty"`
	Extensions  map[string]interface{} `yaml:",inline" json:"-"`
}

// MediaType represents response content type
type MediaType struct {
	Schema  *Schema     `yaml:"schema,omitempty" json:"schema,omitempty"`
	Example interface{} `yaml:"example,omitempty" json:"example,omitempty"`
}

// Components holds reusable objects
type Components struct {
	Schemas    map[string]*Schema    `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Parameters map[string]*Parameter `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Responses  map[string]*Response  `yaml:"responses,omitempty" json:"responses,omitempty"`
}
