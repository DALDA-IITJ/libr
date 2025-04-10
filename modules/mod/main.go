package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DALDA-IITJ/libr/modules/core/config"
	"github.com/DALDA-IITJ/libr/modules/core/crypto"
)

// ---------- STRUCTS ----------
// Updated request struct for message moderation with timestamp.
type ModerateRequest struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ModerateResponse struct {
	Allowed string  `json:"allowed"` // "yes" or "no"
	Score   float64 `json:"score"`
}

type Category struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

func handleModerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request expecting "message" and "timestamp"
	var req ModerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Run moderation check on the message text.
	categories, err := ModerateText(req.Message)
	if err != nil {
		http.Error(w, "Moderation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Compute the moderation score.
	var score float64
	for _, cat := range categories {
		weight := getCategoryWeight(cat.Name)
		score += cat.Confidence * weight
	}

	// If score exceeds threshold then reject the message.
	if score >= Threshold {
		http.Error(w, "Message did not pass moderation", http.StatusBadRequest)
		return
	}

	// combinedMsg := req.Message + req.Timestamp

	payload := map[string]string{
		"message":   req.Message,
		"timestamp": req.Timestamp,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	privateKey, err := crypto.LoadPrivateKey()
	if err != nil {
		http.Error(w, "Could not load private key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	signature, err := crypto.SignMessage(privateKey, string(payloadBytes))
	if err != nil {
		http.Error(w, "Signing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	publicKey, err := crypto.LoadPublicKey()
	if err != nil {
		http.Error(w, "Could not load public key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Response includes the public key and the signature.
	response := map[string]string{
		"public_key": publicKey,
		"sign":       signature,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	config.LoadEnv()

	port := config.GetEnv("PORT", "4000")
	addr := ":" + port

	http.HandleFunc("/moderate", handleModerate)

	fmt.Println("🚀 Server started at http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
