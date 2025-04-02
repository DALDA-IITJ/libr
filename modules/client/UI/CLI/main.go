package main

import (
	"client/UI/CLI/cmd"
	"client/core"
)

func main() {

	core.InitCore()
	cmd.Execute()
}
