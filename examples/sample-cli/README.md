# Sample CLI - Complete Example

This is a complete, working example demonstrating how to use `gospec-cli` to convert a Cobra-based CLI application into OpenCLI Specification format.

## 📁 Directory Structure

```
sample-cli/
├── cmd/                      # CLI command definitions
│   ├── root.go              # Root command with global flags
│   └── server.go            # Server subcommands (start, stop)
├── output/                   # Generated OpenCLI specs
│   ├── sample-cli-spec.yaml # YAML format spec
│   └── sample-cli-spec.json # JSON format spec
├── configspec.yaml          # Configuration for gospec-cli
├── generate_spec.go         # Spec generator script
├── go.mod                   # Go module definition
└── README.md               # This file
```

## Quick Start

### 1. Generate OpenCLI Specification

Run the generator to convert your CLI to OpenCLI spec:

```bash
go run generate_spec.go
```

**Output:**
```
Starting OpenCLI Spec Generation...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Reading configuration from configspec.yaml...
✅ Generated: output/sample-cli-spec.yaml
✅ Generated: output/sample-cli-spec.json

✅ OpenCLI Spec generation completed successfully!
```

### 2. View Generated Specifications

**YAML Format:**
```bash
cat output/sample-cli-spec.yaml
```

**JSON Format:**
```bash
cat output/sample-cli-spec.json
```

## 📝 Configuration File

The `configspec.yaml` file controls how your CLI is converted to OpenCLI spec:

```yaml
# CLI Information
info:
  title: "Sample CLI Application"
  description: "A demonstration CLI application"
  version: "1.0.0"
  contact:
    name: "Sample CLI Team"
    url: "https://github.com/harihs-330/gospec-cli"
  license:
    name: "Apache 2.0"

# Source Configuration
source:
  type: "go-package"
  localPath: "./cmd"
  framework: "cobra"
  rootCommandFunc: "GetRootCmd"

# Output Configuration
output:
  directory: "./output"
  formats:
    - yaml
    - json
  filename: "sample-cli-spec"

# Generation Options
options:
  includeHidden: false
  generateOperationIds: true
  inferResponses: true
  tagStrategy: "auto"
  extractComponents: true
```

## 🔧 How It Works

### Step 1: Define Your CLI Commands

Create Cobra commands in the `cmd/` directory:

**cmd/root.go:**
```go
package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "sample-cli",
    Short: "A sample CLI application",
    Version: "1.0.0",
}

// GetRootCmd exposes root command for gospec-cli
func GetRootCmd() *cobra.Command {
    return rootCmd
}
```

**cmd/server.go:**
```go
package cmd

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "Manage server operations",
}

var serverStartCmd = &cobra.Command{
    Use:   "start",
    Short: "Start the server",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(serverCmd)
    serverCmd.AddCommand(serverStartCmd)
    
    serverStartCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
    serverStartCmd.Flags().StringP("host", "h", "localhost", "Host to bind to")
}
```

### Step 2: Create Configuration File

Create `configspec.yaml` with your CLI metadata and conversion settings.

### Step 3: Generate Specification

Create `generate_spec.go`:

```go
//go:build ignore

package main

import (
    "github.com/harihs-330/gospec-cli"
    "github.com/harihs-330/gospec-cli/examples/sample-cli/cmd"
)

func main() {
    rootCmd := cmd.GetRootCmd()
    gs := gospec.New()
    
    if err := gs.ConvertFromConfig("configspec.yaml", rootCmd); err != nil {
        log.Fatal(err)
    }
}
```

### Step 4: Run Generator

```bash
go run generate_spec.go
```

## 📊 Generated Output

The generator creates two files in the `output/` directory:

### YAML Format (`sample-cli-spec.yaml`)

```yaml
opencli: 1.0.0
info:
  title: Sample CLI Application
  version: 1.0.0
commands:
  sample-cli:
    summary: A sample CLI application
    operationId: sample-cliCommand
    parameters:
      - name: config
        in: flag
        scope: inherited
      - name: verbose
        in: flag
        alias: [v]
        scope: inherited
  /sample-cli/server:
    summary: Manage server operations
    operationId: sample-cliServerCommand
  /sample-cli/server/start:
    summary: Start the server
    operationId: sample-cliServerStartCommand
    parameters:
      - name: port
        in: flag
        alias: [p]
        schema:
          type: integer
          default: "8080"
      - name: host
        in: flag
        alias: [h]
        schema:
          type: string
          default: localhost
```

### JSON Format (`sample-cli-spec.json`)

Complete JSON representation with the same structure.

## 🎯 Key Features Demonstrated

1. **Root Command**: Main CLI entry point with global flags
2. **Subcommands**: Nested command structure (server → start/stop)
3. **Flags**: Both short (-p) and long (--port) flag formats
4. **Default Values**: Flag defaults are preserved
5. **Descriptions**: Command and flag descriptions
6. **Operation IDs**: Auto-generated unique identifiers
7. **Responses**: Inferred success/error responses
8. **Components**: Reusable response definitions

## 🔄 Customization

### Modify CLI Commands

Edit files in `cmd/` directory to add/modify commands.

### Update Configuration

Edit `configspec.yaml` to change:
- CLI metadata (title, version, contact)
- Output formats (YAML, JSON)
- Generation options (hidden commands, operation IDs)
- Platform information

### Regenerate Specs

After any changes, run:
```bash
go run generate_spec.go
```

## 📚 Learn More

- **gospec-cli Documentation**: [README.md](../../README.md)
- **OpenCLI Specification**: [openclispec.org](https://www.openclispec.org)
- **Cobra Framework**: [cobra.dev](https://cobra.dev)

## 💡 Tips

1. **Expose Root Command**: Always provide a `GetRootCmd()` function
2. **Use Descriptive Names**: Clear command and flag descriptions improve the spec
3. **Version Your CLI**: Include version in root command
4. **Test Generation**: Run generator after each CLI change
5. **Review Output**: Check generated specs for accuracy

## 🤝 Contributing

This example is part of the gospec-cli project. Contributions welcome!

## 📄 License

Apache 2.0 - See [LICENSE](../../LICENSE) for details