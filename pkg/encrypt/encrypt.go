package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	mrand "math/rand"

	"time"
)

//GenerateRandom: generate a "random" string of 6 alphanumeric charcaters
func GenerateRandom() string {
	mrand.Seed(time.Now().UnixNano())
	var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = characters[mrand.Intn(len(characters))]
	}
	return string(b)
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

//Encode: base64 encoding of string
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

//Decode: base64 decoding
func Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data, err
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, Secret string) (string, error) {
	block, err := aes.NewCipher([]byte(createHash(Secret)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(text), nil)
	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, Secret string) (string, error) {
	passphrase := []byte(createHash(Secret))
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	cipherText, err := Decode(text)
	if err != nil {
		return "", err
	}
	data := []byte(cipherText)
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext), nil
}
