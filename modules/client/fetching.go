package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func fetchMessages(timestamp string) ([]Message, error) {
	println("timestamp = %s", timestamp)

	// Fetch database nodes
	// relevantTxs, errX := blockchain.FetchBlockchain(msgCert.TS)

	// if errX != nil {
	// 	log.Fatal("❌ Error loading blockchain data:", errX)
	// }

	// dbNodes := fetchDatabaseNodes(RelevantTxs, timestamp, 3)

	// Fetch the database nodes (simulated for now)
	dbNodes := fetchDBTest()

	// Channel to collect responses
	responseChannel := make(chan []Message, len(dbNodes))

	// Channel to collect errors
	errorChannel := make(chan error, len(dbNodes))

	var wg sync.WaitGroup

	// Iterate over each database node and send the GET request
	for _, dbNode := range dbNodes {
		wg.Add(1)

		go func(dbNode DatabaseNode) {
			defer wg.Done()

			url := fmt.Sprintf("http://%s:%s/db/getmsg?time=%s", dbNode.IP, dbNode.Port, timestamp)

			resp, err := http.Get(url)
			if err != nil {
				errorChannel <- fmt.Errorf("error fetching messages from %s: %v", dbNode.IP, err)
				return
			}
			defer resp.Body.Close()

			// If the response is not OK, log an error
			if resp.StatusCode != http.StatusOK {
				errorChannel <- fmt.Errorf("error from %s: %v", dbNode.IP, resp.Status)
				return
			}

			// Read the response body
			var messageResponse struct {
				Data    []Message   `json:"data"`
				Error   interface{} `json:"error"` // Use the appropriate type if needed
				Message string      `json:"message"`
				Status  int         `json:"status"`
			}
			err = json.NewDecoder(resp.Body).Decode(&messageResponse)
			if err != nil {
				errorChannel <- fmt.Errorf("error decoding response from %s: %v", dbNode.IP, err)
				return
			}

			// Check if there's an error in the response
			if messageResponse.Error != nil {
				errorChannel <- fmt.Errorf("error in response from %s: %v", dbNode.IP, messageResponse.Error)
				return
			}

			// Send the messages to the response channel
			responseChannel <- messageResponse.Data
		}(dbNode)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(responseChannel)
	close(errorChannel)

	// Collect all responses
	var allMessages [][]Message
	for msgs := range responseChannel {
		allMessages = append(allMessages, msgs)
	}

	// If there were any errors during fetching, return the last error
	select {
	case err := <-errorChannel:
		if err != nil {
			return nil, err
		}
	default:
		// Continue without error
	}

	// Resolve messages by majority vote
	return resolveMessages(allMessages), nil
}

// resolveMessages resolves the messages by majority vote
func resolveMessages(allMessages [][]Message) []Message {
	var resolvedMessages []Message

	// Iterate over each message index in the responses
	for i := 0; i < len(allMessages[0]); i++ {
		// Count occurrences of each message
		messageCount := make(map[Message]int)

		for _, msgs := range allMessages {
			if i < len(msgs) {
				messageCount[msgs[i]]++
			}
		}

		// Find the message with the majority vote
		var majorityMessage Message
		maxCount := 0
		for msg, count := range messageCount {
			if count > maxCount {
				majorityMessage = msg
				maxCount = count
			}
		}

		// Add the majority message to the resolved list
		resolvedMessages = append(resolvedMessages, majorityMessage)
	}

	return resolvedMessages
}
