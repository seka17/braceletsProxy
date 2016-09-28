package structures

import "fmt"

type Point struct {
	Point    [2]float64 //[lon, lat]
	Accuracy float64
}

func (p Point) Lat() float64 {
	return p.Point[1]
}

func (p Point) Lon() float64 {
	return p.Point[0]
}

// json returns string in json format for bracelet
func (p Point) json() string {
	return fmt.Sprintf(`{"lat"=%v,"lon"=%v}`, p.Lat(), p.Lon())
}
