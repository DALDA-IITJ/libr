package utils

import (
	"crypto/ecdsa"

	"github.com/DALDA-IITJ/libr/modules/core/config"
	"github.com/DALDA-IITJ/libr/modules/core/crypto"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

var PublicKey string
var PrivateKey *ecdsa.PrivateKey

func EnsureKeyPair() error {

	private_key, err := crypto.LoadPrivateKey()
	public_key := config.GetEnv("PUBLIC_KEY", "")

	if err == nil && public_key != "" {
		PrivateKey = private_key
		PublicKey = public_key
		logger.Info("Key pair already exists")
		return nil // Both keys exist
	}

	if err := crypto.GenerateKeys(); err != nil {
		logger.Fatal("Failed to generate keys: " + err.Error())
		return err
	}

	private_key, err = crypto.LoadPrivateKey()
	if err != nil {
		logger.Fatal("Failed to load private key: " + err.Error())
		return err
	}

	public_key = config.GetEnv("PUBLIC_KEY", "")

	PrivateKey = private_key
	PublicKey = public_key
	logger.Info("Key pair generated and saved to .env")
	return nil
}
