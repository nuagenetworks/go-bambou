package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResponse_NewResponse(t *testing.T) {

	Convey("Given I create a new Request", t, func() {
		r := NewResponse()

		Convey("Then Headers should should not be nil", func() {
			So(r.Headers, ShouldNotBeNil)
		})
	})
}

func TestResponse_SetGetHeader(t *testing.T) {

	Convey("Given I create a new Request", t, func() {
		r := NewResponse()

		Convey("When I set the header 'header' to 'value'", func() {
			r.SetHeader("header", "value")

			Convey("Then value of header should value", func() {
				So(r.GetHeader("header"), ShouldEqual, "value")
			})
		})
	})
}
