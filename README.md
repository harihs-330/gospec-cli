# gospec-cli

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**gospec-cli** is an enterprise-level Go package that automatically converts CLI applications (built with frameworks like Cobra) into [OpenCLI Specification](https://openclispec.org) format - similar to how OpenAPI/Swagger works for REST APIs.

## 🚀 Quick Start

### Installation

```bash
go install github.com/harihs-330/gospec-cli/cmd/gospec-cli@latest
```

### Simple Usage - No Code Required!

Just point the tool at your CLI binary or source code:

```bash
# Convert any Cobra-based CLI to OpenCLI spec
gospec-cli convert /path/to/your/cli-binary -o specs/

# Or from source code
gospec-cli convert /path/to/your/cli/cmd/root.go -o specs/

# Customize the output
gospec-cli convert ./my-cli \
  --title "My Awesome CLI" \
  --description "A powerful command-line tool" \
  --version "1.0.0" \
  --format yaml \
  -o my-cli-spec.yaml
```

That's it! No Go code to write, no dependencies to manage.

## 📋 Features

- ✅ **Zero Configuration**: Works out of the box with Cobra CLIs
- ✅ **Multiple Output Formats**: YAML and JSON
- ✅ **Complete Spec Generation**: Commands, flags, arguments, descriptions
- ✅ **Extensible Architecture**: Easy to add support for other CLI frameworks
- ✅ **Enterprise Ready**: Comprehensive error handling and validation
- ✅ **Well Tested**: Full test coverage

## 🎯 Use Cases

- **API Documentation**: Auto-generate CLI documentation
- **Client Generation**: Generate CLI clients in other languages
- **Testing**: Validate CLI structure and behavior
- **CI/CD Integration**: Ensure CLI consistency across versions
- **Developer Onboarding**: Provide clear CLI structure overview

## 📖 Usage Examples

### Basic Conversion

```bash
# Convert and save to default location (./specs/)
gospec-cli convert ./my-cli

# Specify output file
gospec-cli convert ./my-cli -o my-spec.yaml

# Generate JSON instead of YAML
gospec-cli convert ./my-cli --format json -o my-spec.json
```

### With Custom Metadata

```bash
gospec-cli convert ./my-cli \
  --title "My CLI Tool" \
  --description "A comprehensive CLI application" \
  --version "2.0.0" \
  --contact-name "Dev Team" \
  --contact-url "https://github.com/myorg/my-cli" \
  --license "MIT" \
  -o specs/my-cli.yaml
```

### Programmatic Usage (Library)

If you need to integrate into your Go application:

```go
package main

import (
    "os"
    "github.com/harihs-330/gospec-cli"
    "github.com/spf13/cobra"
)

func main() {
    // Your existing Cobra root command
    rootCmd := &cobra.Command{
        Use:   "my-cli",
        Short: "My awesome CLI",
    }
    
    // Convert to OpenCLI spec
    gs := gospec.New()
    options := gospec.DefaultOptions()
    
    // Generate YAML
    file, _ := os.Create("opencli-spec.yaml")
    defer file.Close()
    gs.ConvertToYAML(rootCmd, options, file)
}
```

## 🏗️ Architecture

```
gospec-cli/
├── cmd/gospec-cli/          # CLI tool
├── pkg/
│   ├── spec/                # OpenCLI spec data models
│   ├── parser/              # Parser interface & implementations
│   │   └── cobra/           # Cobra framework parser
│   ├── converter/           # Converts parsed CLI to OpenCLI spec
│   └── generator/           # YAML/JSON generators
└── gospec.go                # Main library API
```

## 🔧 Configuration

### CLI Flags

```
gospec-cli convert [source] [flags]

Flags:
  -o, --output string              Output file path (default "./specs/opencli-spec.yaml")
  -f, --format string              Output format: yaml or json (default "yaml")
      --title string               CLI title
      --description string         CLI description
      --version string             CLI version
      --contact-name string        Contact name
      --contact-url string         Contact URL
      --license string             License name
      --include-hidden             Include hidden commands
      --include-deprecated         Include deprecated commands
  -h, --help                       Help for convert
```

## 📦 OpenCLI Specification Output

The generated specification follows the [OpenCLI Specification 1.0.0](https://openclispec.org) format:

```yaml
opencli: 1.0.0
info:
  title: My CLI Tool
  description: A powerful command-line application
  version: 1.0.0
  contact:
    name: Dev Team
    url: https://github.com/myorg/my-cli
  license:
    name: MIT

commands:
  /my-cli:
    summary: Root command
    description: Main entry point
    parameters:
      - name: config
        in: flag
        description: Config file path
        schema:
          type: string
    responses:
      "0":
        description: Success
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding Support for New CLI Frameworks

1. Implement the `Parser` interface in `pkg/parser/`
2. Register your parser in the registry
3. Add tests
4. Update documentation

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

- Inspired by OpenAPI/Swagger for REST APIs
- Built for the [OpenCLI Specification](https://openclispec.org) standard
- Powered by [Cobra](https://github.com/spf13/cobra)

## 📞 Support

- 📧 Issues: [GitHub Issues](https://github.com/harihs-330/gospec-cli/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/harihs-330/gospec-cli/discussions)
- 📖 Documentation: [Wiki](https://github.com/harihs-330/gospec-cli/wiki)