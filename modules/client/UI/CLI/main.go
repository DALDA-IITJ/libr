package main

import (
	"github.com/DALDA-IITJ/libr/modules/client/UI/CLI/cmd"
	"github.com/DALDA-IITJ/libr/modules/client/core"
)

func main() {

	core.InitCore()
	cmd.Execute()
}
