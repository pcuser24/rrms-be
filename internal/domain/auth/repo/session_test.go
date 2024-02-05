package repo

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/auth/dto"
	"github.com/user2410/rrms-backend/internal/domain/auth/model"
	"github.com/user2410/rrms-backend/internal/utils/random"
)

func newRandomSession(t *testing.T) *model.SessionModel {
	if testUser == nil {
		testUser = NewRandomUserDB(t, testRepo)
	}

	randUUID, err := uuid.NewRandom()
	require.NoError(t, err)

	arg := dto.CreateSession{
		ID:           randUUID,
		UserId:       testUser.ID,
		SessionToken: randUUID.String(),
		Expires:      time.Now(),
		UserAgent:    []byte(random.RandomAlphanumericStr(10)),
		ClientIp:     random.RandomIPv4Address(),
	}

	session, err := testRepo.CreateSession(
		context.Background(),
		&arg,
	)

	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.UserId, session.UserId)
	require.Equal(t, arg.SessionToken, session.SessionToken)
	require.WithinDuration(t, arg.Expires, session.Expires, time.Second)
	require.Equal(t, string(arg.UserAgent), *session.UserAgent)
	require.Equal(t, arg.ClientIp, *session.ClientIp)

	return session
}

func TestCreateSession(t *testing.T) {
	newRandomSession(t)
}

func TestGetSessionById(t *testing.T) {
	session := newRandomSession(t)

	session_1, err := testRepo.GetSessionById(
		context.Background(),
		session.ID,
	)

	require.NoError(t, err)
	require.NotEmpty(t, session_1)
	require.Equal(t, session.ID, session_1.ID)
	require.Equal(t, session.UserId, session_1.UserId)
	require.Equal(t, session.SessionToken, session_1.SessionToken)
	require.Equal(t, session.Expires, session_1.Expires)
	require.Equal(t, session.UserAgent, session_1.UserAgent)
	require.Equal(t, session.ClientIp, session_1.ClientIp)
	require.Equal(t, session.CreatedAt, session_1.CreatedAt)
}

func TestUpdateSessionStatus(t *testing.T) {
	session := newRandomSession(t)
	require.False(t, session.IsBlocked)

	err := testRepo.UpdateSessionStatus(
		context.Background(),
		session.ID,
		true,
	)
	require.NoError(t, err)
}
