package api_test

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-sessions-api/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apiMock "github.com/ONSdigital/dp-sessions-api/api/mock"
	. "github.com/ONSdigital/dp-sessions-api/errors"
	"github.com/ONSdigital/dp-sessions-api/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		mockSession := &apiMock.SessionMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

		req := httptest.NewRequest("POST", "http://localhost:24400/session", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		currentTime := time.Now()
		mockSession := &apiMock.SessionMock{
			NewFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "1234",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}
		mockCache := &apiMock.CacheMock{
			SetFunc: func(s *session.Session) error {
				return nil
			},
		}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

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
				So(mockSession.NewCalls(), ShouldHaveLength, 1)
				So(mockSession.NewCalls()[0].Email, ShouldEqual, "test@test.com")
				So(mockCache.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSession := &apiMock.SessionMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

		req := httptest.NewRequest("POST", "/session", strings.NewReader("this is not json"))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSession := &apiMock.SessionMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a new session", t, func() {
		mockSession := &apiMock.SessionMock{
			NewFunc: func(email string) (*session.Session, error) {
				return nil, errors.New("unable to generate id")
			},
		}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockSession.NewCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		currentTime := time.Now()
		mockSession := &apiMock.SessionMock{
			NewFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "1234",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}
		mockCache := &apiMock.CacheMock{
			SetFunc: func(s *session.Session) error {
				return errors.New("unable to add session to cache")
			},
		}
		sessionHandler := api.CreateSessionHandlerFunc(mockSession, mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockSession.NewCalls(), ShouldHaveLength, 1)
				So(mockCache.SetCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func TestGetByIDSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		sessionID := "123"
		currentTime := time.Now()
		mockCache := &apiMock.CacheMock{
			GetByIDFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "123",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}

		getVars := func(r *http.Request) map[string]string {
			return map[string]string{
				"ID": sessionID,
			}
		}

		sessionHandler := api.GetByIDSessionHandlerFunc(mockCache, getVars)

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
		mockCache := &apiMock.CacheMock{
			GetByIDFunc: func(ID string) (*session.Session, error) {
				return nil, SessionNotFound
			},
		}

		sessionHandler := api.GetByIDSessionHandlerFunc(mockCache, getVars("ID", ""))

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

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByIDFunc: func(ID string) (*session.Session, error) {
				return nil, nil
			},
		}

		sessionHandler := api.GetByIDSessionHandlerFunc(mockCache, getVars("ID", "123"))

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

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByIDFunc: func(ID string) (*session.Session, error) {
				return nil, SessionExpired
			},
		}

		sessionHandler := api.GetByIDSessionHandlerFunc(mockCache, getVars("ID", "123"))

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

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByIDFunc: func(ID string) (*session.Session, error) {
				return nil, errors.New("unexpected error")
			},
		}

		sessionHandler := api.GetByIDSessionHandlerFunc(mockCache, getVars("ID", "123"))

		req := httptest.NewRequest("GET", "/session/123", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
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

func getVars(k string, v string) api.GetVarsFunc {
	return func(r *http.Request) map[string]string {
		return map[string]string{
			k: v,
		}
	}
}
