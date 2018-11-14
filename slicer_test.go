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
var zoom = 15
var bounds = WGS84Bounds{Top: 48.8687073004617, Right: 2.15657022586739, Left: 2.14840505163567, Bottom: 48.8651234503015}
var bounds2 = WGS84Bounds{Top: 48.8687073004617, Right: 2.1647354001, Left: 2.15657022586739, Bottom: 48.8651234503015}

func TestSlicing(t *testing.T) {
	os.RemoveAll("./testdatas/out/")
	os.MkdirAll("./testdatas/out/", 0777)

	src, err := imaging.Open(file1)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	nbSlices, n, m, _ := SliceIt(src, bounds, zoom, out)
	if nbSlices != 2 {
		t.Errorf("Expected  2 tiles  but %v was generated", nbSlices)
	}

	if n != 2 {
		t.Errorf("Expected  2 NEW tiles   but %v was generated", n)
	}

	if m != 0 {
		t.Errorf("Expected  0 MERGED tiles   but %v was generated", m)
	}
}

func TestMergeSlicing(t *testing.T) {
	os.RemoveAll("./testdatas/out/")
	os.MkdirAll("./testdatas/out/", 0777)

	src, err := imaging.Open(file1)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	_, _, _, _ = SliceIt(src, bounds, zoom, out)  // first call for making tiles
	_, n, m, _ := SliceIt(src, bounds, zoom, out) // second call for testin merge

	if n != 0 {
		t.Errorf("Expected  0 NEW tiles   but %v was generated", n)
	}

	if m != 2 {
		t.Errorf("Expected  2 MERGED tiles   but %v was generated", m)
	}

}

func TestJuxtaposeSlicing(t *testing.T) {
	os.RemoveAll("./testdatas/out/")
	os.MkdirAll("./testdatas/out/", 0777)

	src, _ := imaging.Open(file1)
	src2, _ := imaging.Open(file2)

	_, _, _, _ = SliceIt(src, bounds, zoom, out)    // first call for making tiles
	_, n, m, _ := SliceIt(src2, bounds2, zoom, out) // second call for testin merge

	if n != 1 {
		t.Errorf("Expected  1 NEW tiles   but %v was generated", n)
	}

	if m != 1 {
		t.Errorf("Expected  1 MERGED tiles   but %v was generated", m)
	}

}
