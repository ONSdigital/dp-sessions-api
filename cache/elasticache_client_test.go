package cache

import (
	"errors"
	"testing"
	"time"

	"github.com/ONSdigital/dp-sessions-api/session"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testTTL          = 30 * time.Minute
	respLastAccessed = "2020-08-13T08:40:18.652Z"
)

var (
	resp = []byte(`{"id":"1234","email":"user@email.com","start":"2020-08-13T08:40:18.652Z","last_accessed":"2020-08-13T08:40:18.652Z"}`)
)

func TestNewClient(t *testing.T) {
	Convey("Given NewClient returns new redis client", t, func() {

		Convey("When correct redis configuration is provided", func() {
			c, err := New(Config{
				Addr:     "123.0.0.1",
				Password: "1234",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then a new redis client will be returned with no error", func() {
				So(err, ShouldBeNil)
				So(c, ShouldNotBeEmpty)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations address is empty", func() {
			c, err := New(Config{
				Addr:     "",
				Password: "1234",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then the client will not be created and the empty address error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptyAddress)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations password is empty", func() {
			c, err := New(Config{
				Addr:     "123.0.0.1",
				Password: "",
				Database: 0,
				TTL:      testTTL,
			})

			Convey("Then the client will not be created and the empty password error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptyPassword)
			})
		})

	})

	Convey("Given NewClient returns an error", t, func() {

		Convey("When the redis configurations ttl is zero", func() {
			c, err := New(Config{
				Addr:     "123.0.0.1",
				Password: "1234",
				Database: 0,
				TTL:      0,
			})

			Convey("Then the client will not be created and the invalid ttl error is returned", func() {
				So(c, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrInvalidTTL)
			})
		})
	})
}

func TestClient_Set(t *testing.T) {
	Convey("Given a valid sessions and redis client.Set returns no error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusResult("success", nil),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When there is a valid session", func() {
			s := &session.Session{
				ID:           "1234",
				Email:        "user@email.com",
				Start:        time.Now(),
				LastAccessed: time.Now(),
			}

			jsonByes, err := s.MarshalJSON()
			So(err, ShouldBeNil)

			err = client.SetSession(s)

			Convey("Then the session is stored in the cache and no error is returned", func() {
				So(err, ShouldBeNil)
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 2)

				So(mockRedisClient.SetCalls()[0].Key, ShouldEqual, s.ID)
				So(mockRedisClient.SetCalls()[0].Value, ShouldResemble, jsonByes)
				So(mockRedisClient.SetCalls()[0].Expiration, ShouldEqual, testTTL)
			})
		})
	})

	Convey("Given a valid session and redis client.Set returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusResult("fail", errors.New("failed to store session")),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When there is a valid session but redis client.Set errors ", func() {
			s := &session.Session{
				ID:           "1234",
				Email:        "user@email.com",
				Start:        time.Now(),
				LastAccessed: time.Now(),
			}

			jsonByes, err := s.MarshalJSON()
			So(err, ShouldBeNil)

			err = client.SetSession(s)

			Convey("Then the session will not be stored in the cache and an error is returned", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.SetCalls()[0].Key, ShouldEqual, s.ID)
				So(mockRedisClient.SetCalls()[0].Value, ShouldResemble, jsonByes)
				So(mockRedisClient.SetCalls()[0].Expiration, ShouldEqual, testTTL)

				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "redis client.Set returned an unexpected error: failed to store session")
			})
		})
	})

	Convey("Given an invalid session and redis client.Set returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When there is an invalid session", func() {
			var s *session.Session
			err := client.SetSession(s)

			Convey("Then the session will not be stored in the cache and an error is returned", func() {
				So(mockRedisClient.SetCalls(), ShouldHaveLength, 0)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySession)
			})
		})
	})
}

func TestClient_GetByID(t *testing.T) {
	Convey("Given a session ID client.GetByID returns a session and TTL is refreshed", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult(string(resp), nil),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client uses the ID to get the session", func() {
			s, err := client.GetByID("1234")
			So(err, ShouldBeNil)

			Convey("Then redis client.Get is called with the expected parameters", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.GetCalls()[0].Key, ShouldEqual, "1234")

				So(mockRedisClient.ExpireCalls(), ShouldHaveLength, 2) // Expects 2 due to refreshing by ID and Email

				So(mockRedisClient.ExpireCalls()[0].Key, ShouldEqual, "1234")
				So(mockRedisClient.ExpireCalls()[0].Expiration, ShouldEqual, testTTL)

				So(mockRedisClient.ExpireCalls()[1].Key, ShouldEqual, "user@email.com")
				So(mockRedisClient.ExpireCalls()[1].Expiration, ShouldEqual, testTTL)
			})

			Convey("And the expected session is returned", func() {
				So(s, ShouldNotBeEmpty)
				So(s.ID, ShouldEqual, "1234")
				So(s.LastAccessed.String(), ShouldNotEqual, respLastAccessed)
				So(mockRedisClient.ExpireCalls()[0].Expiration, ShouldEqual, testTTL)
			})
		})
	})

	Convey("Given a session ID client.GetByID returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult(string(resp), nil),
			*redis.NewStatusCmd(),
			*redis.NewBoolResult(false, errors.New("unable to refresh expiration")),
		)

		Convey("When client uses the ID to get the session", func() {
			s, err := client.GetByID("1234")

			Convey("Then redis client.Get is called with the expected parameters", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.GetCalls()[0].Key, ShouldEqual, "1234")
				So(mockRedisClient.ExpireCalls(), ShouldHaveLength, 1)
			})

			Convey("And the expected error is returned", func() {
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unable to refresh expiration")
				So(s, ShouldBeNil)
			})
		})
	})

	Convey("Given a blank session ID client.GetByID returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client.GetByID is called has an empty ID", func() {
			s, err := client.GetByID("")

			Convey("Then client.GetByID returns an error and no session is returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 0)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySessionID)
			})
		})
	})

	Convey("Given a session ID client.GetByID returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult("", errors.New("unexpected end of JSON input")),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client.GetByID is called with a valid session ID", func() {
			s, err := client.GetByID("1234")

			Convey("Then the redis client.Get returns an error and no session is returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unexpected end of JSON input")
			})
		})
	})
}

