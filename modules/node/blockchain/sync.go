package blockchain

import "node/utils/logger"

type Blockchain struct {
	State     string
	Timestamp int
}

var BlockchainState *Blockchain

// todo
func Sync(time int) error {
	var err error
	BlockchainState, err = GetBlockchainStateAt(time)
	if err != nil {
		logger.Error("Failed to get blockchain state: " + err.Error())
		return err
	}

	logger.Info("Blockchain state synchronized successfully")
	return nil
}
