package main

import (
	"fmt"
	"os"

	"github.com/DALDA-IITJ/libr/modules/client/UI/CLI/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "libr-cli",
	Short: "CLI for interacting with Libr",
}

func main() {
	// Add all commands here (including addCmd)
	rootCmd.AddCommand(cmd.AddCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
