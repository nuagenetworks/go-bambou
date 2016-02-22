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
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPushCenter_NewPushCenter(t *testing.T) {

	Convey("Given I create a new PushCenter", t, func() {
		p := NewPushCenter()

		Convey("Then the Channel should not be nil", func() {
			So(p.Channel, ShouldNotBeNil)
		})

		Convey("Then the stop channel should not be nil", func() {
			So(p.stop, ShouldNotBeNil)
		})

		Convey("Then the handlers list should not be nil", func() {
			So(p.handlers, ShouldNotBeNil)
		})
	})
}

func TestPushCenter_HandlersRegistration(t *testing.T) {

	Convey("Given I create a new PushCenter", t, func() {

		p := NewPushCenter()
		h := func(*Event) {}

		Convey("When I register a handler for a identity", func() {
			p.RegisterHandlerForIdentity(h, fakeIdentity)

			Convey("Then it should be registered in the list", func() {
				So(p.HasHandlerForIdentity(fakeIdentity), ShouldBeTrue)
			})
		})

		Convey("When I unregister a handler for the identity", func() {
			p.RegisterHandlerForIdentity(h, fakeIdentity)
			p.UnregisterHandlerForIdentity(fakeIdentity)

			Convey("Then it should not be registered in the list anymore", func() {
				So(p.HasHandlerForIdentity(fakeIdentity), ShouldBeFalse)
			})

			Convey("Then the default handler should be nil", func() {
				So(p.defaultHander, ShouldBeNil)
			})
		})

		Convey("When I register handler for the all identity", func() {
			p.RegisterHandlerForIdentity(h, AllIdentity)

			Convey("Then it should not be registered in the list", func() {
				So(p.HasHandlerForIdentity(AllIdentity), ShouldBeTrue)
			})

			Convey("Then it should be set as the defaultHandler", func() {
				So(p.defaultHander, ShouldEqual, h)
			})
		})

		Convey("When I unregister the handler for the all identity", func() {
			p.UnregisterHandlerForIdentity(AllIdentity)

			Convey("Then it should not be registered in the lista anymore", func() {
				So(p.HasHandlerForIdentity(AllIdentity), ShouldBeFalse)
			})

			Convey("Then it should not be set as the defaultHandler anymore", func() {
				So(p.defaultHander, ShouldBeNil)
			})
		})
	})
}

func TestPushCenter_Start(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new PushCenter and resgister a handler", t, func() {

		n := 0
		c := 0
		p := NewPushCenter()
		h1 := func(*Event) { n++ }
		h2 := func(*Event) { n++ }
		p.RegisterHandlerForIdentity(h1, AllIdentity)
		p.RegisterHandlerForIdentity(h2, fakeIdentity)

		Convey("When I start a running push center with a handler", func() {

			defer func() { p.Stop() }()
			defer patch(&sendNativeRequest, func(request *request) *response {
				c++
				if c == 1 {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
						Data:    []byte("{\"uuid\": \"y\", \"events\": [{\"type\": \"CREATE\", \"entityType\": \"fake\", \"updateMechanism\": \"DEFAULT\", \"entities\": [{}]}]}"),
					}
				} else if c == 2 {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
						Data:    []byte("{\"uuid\": \"y\", \"events\": [{\"type\": \"CREATE\", \"entityType\": \"thing\", \"updateMechanism\": \"DEFAULT\", \"entities\": [{}]}]}"),
					}
				} else {
					time.Sleep(1 * time.Second)
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
						Data:    []byte("{\"uuid\": \"y\", \"events\": [{\"type\": \"CREATE\", \"entityType\": \"thing\", \"updateMechanism\": \"DEFAULT\", \"entities\": [{}]}]}"),
					}
				}
			}).restore()

			p.Start()

			time.Sleep(10 * time.Millisecond)

			Convey("Then the number of notification should be 3", func() {
				So(n, ShouldEqual, 3)
			})
		})
	})
}

func TestPushCenter_Stop(t *testing.T) {

	Convey("Given I create a new PushCenter", t, func() {

		p := NewPushCenter()
		p.isRunning = true
		p.lastEventID = "x"

		Convey("When I stop a running push center", func() {

			var stopValue bool
			go p.Stop()

			select {
			case stopValue = <-p.stop:
			case <-time.After(10 * time.Millisecond):
			}

			Convey("Then the stop channel should get a true value", func() {
				So(stopValue, ShouldBeTrue)
			})

			Convey("Then isRunning should false", func() {
				So(p.isRunning, ShouldBeFalse)
			})

			Convey("Then lastEventID should empty", func() {
				So(p.lastEventID, ShouldEqual, "")
			})
		})
	})
}

func TestPushCenter_Listen(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "http://fake.com", nil)
	}).restore()

	Convey("Given I create a new PushCenter", t, func() {

		var notif *Notification
		p := NewPushCenter()

		Convey("Given I receive a push notification", func() {

			Convey("When I receive a push notification and the pushcenter is not runnning", func() {

				defer patch(&sendNativeRequest, func(request *request) *response {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
					}
				}).restore()

				p.lastEventID = "x"

				go p.listen()

				select {
				case notif = <-p.Channel:
				case <-time.After(10 * time.Millisecond):
				}

				Convey("Then notification should not be nil", func() {
					So(notif, ShouldBeNil)
				})

				Convey("Then last Event ID should be the same", func() {
					So(p.lastEventID, ShouldEqual, "x")
				})
			})

			Convey("When I receive a valid push notification", func() {

				defer patch(&sendNativeRequest, func(request *request) *response {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
						Data:    []byte("{\"uuid\": \"y\", \"events\": [{\"type\": \"CREATE\", \"entityType\": \"thing\", \"updateMechanism\": \"DEFAULT\", \"entities\": []}]}"),
					}
				}).restore()

				p.lastEventID = "x"
				p.isRunning = true
				go p.listen()

				select {
				case notif = <-p.Channel:
				case <-time.After(10 * time.Millisecond):
				}

				Convey("Then notification should not be nil", func() {
					So(notif, ShouldNotBeNil)
				})

				Convey("Then last Event ID should be the y", func() {
					So(p.lastEventID, ShouldEqual, "y")
				})
			})

			Convey("When I receive an error", func() {

				defer patch(&sendNativeRequest, func(request *request) *response {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeInternalServerError,
					}
				}).restore()

				p.lastEventID = "x"
				p.isRunning = true
				go p.listen()

				select {
				case notif = <-p.Channel:
				case <-time.After(10 * time.Millisecond):
				}

				Convey("Then notification should be nil", func() {
					So(notif, ShouldBeNil)
				})

				Convey("Then last Event ID should be the same", func() {
					So(p.lastEventID, ShouldEqual, "x")
				})
			})

			Convey("When I receive a push notification with malformed json", func() {

				defer patch(&sendNativeRequest, func(request *request) *response {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeSuccess,
						Data:    []byte("not a valid json"),
					}
				}).restore()

				p.lastEventID = "x"
				p.isRunning = true
				go p.listen()

				select {
				case notif = <-p.Channel:
				case <-time.After(10 * time.Millisecond):
				}

				Convey("Then notification should be nil", func() {
					So(notif, ShouldBeNil)
				})

				Convey("Then last Event ID should be the same", func() {
					So(p.lastEventID, ShouldEqual, "x")
				})
			})
		})
	})
}
