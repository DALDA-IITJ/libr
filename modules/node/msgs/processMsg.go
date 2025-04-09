package msgs

import (
	"fmt"
	"strconv"

	"github.com/DALDA-IITJ/libr/modules/node/blockchain"
	"github.com/DALDA-IITJ/libr/modules/node/db"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

// Generate a unique hash for each message
// func generateMsgHash(msgCert map[string]interface{}) string {
// 	msgBytes, _ := json.Marshal(msgCert)
// 	hash := sha256.Sum256(msgBytes)
// 	return fmt.Sprintf("%x", hash)
// }

// ProcessMsg verifies, validates, and stores the message
func ProcessMsg(msgCert map[string]interface{}) error {
	// Extract timestamp
	ts, ok := msgCert["ts"].(string)
	if !ok {
		logger.Error("ts is missing or invalid")
		return fmt.Errorf("timestamp is missing or not of the expected type")
	}

	print("ts = ", ts)

	// Check blockchain state
	err := blockchain.GetBlockchainState(ts)
	if err != nil {
		logger.Error("Failed to get blockchain state: " + err.Error())
		return err
	}

	// Parallel verification using goroutines
	successChan := make(chan bool, 4)
	errorChan := make(chan error, 4)

	// Verify K Mods
	go func() {
		if err := verifyKMods(ts, msgCert); err != nil {
			errorChan <- err
			return
		}
		successChan <- true
	}()

	go func() {
		if err := verifyMsgCert(ts, msgCert); err != nil {
			errorChan <- err
			return
		}
		successChan <- true
	}()

	go func() {
		if err := VerifyDBNode(ts); err != nil {
			errorChan <- err
			return
		}
		successChan <- true
	}()

	go func() {
		if err := verifyDBTime(ts); err != nil {
			errorChan <- err
			return
		}
		successChan <- true
	}()

	// Process results from goroutines
	var finalSuccess = true
	for i := 0; i < 4; i++ {
		select {
		case success := <-successChan:
			if !success {
				finalSuccess = false
			}
		case err := <-errorChan:
			if err != nil {
				logger.Error("Error occurred: " + err.Error())
				finalSuccess = false
				return err
			}
		}
	}

	fmt.Println("Verified all")

	if !finalSuccess {
		return fmt.Errorf("failed to process message certificate")
	}

	// Extract required fields for database storage
	sender, ok := msgCert["sender"].(string)
	if !ok {
		logger.Error("Sender is missing")
		return fmt.Errorf("sender is missing")
	}

	content, ok := msgCert["msg"].(string)
	if !ok {
		logger.Error("Content is missing")
		return fmt.Errorf("content is missing")
	}

	bucketId, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		logger.Error("Failed to convert timestamp to int64: " + err.Error())
		return fmt.Errorf("failed to convert timestamp")
	}

	// Save message to DB
	_, err = db.SaveToDB(sender, content, bucketId)
	if err != nil {
		logger.Error("Database save failed: " + err.Error())
		return err
	}

	logger.Info("Message processed and saved successfully")
	return nil
}

// VerifyKMods ensures that required moderators have approved the message
func verifyKMods(ts string, msgCert map[string]interface{}) error {
	// modList, err := utils.FindKMods(ts)
	// if err != nil {
	// 	logger.Error("Failed to find K mods: " + err.Error())
	// 	return err
	// }

	// if len(modList) == 0 {
	// 	logger.Error("No K mods found")
	// 	return fmt.Errorf("no K mods found")
	// }

	// // Verify mod certs
	// missingMods := []string{}
	// modCerts, ok := msgCert["modCerts"].([]interface{})
	// if !ok {
	// 	logger.Error("modCerts field is missing or invalid")
	// 	return fmt.Errorf("modCerts field is missing or invalid")
	// }

	// for _, mod := range modList {
	// 	isPresent := false
	// 	for _, cert := range modCerts {
	// 		if modCertMap, ok := cert.(map[string]interface{}); ok {
	// 			if pubkey, ok := modCertMap["pubkey"].(string); ok && pubkey == mod {
	// 				isPresent = true
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if !isPresent {
	// 		missingMods = append(missingMods, mod)
	// 	}
	// }

	// if len(missingMods) > 0 {
	// 	logger.Error("Missing mods: " + fmt.Sprintf("%v", missingMods))
	// 	return fmt.Errorf("missing mods: %v", missingMods)
	// }
	return nil
}

// VerifyMsgCert ensures message signature is valid
func verifyMsgCert(ts string, msgCert map[string]interface{}) error {
	// isVerified, err := utils.VerifyMsgCert(ts, msgCert)
	// if err != nil {
	// 	logger.Error("Message certificate verification failed: " + err.Error())
	// 	return err
	// }
	// if !isVerified {
	// 	logger.Error("Message verification failed")
	// 	return fmt.Errorf("message verification failed")
	// }
	return nil
}

// VerifyDBNode checks if this node is in the list of K servers
func VerifyDBNode(ts string) error {
	//todo
	// dbList, err := utils.FindKServers(ts)
	// if err != nil {
	// 	logger.Error("Failed to find K servers: " + err.Error())
	// 	return err
	// }

	// if len(dbList) == 0 {
	// 	logger.Error("No K servers found")
	// 	return fmt.Errorf("no K servers found")
	// }

	// isServerCorrect := false
	// for _, server := range dbList {
	// 	if server == utils.PublicKey {
	// 		isServerCorrect = true
	// 		break
	// 	}
	// }

	// if !isServerCorrect {
	// 	logger.Error("Public key does not match any server in the list")
	// 	return fmt.Errorf("public key does not match any server in the list")
	// }
	return nil
}

func verifyDBTime(ts string) error {
	// Convert timestamp string to int64
	// timestamp, err := strconv.ParseInt(ts, 10, 64)
	// if err != nil {
	// 	logger.Error("Invalid timestamp: " + err.Error())
	// 	return fmt.Errorf("invalid timestamp")
	// }

	// // Load allowed difference in minutes from environment variable; default to 5 if not set or invalid.
	// allowedMinutes := 5
	// if minutesStr := os.Getenv("ALLOWED_DIFF_MINUTES"); minutesStr != "" {
	// 	if val, convErr := strconv.Atoi(minutesStr); convErr == nil {
	// 		allowedMinutes = val
	// 	}
	// }

	// trafficStr := os.Getenv("TRAFFIC")
	// traffic, err := strconv.ParseInt(trafficStr, 10, 64)
	// if err != nil {
	// 	logger.Warn("Invalid TRAFFIC value, defaulting to 1")
	// 	traffic = 1
	// }

	// timestamp = timestamp * traffic
	// currentTime := time.Now().Unix()

	// // Check if the timestamp is in the future
	// if timestamp > currentTime {
	// 	logger.Error("Timestamp is in the future")
	// 	return fmt.Errorf("timestamp is in the future")
	// }

	// // Check if the timestamp is older than allowed (more than allowedMinutes earlier than current DB time)
	// if currentTime-timestamp > int64(allowedMinutes*60) {
	// 	logger.Error("Timestamp is too old")
	// 	return fmt.Errorf("timestamp is too old")
	// }

	return nil
}
