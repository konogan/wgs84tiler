package wgs84tiler

import (
	"log"
	"testing"

	"github.com/disintegration/imaging"
)

func TestMain(t *testing.T) {
	var testImage = "./testdatas/imagein.jpg"
	var testBounds = WGS84Bounds{top: 48.8687073004617, right: 2.15657022586739, left: 2.14840505163567,
		bottom: 48.8651234503015}
	var testOutputDir = "./testdatas/out/"

	src, err := imaging.Open(testImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	ret := main(src, testBounds, 18, testOutputDir)
	if ret != 24 {
		t.Errorf("Tiles generated are incorrect, got: %d, want: %d.", ret, 24)
	}
}
