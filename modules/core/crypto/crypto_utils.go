package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/joho/godotenv"
)

type ecdsaSignature struct {
	R *big.Int
	S *big.Int
}

// GenerateKeys generates a new ECDSA key pair and stores it in the environment.
func GenerateKeys() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found, creating new keys...")
	}

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	publicKeyHex := os.Getenv("PUBLIC_KEY")

	// Check if keys already exist in .env
	if privateKeyHex != "" && publicKeyHex != "" {
		fmt.Println("🔑 Keys already exist in .env")
		fmt.Println("Private Key:", privateKeyHex)
		fmt.Println("Public Key:", publicKeyHex)
		return nil
	}

	// Generate new keys
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	privateKeyHex = hex.EncodeToString(privateKey.D.Bytes())
	publicKeyHex = fmt.Sprintf("%x%x", privateKey.PublicKey.X, privateKey.PublicKey.Y)

	// Save keys to .env file
	envContent := fmt.Sprintf("PRIVATE_KEY=%s\nPUBLIC_KEY=%s\n", privateKeyHex, publicKeyHex)
	err = os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		return err
	}

	// Set environment variables
	os.Setenv("PRIVATE_KEY", privateKeyHex)
	os.Setenv("PUBLIC_KEY", publicKeyHex)

	fmt.Println("🔑 New Keys Generated and Saved to .env:")
	fmt.Println("Private Key:", privateKeyHex)
	fmt.Println("Public Key:", publicKeyHex)

	return nil
}

// LoadPrivateKey retrieves the private key from environment variables.
func LoadPrivateKey() (string, error) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return "", errors.New("PRIVATE_KEY not found in environment")
	}

	return privateKeyHex, nil
}

// LoadPublicKey retrieves the public key from environment variables.
func LoadPublicKey() (string, error) {
	publicKeyHex := os.Getenv("PUBLIC_KEY")
	if publicKeyHex == "" {
		return "", errors.New("PUBLIC_KEY not found in environment")
	}
	return publicKeyHex, nil
}

func SignMessage(message string, privateKeyHex string) (string, error) {
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
	signatureHex := hex.EncodeToString(signature)
	print("Signature: ", signatureHex, "\n")
	return signatureHex, nil
}

// VerifySignature verifies an ECDSA signature.
func VerifySignature(message string, signatureHex string, publicKeyHex string) (bool, error) {
	fmt.Println("Verifying signature...")
	fmt.Println("Message: ", message)
	fmt.Println("Signature: ", signatureHex)
	fmt.Println("Public Key: ", publicKeyHex)

	// Decode signature from hex
	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Decode the DER-encoded signature
	var sig ecdsaSignature
	_, err = asn1.Unmarshal(signatureBytes, &sig)
	if err != nil {
		return false, fmt.Errorf("failed to decode DER signature: %v", err)
	}

	// Decode public key from hex
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	// Validate public key length
	if len(publicKeyBytes) != 64 {
		return false, fmt.Errorf("expected 64-byte public key (128 hex chars), got %d bytes", len(publicKeyBytes))
	}

	// Extract X and Y coordinates (32 bytes each)
	x := new(big.Int).SetBytes(publicKeyBytes[:32])
	y := new(big.Int).SetBytes(publicKeyBytes[32:])

	// Identify the curve
	curves := []elliptic.Curve{
		elliptic.P256(),
		elliptic.P224(),
		elliptic.P384(),
		elliptic.P521(),
	}

	var validCurve elliptic.Curve
	for _, curve := range curves {
		if curve.IsOnCurve(x, y) {
			validCurve = curve
			break
		}
	}

	if validCurve == nil {
		return false, errors.New("public key point is not on any standard curve")
	}

	publicKey := &ecdsa.PublicKey{Curve: validCurve, X: x, Y: y}

	// Hash the message
	hash := sha256.Sum256([]byte(message))

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], sig.R, sig.S), nil
}
