// Copyright (c) 2015, Alcatel-Lucent Inc.
// All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSession_NewSession(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &testRoot{fakeExposed: fakeExposed{ExposedObject: ExposedObject{Identity: testRootdentity}}}

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

func TestSession_makeAuthorizationHeaders(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &testRoot{fakeExposed: fakeExposed{ExposedObject: ExposedObject{Identity: testRootdentity}}}

		Convey("When I prepare the Headers with a session that doesn't have an APIKey", func() {

			s := NewSession("username", "password", "organization", "http://url.com", r)
			h := s.makeAuthorizationHeaders()

			Convey("Then the header should be 'XREST dXNlcm5hbWU6cGFzc3dvcmQ", func() {
				So(h, ShouldEqual, "XREST dXNlcm5hbWU6cGFzc3dvcmQ=")
			})
		})

		Convey("When I prepare the Headers with a session that already has an APIKey", func() {

			s := NewSession("username", "password", "organization", "http://url.com", r)
			s.APIKey = "api-key"
			h := s.makeAuthorizationHeaders()

			Convey("Then the header should be 'XREST dXNlcm5hbWU6cGFzc3dvcmQ", func() {
				So(h, ShouldEqual, "XREST dXNlcm5hbWU6YXBpLWtleQ==")
			})
		})

		Convey("When I prepare the Headers with a session missing username", func() {

			s := NewSession("", "password", "organization", "http://url.com", r)

			Convey("It should panic", func() {
				So(func() { s.makeAuthorizationHeaders() }, ShouldPanic)
			})
		})

		Convey("When I prepare the Headers with a session missing password", func() {

			s := NewSession("username", "", "organization", "http://url.com", r)

			Convey("It should panic", func() {
				So(func() { s.makeAuthorizationHeaders() }, ShouldPanic)
			})
		})
	})
}

func TestSession_StartStopSession(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &testRoot{fakeExposed: fakeExposed{ExposedObject: ExposedObject{Identity: testRootdentity}}}
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

		Convey("When I start the session and cannot get the root object", func() {

			s.Root = &testFailedRoot{}
			err := s.Start()

			Convey("The err should not be nil", func() {
				So(err, ShouldNotBeNil)
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
