package wgs84tiler

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"time"
)

//TimeTrack debug timer
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func main(image image.Image, wgs84Bounds WGS84Bounds, zoom int, outputdir string) int {
	defer TimeTrack(time.Now(), "tiler")

	report := sliceIt(image, wgs84Bounds, zoom, outputdir)

	fmt.Print(&buf)
	return report
}
