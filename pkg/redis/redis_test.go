package redis

import (
	"testing"

	"github.com/BangNopall/paskihub-be/internal/infra/env"
)

func TestNewRedisClientReturnsErrorWhenRedisUnavailable(t *testing.T) {
	originalHost := env.AppEnv.RedisHost
	originalPort := env.AppEnv.RedisPort
	originalPassword := env.AppEnv.RedisPassword
	t.Cleanup(func() {
		env.AppEnv.RedisHost = originalHost
		env.AppEnv.RedisPort = originalPort
		env.AppEnv.RedisPassword = originalPassword
	})

	env.AppEnv.RedisHost = "127.0.0.1"
	env.AppEnv.RedisPort = "1"
	env.AppEnv.RedisPassword = ""

	client, err := NewRedisClient()
	if err == nil {
		t.Fatal("expected Redis connection error")
	}
	if client != nil {
		t.Fatal("expected nil Redis client on connection failure")
	}
}
