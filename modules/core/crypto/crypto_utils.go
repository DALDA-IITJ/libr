package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
)
 
// GenerateKeys generates a new ECDSA key pair and stores it in the environment.
func GenerateKeys() error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex := fmt.Sprintf("%x%x", privateKey.PublicKey.X, privateKey.PublicKey.Y)

	os.Setenv("PRIVATE_KEY", privateKeyHex)
	os.Setenv("PUBLIC_KEY", publicKeyHex)

	fmt.Println("🔑 New Keys Generated:")
	fmt.Println("Private Key:", privateKeyHex)
	fmt.Println("Public Key:", publicKeyHex)

	return nil
}

// LoadPrivateKey retrieves the private key from environment variables.
func LoadPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return nil, errors.New("PRIVATE_KEY not found in environment")
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKeyBytes)

	return privateKey, nil
}

// LoadPublicKey retrieves the public key from environment variables.
func LoadPublicKey() (string, error) {
	publicKeyHex := os.Getenv("PUBLIC_KEY")
	if publicKeyHex == "" {
		return "", errors.New("PUBLIC_KEY not found in environment")
	}
	return publicKeyHex, nil
}

// SignMessage signs the given message using ECDSA.
func SignMessage(privateKey *ecdsa.PrivateKey, message string) (string, error) {
	hash := sha256.Sum256([]byte(message))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	signature := fmt.Sprintf("%x:%x", r, s)
	return signature, nil
}

// VerifySignature verifies an ECDSA signature.
func VerifySignature(publicKeyHex, message, signature string) (bool, error) {
	hash := sha256.Sum256([]byte(message))

	// Extract r and s from the signature
	var r, s big.Int
	_, err := fmt.Sscanf(signature, "%x:%x", &r, &s)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	// Convert publicKeyHex to ECDSA PublicKey
	if len(publicKeyHex) < 128 {
		return false, errors.New("invalid public key length")
	}

	x := new(big.Int)
	y := new(big.Int)
	x.SetString(publicKeyHex[:64], 16)
	y.SetString(publicKeyHex[64:], 16)

	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	// Verify the signature
	isValid := ecdsa.Verify(&publicKey, hash[:], &r, &s)
	return isValid, nil
}
