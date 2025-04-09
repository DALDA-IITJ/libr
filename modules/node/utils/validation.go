package utils

import (
	"fmt"

	"github.com/DALDA-IITJ/libr/modules/core/crypto"
	"github.com/DALDA-IITJ/libr/modules/node/utils/logger"
)

func VerifyMsgCert(ts string, msgCert map[string]interface{}) (bool, error) {

	errChan := make(chan error, 2) // Buffered channel to collect errors

	go func() {
		err := senderSignCheck(msgCert)
		if err != nil {
			logger.Error("Failed to verify sender sign: " + err.Error())
			errChan <- err
			return
		}
		errChan <- nil // Signal success
	}()

	go func() {
		err := verifyModCerts(msgCert)
		if err != nil {
			logger.Error("Failed to verify module certs: " + err.Error())
			errChan <- err
			return
		}
		errChan <- nil // Signal success
	}()

	// Collect errors from goroutines
	var errs []error
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		logger.Error("Verification failed with errors: " + fmt.Sprintf("%v", errs))
		// Combine errors into a single error (you might want a more sophisticated error aggregation)
		return false, errs[0]
	}

	logger.Info("Message certificate verified successfully")
	return true, nil
}

func senderSignCheck(msgCert map[string]interface{}) error {
	senderSign, ok := msgCert["sign"].(string)
	if !ok {
		err := fmt.Errorf("sender sign is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	sender, ok := msgCert["sender"].(string)
	if !ok {
		err := fmt.Errorf("sender key is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	message := make(map[string]interface{})
	for k, v := range msgCert {
		if k == "sign" || k == "sender" {
			continue
		}
		message[k] = v
	}

	// Convert message to string
	messageStr := ""
	for k, v := range message {
		messageStr += fmt.Sprintf("%s:%v;", k, v)
	}

	isSigned, err := crypto.VerifySignature(sender, messageStr, senderSign)
	if err != nil {
		logger.Error("Failed to verify sender sign: " + err.Error())
		return err
	}

	if !isSigned {
		err := fmt.Errorf("sender sign verification failed")
		logger.Error(err.Error())
		return err
	}

	logger.Info("Sender sign verified successfully")
	return nil
}

func verifyModCerts(msgCert map[string]interface{}) error {
	msg, ok := msgCert["msg"].(string)
	if !ok {
		err := fmt.Errorf("msg is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	ts, ok := msgCert["ts"].(string)
	if !ok {
		err := fmt.Errorf("ts is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	modCerts, ok := msgCert["mod_cert"].([]map[string]string)
	if !ok {
		err := fmt.Errorf("mod_cert is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	for _, modCert := range modCerts {
		modPubKey := modCert["public_key"]
		modSign := modCert["sign"]

		message := map[string]interface{}{
			"msg":       msg,
			"timestamp": ts,
		}

		messageStr := ""
		for k, v := range message {
			messageStr += fmt.Sprintf("%s:%v;", k, v)
		}

		isSigned, err := crypto.VerifySignature(modPubKey, messageStr, modSign)
		if err != nil {
			logger.Error("Failed to verify module sign: " + err.Error())
			return err
		}

		if !isSigned {
			err := fmt.Errorf("mod sign verification failed")
			logger.Error(err.Error())
			return err
		}

		logger.Info("Mod sign verified successfully")
	}

	return nil
}
