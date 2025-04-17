package repl

type CommandHandler func(args []string, timestamp int64) int64

func GetCommandHandlers() map[string]CommandHandler {
	return map[string]CommandHandler{
		"send":  handleSendCommand,
		"fetch": handleFetchCommand,
		"f":     handleFetchCommand,
		"prev":  handlePrevCommand,
		"p":     handlePrevCommand,
		"setup": handleModSetupCommand,
		"help":  handleHelpCommand,
	}
}
