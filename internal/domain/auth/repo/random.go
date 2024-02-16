package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/utils"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

var testUser *model.UserModel = nil

func NewRandomUserDB(t *testing.T, testRepo Repo) *model.UserModel {
	hashedPassword, err := utils.HashPassword(random.RandomAlphanumericStr(10))
	require.NoError(t, err)

	arg := dto.RegisterUser{
		Email:     random.RandomEmail(),
		Password:  hashedPassword,
		FirstName: random.RandomAlphanumericStr(10),
		LastName:  random.RandomAlphanumericStr(10),
	}

	user, err := testRepo.CreateUser(
		context.Background(),
		&arg,
	)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, *user.Password)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.NotZero(t, user.ID)
	require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), time.Second)

	return user
}

func compareUsers(t *testing.T, expected, actual *model.UserModel) {
	require.NotEmpty(t, expected)
	require.NotEmpty(t, actual)
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.Email, actual.Email)
	require.Equal(t, expected.Password, actual.Password)
	require.Equal(t, expected.FirstName, actual.FirstName)
	require.Equal(t, expected.LastName, actual.LastName)
	require.Equal(t, expected.Phone, actual.Phone)
	require.Equal(t, expected.Avatar, actual.Avatar)
	require.Equal(t, expected.Address, actual.Address)
	require.Equal(t, expected.City, actual.City)
	require.Equal(t, expected.District, actual.District)
	require.Equal(t, expected.Ward, actual.Ward)
}

func NewRandomUserModel(t *testing.T) *model.UserModel {
	hashedPassword, err := utils.HashPassword(random.RandomAlphanumericStr(10))
	require.NoError(t, err)

	return &model.UserModel{
		ID:        uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
		Email:     random.RandomEmail(),
		Password:  types.Ptr[string](hashedPassword),
		GroupID:   uuid.Nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.Nil,
		UpdatedBy: uuid.Nil,
		DeletedF:  false,
		FirstName: random.RandomAlphanumericStr(10),
		LastName:  random.RandomAlphanumericStr(10),
		Phone:     types.Ptr[string](random.RandomNumericStr(10)),
		Avatar:    types.Ptr[string](random.RandomAlphanumericStr(10)),
		Address:   types.Ptr[string](random.RandomAlphanumericStr(10)),
		City:      types.Ptr[string](random.RandomAlphanumericStr(10)),
		District:  types.Ptr[string](random.RandomAlphanumericStr(10)),
		Ward:      types.Ptr[string](random.RandomAlphanumericStr(10)),
	}
}