func TestClient_GetByEmail(t *testing.T) {
	Convey("Given a session email client.GetByEmail returns a session and TTL is refreshed", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult(string(resp), nil),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client uses the email to get the session", func() {
			s, err := client.GetByEmail("user@email.com")
			So(err, ShouldBeNil)

			Convey("Then redis client.Get is called with the expected parameters", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.GetCalls()[0].Key, ShouldEqual, "user@email.com")

				So(mockRedisClient.ExpireCalls(), ShouldHaveLength, 2) // Expects 2 due to refreshing by ID and Email
				So(mockRedisClient.ExpireCalls()[0].Key, ShouldEqual, "user@email.com")
				So(mockRedisClient.ExpireCalls()[0].Expiration, ShouldEqual, testTTL)

				So(mockRedisClient.ExpireCalls()[1].Key, ShouldEqual, "1234")
				So(mockRedisClient.ExpireCalls()[1].Expiration, ShouldEqual, testTTL)
			})

			Convey("And the expected session is returned", func() {

				So(s, ShouldNotBeEmpty)
				So(s.Email, ShouldEqual, "user@email.com")
				So(s.LastAccessed.String(), ShouldNotEqual, respLastAccessed)
			})
		})
	})

	Convey("Given a session email client.GetByEmail returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult(string(resp), nil),
			*redis.NewStatusCmd(),
			*redis.NewBoolResult(false, errors.New("unable to refresh expiration")),
		)

		Convey("When client uses the email to get the session", func() {
			s, err := client.GetByEmail("user@email.com")

			Convey("Then redis client.Get is called with the expected parameters", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.GetCalls()[0].Key, ShouldEqual, "user@email.com")

				So(mockRedisClient.ExpireCalls(), ShouldHaveLength, 1) // Expects 2 due to refreshing by ID and Email
				So(mockRedisClient.ExpireCalls()[0].Key, ShouldEqual, "user@email.com")
				So(mockRedisClient.ExpireCalls()[0].Expiration, ShouldEqual, testTTL)
			})

			Convey("Then redis client.Get is called and returns an error", func() {
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unable to refresh expiration")
				So(s, ShouldBeNil)
			})
		})
	})

	Convey("Given a blank session email client.GetByEmail returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client.GetByEmail is called has an empty ID", func() {
			s, err := client.GetByEmail("")

			Convey("Then client.GetByEmail returns an error and no session is returned", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 0)
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err, ShouldEqual, ErrEmptySessionEmail)
			})
		})
	})

	Convey("Given a session ID client.GetByEmail returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringResult("", errors.New("unexpected end of JSON input")),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When client.GetByEmail is called with a valid session ID", func() {
			s, err := client.GetByEmail("user@test.com")

			Convey("Then redis client.Get is called with the expected parameters", func() {
				So(mockRedisClient.GetCalls(), ShouldHaveLength, 1)
				So(mockRedisClient.GetCalls()[0].Key, ShouldEqual, "user@test.com")
				So(mockRedisClient.ExpireCalls(), ShouldHaveLength, 0) // Expects 2 due to refreshing by ID and Email
			})

			Convey("Then the redis client.Get returns an error and no session is returned", func() {
				So(s, ShouldBeNil)
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "unexpected end of JSON input")
			})
		})
	})
}

func TestClient_DeleteAll(t *testing.T) {
	Convey("Given DeleteAll removes all sessions from cache", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringCmd(),
			*redis.NewStatusCmd(),
			*redis.NewBoolCmd(),
		)

		Convey("When DeleteAll is called", func() {
			err := client.DeleteAll()

			Convey("Then all sessions are removed from cache and no error is returned", func() {
				So(err, ShouldBeNil)
				So(mockRedisClient.FlushAllCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given DeleteAll returns an error", t, func() {
		mockRedisClient, client := setUpMocks(
			*redis.NewStatusCmd(),
			*redis.NewStringCmd(),
			*redis.NewStatusResult("fail", errors.New("some redis error")),
			*redis.NewBoolCmd(),
		)

		Convey("When DeleteAll is called", func() {
			err := client.DeleteAll()

			Convey("Then no sessions are removed and a redis error is returned", func() {
				So(err, ShouldNotBeEmpty)
				So(err.Error(), ShouldEqual, "some redis error")
				So(mockRedisClient.FlushAllCalls(), ShouldHaveLength, 1)
			})
		})

	})
}

func setUpMocks(setStatusCmd redis.StatusCmd, getStringCmd redis.StringCmd, flushAllStatusCmd redis.StatusCmd, expireBoolCmd redis.BoolCmd) (*RedisClienterMock, SessionCache) {
	mockRedisClient := &RedisClienterMock{
		PingFunc: nil,
		SetFunc: func(key string, value interface{}, ttl time.Duration) *redis.StatusCmd {
			return &setStatusCmd
		},
		GetFunc: func(key string) *redis.StringCmd {
			return &getStringCmd
		},
		FlushAllFunc: func() *redis.StatusCmd {
			return &flushAllStatusCmd
		},
		ExpireFunc: func(key string, expiration time.Duration) *redis.BoolCmd {
			return &expireBoolCmd
		}}
	return mockRedisClient, &ElasticacheClient{
		client: mockRedisClient,
		ttl:    testTTL,
	}
}
