package client

import (
	"bytes"
	// "client/core/blockchain"
	"encoding/json"
	"fmt"

	// "log"
	"net/http"
	"sync"
)

// storeMessage stores the message in multiple database nodes
func storeMessage(msgCert MsgCert) error {
	// Fetch database nodes
	// relevantTxs, errX := blockchain.FetchBlockchain(msgCert.TS)

	// if errX != nil {
	// 	log.Fatal("❌ Error loading blockchain data:", errX)
	// }

	// dbNodes := fetchDatabaseNodes(relevantTxs) // You will replace this with a real query to blockchain

	dbNodes := fetchDBTest()
	fmt.Printf("Fetched db nodes...")

	// Channel to collect errors
	errorChannel := make(chan error, len(dbNodes))

	var wg sync.WaitGroup

	// Iterate over each database node and send the message certificate
	for _, dbNode := range dbNodes {
		wg.Add(1)

		go func(dbNode DatabaseNode) {
			defer wg.Done()

			fmt.Printf("Fetching URL http://%s:%s/db/savemsg\n", dbNode.IP, dbNode.Port)

			url := fmt.Sprintf("http://%s:%s/db/savemsg", dbNode.IP, dbNode.Port)

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

func fetchDBTest() []DatabaseNode {

	fmt.Printf("Fetching dbs...\n")

	return []DatabaseNode{
		{"172.31.119.163", "8082"},
		// {"172.31.91.201", "8082"},
		{"172.31.87.47", "8082"},
		{"172.31.88.50", "8082"},
	}
}
