package blockchain

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

	nodeIPValue, nodeIPPresent := nodeConfig["node_ip"]
	nodePortValue, nodePortPresent := nodeConfig["node_port"]
	nonceValue, noncePresent := txConfig["nonce"]

	if !nodeIPPresent || !nodePortPresent {
		logger.Error("node_ip or node_port is missing in configuration")
		return
	}
	if !noncePresent {
		logger.Error("nonce is missing in transaction configuration")
		return
	}

	ipStr, ok := nodeIPValue.(string)
	if !ok {
		logger.Error("node_ip must be a non-empty string")
		return
	}
	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		logger.Error("node_ip must be a non-empty string")
		return
	}

	portStr, err := normalizeIntegerField(nodePortValue, "node_port", false)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	nonceStr, err := normalizeIntegerField(nonceValue, "nonce", true)
	if err != nil {
		logger.Error(err.Error())
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

func normalizeIntegerField(value interface{}, fieldName string, allowZero bool) (string, error) {
	var number int64

	switch v := value.(type) {
	case int:
		number = int64(v)
	case int8:
		number = int64(v)
	case int16:
		number = int64(v)
	case int32:
		number = int64(v)
	case int64:
		number = v
	case uint:
		if uint64(v) > math.MaxInt64 {
			return "", fmt.Errorf("%s value is too large", fieldName)
		}
		number = int64(v)
	case uint8:
		number = int64(v)
	case uint16:
		number = int64(v)
	case uint32:
		number = int64(v)
	case uint64:
		if v > math.MaxInt64 {
			return "", fmt.Errorf("%s value is too large", fieldName)
		}
		number = int64(v)
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return "", fmt.Errorf("%s must be a non-empty integer value", fieldName)
		}

		parsed, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return "", fmt.Errorf("%s must be an integer or numeric string", fieldName)
		}
		number = parsed
	default:
		return "", fmt.Errorf("%s must be an integer or numeric string", fieldName)
	}

	if number < 0 {
		return "", fmt.Errorf("%s must be non-negative", fieldName)
	}

	if !allowZero && number == 0 {
		return "", fmt.Errorf("%s must be greater than 0", fieldName)
	}

	return strconv.FormatInt(number, 10), nil
}

func SendTransaction(tx blockchain.Transaction) error {
	return nil
}
