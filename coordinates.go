package wgs84tiler

import "math"

//Long2tile convert a longitude in WSG84 tile coordinate
func Long2tile(lon float64, zoom int) int {
	return int(math.Floor(((lon + 180) / 360) * math.Pow(2, float64(zoom))))
}

//Lat2tile convert a latitude in WSG84 tile coordinate
func Lat2tile(lat float64, zoom int) int {
	return int(((1 - math.Log(math.Tan((lat*math.Pi)/180)+1/math.Cos((lat*math.Pi)/180))/math.Pi) / 2) * math.Pow(2, float64(zoom)))
}

//Tile2long convert a WSG84 tile coordinate in longitude
func Tile2long(x, zoom int) float64 {
	return float64((float64(x)/math.Pow(2, float64(zoom)))*360 - 180)
}

//Tile2lat convert a WSG84 tile coordinate in latitude
func Tile2lat(y, zoom int) float64 {
	var n = math.Pi - (2*math.Pi*float64(y))/math.Pow(2, float64(zoom))
	return float64(180/math.Pi) * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
}
