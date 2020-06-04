package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-sessions-api/api"
	apiMock "github.com/ONSdigital/dp-sessions-api/api/mock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mu          sync.Mutex
	testContext = context.Background()
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		p := &auth.NopHandler{}
		api := GetAPIWithMocks(p)

		Convey("When created the following routes should have been added", func() {
			// Replace the check below with any newly added api endpoints
			So(hasRoute(api.Router, "/sessions", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/sessions/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/sessions", "DELETE"), ShouldBeTrue)
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Given an API instance", t, func() {
		p := &apiMock.AuthHandlerMock{
			RequireFunc: func(required auth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					handler.ServeHTTP(w, r)
				}
			},
		}
		a := GetAPIWithMocks(p)

		Convey("When the api is closed any dependencies are closed also", func() {
			err := a.Close(testContext)
			So(err, ShouldBeNil)
			// Check that dependencies are closed here
		})
	})
}

// GetAPIWithMocks also used in other tests
func GetAPIWithMocks(authMock api.AuthHandler) *api.API {
	mu.Lock()
	defer mu.Unlock()
	return api.Setup(testContext, mux.NewRouter(), authMock)
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
