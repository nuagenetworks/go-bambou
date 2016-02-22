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
		c := NewConnection()

		Convey("Then Timeout should 60s", func() {
			So(c.Timeout, ShouldEqual, time.Duration(60)*time.Second)
		})
	})
}

func TestConnection_Start(t *testing.T) {

	defer patch(&CurrentSession, func() *Session {
		return NewSession("username", "password", "organization", "url", nil)
	}).restore()

	Convey("Given I create a new Connection", t, func() {
		c := NewConnection()

		Convey("When I start a connection with the request that returns ResponseCodeSuccess", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeSuccess,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeCreated", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeCreated,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeEmpty", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeEmpty,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should not be nil", func() {
				So(resp, ShouldNotBeNil)
			})

			Convey("Then error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeMultipleChoices", func() {

			reqCount := 0
			defer patch(&sendNativeRequest, func(request *Request) *Response {
				reqCount++

				if reqCount == 1 {
					return &Response{
						Headers: make(map[string]string),
						Code:    ResponseCodeMultipleChoices,
					}
				}
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeSuccess,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

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

		Convey("When I start a connection with the request that returns ResponseCodeConflict", func() {

			d := "{\"property\": \"prop\", \"type\": \"iznogood\", \"descriptions\": [{\"title\": \"oula\", \"description\": \"pas bon\"}]}"

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeConflict,
					Data:    []byte(d),
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

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

		Convey("When I start a connection with the request that returns ResponseCodeBadRequest", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeBadRequest,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeBadRequest", func() {
				So(err.Code, ShouldEqual, ResponseCodeBadRequest)
			})

			Convey("Then error Message should be 'Bad request.'", func() {
				So(err.Message, ShouldEqual, "Bad request.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeUnauthorized", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeUnauthorized,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeUnauthorized", func() {
				So(err.Code, ShouldEqual, ResponseCodeUnauthorized)
			})

			Convey("Then error Message should be 'Unauthorized.'", func() {
				So(err.Message, ShouldEqual, "Unauthorized.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodePermissionDenied", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodePermissionDenied,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodePermissionDenied", func() {
				So(err.Code, ShouldEqual, ResponseCodePermissionDenied)
			})

			Convey("Then error Message should be 'Permission denied.'", func() {
				So(err.Message, ShouldEqual, "Permission denied.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeMethodNotAllowed", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeMethodNotAllowed,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeMethodNotAllowed", func() {
				So(err.Code, ShouldEqual, ResponseCodeMethodNotAllowed)
			})

			Convey("Then error Message should be 'Not allowed.'", func() {
				So(err.Message, ShouldEqual, "Not allowed.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeNotFound", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeNotFound,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeNotFound", func() {
				So(err.Code, ShouldEqual, ResponseCodeNotFound)
			})

			Convey("Then error Message should be 'Not found.'", func() {
				So(err.Message, ShouldEqual, "Not found.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeConnectionTimeout", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeConnectionTimeout,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeConnectionTimeout", func() {
				So(err.Code, ShouldEqual, ResponseCodeConnectionTimeout)
			})

			Convey("Then error Message should be 'Timeout.'", func() {
				So(err.Message, ShouldEqual, "Timeout.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodePreconditionFailed", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodePreconditionFailed,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodePreconditionFailed", func() {
				So(err.Code, ShouldEqual, ResponseCodePreconditionFailed)
			})

			Convey("Then error Message should be 'Precondition failed.'", func() {
				So(err.Message, ShouldEqual, "Precondition failed.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeInternalServerError", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeInternalServerError,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeInternalServerError", func() {
				So(err.Code, ShouldEqual, ResponseCodeInternalServerError)
			})

			Convey("Then error Message should be 'Internal server error.'", func() {
				So(err.Message, ShouldEqual, "Internal server error.")
			})
		})

		Convey("When I start a connection with the request that returns ResponseCodeServiceUnavailable", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    ResponseCodeServiceUnavailable,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeServiceUnavailable", func() {
				So(err.Code, ShouldEqual, ResponseCodeServiceUnavailable)
			})

			Convey("Then error Message should be 'Service unavailable.'", func() {
				So(err.Message, ShouldEqual, "Service unavailable.")
			})
		})

		Convey("When I start a connection with the request that returns an unknown code", func() {

			defer patch(&sendNativeRequest, func(request *Request) *Response {
				return &Response{
					Headers: make(map[string]string),
					Code:    666,
				}
			}).restore()

			req := NewRequest("http://fake.com/hello")
			resp, err := c.Start(req)

			Convey("Then response should be bil", func() {
				So(resp, ShouldBeNil)
			})

			Convey("Then error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then error Code should be ResponseCodeServiceUnavailable", func() {
				So(err.Code, ShouldEqual, 666)
			})

			Convey("Then error Message should be 'Service unavailable.'", func() {
				So(err.Message, ShouldEqual, "Unknown error.")
			})
		})

	})
}
