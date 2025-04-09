package worker

import (
	"encoding/json"
	"strconv"

	"github.com/DALDA-IITJ/libr/modules/node/msgs"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

// Message queue for processing (buffer size can be adjusted)
var msgProcessingQueue = make(chan map[string]interface{}, 100)

// StartWorkerPool initializes workers that process messages from the queue
func StartWorkerPool(workerCount int) {
	for i := 0; i < workerCount; i++ {
		logger.Info("Starting worker #" + strconv.Itoa(i+1))
		go worker() // Each worker runs in a goroutine
	}
}

// worker continuously processes messages from the queue
func worker() {
	for msgCert := range msgProcessingQueue {
		err := msgs.ProcessMsg(msgCert)
		if err != nil {
			logger.Error("Failed to process message: " + err.Error())
		}
	}
}

// AddToQueue safely pushes messages to the queue for processing
func AddToQueue(msg []byte) error {
	var msgCert map[string]interface{}
	err := json.Unmarshal(msg, &msgCert)
	if err != nil {
		logger.Error("Failed to unmarshal message: " + err.Error())
		return err
	}

	msgProcessingQueue <- msgCert // This triggers workers to wake up if idle
	logger.Info("Message added to queue for processing")
	return nil

}
