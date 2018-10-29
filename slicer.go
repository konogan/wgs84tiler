package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

// WGS84Bounds boundaries of the image in WSG84 world
//
// Example :
// myBounds = WGS84Bounds{top: 48.8687073004617, right: 2.15657022586739, left: 2.14840505163567,bottom: 48.8651234503015}
type WGS84Bounds struct {
	// top holds the top coordinate in WGS84 latitude
	top float64
	// holds the right coordinate in WGS84 longitude
	right float64
	// left  holds the coordinate in WGS84 longitude
	left float64
	// bottom  holds the coordinate in WGS84 latitude
	bottom float64
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
	x1 := tile2long(long2tile(wgs84Bounds.left, zoom), zoom)
	x2 := tile2long(long2tile(wgs84Bounds.left, zoom)+1, zoom)
	lngperpx := (x2 - x1) / float64(TILESIZE)
	fileSizeInLng := wgs84Bounds.right - wgs84Bounds.left
	newWidth := math.Ceil(fileSizeInLng / lngperpx)
	newHeight := math.Ceil(newWidth * float64(imageSouceBounds.Max.Y) / float64(imageSouceBounds.Max.X))
	return dimension{width: int(newWidth), height: int(newHeight)}
}

func getTargetTilesBounds(wgs84Bounds WGS84Bounds, zoom int) (tilebounds, shift) {
	// defer timeTrack(time.Now(), "getTargetTilesBounds")

	// bounds in tiles
	var tilebounds tilebounds
	tilebounds.top = lat2tile(wgs84Bounds.top, zoom)
	tilebounds.bottom = lat2tile(wgs84Bounds.bottom, zoom)
	tilebounds.left = long2tile(wgs84Bounds.left, zoom)
	tilebounds.right = long2tile(wgs84Bounds.right, zoom)

	// tiles in coord
	tileTopLat := tile2lat(tilebounds.top, zoom)
	tileTopLatNext := tile2lat(tilebounds.top+1, zoom)
	tileLeftLng := tile2long(tilebounds.left, zoom)
	tileLeftLngNext := tile2long(tilebounds.left+1, zoom)

	var tileshift shift
	tileshift.top = int(((wgs84Bounds.top - tileTopLat) / (tileTopLatNext - tileTopLat)) * float64(TILESIZE))
	tileshift.left = int(((wgs84Bounds.left - tileLeftLng) / (tileLeftLngNext - tileLeftLng)) * float64(TILESIZE))

	return tilebounds, tileshift
}

func sliceIt(imageSource image.Image, wgs84Bounds WGS84Bounds, zoom int, outputDir string) int {
	targetImageSize := getTargetImageSize(imageSource, wgs84Bounds, zoom)
	targetTilesBounds, targetTileShift := getTargetTilesBounds(wgs84Bounds, zoom)

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
				//premiere et derniere tile de la ligne
				width = targetImageSize.width
				sliceShift.X = targetTileShift.left
			} else if tileX == targetTilesBounds.left {
				//premiere tile de la ligne
				width = TILESIZE - targetTileShift.left
				sliceShift.X = targetTileShift.left
			} else if tileX == targetTilesBounds.right {
				//derniere tile de la ligne
				width = targetImageSize.width - sX*TILESIZE + targetTileShift.left
				left = sX*TILESIZE - targetTileShift.left
			} else {
				//tile intermediaire de la ligne
				left = sX*TILESIZE - targetTileShift.left - 1
			}

			if tileY == targetTilesBounds.top && tileY == targetTilesBounds.bottom {
				//premiere et derniere tile de la colonne
				height = targetImageSize.height
				sliceShift.Y = targetTileShift.top
			} else if tileY == targetTilesBounds.top {
				//premiere tile de la colonne
				height = TILESIZE - targetTileShift.top
				sliceShift.Y = targetTileShift.top
			} else if tileY == targetTilesBounds.bottom {
				//derniere tile de la colonne
				height = targetImageSize.height - sY*TILESIZE + targetTileShift.top
				top = sY*TILESIZE - targetTileShift.top
			} else {
				//tile intermediaire de la colonne
				top = sY*TILESIZE - targetTileShift.top - 1
			}

			sliceExtract.Min = image.Pt(left, top)
			sliceExtract.Max = image.Pt(left+width, top+height)
			makeTheSlice(resizedImage, Tile{x: tileX, y: tileY}, zoom, sliceExtract, sliceShift, outputDir)
			sY++
		}
		sX++
		sY = 0
	}
	nbtiles := (targetTilesBounds.right - targetTilesBounds.left) * (targetTilesBounds.bottom - targetTilesBounds.top)
	return nbtiles
}

func makeTheSlice(imageSource image.Image, tile Tile, zoom int, sliceExtract image.Rectangle, sliceShift image.Point, outputDir string) {
	var path = outputDir + strconv.Itoa(zoom) + "/" + strconv.Itoa(tile.x)
	var file = strconv.Itoa(tile.y) + ".png"
	var fulldest = path + "/" + file
	os.MkdirAll(path, os.ModePerm)
	part := imaging.Crop(imageSource, sliceExtract)
	originalContent, err := imaging.Open(fulldest)
	if err != nil {
		originalContent = virgin
	}
	dst := imaging.Paste(originalContent, part, image.Pt(sliceShift.X, sliceShift.Y))
	err = imaging.Save(dst, fulldest)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
