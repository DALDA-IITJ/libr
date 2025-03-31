package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	logger "libr/core"
	"os"
)

func GetPrivateKey() (ecdsa.PrivateKey, error) {
	file, err := os.Open("key.json")
	if err != nil {
		logger.Error("Failed to open key.json: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read key.json: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}

	var keyData struct {
		PrivateKey string `json:"private_key"`
	}
	if err := json.Unmarshal(data, &keyData); err != nil {
		logger.Error("Failed to unmarshal key.json: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}

	block, _ := pem.Decode([]byte(keyData.PrivateKey))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		err := errors.New("failed to decode PEM block containing private key")
		logger.Error(err.Error())
		return ecdsa.PrivateKey{}, err
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		logger.Error("Failed to parse EC private key: " + err.Error())
		return ecdsa.PrivateKey{}, err
	}

	return *privateKey, nil
}

func GetPublicKey() (ecdsa.PublicKey, error) {
	file, err := os.Open("key.json")
	if err != nil {
		logger.Error("Failed to open key.json: " + err.Error())
		return ecdsa.PublicKey{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Failed to read key.json: " + err.Error())
		return ecdsa.PublicKey{}, err
	}

	var keyData struct {
		PublicKey string `json:"public_key"`
	}
	if err := json.Unmarshal(data, &keyData); err != nil {
		logger.Error("Failed to unmarshal key.json: " + err.Error())
		return ecdsa.PublicKey{}, err
	}

	block, _ := pem.Decode([]byte(keyData.PublicKey))
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
