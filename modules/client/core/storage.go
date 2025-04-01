package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const dbURL = "http://database-service/api/store"

func storeMessage(msg Message) error {
	data, _ := json.Marshal(msg)
	resp, err := http.Post(dbURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending message to database: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to store message")
	}

	return nil
}

func fetchMessages(timestamp string) ([]Message, error) {
	resp, err := http.Get(fmt.Sprintf("%s?timestamp=%s", dbURL, timestamp))
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}
	defer resp.Body.Close()

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("error decoding messages: %v", err)
	}

	return messages, nil
}
