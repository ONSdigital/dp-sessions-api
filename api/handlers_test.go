package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-sessions-api/session"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateSessionHandlerFunc(t *testing.T) {

	Convey("Given a request to /session with no body", t, func() {
		mockSessions := &SessionsMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions)

		req := httptest.NewRequest("POST", "http://localhost:24400/session", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a valid create session request to /session with body with all elements", t, func() {
		mockSessions := &SessionsMock{
			NewFunc: func(email string) (*session.Session, error) {
				return &session.Session{
					ID:    "1234",
					Email: email,
				}, nil
			},
		}
		sessionHandler := CreateSessionHandlerFunc(mockSessions)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the expected success response is returned", func() {
				So(resp.Code, ShouldEqual, 201)
				So(resp.Header().Get("Content-Location"), ShouldEqual, fmt.Sprintf("/session/1234"))
				So(mockSessions.NewCalls(), ShouldHaveLength, 1)
				So(mockSessions.NewCalls()[0].Email, ShouldEqual, "test@test.com")
			})
		})
	})

	Convey("Given a bad request to /session", t, func() {
		mockSessions := &SessionsMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions)

		req := httptest.NewRequest("POST", "/session", strings.NewReader("this is not json"))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a request to /session with a body with missing elements", t, func() {
		mockSessions := &SessionsMock{}
		sessionHandler := CreateSessionHandlerFunc(mockSessions)

		sessJSON, err := newSessionDetailsAndMarshal("")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
				So(mockSessions.NewCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given a create new session returns an error", t, func() {
		mockSessions := &SessionsMock{
			NewFunc: func(email string) (*session.Session, error) {
				return nil, errors.New("unable to generate id")
			},
		}
		sessionHandler := CreateSessionHandlerFunc(mockSessions)

		sessJSON, err := newSessionDetailsAndMarshal("test@test.com")
		So(err, ShouldBeNil)

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then return an error response", func() {
				So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockSessions.NewCalls(), ShouldHaveLength, 1)
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
