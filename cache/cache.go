package cache

import (
	"errors"
	. "github.com/ONSdigital/dp-sessions-api/errors"
	"github.com/ONSdigital/dp-sessions-api/session"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

// Cache defines the structure required for a cache
type Cache struct {
	ttl      time.Duration
	interval time.Duration
	store    map[string]*session.Session
}

// NewCache creates a new cache
func NewCache(interval time.Duration, ttl time.Duration) *Cache {
	return &Cache{
		ttl:      ttl,
		interval: interval,
		store:    map[string]*session.Session{},
	}
}

// Set stores a session into the cache
// error could be returned when using redis
func (c *Cache) Set(s *session.Session) error {
	mutex.Lock()
	defer mutex.Unlock()
	c.store[s.ID] = s
	return nil
}


// GetByID retrieves a session from the cache by ID
func (c *Cache) GetByID(ID string) (*session.Session, error) {
	mutex.Lock()
	defer mutex.Unlock()

	findByID := func(s *session.Session) bool {
		return s.ID == ID
	}

	s := c.findSessionBy(findByID)
	if s == nil {
		return nil, SessionNotFound
	}

	if sinceAccessed := time.Since(s.LastAccessed); sinceAccessed >= c.ttl {
		return nil, SessionExpired
	}

	s.LastAccessed = time.Now()
	c.store[ID] = s

	return s, nil
}

func (c *Cache) DeleteAll() error {
	mutex.Lock()
	defer mutex.Unlock()

	if len(c.store) == 0 {
		return errors.New("no sessions to delete")
	}

	for k := range c.store {
		delete(c.store, k)
	}

	return nil
}

func (c *Cache) findSessionBy(filterFunc func(s *session.Session) bool) *session.Session {
	for _, sess := range c.store {
		if filterFunc(sess) {
			return sess
		}
	}
	return nil
}
