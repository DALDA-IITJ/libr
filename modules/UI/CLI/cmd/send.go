package cmd

import (
	"fmt"

	"github.com/DALDA-IITJ/libr/modules/client/core"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: cli send <message>")
			return
		}

		message := args[0]
		core := core.NewCore()
		err := core.SendMessage(message)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Message sent successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
