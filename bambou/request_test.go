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

func TestRequest_newRequest(t *testing.T) {

	Convey("Given I create a new request", t, func() {
		r := newRequest("https://fake.com")

		Convey("Then URL should https://fake.com", func() {
			So(r.URL, ShouldEqual, "https://fake.com")
		})

		Convey("Then Method should GET", func() {
			So(r.Method, ShouldEqual, requestMethodGet)
		})

		Convey("Then Headers should not be nil", func() {
			So(r.Headers, ShouldNotBeNil)
		})

		Convey("Then Parameters should not be nil", func() {
			So(r.Parameters, ShouldNotBeNil)
		})
	})
}

func TestRequest_SetGetHeader(t *testing.T) {

	Convey("Given I create a new request", t, func() {
		r := newRequest("https://fake.com")

		Convey("When I set the header 'header' to 'value'", func() {
			r.setHeader("header", "value")

			Convey("Then value of header should be value", func() {
				So(r.getHeader("header"), ShouldEqual, "value")
			})
		})
	})
}

func TestRequest_SetGetParameter(t *testing.T) {

	Convey("Given I create a new request", t, func() {
		r := newRequest("https://fake.com")

		Convey("When I set the parameter 'param' to 'value'", func() {
			r.setParameter("param", "value")

			Convey("Then the value of parameter 'param' should 'value", func() {
				So(r.getParameter("param"), ShouldEqual, "value")
			})
		})
	})
}

func TestRequest_toNative(t *testing.T) {

	Convey("Given I create new request with default values", t, func() {
		r := newRequest("https://fake.com")
		r.setHeader("header", "value")
		r.Data = []byte("hello")

		Convey("When I convert the request to the native request", func() {
			n := r.toNative()

			Convey("Then URL should https://fake.com", func() {
				So(n.URL.String(), ShouldEqual, "https://fake.com")
			})

			Convey("Then Header 'header' should be 'value'", func() {
				So(n.Header.Get("header"), ShouldEqual, "value")
			})
		})
	})
}
