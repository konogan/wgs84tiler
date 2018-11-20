package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	wgs "github.com/konogan/wgs84tiler"
)

// command line version of wgs84tiler
// take parameters :
// -f the file to slice
// -zoom the zoom level from 12 to 21
// -out the folder for output (optionnal, default is dir from input)
// -top : WGS84 top latitude
// -bottom : WGS84 bottom latitude
// -left : WGS84 left longitude
// -right : WGS84 right longitude
func main() {

	var file string
	var top float64
	var bottom float64
	var left float64
	var right float64
	var zoom int
	var out string

	flag.StringVar(&file, "f", "", "imgage to slice. (Required)")
	flag.Float64Var(&top, "top", 0, "top, WGS84 latitude. (Required)")
	flag.Float64Var(&bottom, "bottom", 0, "bottom, WGS84 latitude. (Required)")
	flag.Float64Var(&left, "left", 0, "left, WGS84 longitude. (Required)")
	flag.Float64Var(&right, "right", 0, "right, WGS84 longitude. (Required)")
	flag.IntVar(&zoom, "zoom", 15, "WGS84 zoomlevel from [12-21]")
	flag.StringVar(&out, "out", "", "Directory for output. (Optional default is Dir from file param )")
	flag.Parse()

	if file == "" || top == 0 || bottom == 0 || left == 0 || right == 0 {
		flag.PrintDefaults()
		os.Exit(1)

	}

	if out == "" {
		out = filepath.Dir(file) + "/"
	}

	var bounds = wgs.WGS84Bounds{Top: top, Right: right, Left: left, Bottom: bottom}

	src, err := imaging.Open(file)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	nbSlices, time := wgs.SliceIt(src, bounds, zoom, out)

	log.Printf("%s file generate %d tiles at this %d zoom level in %s", file, nbSlices, zoom, time)

}
