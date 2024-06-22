package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func Encrypt(secret string, data []byte) ([]byte, error) {
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	signedData := gcm.Seal(nonce, nonce, data, nil)

	return signedData, nil
}

var (
	ErrInvalidKey = errors.New("invalid key")
)

func Decrypt(secret string, cipherData []byte) ([]byte, error) {
	// decrypt the hash
	// initialize aes cipher
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if err != nil {
		return nil, err
	}
	if len(cipherData) < nonceSize {
		return nil, errors.Join(ErrInvalidKey, err)
	}
	nonce, ciphertext := cipherData[:nonceSize], cipherData[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
