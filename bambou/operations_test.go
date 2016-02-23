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

func TestOperations_FetchEntity(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object", t, func() {

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "xxx"}}

		Convey("When I fetch it with success", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("[{\"ID\": \"xxx\", \"parentType\": \"pedro\", \"parentID\": \"yyy\"}]"),
				}
			}).restore()

			FetchEntity(e)

			Convey("Then parentType should pedro", func() {
				So(e.ParentType, ShouldEqual, "pedro")
			})

			Convey("Then parentID should yyy", func() {
				So(e.ParentID, ShouldEqual, "yyy")
			})
		})

		Convey("When I fetch it and I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			err := FetchEntity(e)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I fetch it and I got bad json", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("definitely a bad json"),
				}
			}).restore()

			Convey("Then it should panic", func() {
				So(func() { FetchEntity(e) }, ShouldPanic)
			})
		})
	})
}

func TestOperations_SaveEntity(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object", t, func() {

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "yyy"}}

		Convey("When I save it with success", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("[{\"ID\": \"zzz\", \"parentType\": \"pedro\", \"parentID\": \"yyy\"}]"),
				}
			}).restore()

			SaveEntity(e)

			Convey("Then ID should zzz", func() {
				So(e.ID, ShouldEqual, "zzz")
			})
		})

		Convey("When I save it and I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			err := SaveEntity(e)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I save it and I got bad json", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("definitely a bad json"),
				}
			}).restore()

			Convey("Then it should panic", func() {
				So(func() { SaveEntity(e) }, ShouldPanic)
			})
		})

		Convey("When I save it and I got no data", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
				}
			}).restore()

			Convey("Then it not should panic", func() {
				So(func() { SaveEntity(e) }, ShouldNotPanic)
			})
		})

	})
}

func TestOperations_DeleteEntity(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object", t, func() {

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "yyy"}}

		Convey("When I delete it with success", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("[{\"ID\": \"zzz\", \"parentType\": \"pedro\", \"parentID\": \"yyy\"}]"),
				}
			}).restore()

			DeleteEntity(e)

			Convey("Then ID should yyy", func() {
				So(e.ID, ShouldEqual, "yyy")
			})
		})

		Convey("When I delete it and I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			err := DeleteEntity(e)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestOperations_FetchChildren(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object", t, func() {

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "yyy"}}

		Convey("When I Fetch its children with success", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("[{\"ID\": \"1\"}, {\"ID\": \"2\"}]"),
				}
			}).restore()

			var l fakeObjectsList
			FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then the lenght of the children list should be 2", func() {
				So(len(l), ShouldEqual, 2)
			})

			Convey("Then the first child ID should be 1", func() {
				So(l[0].ID, ShouldEqual, "1")
			})

			Convey("Then the second child ID should be 2", func() {
				So(l[1].ID, ShouldEqual, "2")
			})
		})

		Convey("When I fetch the children and I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			var l fakeObjectsList
			err := FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I fetch the children I got bad json", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("definitely a bad json"),
				}
			}).restore()

			var l fakeObjectsList

			Convey("Then it should panic", func() {
				So(func() { FetchChildren(e, fakeIdentity, &l, nil) }, ShouldPanic)
			})
		})
	})
}

func TestOperations_CreateChild(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object and a child", t, func() {

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "xxx"}}
		c := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity}}

		Convey("When I create it with success", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 201,
					Data: []byte("[{\"ID\": \"zzz\"}]"),
				}
			}).restore()

			CreateChild(e, c)

			Convey("Then ID of the new children should be zzz", func() {
				So(c.ID, ShouldEqual, "zzz")
			})
		})

		Convey("When I create it I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			err := CreateChild(e, c)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I create it I got bad json", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 200,
					Data: []byte("definitely a bad json"),
				}
			}).restore()

			Convey("Then it should panic", func() {
				So(func() { CreateChild(e, c) }, ShouldPanic)
			})
		})
	})
}

func TestOperations_AssignChildren(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new Exposed object and a list of children to assign", t, func() {

		defer patch(&sendNativeRequest, func(request *request) *response {
			return &response{
				Code: 200,
			}
		}).restore()

		e := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity, ID: "xxx"}}
		c := &fakeObject{ExposedObject: ExposedObject{Identity: fakeIdentity}}
		l := fakeObjectsList{c}

		Convey("When I assign them with success", func() {

			AssignChildren(e, l, fakeIdentity)

			Convey("Then nothing special should happen actually", func() {
			})
		})

		Convey("When I assign them I got an communication error", func() {

			defer patch(&sendNativeRequest, func(request *request) *response {
				return &response{
					Code: 500,
				}
			}).restore()

			err := AssignChildren(e, l, fakeIdentity)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

	})
}
