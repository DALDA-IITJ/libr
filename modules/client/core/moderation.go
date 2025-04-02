package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// FetchModerators simulates retrieving moderator details.
// TODO: Replace with real blockchain data fetching.
func FetchModerators() []Moderator {
	
	fmt.Printf("Fetching mods...")

	return []Moderator{
		{"localhost", "8080", "PublicKey1"},
		{"localhost", "8081", "PublicKey2"},
		// {"localhost", "8082", "PublicKey3"},
	}
}

// SendToModerators sends a message to moderators and collects signatures.
func SendToModerators(msg string, timestamp string) (ModCert, error) {
	fmt.Printf("Sending to mods...")
	moderators := FetchModerators() // requires implementation
	var modCert ModCert
	modCert.Msg = msg
	modCert.Timestamp = timestamp

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel to collect errors
	// Read: https://medium.com/@developer.naren/handling-errors-in-channels-e76ad5cbf3a0
	errorChannel := make(chan error, len(moderators))

	// Iterate over each moderator and send request
	for _, mod := range moderators {
		wg.Add(1)

		go func(mod Moderator) {
			defer wg.Done()

			url := fmt.Sprintf("http://%s:%s/moderate", mod.IP, mod.Port)
			reqBody, _ := json.Marshal(map[string]string{
				"message":   msg,
				"timestamp": timestamp,
			})

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				errorChannel <- fmt.Errorf("error sending to %s: %v", mod.IP, err)
				return
			}
			defer resp.Body.Close()

			// Read response body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errorChannel <- fmt.Errorf("error reading response from %s: %v", mod.IP, err)
				return
			}

			var modSign ModSign
			err = json.Unmarshal(body, &modSign)
			if err != nil {
				errorChannel <- fmt.Errorf("invalid response format from %s", mod.IP)
				return
			}

			// Store valid signatures
			mutex.Lock()
			modCert.Signatures = append(modCert.Signatures, modSign)
			mutex.Unlock()
		}(mod)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChannel)

	// Check for errors
	var finalError error
	for err := range errorChannel {
		fmt.Println("Warning:", err) // Log errors but do not break flow
		finalError = err             // Keep the last error, if any
	}

	return modCert, finalError
}
