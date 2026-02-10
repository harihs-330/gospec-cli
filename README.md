# 🚀 gospec-cli

> Generate an **OpenCLI Specification** for your CLI using a simple configuration file.

`gospec-cli` allows you to convert your CLI into an OpenCLI specification using a configuration-driven approach.
Instead of passing many CLI flags, you define everything in a single `configspec.yaml` file.

---

## ⚡ Quick Start

### Step 1: Create Configuration File

In your project root, create a file named **`configspec.yaml`**:

```yaml
# OpenCLI Spec Configuration
# This file demonstrates how to configure gospec-cli to convert your CLI

# CLI Information - Metadata about your CLI application
info:
  title: "Sample CLI Application"
  description: "A demonstration CLI application showing gospec-cli conversion capabilities"
  version: "1.0.0"
  contact:
    name: "Sample CLI Team"
    url: "https://github.com/harihs-330/gospec-cli"
    email: "support@example.com"
  license:
    name: "Apache 2.0"
    url: "https://www.apache.org/licenses/LICENSE-2.0"

# Source CLI Configuration
source:
  # Type of source: "go-package", "binary", or "go-file"
  type: "go-package"
  
  # Local path to the CLI package
  localPath: "./cmd"
  
  # CLI framework (auto-detected if not specified)
  # Supported: "cobra", "urfave-cli", "flag"
  framework: "cobra"
  
  # Function name that returns the root command
  # Required for go-package type
  rootCommandFunc: "GetRootCmd"

# Output Configuration
output:
  # Directory where spec files will be generated
  directory: "./output"
  
  # Output formats: yaml, json, or both
  formats:
    - yaml
    # - json
  
  # Base filename (extensions added automatically)
  filename: "sample-cli-spec"
```

## Programmatic Usage (Library)

If you need to integrate into your Go application:

```go
//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/harihs-330/gospec-cli"
	"github.com/harihs-330/gospec-cli/examples/sample-cli/cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd()

	gs := gospec.New()
	if err := gs.ConvertFromConfig("configspec.yaml", rootCmd); err != nil {
		log.Fatalf("❌ Error converting CLI to spec: %v", err)
	}

	fmt.Println("\n✅ OpenCLI Spec generation completed successfully!")
	os.Exit(0)
}
```

### Step 2: Generate the Specification

Run the following command:

```bash
go run generate_spec.go
```

### Expected Output

```text
Starting OpenCLI Spec Generation...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Reading configuration from configspec.yaml...
✅ Generated: output/sample-cli-spec.yaml
✅ Generated: output/sample-cli-spec.json

✅ OpenCLI Spec generation completed successfully!
```

### Generated Files

The generated files will be placed in:

```
output/
 ├── sample-cli-spec.yaml
 └── sample-cli-spec.json
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding Support for New CLI Frameworks

1. **Implement** the `Parser` interface in `pkg/parser/`
2. **Register** your parser in the registry
3. **Add tests** to ensure reliability
4. **Update documentation** with examples

---

## 📄 License

Apache 2.0 - See [LICENSE](LICENSE) file for details

---

## 🙏 Acknowledgments

- 🎨 Inspired by **OpenAPI/Swagger** for REST APIs
- 📝 Built for the [**OpenCLI Specification**](https://openclispec.org) standard
- ⚡ Powered by [**Cobra**](https://github.com/spf13/cobra)
- 💙 Special thanks to the [**OpenCLI Specification**](https://www.openclispec.org) team for creating and maintaining this amazing standard

---

## 📞 Support

Need help? We're here for you!

- 📧 **Issues**: [GitHub Issues](https://github.com/harihs-330/gospec-cli/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/harihs-330/gospec-cli/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/harihs-330/gospec-cli/wiki)

