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
		Email:    random.RandomEmail(),
		Password: hashedPassword,
	}

	user, err := testRepo.InsertUser(
		context.Background(),
		&arg,
	)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.ID)
	require.WithinDuration(t, user.CreatedAt, time.Now(), time.Second)
	require.WithinDuration(t, user.UpdatedAt, time.Now(), time.Second)

	return user
}

func NewRandomUserModel(t *testing.T) *model.UserModel {
	hashedPassword, err := utils.HashPassword(random.RandomAlphanumericStr(10))
	require.NoError(t, err)

	id, err := uuid.NewRandom()
	require.NoError(t, err)

	return &model.UserModel{
		ID:        id,
		Email:     random.RandomEmail(),
		Password:  types.Ptr[string](hashedPassword),
		GroupID:   uuid.Nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
