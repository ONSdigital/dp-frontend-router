package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitHelpers(t *testing.T) {
	Convey("test ReturnSecondSegmentFromPath", t, func() {
		Convey("extracts second segment from path", func() {
			second, err := ReturnSecondSegmentFromPath("/first/second")
			So(err, ShouldBeNil)
			So(second, ShouldEqual, "second")

			second, err = ReturnSecondSegmentFromPath("/first/a-different-second/third/forth")
			So(err, ShouldBeNil)
			So(second, ShouldEqual, "a-different-second")
		})

		Convey("returns an error if it is unable to extract the information", func() {
			second, err := ReturnSecondSegmentFromPath("invalid")
			So(err, ShouldBeError, "unable to extract secondSegment from path: invalid")
			So(second, ShouldEqual, "")
		})
	})
}
