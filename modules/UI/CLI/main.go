package main

import (
	"fmt"
	"os"

	"github.com/DALDA-IITJ/libr/modules/UI/config"
	"github.com/DALDA-IITJ/libr/modules/UI/repl"
	"github.com/DALDA-IITJ/libr/modules/client"
)

func main() {
	// Load configuration
	conf := config.GetKeys()

	// Set environment variables
	os.Setenv("PRIVATE_KEY", conf.PrivateKey)
	os.Setenv("PUBLIC_KEY", conf.PublicKey)

	// Print the keys for debugging (optional)
	fmt.Printf("Loaded keys: PRIVATE_KEY=%s, PUBLIC_KEY=%s\n", conf.PrivateKey, conf.PublicKey)

	// Print confirmation for debugging (optional)
	fmt.Println("Configuration loaded successfully.")

	// Initialize core and start REPL
	client.InitCore()
	repl.Start()
}
