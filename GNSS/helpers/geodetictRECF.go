package helpers

import (
	"math"
)

const (
	SemiMajorAxis             = 6378137.0
	SemiMinorAxis             = 6356752.3142
	EccentricitySquared       = 6.69437999014 * 0.001
	SecondEccentricitySquared = 6.73949674228 * 0.001
)

func GeodeticToECEF(geodetic [][]float64, radians bool) [][]float64 {
	inputShape := len(geodetic)
	ratio := 1.0
	if radians {
		ratio = math.Pi / 180.0
	}
	earthCenterEarthFixed := make([][]float64, inputShape)
	for i, coordinate := range geodetic {
		lat := ratio * coordinate[0]
		lon := ratio * coordinate[1]
		alt := coordinate[2]

		//  prime vertical radius of curvature
		xi := math.Sqrt(1 - EccentricitySquared*math.Sin(lat)*math.Sin(lat))
		x := (SemiMajorAxis/xi + alt) * math.Cos(lat) * math.Cos(lon)
		y := (SemiMajorAxis/xi + alt) * math.Cos(lat) * math.Sin(lon)
		z := (SemiMajorAxis/xi*(1-EccentricitySquared) + alt) * math.Sin(lat)

		// calc cords
		earthCenterEarthFixed[i] = []float64{x, y, z}
	}

	return earthCenterEarthFixed
}
