package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// sign takes a message and a private key in hex format, returns a signature in hex format
func sign(message string, privateKeyHex string) (string, error) {
	// Decode private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	// Parse private key
	curve := elliptic.P256() // Using P256 curve
	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(privateKeyBytes)
	privateKey.PublicKey.Curve = curve
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curve.ScalarBaseMult(privateKeyBytes)

	// Hash the message
	hash := sha256.Sum256([]byte(message))

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("signing error: %v", err)
	}

	// Convert signature to hex
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	signature := append(rBytes, sBytes...)
	return hex.EncodeToString(signature), nil
}

// verify takes a message, signature in hex, and public key in hex, returns whether signature is valid
func verify(message string, signatureHex string, publicKeyHex string) (bool, error) {
	// Decode signature from hex
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Split signature into r and s components
	sigLen := len(signatureBytes)
	if sigLen%2 != 0 {
		return false, fmt.Errorf("invalid signature length")
	}

	r := new(big.Int).SetBytes(signatureBytes[:sigLen/2])
	s := new(big.Int).SetBytes(signatureBytes[sigLen/2:])

	// Decode public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	// Parse public key
	curve := elliptic.P256() // Using P256 curve
	x, y := elliptic.Unmarshal(curve, publicKeyBytes)
	if x == nil {
		return false, fmt.Errorf("invalid public key")
	}
	publicKey := &ecdsa.PublicKey{Curve: curve, X: x, Y: y}

	// Hash the message
	hash := sha256.Sum256([]byte(message))

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// Example usage
func main() {
	// Generate a new ECDSA key pair for testing
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		return
	}
	// Convert private key to hex
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	print("Private Key: ", privateKeyHex, "\n")
	// Convert public key to hex
	publicKeyBytes := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	publicKeyHex := hex.EncodeToString(publicKeyBytes)
	print("Public Key: ", publicKeyHex, "\n")

	// Message to sign
	message := "Hello, ECDSA!"

	// Sign the message
	signature, err := sign(message, privateKeyHex)
	if err != nil {
		fmt.Printf("Error signing: %v\n", err)
		return
	}
	fmt.Printf("Signature: %s\n", signature)

	// Verify the signature
	valid, err := verify(message, signature, publicKeyHex)
	if err != nil {
		fmt.Printf("Error verifying: %v\n", err)
		return
	}
	fmt.Printf("Signature valid: %v\n", valid)

	// Try with a tampered message
	tamperedMessage := "Hello, tampered ECDSA!"
	valid, err = verify(tamperedMessage, signature, publicKeyHex)
	if err != nil {
		fmt.Printf("Error verifying: %v\n", err)
		return
	}
	fmt.Printf("Tampered message signature valid: %v\n", valid)
}
