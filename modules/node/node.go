package main

import (
	"net/http"

	"github.com/DALDA-IITJ/libr/modules/node/api"
	"github.com/DALDA-IITJ/libr/modules/node/blockchain"
	"github.com/DALDA-IITJ/libr/modules/node/db"
	"github.com/DALDA-IITJ/libr/modules/node/utils"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
	"github.com/DALDA-IITJ/libr/modules/node/worker"
)

func main() {
	logger.Info("Starting node...")

	utils.LoadConfigAndEnv()

	db.InitDB()

	err := utils.EnsureKeyPair()
	if err != nil {
		logger.Fatal("Failed to ensure key pair: " + err.Error())
	}

	blockchain.RegisterNode()

	worker.StartWorkerPool(1)

	router := api.SetUpRoutes()
	logger.Info("Starting HTTP server on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal("HTTP server failed to start: " + err.Error())
	}

	logger.Info("Node started successfully")
}
