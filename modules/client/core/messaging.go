package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"client/core/blockchain"
	"client/core/crypto" // Import the crypto module
)

var RelevantTxs []blockchain.Transaction

func (c *Core) SendMessage(content string) error {

	fmt.Printf("Sending User Message %s\n", content)

	// Load keys
	privateKey, err := crypto.LoadPrivateKey()
	if err != nil {
		log.Fatal("❌ Error loading private key:", err)
	}
	fmt.Printf("Loaded Private Key\n")
	publicKey, err := crypto.LoadPublicKey()
	if err != nil {
		log.Fatal("❌ Error loading public key:", err)
	}
	fmt.Printf("Loaded Public Key\n")

	msg := UserMessage{
		Content:   content,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	RelevantTxs, err = blockchain.FetchBlockchain(msg.Timestamp)

	if err != nil {
		log.Fatal("❌ Error loading blockchain data:", err)
	}

	// Send to Moderation
	modCert, err := SendToModerators(msg.Content, msg.Timestamp)
	if err != nil {
		return err
	}

	// Create MsgCert object
	msgCert := MsgCert{
		Sender:  publicKey,
		Msg:     msg.Content,
		TS:      msg.Timestamp,
		ModCert: modCert.Signatures, // correct
	}

	// Serialize MsgCert for signing
	msgCertBytes, _ := json.Marshal(msgCert)
	signature, err := crypto.SignMessage(privateKey, string(msgCertBytes))
	if err != nil {
		return err
	}

	// Assign the signature
	msgCert.Sign = signature

	// Print the final message certificate
	fmt.Println("\n✅ Final MsgCert:", msgCert)

	// Send to Storage
	err = storeMessage(msgCert)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) FetchMessages(timestamp string) ([]Message, error) {
	return fetchMessages(timestamp)
}
