package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mcarmonaa/end-of-degree-project/auth-svc/nonces"
	"github.com/mcarmonaa/end-of-degree-project/message"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/pbkdf2"
)

func inDB(db *gorm.DB, mail string) bool {
	var count int
	db.Model(&User{}).Where(&User{Mail: mail}).Count(&count)
	return count > 0
}

func randBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateSalt() (string, error) {
	const saltLen = 32
	salt, err := randBytes(saltLen)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}

func generatePassword(pass, salt string) (string, error) {
	const (
		length = 32
		iters  = 4096
	)

	rawSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", fmt.Errorf("decoding base64 password: %v", err)
	}

	derivatedKey := pbkdf2.Key([]byte(pass), rawSalt, iters, length, sha256.New)
	derivatedKey64 := base64.StdEncoding.EncodeToString(derivatedKey)
	return derivatedKey64, nil
}

func decryptAndValidate(initVector, payload, key, authData string, blackList *nonces.Nonces) (*message.Widespread, error) {
	message, err := decryptPayload(initVector, payload, key, authData)
	if err != nil {
		return nil, err
	}

	if err := validateMessage(message, blackList); err != nil {
		return nil, err
	}

	return message, nil
}

func decryptPayload(ivBase64, payloadBase64, keyBase64, authData string) (*message.Widespread, error) {
	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, err
	}

	buf, err := decryptData(iv, data, key, []byte(authData))
	if err != nil {
		return nil, err
	}

	message := &message.Widespread{}
	if err := json.Unmarshal(buf, message); err != nil {
		return nil, err
	}

	return message, nil
}

func decryptData(iv, data, key, authData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decryptedData, err := aesgcm.Open(nil, iv, data, authData)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func validateMessage(message *message.Widespread, blackList *nonces.Nonces) error {
	const windowTime = 15
	elapsedTS := time.Now().Unix() - message.Timestamp

	if elapsedTS < 0 || elapsedTS > windowTime {
		return fmt.Errorf("invalid timestamp")
	}

	if !blackList.CheckAndAdd(message.Nonce) {
		return fmt.Errorf("invalid nonce")
	}

	return nil
}
