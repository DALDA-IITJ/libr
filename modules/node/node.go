package main

import (
	"net/http"
	"node/api"
	"node/blockchain"
	"node/db"
	"node/utils"
	"node/utils/logger"
)

func main() {
	logger.Info("Starting node...")

	utils.LoadConfig()
	db.InitDB()
	err := utils.EnsureKeyPair()
	if err != nil {
		logger.Fatal("Failed to ensure key pair: " + err.Error())
	}

	blockchain.RegisterNode()

	router := api.SetUpRoutes()
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("Node started successfully")
}
