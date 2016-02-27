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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

/*
	Initialization
*/
func TestSession_NewSession(t *testing.T) {

	Convey("When I create a new Session", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		s := NewSession("username", "password", "organization", "http://url.com", r)

		Convey("Then the property Username should be 'username'", func() {
			So(s.Username, ShouldEqual, "username")
		})

		Convey("Then the property Password should be 'password'", func() {
			So(s.Password, ShouldEqual, "password")
		})

		Convey("Then the property Organization should be 'organization'", func() {
			So(s.Organization, ShouldEqual, "organization")
		})

		Convey("Then the property URL should be 'http://url.com'", func() {
			So(s.URL, ShouldEqual, "http://url.com")
		})

		Convey("Then the property Root should not be nil", func() {
			So(s.Root, ShouldNotBeNil)
		})
	})
}

/*
	Privates
*/
func TestSession_makeAuthorizationHeaders(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		Convey("When I prepare the Authorization with a session that doesn't have an APIKey", func() {

			s := NewSession("username", "password", "organization", "http://url.com", r)
			h := s.makeAuthorizationHeaders()

			Convey("Then the header should be 'XREST dXNlcm5hbWU6cGFzc3dvcmQ=", func() {
				So(h, ShouldEqual, "XREST dXNlcm5hbWU6cGFzc3dvcmQ=")
			})
		})

		Convey("When I prepare the Authorization with a session that already has an APIKey", func() {

			s := NewSession("username", "password", "organization", "http://url.com", r)
			s.Root().SetAPIKey("api-key")
			h := s.makeAuthorizationHeaders()

			Convey("Then the header should be 'XREST dXNlcm5hbWU6YXBpLWtleQ==", func() {
				So(h, ShouldEqual, "XREST dXNlcm5hbWU6YXBpLWtleQ==")
			})
		})

		Convey("When I prepare the Authorization with a session missing username", func() {

			s := NewSession("", "password", "organization", "http://url.com", r)

			Convey("Then it should panic", func() {
				So(func() { s.makeAuthorizationHeaders() }, ShouldPanic)
			})
		})

		Convey("When I prepare the Authorization with a session missing password", func() {

			s := NewSession("username", "", "organization", "http://url.com", r)

			Convey("Then it should panic", func() {
				So(func() { s.makeAuthorizationHeaders() }, ShouldPanic)
			})
		})
	})
}

