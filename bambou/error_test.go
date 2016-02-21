package bambou

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError_NewError(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("Then Code should be 42", func() {
			So(e.Code, ShouldEqual, 42)
		})

		Convey("Then Message should be 'a big error'", func() {
			So(e.Message, ShouldEqual, "a big error")
		})

	})
}

func TestError_Error(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("Then Error() output should be should be '<Error: 42, message: a big error>'", func() {
			So(e.Error(), ShouldEqual, "<Error: 42, message: a big error>")
		})
	})
}

func TestError_FromJSON(t *testing.T) {

	Convey("Given I create a new Error", t, func() {
		e := NewError(42, "a big error")

		Convey("When I unmarshal some data", func() {
			d := "{\"property\": \"prop\", \"type\": \"iznogood\", \"descriptions\": [{\"title\": \"oula\", \"description\": \"pas bon\"}]}"
			json.Unmarshal([]byte(d), e)

			Convey("Then the Message should be 'iznogood'", func() {
				So(e.Message, ShouldEqual, "iznogood")
			})

			Convey("Then the Property should be 'prop'", func() {
				So(e.Property, ShouldEqual, "prop")
			})

			Convey("Then the lenght of Descriptions should be 1", func() {
				So(len(e.Descriptions), ShouldEqual, 1)
			})

			Convey("When I retrieve the content of Descriptions", func() {
				d := e.Descriptions[0]

				Convey("Then the Title should be 'oula'", func() {
					So(d.Title, ShouldEqual, "oula")
				})

				Convey("Then the Description should be ' bon'", func() {
					So(d.Description, ShouldEqual, "pas bon")
				})

			})
		})
	})
}
