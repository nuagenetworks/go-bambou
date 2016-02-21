package bambou

import (
	"testing"

	"github.com/ccding/go-logging/logging"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogger_Logger(t *testing.T) {

	Convey("Given I retrieve the Logger", t, func() {
		l := Logger()

		Convey("Then the Level should be logging.ERROR", func() {
			So(l.Level(), ShouldEqual, logging.ERROR)
		})

		Convey("Then the Name should be 'bambou", func() {
			So(l.Name(), ShouldEqual, "bambou")
		})
	})
}
