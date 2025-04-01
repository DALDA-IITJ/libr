package cmd

import (
	"fmt"
	"strconv"

	"github.com/DALDA-IITJ/libr/modules/client/core/calc"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add [num1] [num2]",
	Short: "Adds two numbers",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		num1, _ := strconv.Atoi(args[0])
		num2, _ := strconv.Atoi(args[1])
		result := calc.Add(num1, num2)
		fmt.Println("Result:", result)
	},
}