func TestSession_prepareHeaders(t *testing.T) {

	r := &fakeRootObject{}
	r.SetIdentity(fakeRootIdentity)

	session := NewSession("username", "password", "organization", "http://fake.com", r)

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()
		r, _ := http.NewRequest("GET", "http://fake.com", nil)

		Convey("When I prepareHeaders with a no fetching info", func() {
			session.prepareHeaders(r, nil)

			Convey("Then I should not have a value for X-Nuage-Page", func() {
				So(r.Header.Get("X-Nuage-Page"), ShouldEqual, "")
			})

			Convey("Then I should have a the X-Nuage-PageSize set to 50", func() {
				So(r.Header.Get("X-Nuage-PageSize"), ShouldEqual, "50")
			})

			Convey("Then I should not have a value for X-Nuage-Filter", func() {
				So(r.Header.Get("X-Nuage-Filter"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-OrderBy", func() {
				So(r.Header.Get("X-Nuage-OrderBy"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-GroupBy", func() {
				So(r.Header.Get("X-Nuage-GroupBy"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-Attributes", func() {
				So(r.Header.Get("X-Nuage-Attributes"), ShouldEqual, "")
			})
		})

		Convey("When I prepareHeaders witha fetching info that has a all fields", func() {
			f.Page = 2
			f.PageSize = 42
			f.Filter = "filter"
			f.OrderBy = "orderby"
			f.GroupBy = []string{"group1", "group2"}

			session.prepareHeaders(r, f)

			Convey("Then I should have a the X-Nuage-Page set to 2", func() {
				So(r.Header.Get("X-Nuage-Page"), ShouldEqual, "2")
			})

			Convey("Then I should have a the X-Nuage-PageSize set to 42", func() {
				So(r.Header.Get("X-Nuage-PageSize"), ShouldEqual, "42")
			})

			Convey("Then I should have a value for X-Nuage-Filter set to 'filter'", func() {
				So(r.Header.Get("X-Nuage-Filter"), ShouldEqual, "filter")
			})

			Convey("Then I should have a value for X-Nuage-OrderBy set to 'orderby'", func() {
				So(r.Header.Get("X-Nuage-OrderBy"), ShouldEqual, "orderby")
			})

			Convey("Then I should have a value for X-Nuage-GroupBy set to true", func() {
				So(r.Header.Get("X-Nuage-GroupBy"), ShouldEqual, "true")
			})

			Convey("Then I should have a value for X-Nuage-Attributes contains group1 and group2", func() {
				So(r.Header.Get("X-Nuage-Attributes"), ShouldEqual, "group1, group2")
			})
		})

	})
}

func TestSession_readHeaders(t *testing.T) {

	Convey("Given I create a new session an a FetchingInfo", t, func() {

		session := NewSession("username", "password", "organization", "http://fake.com", nil)

		f := NewFetchingInfo()
		r := &http.Response{Header: http.Header{}}

		Convey("When I readHeaders with a no fetching info", func() {

			session.readHeaders(r, nil)

			Convey("Then nothing should happen", func() {
			})
		})

		Convey("When I readHeaders with a request that has informations", func() {

			r.Header.Set("X-Nuage-Page", "3")
			r.Header.Set("X-Nuage-PageSize", "42")
			r.Header.Set("X-Nuage-Filter", "filter")
			r.Header.Set("X-Nuage-FilterType", "text")
			r.Header.Set("X-Nuage-OrderBy", "value")
			r.Header.Set("X-Nuage-Count", "666")

			session.readHeaders(r, f)

			Convey("Then FecthingInfo.Page should be 3", func() {
				So(f.Page, ShouldEqual, 3)
			})

			Convey("Then FecthingInfo.PageSize should be 42", func() {
				So(f.PageSize, ShouldEqual, 42)
			})

			Convey("Then FecthingInfo.Filter should be filter", func() {
				So(f.Filter, ShouldEqual, "filter")
			})

			Convey("Then FecthingInfo.FilterType should be text", func() {
				So(f.FilterType, ShouldEqual, "text")
			})

			Convey("Then FecthingInfo.OrderBy should be value", func() {
				So(f.OrderBy, ShouldEqual, "value")
			})

			Convey("Then FecthingInfo.TotalCount should be 666", func() {
				So(f.TotalCount, ShouldEqual, 666)
			})
		})
	})
}

func TestSession_URI(t *testing.T) {

	Convey("Given I create a new Session", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)
		r.SetIdentifier("yyy")

		s := NewSession("username", "password", "organization", "http://url.com", r)

		Convey("When I check the URI the root object", func() {

			Convey("Then personal URL should be http://url.com/root", func() {
				So(s.getPersonalURL(r), ShouldEqual, "http://url.com/root")
			})

			Convey("Then general URL should be http://url.com/root", func() {
				So(s.getGeneralURL(r), ShouldEqual, "http://url.com/root")
			})

			Convey("Then URL of the children with fakeIdentity should be http://url.com/fakes", func() {
				So(s.getURLForChildrenIdentity(r, fakeIdentity), ShouldEqual, "http://url.com/fakes")
			})
		})

		Convey("When I check the URI of a standard object with an ID", func() {

			e := &fakeObject{}
			e.SetIdentity(fakeIdentity)
			e.SetIdentifier("xxx")

			Convey("Then personal URL should be http://url.com/fakes/xxx", func() {
				So(s.getPersonalURL(e), ShouldEqual, "http://url.com/fakes/xxx")
			})

			Convey("Then general URL should be http://url.com/fakes", func() {
				So(s.getGeneralURL(e), ShouldEqual, "http://url.com/fakes")
			})

			Convey("Then URL of the children with fakeRootIdentity should be http://url.com/fakes/xxx/root", func() {
				So(s.getURLForChildrenIdentity(e, fakeRootIdentity), ShouldEqual, "http://url.com/fakes/xxx/root")
			})
		})

		Convey("When I check the URI of a standard object without an ID", func() {

			e := &fakeObject{}
			e.SetIdentity(fakeIdentity)

			Convey("Then getting general URL should not panic", func() {
				So(func() { s.getGeneralURL(e) }, ShouldNotPanic)
			})

			Convey("Then getting the personal URL should panic", func() {
				So(func() { s.getPersonalURL(e) }, ShouldPanic)
			})
		})

		Convey("When I check the URI of an object without an Identity", func() {

			e := &fakeObject{}
			e.SetIdentifier("xxx")

			Convey("Then getting general URL should panic", func() {
				So(func() { s.getGeneralURL(e) }, ShouldPanic)
			})

			Convey("Then getting the personal URL should panic", func() {
				So(func() { s.getPersonalURL(e) }, ShouldPanic)
			})
		})
	})
}

/*
	Operations
*/
func TestSession_StartStopSession(t *testing.T) {

	Convey("GivenI create a new session", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `[{"ID": "xxx", "APIKey": "api-key"}]`)
		}))
		defer ts.Close()
		session := NewSession("username", "password", "organization", ts.URL, r)

		Convey("When I start the session and I can get the root object", func() {

			session.Start()

			Convey("Then the session should be current", func() {
				So(CurrentSession(), ShouldEqual, session)
			})

			Convey("Then the Root User APIKey should be 'api-key'", func() {
				So(CurrentSession().Root().APIKey(), ShouldEqual, "api-key")
			})

			Convey("When I reset the session everything should be back to nil", func() {

				session.Reset()

				Convey("Then the session API key should be ''", func() {
					So(session.Root().APIKey(), ShouldEqual, "")
				})

				Convey("Then the current session should be nil", func() {
					So(CurrentSession(), ShouldBeNil)
				})

			})
		})
	})

	Convey("When I start the session and I cannot get the root object", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "woops", 500)
		}))
		defer ts.Close()

		session := NewSession("username", "password", "organization", ts.URL, r)
		err := session.Start()

		Convey("Then err should not be nil", func() {
			So(err, ShouldNotBeNil)
		})
	})
}

