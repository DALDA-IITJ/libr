package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/DALDA-IITJ/libr/modules/node/db"
	"github.com/DALDA-IITJ/libr/modules/node/msgs"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
	"github.com/DALDA-IITJ/libr/modules/node/utils/resp"
	"github.com/DALDA-IITJ/libr/modules/node/worker"
)

func SaveMsgHandler(w http.ResponseWriter, r *http.Request) {

	logger.Info("got req")

	body, err := io.ReadAll((r.Body))
	if err != nil {
		logger.Error("Failed to read request body: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Failed to read request body", err, nil)
		return
	}

	logger.Info(string(body))

	err = worker.AddToQueue(body)
	if err != nil {
		logger.Error("Failed to unwrap message certificate: " + err.Error())
		resp.SendResponse(w, http.StatusBadRequest, "Failed to unwrap message certificate", err, nil)
		return
	}

	logger.Info("Message certificate added to queue for processing")
	// resp.SendResponse(w, http.StatusOK, "Message certificate added to queue for processing", nil, nil)
}

func GetMsgHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	time := query.Get("time")
	if time == "" {
		logger.Error("Missing time parameter")
		resp.SendResponse(w, http.StatusBadRequest, "Missing time parameter", nil, nil)
		return
	}

	err := msgs.VerifyDBNode(time)
	if err != nil {
		logger.Error("Failed to verify if the server id correct DB node : " + err.Error())
		resp.SendResponse(w, http.StatusInternalServerError, "DB node is not the correct server", err, nil)
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
