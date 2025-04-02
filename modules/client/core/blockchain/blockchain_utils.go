package blockchain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// BlockchainResponse represents the JSON structure of blockchain data.
type BlockchainResponse struct {
	BlockchainHeader struct {
		BlockchainLength int `json:"blockchainLength"`
	} `json:"blockchainHeader"`
	Blocks []Block `json:"blocks"`
}

// Block represents a block in the blockchain.
type Block struct {
	PrevBlockHash string        `json:"prevBlockHash"`
	Transactions  []Transaction `json:"transactions"`
	BlockNumber   int           `json:"blockNumber"`
	Nonce         string        `json:"nonce"`
	BlockHash     string        `json:"blockHash"`
}

// Transaction represents a transaction in a block.
type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amt"`
	Nonce     string `json:"nonce"`
	Data      TxData `json:"data"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

// TxData represents the transaction's data field.
type TxData struct {
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
	Leaver   string            `json:"leaver,omitempty"`
}

// FetchDatabaseNode fetches blockchain transactions at a given timestamp.
func FetchBlockchain(timestamp string) ([]Transaction, error) {
	// Read from config (replace with your config handling)
	blockchainURL := os.Getenv("BLOCKCHAIN_URL") // Example: "http://127.0.0.1:5000"
	communityPublicKey := os.Getenv("COMMUNITY_PUBLIC_KEY")

	// Construct request URL
	url := fmt.Sprintf("%s/blockchain/get-at/%s", blockchainURL, timestamp)

	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blockchain: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var blockchainData BlockchainResponse
	if err := json.Unmarshal(body, &blockchainData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Filter relevant transactions
	var relevantTxs []Transaction
	for _, block := range blockchainData.Blocks {
		for _, tx := range block.Transactions {
			if tx.Recipient == communityPublicKey && (tx.Data.Type == "DB_JOINED" || tx.Data.Type == "DB_LEFT" || tx.Data.Type == "MOD_JOINED" || tx.Data.Type == "MOD_LEFT") {
				relevantTxs = append(relevantTxs, tx)
			}
		}
	}

	return relevantTxs, nil
}
