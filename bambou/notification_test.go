package bambou

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotification_NewNotification(t *testing.T) {

	Convey("Given I create a new Notification", t, func() {
		n := NewNotification()

		Convey("Then Events should not be nil", func() {
			So(n.Events, ShouldNotBeNil)
		})

		Convey("Then UUID should not be nil", func() {
			So(n.UUID, ShouldNotBeNil)
		})
	})
}

func TestNotification_FromJSON(t *testing.T) {

	Convey("Given I create a new notification", t, func() {
		n := NewNotification()

		Convey("When I unmarshal son json data", func() {
			d := "{\"uuid\": \"007\", \"events\": [{\"entityType\": \"cat\", \"type\": \"UPDATE\", \"updateMechanism\": \"useless\", \"entities\":[{\"name\": \"hello\"}]}]}"
			json.Unmarshal([]byte(d), &n)

			Convey("Then UUI should be '007'", func() {
				So(n.UUID, ShouldEqual, "007")
			})

			Convey("Then lenght of Events should be 1", func() {
				So(len(n.Events), ShouldEqual, 1)
			})

			Convey("When I retrieve the Events", func() {
				e := n.Events[0]

				Convey("Then EntityType should be cat", func() {
					So(e.EntityType, ShouldEqual, "cat")
				})

				Convey("Then Type should UPDATE", func() {
					So(e.Type, ShouldEqual, "UPDATE")
				})

				Convey("Then UpdateMechanism should useless", func() {
					So(e.UpdateMechanism, ShouldEqual, "useless")
				})

				Convey("Then the lenght of DataMap should be 1", func() {
					So(len(e.DataMap), ShouldEqual, 1)
				})

				Convey("Then the value of item 0 of DataMap should hello", func() {
					So(e.DataMap[0]["name"], ShouldEqual, "hello")
				})
			})
		})
	})
}
