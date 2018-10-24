package wgs84tiler

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
type WGS84Bounds struct {
	top    float64
	right  float64
	left   float64
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

func getTargetImageSize(imageSource image.Image, wgs84Bounds WGS84Bounds, zoom int) dimension {
	// defer TimeTrack(time.Now(), "getTargetImageSize")
	imageSouceBounds := imageSource.Bounds()
	x1 := Tile2long(Long2tile(wgs84Bounds.left, zoom), zoom)
	x2 := Tile2long(Long2tile(wgs84Bounds.left, zoom)+1, zoom)
	lngperpx := (x2 - x1) / float64(TILESIZE)
	fileSizeInLng := wgs84Bounds.right - wgs84Bounds.left
	newWidth := math.Ceil(fileSizeInLng / lngperpx)
	newHeight := math.Ceil(newWidth * float64(imageSouceBounds.Max.Y) / float64(imageSouceBounds.Max.X))
	return dimension{width: int(newWidth), height: int(newHeight)}
}

func getTargetTilesBounds(wgs84Bounds WGS84Bounds, zoom int) (tilebounds, shift) {
	// defer TimeTrack(time.Now(), "getTargetTilesBounds")

	// bounds in tiles
	var tilebounds tilebounds
	tilebounds.top = Lat2tile(wgs84Bounds.top, zoom)
	tilebounds.bottom = Lat2tile(wgs84Bounds.bottom, zoom)
	tilebounds.left = Long2tile(wgs84Bounds.left, zoom)
	tilebounds.right = Long2tile(wgs84Bounds.right, zoom)

	// tiles in coord
	tileTopLat := Tile2lat(tilebounds.top, zoom)
	tileTopLatNext := Tile2lat(tilebounds.top+1, zoom)
	tileLeftLng := Tile2long(tilebounds.left, zoom)
	tileLeftLngNext := Tile2long(tilebounds.left+1, zoom)

	var tileshift shift
	tileshift.top = int(((wgs84Bounds.top - tileTopLat) / (tileTopLatNext - tileTopLat)) * float64(TILESIZE))
	tileshift.left = int(((wgs84Bounds.left - tileLeftLng) / (tileLeftLngNext - tileLeftLng)) * float64(TILESIZE))

	return tilebounds, tileshift
}

func sliceIt(imageSource image.Image, wgs84Bounds WGS84Bounds, zoom int, outputDir string) int {
	targetImageSize := getTargetImageSize(imageSource, wgs84Bounds, zoom)
	targetTilesBounds, targetTileShift := getTargetTilesBounds(wgs84Bounds, zoom)

	sX := 0
	sY := 0
	for tileX := targetTilesBounds.left; tileX <= targetTilesBounds.right; tileX++ {
		for tileY := targetTilesBounds.top; tileY <= targetTilesBounds.bottom; tileY++ {

			// for each slice initialize extract and shift
			var sliceExtract extract
			var sliceShift shift
			sliceExtract.width = TILESIZE
			sliceExtract.height = TILESIZE
			// logger.Print(" ")
			// logger.Print("-----", tileX, tileY, "-------")

			if tileX == targetTilesBounds.left && tileX == targetTilesBounds.right {
				// logger.Print("premiere et derniere tile de la ligne")
				sliceExtract.width = targetImageSize.width
				sliceShift.left = targetTileShift.left
			} else if tileX == targetTilesBounds.left {
				// logger.Print("premiere tile de la ligne")
				sliceExtract.width = TILESIZE - targetTileShift.left
				sliceShift.left = targetTileShift.left
			} else if tileX == targetTilesBounds.right {
				// logger.Print("derniere tile de la ligne")
				sliceExtract.width = targetImageSize.width - sX*TILESIZE + targetTileShift.left
				sliceExtract.left = sX*TILESIZE - targetTileShift.left
			} else {
				// logger.Print("tile intermediaire de la ligne")
				sliceExtract.left = sX*TILESIZE - targetTileShift.left - 1
			}

			if tileY == targetTilesBounds.top && tileY == targetTilesBounds.bottom {
				// logger.Print("premiere et derniere tile de la colonne")
				sliceExtract.height = targetImageSize.height
				sliceShift.top = targetTileShift.top
			} else if tileY == targetTilesBounds.top {
				// logger.Print("premiere tile de la colonne")
				sliceExtract.height = TILESIZE - targetTileShift.top
				sliceShift.top = targetTileShift.top
			} else if tileY == targetTilesBounds.bottom {
				// logger.Print("derniere tile de la colonne")
				sliceExtract.height = targetImageSize.height - sY*TILESIZE + targetTileShift.top
				sliceExtract.top = sY*TILESIZE - targetTileShift.top
			} else {
				// logger.Print("tile intermediaire de la colonne")
				sliceExtract.top = sY*TILESIZE - targetTileShift.top - 1
			}

			// logger.Print("targetTilesBounds", targetTilesBounds)
			// logger.Print("sliceExtract", sliceExtract)
			// logger.Print("sliceShift", sliceShift)

			//logger.Print(tileX, tileY, sliceExtract, sliceShift)

			makeTheSlice(imageSource, tileX, tileY, zoom, sliceExtract, sliceShift, outputDir)

			sY++

		}
		sX++
		sY = 0
	}
	nbtiles := (targetTilesBounds.right - targetTilesBounds.left) * (targetTilesBounds.bottom - targetTilesBounds.top)
	return nbtiles
}

func makeTheSlice(imageSource image.Image, tileX int, tileY int, zoom int, sliceExtract extract, sliceShift shift, outputDir string) {
	// defer TimeTrack(time.Now(), "   makeTheSlice ("+strconv.Itoa(tileX)+"-"+strconv.Itoa(tileY)+")")
	var path = outputDir + strconv.Itoa(zoom) + "/" + strconv.Itoa(tileX)
	var file = strconv.Itoa(tileY) + ".jpg"
	var fulldest = path + "/" + file
	os.MkdirAll(path, os.ModePerm)

	var rect image.Rectangle
	rect.Min = image.Pt(sliceExtract.left, sliceExtract.top)
	rect.Max = image.Pt(sliceExtract.left+sliceExtract.width, sliceExtract.top+sliceExtract.height)

	part := imaging.Crop(imageSource, rect)
	// todo :check if fulldest exist if yes get it
	vierge := imaging.New(TILESIZE, TILESIZE, color.NRGBA{128, 128, 128, 255})
	dst := imaging.Paste(vierge, part, image.Pt(sliceShift.left, sliceShift.top))
	var err = imaging.Save(dst, fulldest)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	//logger.Printf("tile %d/%d %+v %+v", x, y, sliceExtract, sliceShift)
}
