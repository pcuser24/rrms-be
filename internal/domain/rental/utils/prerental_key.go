package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"time"

	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	aes_util "github.com/user2410/rrms-backend/internal/utils/encryption/aes"
)

func stringify(b []byte) string {
	return fmt.Sprintf("%x", b)
}

func parse(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func CreatePreRentalKey(secret string, prerental *rental_model.PreRental) (string, error) {
	urlValues := url.Values{}
	urlValues.Add("id", fmt.Sprintf("%d", prerental.ID))
	urlValues.Add("creatorId", prerental.CreatorID.String())
	urlValues.Add("propertyId", prerental.PropertyID.String())
	urlValues.Add("unitId", prerental.UnitID.String())
	urlValues.Add("tenantName", prerental.TenantName)
	urlValues.Add("tenantEmail", prerental.TenantEmail)
	urlValues.Add("tenantPhone", prerental.TenantPhone)
	urlValues.Add("createdAt", prerental.CreatedAt.Format(time.RFC3339Nano))
	signData := urlValues.Encode()

	signedData, err := aes_util.Encrypt(secret, []byte(signData))
	if err != nil {
		return "", err
	}

	return stringify(signedData), nil
}

var ErrMismatchField = errors.New("mismatch field")

func VerifyPreRentalKey(prerental *rental_model.PreRental, key, secret string) error {
	cipherData, err := parse(key)
	if err != nil {
		return err
	}
	decryptedData, err := aes_util.Decrypt(secret, cipherData)
	if err != nil {
		return err
	}

	// validate the data
	params, err := url.ParseQuery(string(decryptedData))
	if err != nil {
		return err
	}
	if params.Get("id") != fmt.Sprintf("%d", prerental.ID) ||
		params.Get("creatorId") != prerental.CreatorID.String() ||
		params.Get("propertyId") != prerental.PropertyID.String() ||
		params.Get("unitId") != prerental.UnitID.String() ||
		params.Get("tenantName") != prerental.TenantName ||
		params.Get("tenantEmail") != prerental.TenantEmail ||
		params.Get("tenantPhone") != prerental.TenantPhone ||
		params.Get("createdAt") != prerental.CreatedAt.Format(time.RFC3339Nano) {
		return ErrMismatchField
	}

	return nil
}
