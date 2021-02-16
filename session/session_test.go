package session

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("New should return expected error if email address is empty", t, func() {
		s, err := New("")
		So(err, ShouldEqual, EmailEmptyErr)
		So(s, ShouldBeNil)
	})

	Convey("New should return valid session if email address is not empty", t, func() {
		s, err := New("test@test.com")

		So(err, ShouldBeNil)
		So(s.Email, ShouldEqual, "test@test.com")
		So(s.ID, ShouldNotBeEmpty)
		So(s.LastAccessed, ShouldNotBeNil)
		So(s.Start, ShouldNotBeNil)
		So(s.Start, ShouldEqual, s.LastAccessed)
	})

}

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
				expectedStartVal := start.Format(DateTimeFMT)
				expectedLastAccessedVal := lastAccess.Format(DateTimeFMT)

				assertJSONFieldValue("id", "123", jsonMap)
				assertJSONFieldValue("email", "test@test.com", jsonMap)
				assertJSONFieldValue("start", expectedStartVal, jsonMap)
				assertJSONFieldValue("last_accessed", expectedLastAccessedVal, jsonMap)
			})
		})
	})
}

func TestSession_UnmarshalJSON(t *testing.T) {
	Convey("Should unmarshal valid session from JSON", t, func() {
		input, err := New("test@ons.gov.uk")
		So(err, ShouldBeNil)

		jsonBytes, err := json.Marshal(input)
		So(err, ShouldBeNil)

		var output Session
		err = json.Unmarshal(jsonBytes, &output)
		So(err, ShouldBeNil)

		So(input, ShouldResemble, &output)
	})

	Convey("Should return expected error is session.Start is blank/empty ", t, func() {
		input := `{"id":"123","email":"test@ons.gov.uk","last_accessed":"2021-02-02T11:51:48.300Z"}`

		var output Session
		err := json.Unmarshal([]byte(input), &output)
		So(err, ShouldEqual, StartEmptyErr)
	})

	Convey("Should return expected error is session.Start is blank/empty ", t, func() {
		input := `{"id":"123","email":"test@ons.gov.uk","start":"2021-02-02T11:51:48.300Z"}`

		var output Session
		err := json.Unmarshal([]byte(input), &output)
		So(err, ShouldEqual, LastAccessedEmptyErr)

		err = json.Unmarshal([]byte(input), &output)
		if err != nil {

		}
	})

	Convey("Should return expected error if session.Start invalid time value", t, func() {
		input := `{"id":"123","email":"test@ons.gov.uk","start":"bob", "last_accessed":"2021-02-02T11:51:48.300Z"}`

		var output Session
		err := json.Unmarshal([]byte(input), &output)
		So(err.Error(), ShouldStartWith, "error parsing session.Start as time.Time value")
	})

	Convey("Should return expected error if session.Start invalid time value", t, func() {
		input := `{"id":"123","email":"test@ons.gov.uk","start":"2021-02-02T11:51:48.300Z", "last_accessed":"wibble"}`

		var output Session
		err := json.Unmarshal([]byte(input), &output)
		So(err.Error(), ShouldStartWith, "error parsing session.LastAccessed as time.Time value")
	})
}

func assertJSONFieldValue(key, expectedValue string, jsonMap map[string]interface{}) {
	actualValue, exists := jsonMap[key]
	So(exists, ShouldBeTrue)
	So(actualValue, ShouldEqual, expectedValue)
}
