package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/utils/random"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestInsertUser(t *testing.T) {
	NewRandomUserDB(t, testRepo)
}

func TestGetUserById(t *testing.T) {
	if testUser == nil {
		testUser = NewRandomUserDB(t, testRepo)
	}

	user, err := testRepo.GetUserById(context.Background(), testUser.ID)
	require.NoError(t, err)
	compareUsers(t, testUser, user)
}

func TestGetUserByEmail(t *testing.T) {
	if testUser == nil {
		testUser = NewRandomUserDB(t, testRepo)
	}

	user, err := testRepo.GetUserByEmail(
		context.Background(),
		testUser.Email,
	)
	require.NoError(t, err)
	compareUsers(t, testUser, user)
}

func TestUpdateUser(t *testing.T) {
	user := NewRandomUserDB(t, testRepo)

	arg := dto.UpdateUser{
		UpdatedBy: user.ID,
		Email:     types.Ptr[string](random.RandomEmail()),
		Password:  types.Ptr[string](random.RandomAlphabetStr(15)),
		FirstName: types.Ptr[string](random.RandomAlphabetStr(15)),
		LastName:  types.Ptr[string](random.RandomAlphabetStr(15)),
		Phone:     types.Ptr[string](random.RandomNumericStr(10)),
		Address:   types.Ptr[string](random.RandomAlphabetStr(20)),
		City:      types.Ptr[string](random.RandomAlphabetStr(10)),
		District:  types.Ptr[string](random.RandomAlphabetStr(10)),
		Ward:      types.Ptr[string](random.RandomAlphabetStr(10)),
	}

	user.Email = *arg.Email
	user.Password = arg.Password
	user.FirstName = *arg.FirstName
	user.LastName = *arg.LastName
	user.Phone = arg.Phone
	user.Address = arg.Address
	user.City = arg.City
	user.District = arg.District
	user.Ward = arg.Ward

	err := testRepo.UpdateUser(context.Background(), user.ID, &arg)
	require.NoError(t, err)

	_user, err := testRepo.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)

	compareUsers(t, user, _user)
}
