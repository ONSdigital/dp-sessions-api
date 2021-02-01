package session

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSession_MarshalJSON(t *testing.T) {

	Convey("Given valid session", t, func() {
		start := time.Now()
		lastAccess := time.Now()

		s := &Session{
			ID:           "123",
			Email:        "test@test.com",
			Start:        start,
			LastAccessed: lastAccess,
		}

		Convey("When MarshalJSON is invoked", func() {
			jsonBytes, err := s.MarshalJSON()
			So(err, ShouldBeNil)
			So(jsonBytes, ShouldNotBeNil)

			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonBytes, &jsonMap)
			So(err, ShouldBeNil)

			Convey("Then session JSON has the expected field values", func() {
				expectedStartVal := start.Format(dateTimeFMT)
				expectedLastAccessedVal := lastAccess.Format(dateTimeFMT)

				assertJSONFieldValue("id", "123", jsonMap)
				assertJSONFieldValue("email", "test@test.com", jsonMap)
				assertJSONFieldValue("start", expectedStartVal, jsonMap)
				assertJSONFieldValue("last_accessed", expectedLastAccessedVal, jsonMap)
			})
		})
	})
}

func assertJSONFieldValue(key, expectedValue string, jsonMap map[string]interface{}) {
	actualValue, exists := jsonMap[key]
	So(exists, ShouldBeTrue)
	So(actualValue, ShouldEqual, expectedValue)
}
