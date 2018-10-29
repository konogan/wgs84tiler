package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

//Slice execute the slicing process on the image
//
// Examples:
//
//	countSlice = wgs84tiler.Slice(image, wgs84Bounds, 18, "./out")
//
func Slice(image image.Image, wgs84Bounds WGS84Bounds, zoom int, outputdir string) int {
	defer timeTrack(time.Now(), "tiler")

	report := sliceIt(image, wgs84Bounds, zoom, outputdir)

	fmt.Print(&buf)
	return report
}

func main() {

}
