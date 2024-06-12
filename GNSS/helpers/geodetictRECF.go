package helpers

import "math"

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

		// todo LOOKUP THE FORMULA
		//  prime vertical radius of curvature

		// calc cords
		earthCenterEarthFixed[i] = []float64{}
	}

	return earthCenterEarthFixed
}
