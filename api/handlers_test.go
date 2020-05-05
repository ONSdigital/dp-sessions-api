package api

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-sessions-api/session"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		mockSessions := &SessionsMock{}
		mockCache := &CacheMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		req := httptest.NewRequest("POST", "http://localhost:24400/session", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		currentTime := time.Now()
		mockSessions := &SessionsMock{
			NewFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "1234",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}
		mockCache := &CacheMock{
			SetFunc: func(s *session.Session) error {
				return nil
			},
		}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the expected success response is returned", func() {
				expected := &session.Session{
					ID:    "1234",
					Email: "test@test.com",
					Start: currentTime,
				}

				b, err := json.Marshal(expected)
				So(err, ShouldBeNil)
				So(resp.Body.String(), ShouldEqual, string(b))
				So(resp.Code, ShouldEqual, 201)
				So(mockSessions.NewCalls(), ShouldHaveLength, 1)
				So(mockSessions.NewCalls()[0].Email, ShouldEqual, "test@test.com")
				So(mockCache.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSessions := &SessionsMock{}
		mockCache := &CacheMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		req := httptest.NewRequest("POST", "/session", strings.NewReader("this is not json"))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSessions := &SessionsMock{}
		mockCache := &CacheMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a new session", t, func() {
		mockSessions := &SessionsMock{
			NewFunc: func(email string) (*session.Session, error) {
				return nil, errors.New("unable to generate id")
			},
		}
		mockCache := &CacheMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockSessions.NewCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		currentTime := time.Now()
		mockSessions := &SessionsMock{
			NewFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "1234",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}
		mockCache := &CacheMock{
			SetFunc: func(s *session.Session) error {
				return errors.New("unable to add session to cache")
			},
		}
		sessionHandler := CreateSessionHandlerFunc(mockSessions, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockSessions.NewCalls(), ShouldHaveLength, 1)
				So(mockCache.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func TestGetByIDSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		sessionID := "123"
		currentTime := time.Now()
		mockCache := &CacheMock{
			GetByIDFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "123",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}

		getVars := func(r *http.Request)map[string]string {
			return map[string]string{
				"ID": sessionID,
			}
		}

		sessionHandler := GetByIDSessionHandlerFunc(mockCache, getVars)

		req := httptest.NewRequest("GET", "/session/123", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then expected session is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				So(mockCache.GetByIDCalls(), ShouldHaveLength, 1)
				So(mockCache.GetByIDCalls()[0].ID, ShouldEqual, "123")
			})
		})
	})

	Convey("Given a request to retrieve a session", t, func() {
		mockCache := &CacheMock{
			GetByIDFunc: func(ID string) (*session.Session, error) {
				return nil, errors.New("unable to get session by id")
			},
		}

		sessionHandler := GetByIDSessionHandlerFunc(mockCache, getVars("ID", ""))

		req := httptest.NewRequest("GET", "/session/123", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(mockCache.GetByIDCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func newSessionDetailsAndMarshal(email string) ([]byte, error) {
	sess := session.NewSessionDetails{
		Email: email,
	}
	sessJSON, err := json.Marshal(sess)
	if err != nil {
		return nil, err
	}

	return sessJSON, nil
}

func getVars(k string, v string) GetVarsFunc {
	return func(r *http.Request) map[string]string {
		return map[string]string{
			k: v,
		}
	}
}
