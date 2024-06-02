package redisd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

var (
	clusterClient *redis.Client
)

func TestMain(m *testing.M) {
	opt := redis.Options{
		Addr:     "localhost:36379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	clusterClient = redis.NewClient(&opt)

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res := clusterClient.Ping(ctx)
	require.NoError(t, res.Err())
}
