package main

import "math"

// Latlong define a latitude longitude point
type Latlong struct {
	lat  float64
	long float64
}

//Tile define a WSG84 Tile coordinate
type Tile struct {
	x int
	y int
}

//ToTile convert a Latlong coordinate to a  WSG84 Tile coordinate for a zoom level
func (ll *Latlong) ToTile(zoom int) Tile {
	return Tile{lat2tile(ll.lat, zoom), long2tile(ll.long, zoom)}
}

//ToLatLong convert a WSG84 Tile coordinate to a Latlong coordinate for a zoom level
func (t *Tile) ToLatLong(zoom int) Latlong {
	return Latlong{tile2lat(t.x, zoom), tile2long(t.y, zoom)}
}

func long2tile(lon float64, zoom int) int {
	return int(math.Floor(((lon + 180) / 360) * math.Pow(2, float64(zoom))))
}

func lat2tile(lat float64, zoom int) int {
	return int(((1 - math.Log(math.Tan((lat*math.Pi)/180)+1/math.Cos((lat*math.Pi)/180))/math.Pi) / 2) * math.Pow(2, float64(zoom)))
}

func tile2long(x, zoom int) float64 {
	return float64((float64(x)/math.Pow(2, float64(zoom)))*360 - 180)
}

func tile2lat(y, zoom int) float64 {
	var n = math.Pi - (2*math.Pi*float64(y))/math.Pow(2, float64(zoom))
	return float64(180/math.Pi) * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
}
