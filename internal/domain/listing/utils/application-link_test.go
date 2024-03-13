package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/utils/random"
)

func TestApplicationLink(t *testing.T) {
	data := dto.CreateApplicationLink{
		ListingId: uuid.MustParse("978fe220-663f-464c-9afd-fe05de7be44b"),
		FullName:  "John Doe",
		Email:     "abc@email.com",
		Phone:     "1234567890",
	}
	secret := random.RandomAlphabetStr(32)
	encrypted, err := EncryptApplicationLink(secret, &data)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)
	t.Log(encrypted)

	res, err := VerifyApplicationLink(&dto.VerifyApplicationLink{
		CreateApplicationLink: data,
		Key:                   encrypted,
	}, secret)
	require.NoError(t, err)
	require.True(t, res)

	res, err = VerifyApplicationLink(&dto.VerifyApplicationLink{
		CreateApplicationLink: data,
		Key:                   encrypted[1 : len(encrypted)-1],
	}, secret)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidKey)
	require.False(t, res)

	res, err = VerifyApplicationLink(&dto.VerifyApplicationLink{
		CreateApplicationLink: dto.CreateApplicationLink{
			ListingId: uuid.MustParse("d2099b7d-c72f-4c11-aa64-630b836d750f"),
			FullName:  "John Doe",
			Email:     "abc@email.com",
			Phone:     "1234567890",
		},
		Key: encrypted,
	}, secret)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMismatchField)
	require.False(t, res)

	res, err = VerifyApplicationLink(&dto.VerifyApplicationLink{
		CreateApplicationLink: dto.CreateApplicationLink{
			ListingId: uuid.MustParse("978fe220-663f-464c-9afd-fe05de7be44b"),
			FullName:  "John Does",
			Email:     "abc@email.com",
			Phone:     "1234567890",
		},
		Key: encrypted,
	}, secret)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMismatchField)
	require.False(t, res)

}
