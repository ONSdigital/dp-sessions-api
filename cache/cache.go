package cache

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/ONSdigital/log.go/log"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type Cache struct {
	ttl      time.Duration
	interval time.Duration
	store    map[string]*session.Session
}

func NewCache(interval time.Duration, ttl time.Duration) *Cache {
	return &Cache{
		ttl:      ttl,
		interval: interval,
		store:    map[string]*session.Session{},
	}
}

func (c *Cache) Set(s *session.Session) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Event(nil, "adding session to cache")
	c.store[s.ID] = s
}

func (c *Cache) GetByID(ID string) (*session.Session, error) {
	mutex.Lock()
	defer mutex.Unlock()

	findByID := func(s *session.Session) bool {
		return s.ID == ID
	}

	s := c.findSessionBy(findByID)
	if s == nil {
		return nil, nil
	}

	return s, nil
}

func (c *Cache) findSessionBy(filterFunc func(s *session.Session) bool) *session.Session {
	for _, sess := range c.store {
		if filterFunc(sess) {
			return sess
		}
	}
	return nil
}
