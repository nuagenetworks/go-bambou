// Copyright (c) 2015, Alcatel-Lucent Inc.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
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

func TestExposedObject_GetIdentity(t *testing.T) {

	Convey("Given I create a new ExposedObject", t, func() {
		e := &fakeExposed{ExposedObject: ExposedObject{Identity: fakeIdentity}}

		Convey("Then Identity should fake", func() {
			So(e.GetIdentity(), ShouldResemble, fakeIdentity)
		})
	})
}

func TestExposedObject_URL(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new ExposedObject", t, func() {

		e := &fakeExposed{ExposedObject: ExposedObject{Identity: fakeIdentity}}

		Convey("When I don't set the ID", func() {
			Convey("Then URL should be http://fake.com/fakes", func() {
				So(e.GetURL(), ShouldEqual, "http://fake.com/fakes")
			})
		})

		Convey("When I set the ID", func() {

			e.ID = "xxx"

			Convey("Then URL should be http://fake.com/fakes/xxx", func() {
				So(e.GetURL(), ShouldEqual, "http://fake.com/fakes/xxx")
			})

		})
	})
}

func TestExposedObject_GetURLForChildrenIdentity(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new ExposedObject", t, func() {

		e := &fakeExposed{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "xxx"}}

		Convey("When I retrieve the URL for an identity", func() {

			i := Identity{"child", "children"}

			Convey("Then children URL should http://fake.com/fakes/xxx/children", func() {
				So(e.GetURLForChildrenIdentity(i), ShouldEqual, "http://fake.com/fakes/xxx/children")
			})
		})
	})
}

func TestExposedObject_String(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new ExposedObject", t, func() {

		e := &fakeExposed{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "xxx"}}

		Convey("Then the string representation should be <fake:xxx>", func() {
			So(e.String(), ShouldEqual, "<fake:xxx>")
		})
	})
}

var fakeIdentity = Identity{"fake", "fakes"}

type fakeExposed struct {
	ExposedObject
}

func (o *fakeExposed) Save() *Error   { return nil }
func (o *fakeExposed) Delete() *Error { return nil }
func (o *fakeExposed) Fetch() *Error  { return nil }
