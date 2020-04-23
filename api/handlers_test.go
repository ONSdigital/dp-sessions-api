package api

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateSessionHandlerFunc(t *testing.T) {

	sessionHandler := CreateSessionHandlerFunc()

	Convey("Given a request to /session with no body", t, func() {
		req := httptest.NewRequest("POST", "http://localhost:24400/session", nil)
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the response should be a 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})

	Convey("Give a valid create session request to /session with body with all elements", t, func() {
		sess := session.Session{
			ID:    "123",
			Email: "me@me.com",
			Start: time.Now(),
		}
		sessJSON, err := sess.MarshalJSON()

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the response should be 201", func() {
				So(resp.Code, ShouldEqual, 201)
				So(err, ShouldBeNil)
			})
		})
	})

	/*
	Convey("Give a request to /session with a body with missing elements", t, func() {
		sess := session.Session{
			Email: "me@me.com",
			Start: time.Now(),
		}
		sessJSON, _ := sess.MarshalJSON()

		req := httptest.NewRequest("POST", "/session", strings.NewReader(string(sessJSON)))
		resp := httptest.NewRecorder()

		Convey("When the request is handled by the router", func() {
			sessionHandler.ServeHTTP(resp, req)

			Convey("Then the response should be 400", func() {
				So(resp.Code, ShouldEqual, 400)
			})
		})
	})
	*/
}
