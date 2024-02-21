package assets

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAsset(t *testing.T) {
	Convey("Given an asset exists", t, func() {

		assetPath := "redirects/redirects.csv"

		// get the file details (largely, size)
		fileInfo, err := os.Stat(assetPath)
		So(err, ShouldBeNil)

		Convey("When we request the asset", func() {

			asset, err := Asset(assetPath)

			Convey("Then the asset should have been obtained successfully", func() {
				So(err, ShouldBeNil)
				So(len(asset), ShouldEqual, fileInfo.Size())
			})
		})
	})
}
