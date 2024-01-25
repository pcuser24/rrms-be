package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserById(t *testing.T) {
	if testUser == nil {
		testUser = NewRandomUser(t, testRepo)
	}

	user, err := testRepo.GetUserById(
		context.Background(),
		testUser.ID,
	)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.ID, user.ID)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.CreatedAt, user.CreatedAt)
	require.Equal(t, testUser.UpdatedAt, user.UpdatedAt)
	require.Equal(t, testUser.CreatedBy, user.CreatedBy)
	require.Equal(t, testUser.UpdatedBy, user.UpdatedBy)
	require.Equal(t, testUser.DeletedF, user.DeletedF)
}

func TestGetUserByEmail(t *testing.T) {
	if testUser == nil {
		testUser = NewRandomUser(t, testRepo)
	}

	user, err := testRepo.GetUserByEmail(
		context.Background(),
		testUser.Email,
	)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.ID, user.ID)
	require.Equal(t, testUser.Email, user.Email)
	require.Equal(t, testUser.CreatedAt, user.CreatedAt)
	require.Equal(t, testUser.UpdatedAt, user.UpdatedAt)
	require.Equal(t, testUser.CreatedBy, user.CreatedBy)
	require.Equal(t, testUser.UpdatedBy, user.UpdatedBy)
	require.Equal(t, testUser.DeletedF, user.DeletedF)
}
