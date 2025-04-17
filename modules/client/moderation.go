package client

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
		{"172.31.119.163", "8080", "0449fa38bd8e4c31e49511e2a6475012d0e8fea28f7a882fe8755ef916b0fed71ad5b3e454839d192dbc76e995b152ec1658cbd621f9a7088e7e5d96a5aeaa5d24"},
		// {"172.31.91.201", "8080", "043aa3676f77bd5af04af81e7aff09e53433b9d4a4db259867be7da5a54e1cf4441831604e0a2639505feab3dab3113d5648929252d8398ea67f71843fd33a0260"},
		{"172.31.87.47", "8080", "04227547a108c40745cf479c54e430da76a75f2350772370020c10832d4de14409f6996cfff613048c4b6efe1251b4bbda86c0222a983b35c941826f2cf965c77c"},
		{"172.31.88.50", "8080", "0449fa38bd8e4c31e49511e2a6475012d0e8fea28f7a882fe8755ef916b0fed71ad5b3e454839d192dbc76e995b152ec1658cbd621f9a7088e7e5d96a5aeaa5d24"},
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
			print("body = ", string(body), "\n")
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
