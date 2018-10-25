/*
Package `wgs84tiler` provides slicing functions for slippy map (tiled map) generation.

It output out/Z/X/Y.png files , where Z is the zoom level, and X and Y identify the tile.

Example:

	package main

	import (
		"image/jpeg"
		"log"
		"os"

		"github.com/konogan/wgs84tiler"
	)

	func main() {
		// Open the test image.
		f, err := os.Open("testdata/imagein.jpg")
		if err != nil {
			log.Fatalf("os.Open failed: %v", err)
		}

  	// define the boundaries of this image
  	var b = WGS84Bounds{top: 48.8687073004617, right: 2.15657022586739, left: 2.14840505163567,
    	bottom: 48.8651234503015}

		// Slice the image at zoom level 14 and put all slices in out folfer
		nbSlices := wgs84tiler.Slice(f,b,14,"./out/")
	}


*/
package wgs84tiler
