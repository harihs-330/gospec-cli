package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	port int
	host string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage server operations",
	Long:  `Start, stop, or manage the server instance.`,
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long:  `Start the server on the specified host and port.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting server on %s:%d\n", host, port)
	},
}

var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the server",
	Long:  `Stop the running server instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping server...")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverStartCmd)
	serverCmd.AddCommand(serverStopCmd)

	serverStartCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	serverStartCmd.Flags().StringVarP(&host, "host", "h", "localhost", "Host to bind to")
}
