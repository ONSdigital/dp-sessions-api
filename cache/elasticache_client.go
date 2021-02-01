package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
)

var (
	HealthyMessage       = "elasticache is OK"
	ErrEmptySessionID    = errors.New("session id required but was empty")
	ErrEmptySessionEmail = errors.New("session email required but was empty")
	ErrEmptySession      = errors.New("session is empty")
	ErrEmptyAddress      = errors.New("address is empty")
	ErrEmptyPassword     = errors.New("password is empty")
	ErrInvalidTTL        = errors.New("ttl should not be zero")
)

type ElasticacheClient struct {
	client RedisClienter
	ttl    time.Duration
}

// Config - config options for the elasticache client
type Config struct {
	Addr     string
	Password string `json:"-"`
	Database int
	TTL      time.Duration
	TLS      *tls.Config
}

// New - create new session cache client instance
func New(c Config) (*ElasticacheClient, error) {
	if c.Addr == "" {
		return nil, ErrEmptyAddress
	}

	if c.Password == "" {
		return nil, ErrEmptyPassword
	}

	if c.TTL == 0 {
		return nil, ErrInvalidTTL
	}

	return &ElasticacheClient{
		client: redis.NewClient(&redis.Options{
			Addr:      c.Addr,
			Password:  c.Password,
			DB:        c.Database,
			TLSConfig: c.TLS,
		}),
		ttl: c.TTL,
	}, nil
}

// SetSession - add session to elasticache
func (c *ElasticacheClient) SetSession(s *session.Session) error {
	if s == nil {
		return ErrEmptySession
	}

	sJSON, err := s.MarshalJSON()
	if err != nil {
		return err
	}

	// Add session using ID as key
	err = c.client.Set(s.ID, sJSON, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("elasticache client.Set returned an unexpected error: %w", err)
	}

	// Add session using email as key
	err = c.client.Set(s.Email, sJSON, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("elasticache client.Set returned an unexpected error: %w", err)
	}

	return nil
}

// GetByID - gets a session from elasticache using its ID
func (c *ElasticacheClient) GetByID(id string) (*session.Session, error) {
	if id == "" {
		return nil, ErrEmptySessionID
	}

	msg, err := c.client.Get(id).Result()
	if err != nil {
		return nil, err
	}

	var s *session.Session

	err = json.Unmarshal([]byte(msg), &s)
	if err != nil {
		return nil, err
	}

	// Refresh TTL on access and update LastAccessed in session
	s.LastAccessed = time.Now()
	err = c.Expire(s.ID, c.ttl)
	if err != nil {
		return nil, err
	}

	err = c.Expire(s.Email, c.ttl)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetByEmail - gets a session from elasticache using its ID
func (c *ElasticacheClient) GetByEmail(email string) (*session.Session, error) {
	if email == "" {
		return nil, ErrEmptySessionEmail
	}

	msg, err := c.client.Get(email).Result()
	if err != nil {
		return nil, err
	}

	var s *session.Session

	err = json.Unmarshal([]byte(msg), &s)
	if err != nil {
		return nil, err
	}

	// Refresh TTL on access and update LastAccessed in session
	s.LastAccessed = time.Now()
	err = c.Expire(s.Email, c.ttl)
	if err != nil {
		return nil, err
	}

	err = c.Expire(s.ID, c.ttl)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteAll - removes all items from elasticache
func (c *ElasticacheClient) DeleteAll() error {
	return c.client.FlushAll().Err()
}

// Ping - checks the connection to elasticache
func (c *ElasticacheClient) Ping() error {
	return c.client.Ping().Err()
}

// Expire - sets the expiration of key
func (c *ElasticacheClient) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(key, expiration).Err()
}

func (c *ElasticacheClient) Checker(ctx context.Context, state *health.CheckState) error {
	err := c.Ping()
	if err != nil {
		// Generic error
		return state.Update(health.StatusCritical, err.Error(), 0)
	}
	// Success
	return state.Update(health.StatusOK, HealthyMessage, 0)
}
