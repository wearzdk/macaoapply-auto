package encode

import (
	"crypto/rc4"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

// encrypt tool
var pass = "1234"
var key1 = "123"

func Rc4Encrypt(plaintext []byte, key []byte) []byte {
	// generate cipher.Block
	block, err := rc4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(plaintext))
	block.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func Encode(url, data string) string {
	// PBKDF2
	firstKey := pbkdf2.Key([]byte(pass), []byte(key1+url), 721, 32, sha256.New)

	// encrypt
	ciphertext := Rc4Encrypt([]byte(data), firstKey)
	return base64.StdEncoding.EncodeToString(ciphertext)
}
