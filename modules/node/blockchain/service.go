package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DALDA-IITJ/libr/modules/core/blockchain"
	"github.com/DALDA-IITJ/libr/modules/core/crypto"
	"github.com/DALDA-IITJ/libr/modules/node/utils"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

var BlockchainStates map[string][]blockchain.Transaction

func GetBlockchainState(time string) error {
	// BlockchainStates[time] = []blockchain.Transaction{}

	// txs, err := blockchain.FetchBlockchain(time)
	// if err != nil {
	// 	logger.Error("Failed to fetch blockchain: " + err.Error())
	// 	return err
	// }

	// BlockchainStates[time] = txs
	// logger.Info("Blockchain state fetched successfully for time: " + string(time))
	return nil
}

func RegisterNode() {
	communityPublicKey := os.Getenv("COMMUNITY_PUBLIC_KEY")
	txConfig := utils.GetTransactionConfig()
	nodeConfig := utils.GetNodeConfig()

	if nodeConfig["node_ip"] == nil || nodeConfig["node_port"] == nil {
		logger.Error("node_ip or node_port is missing in configuration")
		return
	}
	if txConfig["nonce"] == nil {
		logger.Error("nonce is missing in transaction configuration")
		return
	}

	ipStr := fmt.Sprintf("%v", nodeConfig["node_ip"])
	portStr := fmt.Sprintf("%v", nodeConfig["node_port"])
	if ipStr == "" || portStr == "" {
		logger.Error("node_ip or node_port cannot be empty")
		return
	}

	data := blockchain.TxData{
		Type: "DB_JOINED",
		Metadata: map[string]string{
			"ip":   ipStr,
			"port": portStr,
			// "other": "",
		},
	}

	messageBytes, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal data: " + err.Error())
		return
	}

	dataStr := string(messageBytes)

	sign, err := crypto.SignMessage(utils.PrivateKey, dataStr)
	if err != nil {
		logger.Error("Failed to sign message: " + err.Error())
		return
	}

	nonceStr := fmt.Sprintf("%v", txConfig["nonce"])
	timeStampStr := strconv.Itoa(int(time.Now().Unix()))

	var tx = blockchain.Transaction{
		Sender:    utils.PublicKey,
		Recipient: communityPublicKey,
		Nonce:     nonceStr,
		Data:      data,
		Timestamp: timeStampStr,
		Sign:      sign,
	}

	err = SendTransaction(tx)
	if err != nil {
		logger.Error("Failed to send transaction: " + err.Error())
		return
	}

	logger.Info("DB JOIN Transaction sent successfully")
}

func SendTransaction(tx blockchain.Transaction) error {
	return nil
}
