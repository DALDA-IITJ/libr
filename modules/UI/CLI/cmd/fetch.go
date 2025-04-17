package cmd

import (
	"fmt"

	"github.com/DALDA-IITJ/libr/modules/client"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch messages by timestamp",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: cli fetch <timestamp>")
			return
		}

		timestamp := args[0]
		core := client.NewCore()
		messages, err := core.FetchMessages(timestamp)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			for _, msg := range messages {
				fmt.Printf("%s\n", msg.Content)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
