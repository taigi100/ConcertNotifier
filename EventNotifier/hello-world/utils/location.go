package utils

import "math"

//Location describes a lat/long coordinate
type Location struct {
	Lat  float64
	Long float64
}

func toRadians(val float64) float64 {
	return val * math.Pi / 180
}

//DistanceTo returns the distance between coordonates in meters using haversine formula
func (loc Location) DistanceTo(other Location) float64 {
	R := 6371e3
	omega1 := toRadians(loc.Lat)
	omega2 := toRadians(other.Lat)
	deltaomega := toRadians(other.Lat - loc.Lat)
	deltaalpha := toRadians(other.Long - loc.Long)

	a := math.Sin(deltaomega/2)*math.Sin(deltaomega/2) +
		math.Cos(omega1)*math.Cos(omega2)*math.Sin(deltaalpha/2)*math.Sin(deltaalpha/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return math.Round(d)
}
