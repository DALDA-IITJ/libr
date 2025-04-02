package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"math/rand"
	"node/utils/logger"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var PublicKey ecdsa.PublicKey
var PrivateKey ecdsa.PrivateKey
// var PrivateKeyHex string
var PublicKeyHex string

func EnsureKeyPair() error {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Failed to load .env: " + err.Error())
	}

	private_key, privErr := getPrivateKey()
	public_key, pubErr := getPublicKey()

	if privErr == nil && pubErr == nil {
		PrivateKey = private_key
		PublicKey = public_key
		// PrivateKeyHex = hex.EncodeToString(PrivateKey.D.Bytes())
		PublicKeyHex = hex.EncodeToString(PublicKey.X.Bytes()) + hex.EncodeToString(PublicKey.Y.Bytes())
		logger.Info("Key pair already exists")
		return nil // Both keys exist
	}

	// Generate a new key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		logger.Error("Failed to generate key pair: " + err.Error())
		return err
	}

	// Encode private key to hex
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		logger.Error("Failed to marshal private key: " + err.Error())
		return err
	}
	privateKeyHex := hex.EncodeToString(privBytes)

	// Encode public key to PEM
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		logger.Error("Failed to marshal public key: " + err.Error())
		return err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	// Save private key to .env
	err = godotenv.Write(map[string]string{"PRIVATE_KEY": privateKeyHex}, ".env")
	if err != nil {
		logger.Error("Failed to write private key to .env: " + err.Error())
		return err
	}

	// Save public key to a separate file (optional)
	if err := os.WriteFile("public_key.pem", pubPem, 0644); err != nil {
		logger.Error("Failed to write public key to file: " + err.Error())
		return err
	}

	PublicKey = privateKey.PublicKey
	PrivateKey = *privateKey
	PublicKeyHex = hex.EncodeToString(PublicKey.X.Bytes()) + hex.EncodeToString(PublicKey.Y.Bytes())
	
	logger.Info("New key pair generated and saved to .env")
	return nil
}

func getPrivateKey() (ecdsa.PrivateKey, error) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		err := errors.New("private key not found in .env")
		logger.Error(err.Error())
		return ecdsa.PrivateKey{}, err
	}

	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		logger.Error("Failed to decode private key hex: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}

	privateKey, err := x509.ParseECPrivateKey(privBytes)
	if err != nil {
		logger.Error("Failed to parse EC private key: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}

	return *privateKey, nil
}

func getPublicKey() (ecdsa.PublicKey, error) {
	pubPem, err := os.ReadFile("public_key.pem")
	if err != nil {
		logger.Error("Failed to read public_key.pem: " + err.Error())
		return ecdsa.PublicKey{}, err
	}

	block, _ := pem.Decode(pubPem)
	if block == nil || block.Type != "PUBLIC KEY" {
		err := errors.New("failed to decode PEM block containing public key")
		logger.Error(err.Error())
		return ecdsa.PublicKey{}, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Error("Failed to parse public key: " + err.Error())
		return ecdsa.PublicKey{}, err
	}

	return *publicKey.(*ecdsa.PublicKey), nil
}
