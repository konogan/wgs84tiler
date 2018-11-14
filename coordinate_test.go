package wgs84tiler

import (
	"testing"
)

func TestIsoConversion(t *testing.T) {

	zoom := 15
	x := 16850
	y := 11272

	tile := Tile{x, y}

	tileToLatLong := tile.ToLatLong(zoom)
	latLong2Tile := tileToLatLong.ToTile(zoom)

	if latLong2Tile.x != x || latLong2Tile.y != y {
		t.Errorf("Expected  %v = %v ", tile, latLong2Tile)
	}
}