func TestSession_FetchEntity(t *testing.T) {

	r := &fakeRootObject{}
	r.SetIdentity(fakeRootIdentity)

	Convey("Given I create a new session", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `[{"ID": "xxx", "parentType": "pedro", "parentID": "yyy"}]`)
		}))
		defer ts.Close()
		session := NewSession("username", "password", "organization", ts.URL, r)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("xxx")

		Convey("When I fetch an entity with success", func() {

			session.FetchEntity(e)

			Convey("Then parentType should pedro", func() {
				So(e.ParentType, ShouldEqual, "pedro")
			})

			Convey("Then parentID should yyy", func() {
				So(e.ParentID, ShouldEqual, "yyy")
			})
		})

		Convey("When I fetch an entity and I got an communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "bad comm", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			err := session.FetchEntity(e)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I fetch an entity and I got a bad json", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `not good at all`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			Convey("Then it should panic", func() {
				So(func() { session.FetchEntity(e) }, ShouldPanic)
			})
		})
	})
}

func TestSession_SaveEntity(t *testing.T) {

	Convey("Given I create a new object", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("yyy")

		Convey("When I save it with success", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[{"ID": "zzz", "parentType": "pedro", "parentID": "yyy"}]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			session.SaveEntity(e)

			Convey("Then ID should 'zzz'", func() {
				So(e.Identifier(), ShouldEqual, "zzz")
			})
		})

		Convey("When I save it and I got an communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "nope", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			err := session.SaveEntity(e)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I save it and I got a bad json", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `bad json`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			Convey("Then it should panic", func() {
				So(func() { session.SaveEntity(e) }, ShouldPanic)
			})
		})

		Convey("When I save it and I got no data", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			Convey("Then it not should panic", func() {
				So(func() { session.SaveEntity(e) }, ShouldNotPanic)
			})
		})
	})
}

func TestSession_DeleteEntity(t *testing.T) {

	Convey("Given I have an existing object", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("yyy")

		Convey("When I delete it with success", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[{"ID": "zzz", "parentType": "pedro", "parentID": "yyy"}]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			session.DeleteEntity(e)

			Convey("Then ID should 'yyy'", func() {
				So(e.Identifier(), ShouldEqual, "yyy")
			})
		})

		Convey("When I delete it and I got a communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "nope", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			err := session.DeleteEntity(e)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSession_FetchChildren(t *testing.T) {

	Convey("Given I have an existing object", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("yyy")

		Convey("When I fetch its children with success", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[{"ID": "1"}, {"ID": "2"}]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			var l fakeObjectsList
			session.FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then the lenght of the children list should be 2", func() {
				So(len(l), ShouldEqual, 2)
			})

			Convey("Then the first child ID should be 1", func() {
				So(l[0].Identifier(), ShouldEqual, "1")
			})

			Convey("Then the second child ID should be 2", func() {
				So(l[1].Identifier(), ShouldEqual, "2")
			})
		})

		Convey("When I fetch its children while there is no data", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			var l fakeObjectsList
			session.FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then the lenght of the children list should be 0", func() {
				So(l, ShouldBeNil)
			})
		})

		Convey("When I fetch its children while there is none, but I still get an empty array", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			var l fakeObjectsList
			session.FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then the lenght of the children list should be 0", func() {
				So(len(l), ShouldEqual, 0)
			})
		})

		Convey("When I fetch the children and I got a communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "woops", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			var l fakeObjectsList
			err := session.FetchChildren(e, fakeIdentity, &l, nil)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I fetch the children I got a bad json", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `[not good]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			var l fakeObjectsList

			Convey("Then it should panic", func() {
				So(func() { session.FetchChildren(e, fakeIdentity, &l, nil) }, ShouldPanic)
			})
		})
	})
}

