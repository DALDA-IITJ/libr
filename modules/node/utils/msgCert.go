package utils

import (
	"encoding/json"
	"fmt"

	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

var MsgCert = make(map[string]interface{})
var ModCerts = make([]map[string]interface{}, 0)

func UnWrapMsgCert(msgCert []byte) error {
	err := json.Unmarshal(msgCert, &MsgCert)
	if err != nil {
		logger.Error("Failed to unmarshal message certificate: " + err.Error())
		return err
	}

	err = UnWrapModCerts()
	if err != nil {
		logger.Error("Failed to unwrap module certificates: " + err.Error())
		return err
	}

	logger.Info("Message certificate unmarshalled successfully")
	logger.Debug("Message certificate: " + fmt.Sprintf("%v", MsgCert))

	return nil
}

func UnWrapModCerts() error {
	modCerts, ok := MsgCert["modCerts"].([]interface{})
	if !ok {
		err := fmt.Errorf("modCerts key is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	for _, modCert := range modCerts {
		modCertMap, ok := modCert.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("one of the modCerts entries is not of the expected type")
			logger.Error(err.Error())
			return err
		}
		ModCerts = append(ModCerts, modCertMap)
	}

	logger.Info("Module certificates unmarshalled successfully")
	logger.Debug("Module certificates: " + fmt.Sprintf("%v", ModCerts))
	return nil
}
