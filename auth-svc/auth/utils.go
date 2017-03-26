package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"

	"fmt"

	"encoding/json"

	"github.com/jinzhu/gorm"
	"github.com/mcarmonaa/eodp/message"
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

func decryptLogin(initVector, payload string, user *User) (*message.Login, error) {
	iv, err := base64.StdEncoding.DecodeString(initVector)
	if err != nil {
		return nil, fmt.Errorf("login: couldn't decode base64 iv: %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("login: couldn't decode base64 payload: %v", err)
	}

	pass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		return nil, fmt.Errorf("login: couldn't decode base64 password: %v", err)
	}

	buf, err := decryptPayload(iv, data, pass, []byte(user.Mail))
	if err != nil {
		return nil, fmt.Errorf("login: couldn't decrypt payload: %v", err)
	}

	fmt.Println(string(buf))

	message := &message.Login{}
	if err := json.Unmarshal(buf, message); err != nil {
		return nil, fmt.Errorf("login: couldn't unmarshal login message: %v", err)
	}

	return message, nil
}

func decryptPayload(iv, data, key, authData []byte) ([]byte, error) {
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
