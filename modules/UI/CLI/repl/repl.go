package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func Start() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("=== Welcome to Libr ===")
	fmt.Println("Type 'help' for commands. Type '\\q' to quit.")
	fmt.Println("===========================")

	currentTimestamp := divideByHundred(time.Now().Unix())

	commands := GetCommandHandlers()

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "\\q" {
			fmt.Println("Exiting Libr. Goodbye!")
			break
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		if handler, exists := commands[command]; exists {
			currentTimestamp = handler(args, currentTimestamp)
		} else {
			fmt.Println("Unknown command. Type 'help' for commands.")
		}
	}
}

func divideByHundred(ts int64) int64 {
	return ts / 100
}
