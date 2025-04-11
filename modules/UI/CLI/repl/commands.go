package repl

import "fmt"

type CommandHandler func(args []string, timestamp int64) int64

func GetCommandHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		"send":      handleSendCommand,
		"fetch":     handleFetchCommand,
		"f":         handleFetchCommand,
		"prev":      handlePrevCommand,
		"p":         handlePrevCommand,
		"setup mod": handleModSetupCommand,
		"help":      handleHelpCommand,
	}
}

func handleSendCommand(args []string, timestamp int64) int64 {
	fmt.Println("Send command executed with args:", args)
	return timestamp
}

func handleFetchCommand(args []string, timestamp int64) int64 {
	fmt.Println("Fetch command executed.")
	return timestamp
}

func handlePrevCommand(args []string, timestamp int64) int64 {
	fmt.Println("Previous command executed.")
	return timestamp - 1
}

func handleHelpCommand(args []string, timestamp int64) int64 {
	fmt.Println("Available commands: send, fetch, prev, setup mod, help")
	return timestamp
}
