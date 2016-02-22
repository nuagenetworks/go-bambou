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

func TestConnection_NewError(t *testing.T) {

	Convey("Given I create a new Connection", t, func() {
		c := newConnection()

		Convey("Then Timeout should 60s", func() {
			So(c.Timeout, ShouldEqual, time.Duration(60)*time.Second)
		})
	})
}

func TestConnection_Start(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "url", &testRoot{})
	}).restore()

	defer func() {
		CurrentSession().Reset()
	}()

	Convey("Given I create a new Connection", t, func() {
		c := newConnection()

		Convey("When I start a connection with the request that returns responseCodeSuccess", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeSuccess,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns responseCodeCreated", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeCreated,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns responseCodeEmpty", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeEmpty,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns responseCodeMultipleChoices", func() {

			reqCount := 0
			defer patch(&sendNativerequest, func(request *request) *response {
				reqCount++

				if reqCount == 1 {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeMultipleChoices,
					}
				}
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeSuccess,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the request URL should be http://fake.com/hello?responseChoice=1", func() {
				So(req.URL, ShouldEqual, "http://fake.com/hello?responseChoice=1")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeAuthenticationExpired", func() {

			reqCount := 0
			defer patch(&sendNativerequest, func(request *request) *response {
				reqCount++

				if reqCount == 1 {
					return &response{
						Headers: make(map[string]string),
						Code:    responseCodeAuthenticationExpired,
					}
				}
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeSuccess,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns responseCodeConflict", func() {

			d := "{\"property\": \"prop\", \"type\": \"iznogood\", \"descriptions\": [{\"title\": \"oula\", \"description\": \"pas bon\"}]}"

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeConflict,
					Data:    []byte(d),
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be nil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the error Message should 'iznogood'", func() {
				So(string(err.Message), ShouldEqual, "iznogood")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeBadrequest", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeBadrequest,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeBadrequest", func() {
				So(err.Code, ShouldEqual, responseCodeBadrequest)
			})

			Convey("Then error Message should be 'Bad request.'", func() {
				So(err.Message, ShouldEqual, "Bad request.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeUnauthorized", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeUnauthorized,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeUnauthorized", func() {
				So(err.Code, ShouldEqual, responseCodeUnauthorized)
			})

			Convey("Then error Message should be 'Unauthorized.'", func() {
				So(err.Message, ShouldEqual, "Unauthorized.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodePermissionDenied", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodePermissionDenied,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodePermissionDenied", func() {
				So(err.Code, ShouldEqual, responseCodePermissionDenied)
			})

			Convey("Then error Message should be 'Permission denied.'", func() {
				So(err.Message, ShouldEqual, "Permission denied.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeMethodNotAllowed", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeMethodNotAllowed,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeMethodNotAllowed", func() {
				So(err.Code, ShouldEqual, responseCodeMethodNotAllowed)
			})

			Convey("Then error Message should be 'Not allowed.'", func() {
				So(err.Message, ShouldEqual, "Not allowed.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeNotFound", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeNotFound,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeNotFound", func() {
				So(err.Code, ShouldEqual, responseCodeNotFound)
			})

			Convey("Then error Message should be 'Not found.'", func() {
				So(err.Message, ShouldEqual, "Not found.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeConnectionTimeout", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeConnectionTimeout,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeConnectionTimeout", func() {
				So(err.Code, ShouldEqual, responseCodeConnectionTimeout)
			})

			Convey("Then error Message should be 'Timeout.'", func() {
				So(err.Message, ShouldEqual, "Timeout.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodePreconditionFailed", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodePreconditionFailed,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodePreconditionFailed", func() {
				So(err.Code, ShouldEqual, responseCodePreconditionFailed)
			})

			Convey("Then error Message should be 'Precondition failed.'", func() {
				So(err.Message, ShouldEqual, "Precondition failed.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeInternalServerError", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeInternalServerError,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeInternalServerError", func() {
				So(err.Code, ShouldEqual, responseCodeInternalServerError)
			})

			Convey("Then error Message should be 'Internal server error.'", func() {
				So(err.Message, ShouldEqual, "Internal server error.")
			})
		})

		Convey("When I start a connection with the request that returns responseCodeServiceUnavailable", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    responseCodeServiceUnavailable,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeServiceUnavailable", func() {
				So(err.Code, ShouldEqual, responseCodeServiceUnavailable)
			})

			Convey("Then error Message should be 'Service unavailable.'", func() {
				So(err.Message, ShouldEqual, "Service unavailable.")
			})
		})

		Convey("When I start a connection with the request that returns an unknown code", func() {

			defer patch(&sendNativerequest, func(request *request) *response {
				return &response{
					Headers: make(map[string]string),
					Code:    666,
				}
			}).restore()

			req := newRequest("http://fake.com/hello")
			resp, err := c.start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be responseCodeServiceUnavailable", func() {
				So(err.Code, ShouldEqual, 666)
			})

			Convey("Then error Message should be 'Service unavailable.'", func() {
				So(err.Message, ShouldEqual, "Unknown error.")
			})
		})

	})
}

func TestConnection_sendNativerequest(t *testing.T) {

	Convey("Given I create a sucessful request to github", t, func() {

		req := newRequest("http://jsonplaceholder.typicode.com/post/1")
		resp := sendNativerequest(req)

		Convey("Then I should get a response with code 404", func() {
			So(resp.Code, ShouldEqual, 404)
		})
	})

	Convey("Given I create a bad request", t, func() {

		req := newRequest("https:///nope/nope/nope")

		Convey("Then I it should panic", func() {
			So(func() { sendNativerequest(req) }, ShouldPanic)
		})
	})
}
