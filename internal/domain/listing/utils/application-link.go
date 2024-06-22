package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	aes_util "github.com/user2410/rrms-backend/internal/utils/encryption/aes"
)

func stringify(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func parse(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func EncryptApplicationLink(secret string, data *dto.CreateApplicationLink) (string, error) {
	// encode data to string
	urlValues := url.Values{}
	urlValues.Add("listingId", data.ListingId.String())
	urlValues.Add("fullName", data.FullName)
	urlValues.Add("email", data.Email)
	urlValues.Add("phone", data.Phone)
	urlValues.Add("createdAt", time.Now().Format(time.RFC3339Nano))
	signData := urlValues.Encode()

	// encrypt data
	signedData, err := aes_util.Encrypt(secret, []byte(signData))
	if err != nil {
		return "", err
	}

	return stringify(signedData), nil
}

var (
	ErrMismatchField = errors.New("mismatch field")
	ErrExpired       = errors.New("expired")
	ErrInvalidKey    = errors.New("invalid key")
)

func VerifyApplicationLink(query *dto.VerifyApplicationLink, secret string) (bool, error) {
	cipherData, err := parse(query.Key)
	if err != nil {
		return false, err
	}
	decryptedData, err := aes_util.Decrypt(secret, cipherData)
	if err != nil {
		return false, err
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
