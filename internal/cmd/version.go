package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version Variable Structure
var Version = &cobra.Command{
	Use:   "version",
	Short: "Show Object Storage Proxy Version",
	Long:  "Show Object Storage Proxy Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Object Storage Proxy v1.0.0")
		fmt.Println("Initial Source : https://github.com/Kriechi")
		fmt.Println("Refactorer     : https://github.com/dimaskiddo (Dimas Restu Hidayanto)")
	},
}
