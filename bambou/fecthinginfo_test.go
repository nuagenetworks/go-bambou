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

func TestFetchingInfo_NewFetchingInfo(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()

		Convey("Then Page should be -1", func() {
			So(f.Page, ShouldEqual, -1)
		})

		Convey("Then PageSize should -1", func() {
			So(f.PageSize, ShouldEqual, -1)
		})
	})
}

func TestFetchingInfo_String(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()

		Convey("When I set some values", func() {

			f.Filter = "filer"
			f.Page = 2
			f.PageSize = 50

			Convey("Then string representation should <FetchingInfo page: 2, pagesize: 50, totalcount: 0>", func() {
				So(f.String(), ShouldEqual, "<FetchingInfo page: 2, pagesize: 50, totalcount: 0>")
			})
		})
	})
}

func TestFetchingInfo_prepareHeaders(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()
		r := newRequest("http://fake.com")

		Convey("When I prepareHeaders with a no fetching info", func() {
			prepareHeaders(r, nil)

			Convey("Then I should not have a value for X-Nuage-Page", func() {
				So(r.getHeader("X-Nuage-Page"), ShouldEqual, "")
			})

			Convey("Then I should have a the X-Nuage-PageSize set to 50", func() {
				So(r.getHeader("X-Nuage-PageSize"), ShouldEqual, "50")
			})

			Convey("Then I should not have a value for X-Nuage-Filter", func() {
				So(r.getHeader("X-Nuage-Filter"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-OrderBy", func() {
				So(r.getHeader("X-Nuage-OrderBy"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-GroupBy", func() {
				So(r.getHeader("X-Nuage-GroupBy"), ShouldEqual, "")
			})

			Convey("Then I should not have a value for X-Nuage-Attributes", func() {
				So(r.getHeader("X-Nuage-Attributes"), ShouldEqual, "")
			})
		})

		Convey("When I prepareHeaders witha fetching info that has a all fields", func() {
			f.Page = 2
			f.PageSize = 42
			f.Filter = "filter"
			f.OrderBy = "orderby"
			f.GroupBy = []string{"group1", "group2"}

			prepareHeaders(r, f)

			Convey("Then I should have a the X-Nuage-Page set to 2", func() {
				So(r.getHeader("X-Nuage-Page"), ShouldEqual, "2")
			})

			Convey("Then I should have a the X-Nuage-PageSize set to 42", func() {
				So(r.getHeader("X-Nuage-PageSize"), ShouldEqual, "42")
			})

			Convey("Then I should have a value for X-Nuage-Filter set to 'filter'", func() {
				So(r.getHeader("X-Nuage-Filter"), ShouldEqual, "filter")
			})

			Convey("Then I should have a value for X-Nuage-OrderBy set to 'orderby'", func() {
				So(r.getHeader("X-Nuage-OrderBy"), ShouldEqual, "orderby")
			})

			Convey("Then I should have a value for X-Nuage-GroupBy set to true", func() {
				So(r.getHeader("X-Nuage-GroupBy"), ShouldEqual, "true")
			})

			Convey("Then I should have a value for X-Nuage-Attributes contains group1 and group2", func() {
				So(r.getHeader("X-Nuage-Attributes"), ShouldEqual, "group1, group2")
			})
		})

	})
}

func TestFetchingInfo_readHeaders(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()
		r := newResponse()

		r.setHeader("X-Nuage-Page", "3")
		r.setHeader("X-Nuage-PageSize", "42")
		r.setHeader("X-Nuage-Filter", "filter")
		r.setHeader("X-Nuage-FilterType", "text")
		r.setHeader("X-Nuage-OrderBy", "value")
		r.setHeader("X-Nuage-Count", "666")

		Convey("When I readHeaders with a no fetching info", func() {
			readHeaders(r, nil)

			Convey("Then nothing should happen", func() {
			})
		})

		Convey("When I readHeaders with a request", func() {
			readHeaders(r, f)

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
