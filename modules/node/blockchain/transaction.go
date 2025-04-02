package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"node/utils"
	"node/utils/logger"
)

type Txn struct {
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
	signature, err := makeSign(data)
	if err != nil {
		logger.Error("Failed to create signature: " + err.Error())
		return err
	}
	txn := Txn{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Data:      data,
		Nonce:     nonce,
		Signature: signature,
	}

	err = SendTxn(txn)
	if err != nil {
		logger.Error("Failed to send transaction: " + err.Error())
		return err
	}

	logger.Info("Transaction sent successfully")
	return nil
}

// todo
func getSender() string {
	publicKey := utils.PublicKey
	return string(publicKey.X.Bytes()) + string(publicKey.Y.Bytes())
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
func makeSign(data string) (string, error) {

	privateKey := utils.PrivateKey
	hash := sha256.Sum256([]byte(data))

	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, hash[:])
	if err != nil {
		logger.Error("Failed to sign data: " + err.Error())
		return "", err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return string(signature), nil
}
