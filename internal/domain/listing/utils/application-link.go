package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
)

func stringify(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func parse(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func EncryptApplicationLink(secret string, data *dto.CreateApplicationLink) (string, error) {
	// initialize aes cipher
	c, err := aes.NewCipher([]byte(secret))
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

	// encode data to string
	urlValues := url.Values{}
	urlValues.Add("listingId", data.ListingId.String())
	urlValues.Add("fullName", data.FullName)
	urlValues.Add("email", data.Email)
	urlValues.Add("phone", data.Phone)
	urlValues.Add("createdAt", time.Now().Format(time.RFC3339Nano))
	signData := urlValues.Encode()

	// encrypt data
	signedData := gcm.Seal(nonce, nonce, []byte(signData), nil)

	return stringify(signedData), nil
}

var (
	ErrMismatchField = errors.New("mismatch field")
	ErrExpired       = errors.New("expired")
	ErrInvalidKey    = errors.New("invalid key")
)

func VerifyApplicationLink(query *dto.VerifyApplicationLink, secret string) (bool, error) {
	// decrypt the hash
	// initialize aes cipher
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return false, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return false, err
	}
	nonceSize := gcm.NonceSize()
	ciphertext, err := parse(query.Key)
	if err != nil {
		return false, err
	}
	if len(ciphertext) < nonceSize {
		return false, errors.Join(ErrInvalidKey, err)
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, errors.Join(ErrInvalidKey, err)
	}

	// validate the data
	params, err := url.ParseQuery(string(decryptedData))
	if err != nil {
		return false, err
	}
	if params.Get("listingId") != query.ListingId.String() ||
		params.Get("fullName") != query.FullName ||
		params.Get("email") != query.Email ||
		params.Get("phone") != query.Phone {
		return false, ErrMismatchField
	}
	t, err := time.Parse(time.RFC3339Nano, params.Get("createdAt"))
	if err != nil {
		return false, err
	}
	if time.Since(t) > 30*24*time.Hour {
		return false, ErrExpired
	}

	return true, nil
}
