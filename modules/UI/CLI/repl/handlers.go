package repl

import (
	"fmt"
	"time"

	"github.com/DALDA-IITJ/libr/modules/client/core"
)

func handleSendCommand(args []string, timestamp int64) int64 {
	if len(args) < 1 {
		fmt.Println("Usage: send <message>")
		return timestamp
	}

	message := args[0]
	core := core.NewCore()
	err := core.SendMessage(message)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("✔ Message sent successfully")
	}
	return divideByHundred(time.Now().Unix())
}

func handleFetchCommand(args []string, timestamp int64) int64 {
	core := core.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(timestamp))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return timestamp
	}

	fmt.Println("\n=== Messages After Timestamp ===")
	if len(messages) == 0 {
		fmt.Println("No newer messages found.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] [%s] => %s\n", time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	fmt.Println("==============================\n")
	return timestamp
}

func handlePrevCommand(args []string, currentBucket int64) int64 {
	prevBucket := currentBucket - 1
	core := core.NewCore()
	messages, err := core.FetchMessages(fmt.Sprint(prevBucket))
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		return currentBucket // Don't update timestamp on error
	}

	fmt.Println("\n=== Messages in Previous Bucket ===")
	if len(messages) == 0 {
		fmt.Println("No messages found in previous bucket.")
	} else {
		for _, msg := range messages {
			fmt.Printf("[%s] [%s] => %s\n", time.Unix(currentBucket*100, 0).Format("2006-01-02 15:04:05"), msg.Sender, msg.Content)
		}
	}
	fmt.Println("===============================\n")
	return prevBucket
}

func handleHelpCommand(args []string, timestamp int64) int64 {
	fmt.Println("\n=== CLI Commands ===")
	fmt.Println("send <message>     - Send a message.")
	fmt.Println("f, fetch           - Fetch messages after the current timestamp.")
	fmt.Println("p, prev            - Fetch messages before the current timestamp.")
	fmt.Println("\\q                 - Quit the CLI.")
	fmt.Println("====================\n")
	return timestamp
}
