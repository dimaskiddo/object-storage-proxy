package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dimaskiddo/object-storage-proxy/internal/cmd"
)

// Root Variable Structure
var r = &cobra.Command{
	Use:   "object-storage-proxy",
	Short: "Object Storage Proxy",
	Long:  "Object Storage Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Init Function
func init() {
	// Add Child for Root Command
	r.AddCommand(cmd.Proxy)
	r.AddCommand(cmd.Version)
}

// Main Function
func main() {
	err := r.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
