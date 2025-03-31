package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	logger "libr/core"
	"libr/modules/node/utils"
)

type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amt"`
	Data      string `json:"data"`
	Nonce     int    `json:"nonce"`
	Signature string `json:"sign"`
}

func RegisterNode() error {
	sender := getSender()
	recipient := getRecipient()
	amount := getAmount()
	data := getData()
	nonce := getNonce()
	signature, err := makeSignature(data)
	if err != nil {
		logger.Error("Failed to create signature: " + err.Error())
		return err
	}
	transaction := Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Data:      data,
		Nonce:     nonce,
		Signature: signature,
	}

	err = SendTransaction(transaction)
	if err != nil {
		logger.Error("Failed to send transaction: " + err.Error())
		return err
	}

	logger.Info("Transaction sent successfully")
	return nil
}

// todo
func getSender() string {
	// Get public key from key.json and return it

	return "sender_public_key"
}

// todo
func getRecipient() string {
	// Get recipient public key from config.yaml and return it
	return "recipient_public_key"
}

// todo
func getAmount() int {
	return 0
}

// todo
func getData() string {
	return "data"
}

// todo
func getNonce() int {
	return 0
}

// todo
func makeSignature(data string) (string, error) {

	privateKey, err := utils.GetPrivateKey()
	if err != nil {
		logger.Error("Failed to retrieve private key: " + err.Error())
		return "", err
	}

	hash := sha256.Sum256([]byte(data))

	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, hash[:])
	if err != nil {
		logger.Error("Failed to sign data: " + err.Error())
		return "", err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return string(signature), nil
}
