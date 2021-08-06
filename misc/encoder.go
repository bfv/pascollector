package misc

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const Key = "1wLWsuc8jlcy5ThjcznWwSChDvH"

var iv = []byte{68, 105, 116, 73, 115, 87, 105, 108, 108, 101, 107, 101, 117, 114, 105, 103}

func Encrypt(text string) string {
	block, err := aes.NewCipher([]byte(getKey()))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func Decrypt(text string) string {

	block, err := aes.NewCipher([]byte(getKey()))
	if err != nil {
		panic(err)
	}
	ciphertext, _ := base64.StdEncoding.DecodeString(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	cfb.XORKeyStream(plaintext, ciphertext)
	return string(plaintext)
}

func getKey() string {
	return Key[5:21]
}
