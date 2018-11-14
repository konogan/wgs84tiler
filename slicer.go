package wgs84tiler

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
)

// WGS84Bounds boundaries of the image in WSG84 world
//
// Example :
// myBounds = WGS84Bounds{top: 48.8687073004617, right: 2.15657022586739, left: 2.14840505163567,bottom: 48.8651234503015}
type WGS84Bounds struct {
	Top    float64 // top holds the top coordinate in WGS84 latitude
	Right  float64 // holds the right coordinate in WGS84 longitude
	Left   float64 // left  holds the coordinate in WGS84 longitude
	Bottom float64 // bottom  holds the coordinate in WGS84 latitude
}

type dimension struct {
	width  int
	height int
}

// tilebounds wgs84Bounds in tile coordinate
type tilebounds struct {
	top    int
	bottom int
	left   int
	right  int
}

// shift offset image in pixels from top and left tiles
type shift struct {
	left int
	top  int
}

// portion of the original image to extract
type extract struct {
	width  int
	height int
	left   int
	top    int
}

// TILESIZE the slippy tilesize value in pixels
const TILESIZE = 256

var virgin = imaging.New(TILESIZE, TILESIZE, color.NRGBA{128, 128, 128, 0})

func getTargetImageSize(imageSource image.Image, wgs84Bounds WGS84Bounds, zoom int) dimension {
	// defer timeTrack(time.Now(), "getTargetImageSize")
	imageSouceBounds := imageSource.Bounds()
	x1 := tile2long(long2tile(wgs84Bounds.Left, zoom), zoom)
	x2 := tile2long(long2tile(wgs84Bounds.Left, zoom)+1, zoom)
	lngperpx := (x2 - x1) / float64(TILESIZE)
	fileSizeInLng := wgs84Bounds.Right - wgs84Bounds.Left
	newWidth := math.Ceil(fileSizeInLng / lngperpx)
	newHeight := math.Ceil(newWidth * float64(imageSouceBounds.Max.Y) / float64(imageSouceBounds.Max.X))
	return dimension{width: int(newWidth), height: int(newHeight)}
}

func getTargetTilesBounds(wgs84Bounds WGS84Bounds, zoom int) (tilebounds, shift) {
	// defer timeTrack(time.Now(), "getTargetTilesBounds")

	// bounds in tiles
	var tilebounds tilebounds
	tilebounds.top = lat2tile(wgs84Bounds.Top, zoom)
	tilebounds.bottom = lat2tile(wgs84Bounds.Bottom, zoom)
	tilebounds.left = long2tile(wgs84Bounds.Left, zoom)
	tilebounds.right = long2tile(wgs84Bounds.Right, zoom)

	// tiles in coord
	tileTopLat := tile2lat(tilebounds.top, zoom)
	tileTopLatNext := tile2lat(tilebounds.top+1, zoom)
	tileLeftLng := tile2long(tilebounds.left, zoom)
	tileLeftLngNext := tile2long(tilebounds.left+1, zoom)

	var tileshift shift
	tileshift.top = int(((wgs84Bounds.Top - tileTopLat) / (tileTopLatNext - tileTopLat)) * float64(TILESIZE))
	tileshift.left = int(((wgs84Bounds.Left - tileLeftLng) / (tileLeftLngNext - tileLeftLng)) * float64(TILESIZE))

	return tilebounds, tileshift
}

// SliceIt : analyse image and command the slicing
func SliceIt(imageSource image.Image, wgs84Bounds WGS84Bounds, zoom int, outputDir string) (nbtiles, statNew, statMerge int, tps time.Duration) {
	start := time.Now()
	targetImageSize := getTargetImageSize(imageSource, wgs84Bounds, zoom)
	targetTilesBounds, targetTileShift := getTargetTilesBounds(wgs84Bounds, zoom)

	statNew = 0
	statMerge = 0

	resizedImage := imaging.Resize(imageSource, targetImageSize.width, targetImageSize.height, imaging.Lanczos)
	sX := 0
	sY := 0
	for tileX := targetTilesBounds.left; tileX <= targetTilesBounds.right; tileX++ {
		for tileY := targetTilesBounds.top; tileY <= targetTilesBounds.bottom; tileY++ {

			// for each slice initialize extract and shift
			var sliceExtract image.Rectangle
			var sliceShift image.Point
			var width = TILESIZE
			var height = TILESIZE
			var top = 0
			var left = 0

			if tileX == targetTilesBounds.left && tileX == targetTilesBounds.right {
				//first and last tile of the row
				width = targetImageSize.width
				sliceShift.X = targetTileShift.left
			} else if tileX == targetTilesBounds.left {
				//first tile of the row
				width = TILESIZE - targetTileShift.left
				sliceShift.X = targetTileShift.left
			} else if tileX == targetTilesBounds.right {
				//last tile of the row
				width = targetImageSize.width - sX*TILESIZE + targetTileShift.left
				left = sX*TILESIZE - targetTileShift.left
			} else {
				//tile intermediaire of the row
				left = sX*TILESIZE - targetTileShift.left - 1
			}

			if tileY == targetTilesBounds.top && tileY == targetTilesBounds.bottom {
				//first and last tile of the column
				height = targetImageSize.height
				sliceShift.Y = targetTileShift.top
			} else if tileY == targetTilesBounds.top {
				//first tile of the column
				height = TILESIZE - targetTileShift.top
				sliceShift.Y = targetTileShift.top
			} else if tileY == targetTilesBounds.bottom {
				//last tile of the column
				height = targetImageSize.height - sY*TILESIZE + targetTileShift.top
				top = sY*TILESIZE - targetTileShift.top
			} else {
				//tile intermediaire of the column
				top = sY*TILESIZE - targetTileShift.top - 1
			}

			sliceExtract.Min = image.Pt(left, top)
			sliceExtract.Max = image.Pt(left+width, top+height)
			isNew := makeTheSlice(resizedImage, Tile{x: tileX, y: tileY}, zoom, sliceExtract, sliceShift, outputDir)
			if isNew {
				statNew++
			} else {
				statMerge++
			}

			sY++
		}
		sX++
		sY = 0
	}
	nbtiles = statMerge + statNew
	elapsed := time.Since(start)
	return nbtiles, statNew, statMerge, elapsed
}

func makeTheSlice(imageSource image.Image, tile Tile, zoom int, sliceExtract image.Rectangle, sliceShift image.Point, outputDir string) bool {
	var path = outputDir + strconv.Itoa(zoom) + "/" + strconv.Itoa(tile.x)
	var file = strconv.Itoa(tile.y) + ".png"
	var fulldest = path + "/" + file
	os.MkdirAll(path, os.ModePerm)
	part := imaging.Crop(imageSource, sliceExtract)
	isNew := false
	originalContent, err := imaging.Open(fulldest)
	if err != nil {
		isNew = true
		originalContent = virgin
	}
	dst := imaging.Paste(originalContent, part, image.Pt(sliceShift.X, sliceShift.Y))
	err = imaging.Save(dst, fulldest)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	return isNew
}
