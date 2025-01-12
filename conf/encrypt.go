package conf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var cipherKey = []byte("AmH0lGOt7S07N0QrUwgMKjXNC0dxcJPZ")

func encryptAES(text string) (string, error) {
	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	out := gcm.Seal(nonce, nonce, []byte(text), nil)
	str := base64.StdEncoding.EncodeToString(out)
	return str, nil
}

func decryptAES(b64 string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ct) < nonceSize {
		fmt.Println(err)
	}

	nonce, ct := ct[:nonceSize], ct[nonceSize:]
	text, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(text), nil
}
