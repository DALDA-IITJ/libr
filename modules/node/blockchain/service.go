package blockchain

import (
	"encoding/json"
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

	data := blockchain.TxData{
		Type: "DB_JOINED",
		Metadata: map[string]string{
			"ip":   nodeConfig["node_ip"].(string),
			"port": nodeConfig["node_port"].(string),
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

	nonceStr := strconv.Itoa(txConfig["nonce"].(int))
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
