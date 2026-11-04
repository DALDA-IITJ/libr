package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestVerifySignatureWithSEC1AndRawKeys(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	msg := "test moderation message"
	sig, err := SignMessage(priv, msg)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	raw128 := fmt.Sprintf("%064x%064x", priv.PublicKey.X, priv.PublicKey.Y)
	sec1_130 := "04" + raw128

	// Test 1: Raw 128-char key
	valid, err := VerifySignature(raw128, msg, sig)
	if err != nil || !valid {
		t.Errorf("Expected raw 128-char key to verify, got valid=%v, err=%v", valid, err)
	}

	// Test 2: SEC1 130-char key with 04 prefix
	validSec1, err := VerifySignature(sec1_130, msg, sig)
	if err != nil || !validSec1 {
		t.Errorf("Expected SEC1 130-char key to verify, got valid=%v, err=%v", validSec1, err)
	}

	// Test 3: Invalid length key
	_, err = VerifySignature("04abcd", msg, sig)
	if err == nil {
		t.Error("Expected error for invalid key length, got nil")
	}

	// Test 4: Invalid hex coordinates
	badHex := "04" + "zz" + raw128[2:]
	_, err = VerifySignature(badHex, msg, sig)
	if err == nil {
		t.Error("Expected error for invalid hex digits, got nil")
	}
}
