package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSession_NewSession(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &TestRoot{ExposedObject: ExposedObject{Identity: TestRooIdentity}}
		s := NewSession("username", "password", "organization", "http://url.com", r)

		Convey("Then the properties Username should be 'username'", func() {
			So(s.Username, ShouldEqual, "username")
		})

		Convey("Then the properties Password should be 'password'", func() {
			So(s.Password, ShouldEqual, "password")
		})

		Convey("Then the properties Organization should be 'organization'", func() {
			So(s.Organization, ShouldEqual, "organization")
		})

		Convey("Then the properties URL should be 'http://url.com'", func() {
			So(s.URL, ShouldEqual, "http://url.com")
		})

		Convey("Then the properties Root should not be nil", func() {
			So(s.Root, ShouldNotBeNil)
		})
	})
}

func TestSession_MakeAuthorizationHeaders(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &TestRoot{ExposedObject: ExposedObject{Identity: TestRooIdentity}}
		s := NewSession("username", "password", "organization", "http://url.com", r)

		Convey("When I prepare the Header", func() {
			h := s.MakeAuthorizationHeaders()

			Convey("Then the header should be 'XREST dXNlcm5hbWU6cGFzc3dvcmQ", func() {
				So(h, ShouldEqual, "XREST dXNlcm5hbWU6cGFzc3dvcmQ=")
			})
		})
	})
}

func TestSession_StartStopSession(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &TestRoot{ExposedObject: ExposedObject{Identity: TestRooIdentity}}
		s := NewSession("username", "password", "organization", "http://url.com", r)

		Convey("Then the CurrentSession() should be nil", func() {
			So(CurrentSession(), ShouldBeNil)
		})

		Convey("When I start the session and retrieve the CurrentSession", func() {
			s.Start()
			c := CurrentSession()

			Convey("Then the CurrentSession should be equal to session", func() {
				So(c, ShouldEqual, s)
			})

			Convey("Then the session APIKey should be 'api-key'", func() {
				So(c.APIKey, ShouldEqual, "api-key")
			})

			Convey("Then the Root User APIKey should be 'api-key'", func() {
				So(c.Root.GetAPIKey(), ShouldEqual, "api-key")
			})

		})

		Convey("When I reset the session and retrieve the CurrentSession", func() {
			s.Reset()
			c := CurrentSession()

			Convey("Then the CurrentSession should be nil", func() {
				So(c, ShouldBeNil)
			})
		})
	})

}

var TestRooIdentity = Identity{
	RESTName:     "root",
	ResourceName: "root",
}

type TestRoot struct {
	ExposedObject

	UserName     string `json:"userName,omitempty"`
	Password     string `json:"password,omitempty"`
	APIKey       string `json:"APIKey,omitempty"`
	Organization string `json:"enterprise,omitempty"`
}

func (o *TestRoot) GetAPIKey() string                           { return o.APIKey }
func (o *TestRoot) SetAPIKey(key string)                        { o.APIKey = key }
func (o *TestRoot) GetURL() string                              { return CurrentSession().URL + "/" + o.Identity.ResourceName }
func (o *TestRoot) Save() *Error                                { return nil }
func (o *TestRoot) Delete() *Error                              { return nil }
func (o *TestRoot) GetURLForChildrenIdentity(i Identity) string { return "" }

func (o *TestRoot) Fetch() *Error {
	o.APIKey = "api-key"
	return nil
}
