package bambou

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIdentity_AllIdentity(t *testing.T) {

	Convey("Given I retrieve the AllIdentity", t, func() {
		i := AllIdentity

		Convey("Then RESTName should __all__", func() {
			So(i.RESTName, ShouldEqual, "__all__")
		})

		Convey("Then ResourceName should __all__", func() {
			So(i.ResourceName, ShouldEqual, "__all__")
		})
	})
}

func TestIdentity_NewIdentity(t *testing.T) {

	Convey("Given I create a new identity", t, func() {
		i := Identity{"restname", "resourcename"}

		Convey("Then RESTName should restname", func() {
			So(i.RESTName, ShouldEqual, "restname")
		})

		Convey("Then ResourceName should resourcename", func() {
			So(i.ResourceName, ShouldEqual, "resourcename")
		})
	})
}

func TestIdentity_String(t *testing.T) {

	Convey("Given I create a new identity", t, func() {
		i := Identity{"restname", "resourcename"}

		Convey("Then String should <Identity restname|resourcename>", func() {
			So(i.String(), ShouldEqual, "<Identity restname|resourcename>")
		})
	})
}

func TestIdentity_Identify(t *testing.T) {

	Convey("Given I create a new identity", t, func() {

		i := Identity{
			RESTName:     "restname",
			ResourceName: "resourcename",
		}

		Convey("When I Identity some objects with the new Identity", func() {

			l := []*RemoteObject{&RemoteObject{}, &RemoteObject{}}
			identify(&l, i)

			Convey("Then all objects should have the correct identity", func() {
				So(l[0].Identity(), ShouldResemble, i)
				So(l[1].Identity(), ShouldResemble, i)
			})

		})
	})
}
