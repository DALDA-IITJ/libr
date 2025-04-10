package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/DALDA-IITJ/libr/modules/client/core"
)

func main() {
	core.InitCore()
	repl()
}

func roundToHundred(ts int64) int64 {
	return ts - (ts % 100)
}

// Read-Eval-Print Loop
func repl() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("=== Welcome to Libr ===")
	fmt.Println("Type 'help' for commands. Type '\\q' to quit.")
	fmt.Println("===========================")

	currentTimestamp := roundToHundred(time.Now().Unix())

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Process input
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

		switch command {
		case "send":
			handleSendCommand(args)
			// currentTimestamp = time.Now().Unix() // Update timestamp
			currentTimestamp = roundToHundred(time.Now().Unix())
		case "fetch", "f":
			handleFetchCommand(currentTimestamp)
		case "prev", "p":
			currentTimestamp = handlePrevCommand(currentTimestamp)
		case "help":
			printHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for commands.")
		}
	}
}

func handleSendCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: send <message>")
		return
	}

	message := strings.Join(args, " ")
	core := core.NewCore()
	err := core.SendMessage(message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✔ Message sent successfully")
	}
}

func handleFetchCommand(timestamp int64) {
	core := core.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(timestamp))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return
	}

	fmt.Println("\n=== Messages After Timestamp ===")
	if len(messages) == 0 {
		fmt.Println("No newer messages found.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] [%s]=> %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	fmt.Println("==============================\n")
}

func handlePrevCommand(timestamp int64) int64 {
	core := core.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(timestamp - 100)) // Assuming FetchMessages handles older timestamps //100
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return timestamp
	}

	fmt.Println("\n=== Messages Before Timestamp ===")
	if len(messages) == 0 {
		fmt.Println("No older messages found.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"), msg.Content)
		}
		timestamp = timestamp - 100 // Move to the previous timestamp
	}
	fmt.Println("===============================\n")
	return timestamp
}

func printHelp() {
	fmt.Println("\n=== CLI Commands ===")
	fmt.Println("send <message>     - Send a message.")
	fmt.Println("f, fetch           - Fetch messages after the current timestamp.")
	fmt.Println("p, prev            - Fetch messages before the current timestamp.")
	fmt.Println("\\q                 - Quit the CLI.")
	fmt.Println("====================\n")
}
