package api_test

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-sessions-api/api"
	"io"
	"io/ioutil"
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

var timeBefore = time.Now().UTC()

func TestCreateSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		mockSession := &apiMock.SessionUpdaterMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:24400/session", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.UpdateCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			SetSessionFunc: func(s *session.Session) error {
				return nil
			},
		}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			sessionResp, err := unmarshalJSON(resp.Body)
			if err != nil {
				t.Errorf("failed to retreive response body")
			}

			timeAfter := time.Now().UTC()

			Convey("Then the expected success response is returned", func() {
				So(err, ShouldBeNil)
				So(resp.Code, ShouldEqual, 201)
				So(mockCache.SetSessionCalls(), ShouldHaveLength, 1)
				So(sessionResp.Email, ShouldEqual, "test@test.com")
				So(sessionResp.ID, ShouldNotBeNil)
				So(sessionResp.Start, ShouldHappenBetween, timeBefore, timeAfter)
				So(sessionResp.LastAccessed, ShouldHappenBetween, timeBefore, timeAfter)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSession := &apiMock.SessionUpdaterMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader("this is not json"))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.UpdateCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a bad request", t, func() {
		mockSession := &apiMock.SessionUpdaterMock{}
		mockCache := &apiMock.CacheMock{}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("")
		So(err, ShouldBeNil)

		req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSession.UpdateCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a new session", t, func() {
		mockCache := &apiMock.CacheMock{
			SetSessionFunc: func(s *session.Session) error {
				return errors.New("unable to store session in cache")
			}}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			SetSessionFunc: func(s *session.Session) error {
				return errors.New("unable to add session to cache")
			},
		}
		sessionHandler := api.CreateSessionHandlerFunc(mockCache)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest(http.MethodPost, "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockCache.SetSessionCalls(), ShouldHaveLength, 1)
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

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
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

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
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

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
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

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
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

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
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

func TestGetByEmailSessionHandlerFunc(t *testing.T) {
	Convey("Given a valid request", t, func() {
		sessionEmail := "user@test.com"
		currentTime := time.Now()
		mockCache := &apiMock.CacheMock{
			GetByEmailFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "123",
					Email: email,
					Start: currentTime,
				}, nil
			},
		}

		getVars := func(r *http.Request) map[string]string {
			return map[string]string{
				"Email": sessionEmail,
			}
		}

		sessionHandler := api.GetByEmailSessionHandlerFunc(mockCache, getVars)

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then expected session is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				So(mockCache.GetByEmailCalls(), ShouldHaveLength, 1)
				So(mockCache.GetByEmailCalls()[0].Email, ShouldEqual, "user@test.com")
			})
		})
	})

	Convey("Given a request to retrieve a session", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByEmailFunc: func(email string) (*session.Session, error) {
				return nil, SessionNotFound
			},
		}

		sessionHandler := api.GetByEmailSessionHandlerFunc(mockCache, getVars("Email", ""))

		req := httptest.NewRequest(http.MethodGet, "/session/user@test.com", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(mockCache.GetByEmailCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByEmailFunc: func(email string) (*session.Session, error) {
				return nil, nil
			},
		}

		sessionHandler := api.GetByEmailSessionHandlerFunc(mockCache, getVars("Email", "user@test.com"))

		req := httptest.NewRequest(http.MethodGet, "/session/user@test.com", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(mockCache.GetByEmailCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByEmailFunc: func(email string) (*session.Session, error) {
				return nil, SessionExpired
			},
		}

		sessionHandler := api.GetByEmailSessionHandlerFunc(mockCache, getVars("Email", "123"))

		req := httptest.NewRequest(http.MethodGet, "/session/123", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(mockCache.GetByEmailCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given a valid request", t, func() {
		mockCache := &apiMock.CacheMock{
			GetByEmailFunc: func(email string) (*session.Session, error) {
				return nil, errors.New("unexpected error")
			},
		}

		sessionHandler := api.GetByEmailSessionHandlerFunc(mockCache, getVars("Email", "user@test.com"))

		req := httptest.NewRequest(http.MethodGet, "/session/user@test.com", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then an error response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockCache.GetByEmailCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func TestDeleteAllSessionsHandlerFunc(t *testing.T) {
	Convey("Give a valid request", t, func() {
		mockCache := &apiMock.CacheMock{DeleteAllFunc: func() error {
			return nil
		}}

		sessionHandler := api.DeleteAllSessionsHandlerFunc(mockCache)

		req := httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the correct success response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				So(mockCache.DeleteAllCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Give a valid request", t, func() {
		mockCache := &apiMock.CacheMock{DeleteAllFunc: func() error {
			return errors.New("no sessions to delete")
		}}

		sessionHandler := api.DeleteAllSessionsHandlerFunc(mockCache)

		req := httptest.NewRequest(http.MethodDelete, "/sessions", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is received", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the correct success response is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(mockCache.DeleteAllCalls(), ShouldHaveLength, 1)
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

func unmarshalJSON(body io.Reader) (*session.Session, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var s *session.Session

	if err = json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	return s, err
}
