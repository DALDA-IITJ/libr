package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)


// storeMessage stores the message in multiple database nodes
func storeMessage(msgCert MsgCert) error {
	// Fetch simulated database nodes
	dbNodes := fetchDatabaseNodes() // You will replace this with a real query to blockchain

	// Channel to collect errors
	errorChannel := make(chan error, len(dbNodes))

	var wg sync.WaitGroup

	// Iterate over each database node and send the message certificate
	for _, dbNode := range dbNodes {
		wg.Add(1)

		go func(dbNode DatabaseNode) {
			defer wg.Done()

			url := fmt.Sprintf("http://%s:%s/store", dbNode.IP, dbNode.Port)

			// Marshal the msgCert to JSON
			reqBody, err := json.Marshal(msgCert)
			if err != nil {
				errorChannel <- fmt.Errorf("failed to marshal msgCert for %s: %v", dbNode.IP, err)
				return
			}

			// Send POST request to store the message cert
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				errorChannel <- fmt.Errorf("error sending to %s: %v", dbNode.IP, err)
				return
			}
			defer resp.Body.Close()

			// If the response status code is not 200 OK, it's an error
			if resp.StatusCode != http.StatusOK {
				errorChannel <- fmt.Errorf("error from %s: %v", dbNode.IP, resp.Status)
				return
			}

			fmt.Printf("Message cert successfully sent to %s\n", dbNode.IP)
		}(dbNode)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChannel)

	// Check for errors
	for err := range errorChannel {
		if err != nil {
			return err
		}
	}

	return nil
}

// fetchDatabaseNodes simulates fetching the database nodes.
// TODO: Replace this with real blockchain query.
func fetchDatabaseNodes() []DatabaseNode {
	return []DatabaseNode{
		{"192.168.1.20", "8080"},
		{"192.168.1.21", "8081"},
		{"192.168.1.22", "8082"},
	}
}
