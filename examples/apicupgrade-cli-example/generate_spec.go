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
	fmt.Println("Starting OpenCLI Spec Generation...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get the root command from your CLI
	rootCmd := cmd.GetRootCmd()

	// Create a new gospec converter
	gs := gospec.New()

	// Convert using config file
	// This reads configspec.yaml and generates the spec
	fmt.Println("\nReading configuration from configspec.yaml...")
	if err := gs.ConvertFromConfig("configspec.yaml", rootCmd); err != nil {
		log.Fatalf("❌ Error converting CLI to spec: %v", err)
	}

	fmt.Println("\n✅ OpenCLI Spec generation completed successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nGenerated files:")
	fmt.Println("   • output/sample-cli-spec.yaml")
	fmt.Println("   • output/sample-cli-spec.json")
	fmt.Println("\nUsage:")
	fmt.Println("   go run generate_spec.go")
	fmt.Println("\nLearn more:")
	fmt.Println("   https://github.com/harihs-330/gospec-cli")

	os.Exit(0)
}
