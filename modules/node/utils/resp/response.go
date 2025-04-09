package resp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

// SendResponse sends a standardized JSON response to the client.
func SendResponse(w http.ResponseWriter, statusCode int, message string, err error, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Construct the response payload
	resp := map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"error":   nil,
		"data":    nil,
	}

	if err != nil {
		resp["error"] = err.Error()
	}
	if data != nil {
		resp["data"] = data
	}

	// Encode the response as JSON and handle encoding errors
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.Error(fmt.Sprintf("Failed to encode response: %v", encodeErr))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("Response sent: %v", resp))
}
