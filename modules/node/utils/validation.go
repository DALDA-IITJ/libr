package utils

import (
	"encoding/json"
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

	message := map[string]interface{}{
		"sender":   msgCert["sender"],
		"msg":      msgCert["msg"],
		"ts":       msgCert["ts"],
		"mod_cert": msgCert["mod_cert"],
		"sign":     "",
	}

	// Convert message to string
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Error("Failed to marshal message to JSON: " + err.Error())
		return err
	}
	messageStr := string(messageBytes)

	logger.Debug(fmt.Sprintf("Verifying sender sign for public key: %s", sender))
	logger.Debug(fmt.Sprintf("Message string for verification: %s", messageStr))
	logger.Debug(fmt.Sprintf("Sender signature: %s", senderSign))

	isSigned, err := crypto.VerifySignature(messageStr, senderSign, sender)
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
	logger.Debug(fmt.Sprintf("Verifying module certificates with msgCert: %v", msgCert))

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

	logger.Debug(fmt.Sprintf("Type of mod_cert: %T", msgCert["mod_cert"]))
	rawModCerts, ok := msgCert["mod_cert"].([]interface{})
	if !ok {
		err := fmt.Errorf("mod_cert is missing or not of the expected type")
		logger.Error(err.Error())
		return err
	}

	modCerts := make([]map[string]interface{}, len(rawModCerts))
	for i, rawCert := range rawModCerts {
		modCert, ok := rawCert.(map[string]interface{})
		if !ok {
			err := fmt.Errorf("mod_cert[%d] is not of the expected type", i)
			logger.Error(err.Error())
		}
		modCerts[i] = modCert
	}

	logger.Debug(fmt.Sprintf("Verifying %d module certificates", len(modCerts)))

	for _, modCert := range modCerts {
		modPubKey := modCert["public_key"].(string)
		modSign := modCert["sign"].(string)

		message := map[string]interface{}{
			"msg":       msg,
			"timestamp": ts,
		}

		messageStr := ""
		for k, v := range message {
			messageStr += fmt.Sprintf("%s:%v;", k, v)
		}

		logger.Debug(fmt.Sprintf("Verifying module sign for public key: %s", modPubKey))
		logger.Debug(fmt.Sprintf("Message string for verification: %s", messageStr))
		logger.Debug(fmt.Sprintf("Module signature: %s", modSign))

		isSigned, err := crypto.VerifySignature(messageStr, modSign, modPubKey)
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
