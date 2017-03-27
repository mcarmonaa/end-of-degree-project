package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mcarmonaa/end-of-degree-project/message"

	"golang.org/x/crypto/pbkdf2"
)

const ()

func main() {
	var mail string
	var pass string
	flag.StringVar(&mail, "u", "manu@mail.cool", "user to login")
	flag.StringVar(&pass, "p", "pass.123", "user's password")
	flag.Parse()

	nonce, err := generateNonce()
	if err != nil {
		log.Fatal(err)
	}

	sharedKey, err := generateSharedKey()
	if err != nil {
		log.Fatal(err)
	}

	loginMessage := &message.Login{SharedKey: sharedKey}
	loginJSON, err := json.Marshal(loginMessage)
	if err != nil {
		log.Fatal(err)
	}

	widespreadMessage := &message.Widespread{
		Timestamp: time.Now().Unix(),
		Nonce:     nonce,
		Content:   string(loginJSON),
	}

	widespreadJSON, err := json.Marshal(widespreadMessage)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(widespreadJSON))

	keyDF, err := getKDF(mail, pass)
	if err != nil {
		log.Fatal(err)
	}

	encryptedMessage, err := encryptLoginMessage(widespreadJSON, keyDF, []byte(mail))
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	if err := enc.Encode(encryptedMessage); err != nil {
		log.Fatal(err)
	}
}

func generateNonce() (uint64, error) {
	buf, err := randBytes(8)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(buf), nil
}

func randBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateSharedKey() (string, error) {
	const length = 32 // AES256-GCM needs a key length of 32 bytes (256bits)

	key, err := randBytes(length)
	if err != nil {
		return "", err
	}

	sharedKey := base64.StdEncoding.EncodeToString(key)
	return sharedKey, nil
}

func getKDF(mail, pass string) ([]byte, error) {
	const (
		url    = "http://localhost:8080/salt"
		length = 32
		iters  = 4096
	)

	query := "?mail=" + mail
	resp, err := http.Get(url + query)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response status: %q", resp.Status)
	}

	salt, err := getSalt(resp.Body)
	if err != nil {
		return nil, err
	}

	derivatedKey := pbkdf2.Key([]byte(pass), salt, iters, length, sha256.New)
	return derivatedKey, nil
}

func getSalt(in io.Reader) ([]byte, error) {
	user := &message.AuthUser{}
	dec := json.NewDecoder(in)
	if err := dec.Decode(user); err != nil {
		return nil, err
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func encryptLoginMessage(data, key, authData []byte) (*message.Encrypted, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv, err := randBytes(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}

	payload := aesgcm.Seal(nil, iv, data, authData)

	message := &message.Encrypted{
		Mail:    string(authData),
		IVector: base64.StdEncoding.EncodeToString(iv),
		Payload: base64.StdEncoding.EncodeToString(payload),
	}

	return message, nil
}
