package cache

//go:generate moq -out mock_redisclienter.go . RedisClienter

import (
	"time"

	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
)

// SessionCache interface for storing and retrieving sessions
type SessionCache interface {
	SetSession(s *session.Session) error
	GetByID(ID string) (*session.Session, error)
	GetByEmail(email string) (*session.Session, error)
	DeleteAll() error
}

// RedisClienter - interface for redis
type RedisClienter interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	FlushAll() *redis.StatusCmd
	Ping() *redis.StatusCmd
}