func TestSession_CreateChild(t *testing.T) {

	Convey("Given I have an existing object", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("xxx")

		Convey("When I create a child with success", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `[{"ID": "zzz"}]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			c := &fakeObject{}
			c.SetIdentity(fakeIdentity)
			session.CreateChild(e, c)

			Convey("Then ID of the new children should be zzz", func() {
				So(c.Identifier(), ShouldEqual, "zzz")
			})
		})

		Convey("When I create a child and I got a communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "woops", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			c := &fakeObject{}
			c.SetIdentity(fakeIdentity)
			err := session.CreateChild(e, c)

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I create a child I got a bad json", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, `[{"bad"}]`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			c := &fakeObject{}
			c.SetIdentity(fakeIdentity)

			Convey("Then it should panic", func() {
				So(func() { session.CreateChild(e, c) }, ShouldPanic)
			})
		})
	})
}

func TestSession_AssignChildren(t *testing.T) {

	Convey("Given I have two existing objects", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("xxx")

		c := &fakeObject{}
		c.SetIdentity(fakeIdentity)
		c.SetIdentifier("yyy")

		Convey("When I assign them with success", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			l := fakeObjectsList{c}
			session.AssignChildren(e, l, fakeIdentity)

			Convey("Then nothing special should happen", func() {
			})
		})

		Convey("When I assign them I got an communication error", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "woops", 500)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			l := fakeObjectsList{c}
			err := session.AssignChildren(e, l, fakeIdentity)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

/*
	Events
*/
func TestSession_NextEvent(t *testing.T) {

	Convey("When I use NextEvent and I receive a valid push notification", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"uuid": "y", "events": [{"type": "CREATE", "entityType": "thing", "updateMechanism": "DEFAULT", "entities": []}]}`)
		}))
		defer ts.Close()

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		session := NewSession("username", "password", "organization", ts.URL, r)

		lID := "x"
		var notif *Notification
		c := make(NotificationsChannel)
		go session.NextEvent(c, &lID)

		select {
		case notif = <-c:
		case <-time.After(10 * time.Millisecond):
		}

		Convey("Then notification should not be nil", func() {
			So(notif, ShouldNotBeNil)
		})

		Convey("Then last Event ID should be y", func() {
			So(lID, ShouldEqual, "y")
		})
	})

	Convey("When I use NextEvent and I receive an error", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "woops", 500)
		}))
		defer ts.Close()

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		session := NewSession("username", "password", "organization", ts.URL, r)

		lID := "x"
		var notif *Notification
		c := make(NotificationsChannel)
		go session.NextEvent(c, &lID)

		select {
		case notif = <-c:
		case <-time.After(10 * time.Millisecond):
		}

		Convey("Then notification should be nil", func() {
			So(notif, ShouldBeNil)
		})

		Convey("Then last Event ID should be the same", func() {
			So(lID, ShouldEqual, "x")
		})
	})

	Convey("When I use NextEvent and I receive a push notification with malformed json", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"not good}`)
		}))
		defer ts.Close()

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		session := NewSession("username", "password", "organization", ts.URL, r)

		lID := "x"
		var notif *Notification
		c := make(NotificationsChannel)
		go session.NextEvent(c, &lID)

		select {
		case notif = <-c:
		case <-time.After(10 * time.Millisecond):
		}

		Convey("Then notification should be nil", func() {
			So(notif, ShouldBeNil)
		})

		Convey("Then last Event ID should be the same", func() {
			So(lID, ShouldEqual, "x")
		})
	})
}

/*
	Send
*/
func TestSession_Send(t *testing.T) {

	Convey("Given I am authenticated", t, func() {

		r := &fakeRootObject{}
		r.SetIdentity(fakeRootIdentity)

		Convey("When I send a request that returns OK", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response status code should be 200", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I send a request that returns Created", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response status code should be 201", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusCreated)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I send a request that returns No Content", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response status code should be 204", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I send a request that returns Multiple Choices", func() {

			choiceMade := false
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				if choiceMade {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusMultipleChoices)
					choiceMade = true
				}
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response status code should be 200", func() {
				So(resp.StatusCode, ShouldEqual, http.StatusOK)
			})

			Convey("The choice should have been made", func() {
				So(choiceMade, ShouldBeTrue)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I send a request that returns Conflict", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusConflict)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"property": "prop", "type": "iznogood", "descriptions": [{"title": "oula", "description": "pas bon"}]}`)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response should be nil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the error Message should 'iznogood' and the Code should be Conflict", func() {
				So(string(err.Message), ShouldEqual, "iznogood")
				So(err.Code, ShouldEqual, http.StatusConflict)
			})
		})

		Convey("When I send a request that returns any other code", func() {

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer ts.Close()
			session := NewSession("username", "password", "organization", ts.URL, r)

			req, _ := http.NewRequest("GET", ts.URL, nil)
			resp, err := session.send(req, nil)

			Convey("Then response should be nil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the error Message should 'iznogood' and the Code should be StatusInternalServerError", func() {
				So(err.Message, ShouldEqual, "500 Internal Server Error")
				So(err.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
