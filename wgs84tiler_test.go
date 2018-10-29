package main

import (
	"log"
	"os"
	"testing"

	"github.com/disintegration/imaging"
)

func TestSlice(t *testing.T) {
	var testImage = "./testdatas/imagein.jpg"
	var testBounds = WGS84Bounds{top: 48.8795479494599, right: 2.18925349493902, left: 2.18108416654134, bottom: 48.8759617549741}

	var testOutputDir = "./testdatas/app/out/"

	err := os.RemoveAll(testOutputDir)
	if err != nil {
		log.Fatalf("failed to clean dir: %v", err)
	}

	src, err := imaging.Open(testImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	var ret int

	ret = Slice(src, testBounds, 18, testOutputDir)
	if ret != 24 {
		t.Errorf("At zomm %d Tiles generated are incorrect, got: %d, want: %d.", 24, ret, 24)
	}

	Slice(src, testBounds, 14, testOutputDir)
}

func TestZooms(t *testing.T) {
	var testImage = "./testdatas/imagein.jpg"
	var testBounds = WGS84Bounds{top: 48.8795479494599, right: 2.18925349493902, left: 2.18108416654134, bottom: 48.8759617549741}

	var testOutputDir = "./testdatas/app/out/"
	src, err := imaging.Open(testImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	Slice(src, testBounds, 11, testOutputDir)
	Slice(src, testBounds, 12, testOutputDir)
	Slice(src, testBounds, 13, testOutputDir)
	Slice(src, testBounds, 14, testOutputDir)
	Slice(src, testBounds, 15, testOutputDir)
	Slice(src, testBounds, 16, testOutputDir)
	Slice(src, testBounds, 17, testOutputDir)
	Slice(src, testBounds, 18, testOutputDir)
	Slice(src, testBounds, 19, testOutputDir)
	Slice(src, testBounds, 20, testOutputDir)

}

func TestMerge(t *testing.T) {
	var testOutputDir = "./testdatas/app/out/"

	var testImage1 = "./testdatas/imagein.jpg"
	var testBounds1 = WGS84Bounds{top: 48.8795479494599, right: 2.18925349493902, left: 2.18108416654134, bottom: 48.8759617549741}
	var testImage2 = "./testdatas/imagein2.jpg"
	var testBounds2 = WGS84Bounds{top: 48.8759508670196, right: 2.18926420416061, left: 2.18109546619253, bottom: 48.8723646719228}

	src1, err := imaging.Open(testImage1)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	src2, err := imaging.Open(testImage2)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	Slice(src1, testBounds1, 18, testOutputDir)
	Slice(src2, testBounds2, 18, testOutputDir)

	Slice(src1, testBounds1, 14, testOutputDir)
	Slice(src2, testBounds2, 14, testOutputDir)

	Slice(src1, testBounds1, 15, testOutputDir)
	Slice(src2, testBounds2, 15, testOutputDir)
}
