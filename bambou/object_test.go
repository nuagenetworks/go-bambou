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

func TestExposedObject_Identifier(t *testing.T) {

	Convey("Given I create a new object", t, func() {

		e := &fakeObject{}

		Convey("When I set the ID", func() {
			e.SetIdentifier("xxx")

			Convey("Then ID should return 'xxx'", func() {
				So(e.Identifier(), ShouldEqual, "xxx")
			})
		})

		Convey("When I don't set the ID", func() {

			Convey("Then ID should ''", func() {
				So(e.Identifier(), ShouldEqual, "")
			})
		})
	})
}

func TestExposedObject_SetGetIdentity(t *testing.T) {

	Convey("Given I create a new object", t, func() {

		e := &fakeObject{}

		Convey("When I set the identity", func() {
			e.SetIdentity(fakeIdentity)

			Convey("Then Identity should fake", func() {
				So(e.Identity(), ShouldResemble, fakeIdentity)
			})
		})
	})
}

func TestExposedObject_String(t *testing.T) {

	Convey("Given I create a new object", t, func() {

		e := &fakeObject{}
		e.SetIdentity(fakeIdentity)
		e.SetIdentifier("xxx")

		Convey("Then the string representation should be <ExposedObject fake:xxx>", func() {
			So(e.String(), ShouldEqual, "<ExposedObject fake:xxx>")
		})
	})
}
