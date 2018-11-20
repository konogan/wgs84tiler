package wgs84tiler

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
)

var file1 = "./testdatas/imagein.jpg"
var file2 = "./testdatas/imagein2.jpg"
var out = filepath.Dir(file1) + "/out/"
var zoom = 18
var bounds = WGS84Bounds{Top: 48.8687073004617, Right: 2.15657022586739, Left: 2.14840505163567, Bottom: 48.8651234503015}
var bounds2 = WGS84Bounds{Top: 48.8687073004617, Right: 2.1647354001, Left: 2.15657022586739, Bottom: 48.8651234503015}

func TestSlicing(t *testing.T) {
	os.RemoveAll("./testdatas/out/")
	os.MkdirAll("./testdatas/out/", 0777)

	src, err := imaging.Open(file1)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	nbSlices, _ := SliceIt(src, bounds, zoom, out)
	if nbSlices != 35 {
		t.Errorf("At zoom %d : expected  35 tiles  but %v was generated", zoom, nbSlices)
	}
}
