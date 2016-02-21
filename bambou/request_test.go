package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRequest_NewRequest(t *testing.T) {

	Convey("Given I create a new Request", t, func() {
		r := NewRequest("https://fake.com")

		Convey("Then URL should https://fake.com", func() {
			So(r.URL, ShouldEqual, "https://fake.com")
		})

		Convey("Then Method should GET", func() {
			So(r.Method, ShouldEqual, RequestMethodGet)
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

	Convey("Given I create a new Request", t, func() {
		r := NewRequest("https://fake.com")

		Convey("When I set the header 'header' to 'value'", func() {
			r.SetHeader("header", "value")

			Convey("Then value of header should be value", func() {
				So(r.GetHeader("header"), ShouldEqual, "value")
			})
		})
	})
}

func TestRequest_SetGetParameter(t *testing.T) {

	Convey("Given I create a new request", t, func() {
		r := NewRequest("https://fake.com")

		Convey("When I set the parameter 'param' to 'value'", func() {
			r.SetParameter("param", "value")

			Convey("Then the value of parameter 'param' should 'value", func() {
				So(r.GetParameter("param"), ShouldEqual, "value")
			})
		})
	})
}

func TestReques_ToNative(t *testing.T) {

	Convey("Given I create new request with default values", t, func() {
		r := NewRequest("https://fake.com")
		r.SetHeader("header", "value")
		r.Data = []byte("hello")

		Convey("When I convert the request to the native request", func() {
			n := r.ToNative()

			Convey("Then URL should https://fake.com", func() {
				So(n.URL.String(), ShouldEqual, "https://fake.com")
			})

			Convey("Then Header 'header' should be 'value'", func() {
				So(n.Header.Get("header"), ShouldEqual, "value")
			})
		})
	})
}
