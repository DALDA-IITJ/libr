package api

import (
	"fmt"
	"io"
	"net/http"
	"node/db"
	"node/utils"
	"node/utils/logger"
	"node/utils/resp"
	"strconv"
)

func SaveMsgHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll((r.Body))
	if err != nil {
		logger.Error("Failed to read request body: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Failed to read request body", err, nil)
		return
	}

	err = utils.UnWrapMsgCert(body)
	if err != nil {
		logger.Error("Failed to unwrap message certificate: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Failed to unwrap message certificate", err, nil)
		return
	}

	successChan := make(chan bool, 3)
	errorChan := make(chan error, 3)

	go func() {
		ts, ok := utils.MsgCert["timestamp"].(int64)
		if !ok {
			logger.Error("Timestamp is missing or not of the expected type.")
			errorChan <- nil
			successChan <- false
			return
		}

		modList, err := utils.FindKMods(ts)
		if err != nil {
			logger.Error("Failed to find K mods: " + err.Error())
			errorChan <- err
			successChan <- false
			return
		}

		if len(modList) == 0 {
			logger.Error("No K mods found.")
			errorChan <- nil
			successChan <- false
			return
		}

		missingMods := []string{}
		for _, mod := range modList {
			isPresent := false
			for _, modCert := range utils.ModCerts {
				if mod == modCert["pubkey"] {
					isPresent = true
					break
				}
			}
			if !isPresent {
				missingMods = append(missingMods, mod)
			}
		}

		if len(missingMods) > 0 {
			logger.Error("Some mods are not present in the mod certificates: " + fmt.Sprintf("%v", missingMods))
			errorChan <- fmt.Errorf("missing mods: %v", missingMods)
			successChan <- false
			return
		}
		successChan <- true
	}()

	go func() {
		isVerified, err := utils.VerifyMsgCert()
		if err != nil {
			logger.Error("Failed to verify message certificate: " + err.Error())
			errorChan <- err
			successChan <- false
			return
		}
		if !isVerified {
			logger.Error("Message certificate verification failed.")
			errorChan <- nil
			successChan <- false
			return
		}
		successChan <- true
	}()

	go func() {
		ts, ok := utils.MsgCert["timestamp"].(int64)
		if !ok {
			logger.Error("Timestamp is missing or not of the expected type.")
			errorChan <- nil
			successChan <- false
			return
		}

		dbList, err := utils.FindKServers(ts)
		if err != nil {
			logger.Error("Failed to find K servers: " + err.Error())
			errorChan <- err
			successChan <- false
			return
		}

		if len(dbList) == 0 {
			logger.Error("No K servers found.")
			errorChan <- nil
			successChan <- false
			return
		}

		isServerCorrect := false
		for _, server := range dbList {
			if server == utils.PublicKeyHex {
				isServerCorrect = true
				break
			}
		}

		if !isServerCorrect {
			logger.Error("Public key does not match any server in the list.")
			errorChan <- nil
			successChan <- false
			return
		}
		successChan <- true
	}()

	var finalSuccess bool = true
	for i := 0; i < 2; i++ {
		select {
		case success := <-successChan:
			if !success {
				finalSuccess = false
			}
		case err := <-errorChan:
			if err != nil {
				resp.SendResponse(w, http.StatusBadRequest, "An error occurred", err, nil)
				return
			}
		}
	}

	if !finalSuccess {
		resp.SendResponse(w, http.StatusUnauthorized, "Validation failed", nil, nil)
		return
	}

	sender, ok := utils.MsgCert["sender"].(string)
	if !ok {
		logger.Error("Sender is missing or not of the expected type.")
		resp.SendResponse(w, http.StatusBadRequest, "Sender is missing or not of the expected type", nil, nil)
		return
	}

	content, ok := utils.MsgCert["msg"].(string)
	if !ok {
		logger.Error("Content is missing or not of the expected type.")
		resp.SendResponse(w, http.StatusBadRequest, "Content is missing or not of the expected type", nil, nil)
		return
	}

	timestampStr, ok := utils.MsgCert["timestamp"].(string)
	if !ok {
		logger.Error("Timestamp is not a string.")
		resp.SendResponse(w, http.StatusBadRequest, "Invalid timestamp format", nil, nil)
		return
	}

	bucketId, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		logger.Error("Failed to convert timestamp to int64: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Invalid timestamp format", err, nil)
		return
	}

	_, err = db.SaveToDB(sender, content, bucketId)
	if err != nil {
		logger.Error("Failed to save message to database: " + err.Error())
		resp.SendResponse(w, http.StatusInternalServerError, "Failed to save message to database", err, nil)
		return
	}

	logger.Info("Message saved successfully")
	resp.SendResponse(w, http.StatusOK, "Message saved successfully", nil, nil)

}

func GetMsgHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	time := query.Get("time")
	if time == "" {
		logger.Error("Missing time parameter")
		resp.SendResponse(w, http.StatusBadRequest, "Missing time parameter", nil, nil)
		return
	}

	timestamp, err := strconv.ParseInt(time, 10, 64)
	if err != nil {
		logger.Error("Failed to convert time parameter to int64: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Invalid time parameter", err, nil)
		return
	}

	messages, err := db.GetFromDB(timestamp)
	if err != nil {
		logger.Error("Failed to retrieve messages from database: " + err.Error())
		resp.SendResponse(w, http.StatusInternalServerError, "Failed to retrieve messages", err, nil)
		return
	}

	logger.Info("Messages retrieved successfully")
	resp.SendResponse(w, http.StatusOK, "Messages retrieved successfully", nil, messages)
}

func IsAliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.Error("Invalid request method")
		resp.SendResponse(w, http.StatusMethodNotAllowed, "Invalid request method", nil, nil)
		return
	}

	resp.SendResponse(w, http.StatusOK, "Node is alive", nil, nil)
}

// func TimeSignHandler(w http.ResponseWriter, r *http.Request) {
// 	//todo
// }
