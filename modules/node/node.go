package main

import (
	logger "libr/core"
	"libr/modules/node/api"
	"libr/modules/node/blockchain"
	"libr/modules/node/db"
	"net/http"
)

func main() {
	logger.Info("Starting node...")

	db.InitDB()
	blockchain.RegisterNode()

	router := api.SetUpRoutes()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Node started successfully")
}