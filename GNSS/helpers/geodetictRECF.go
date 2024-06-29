package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type LocalCoordinates struct {
	initECEF       []float64
	ned2ecefMatrix [][]float64
	ecef2nedMatrix [][]float64
}

const (
	// SemiMajorAxis             = 6378137.0
	// SemiMinorAxis             = 6356752.3142
	EccentricitySquared       = 0.00669437999
	SecondEccentricitySquared = 6.73949674228 * 0.001
)

const (
	// semi-major axis
	a  float64 = 6378137
	e2 float64 = 6.69437999014 * 0.001
	// semi-minor axis
	b float64 = 6356752.3142

	// first eccentricity squared
	esq float64 = 6.69437999014 * 0.001

	// second eccentricity squared
	e1sq float64 = 6.73949674228 * 0.001

	SemiMajorAxis            = 6378137.0                                                       // Semi-major axis (a)
	Flattening               = 1 / 298.257223563                                               // Flattening (f)
	SemiMinorAxis            = SemiMajorAxis * (1 - Flattening)                                // Semi-minor axis (b)
	FirstEccentricitySquared = 1 - (SemiMinorAxis*SemiMinorAxis)/(SemiMajorAxis*SemiMajorAxis) // First eccentricity squared (e^2)
)

// ======================================

// ======================================

func GeodeticToECEF(geodetic [][]float64, radians bool) [][]float64 {
	inputShape := len(geodetic)
	ratio := 1.0
	if !radians {
		ratio = math.Pi / 180.0
	}
	earthCenterEarthFixed := make([][]float64, inputShape)
	for i, coordinate := range geodetic {
		lat := ratio * coordinate[0]
		lon := ratio * coordinate[1]
		alt := coordinate[2]

		// Prime vertical radius of curvature
		N := SemiMajorAxis / math.Sqrt(1-FirstEccentricitySquared*math.Sin(lat)*math.Sin(lat))
		x := (N + alt) * math.Cos(lat) * math.Cos(lon)
		y := (N + alt) * math.Cos(lat) * math.Sin(lon)
		z := (N*(1-FirstEccentricitySquared) + alt) * math.Sin(lat)

		// Calculate coordinates
		earthCenterEarthFixed[i] = []float64{x, y, z}
	}

	return earthCenterEarthFixed
}

// ======================================

// ======================================

func FlattenMatrix(slice [][]float64) []float64 {
	result := make([]float64, 0, len(slice)*len(slice[0]))
	for _, row := range slice {
		result = append(result, row...)
	}
	return result
}

// ======================================

// ======================================

func Reshape(slice []float64, rows, cols int) [][]float64 {
	result := make([][]float64, rows)
	for i := range result {
		result[i] = slice[i*cols : (i+1)*cols]
	}
	return result
}

// ======================================

// ======================================

// https://gis.stackexchange.com/questions/28446/computational-most-efficient-way-to-convert-cartesian-to-geodetic-coordinates
// SEEMS LIKE OLSO METHOD IS THE BEST
func ECEFToGeodetic(ecef [][]float64, radians bool) [][]float64 {
	output := make([][]float64, len(ecef))

	ratio := 1.0
	if !radians {
		ratio = 180.0 / math.Pi
	}

	for i := range ecef {
		x, y, z := ecef[i][0], ecef[i][1], ecef[i][2]

		a1 := 4.2697672707157535e+4
		a2 := 1.8230912546075455e+9
		a3 := 1.4291722289812413e+2
		a4 := 4.5577281365188637e+9
		a5 := 4.2840589930055659e+4
		a6 := 9.9330562000986220e-1

		zp := math.Abs(z)
		w2 := x*x + y*y
		w := math.Sqrt(w2)
		z2 := z * z
		r2 := w2 + z2
		r := math.Sqrt(r2)
		if r < 100000 {
			output[i] = []float64{0.0, 0.0, -1e7}
			continue
		}
		lon := math.Atan2(y, x)
		s2 := z2 / r2
		c2 := w2 / r2
		u := a2 / r
		v := a3 - a4/r
		var s, c, lat float64
		if c2 > 0.3 {
			s = (zp / r) * (1.0 + c2*(a1+u+s2*v)/r)
			lat = math.Asin(s)
			ss := s * s
			c = math.Sqrt(1.0 - ss)
		} else {
			c = (w / r) * (1.0 - s2*(a5-u-c2*v)/r)
			lat = math.Acos(c)
			ss := 1.0 - c*c
			s = math.Sqrt(ss)
		}
		g := 1.0 - e2*s*s
		rg := a / math.Sqrt(g)
		rf := a6 * rg
		u = w - rg*c
		v = zp - rf*s
		f := c*u + s*v
		m := c*v - s*u
		p := m / (rf/g + f)
		lat += p
		h := f + m*p/2.0
		if z < 0.0 {
			lat = -lat
		}

		output[i] = []float64{ratio * lat, ratio * lon, h}
	}

	return output
}

// ======================================

// ======================================

func NewLocalCoordinates(Geodetic []float64) *LocalCoordinates {

	if len(Geodetic) != 3 {
		panic("Geodetic coordinates must be of length 3")
	}

	fmt.Printf("Input Geodetic: %v\n", Geodetic)

	initECEF := GeodeticToECEF([][]float64{Geodetic}, false)[0]

	fmt.Printf("Initial ECEF: %v\n", initECEF)

	lat := (math.Pi / 180) * Geodetic[0]
	lon := (math.Pi / 180) * Geodetic[1]

	ned2ecefMatrix := [][]float64{
		{-math.Sin(lat) * math.Cos(lon), -math.Sin(lon), -math.Cos(lat) * math.Cos(lon)},
		{-math.Sin(lat) * math.Sin(lon), math.Cos(lon), -math.Cos(lat) * math.Sin(lon)},
		{math.Cos(lat), 0, -math.Sin(lat)},
	}

	ecef2nedMatrix := Transpose(ned2ecefMatrix)
	fmt.Println("====================================")
	geodeticFromECEF := ECEFToGeodetic([][]float64{initECEF}, false)
	fmt.Printf("Geodetic from ECEF: %v\n", geodeticFromECEF)
	fmt.Println("+++++++++++++++++++++++++++++++")
	return &LocalCoordinates{
		initECEF:       initECEF,
		ned2ecefMatrix: ned2ecefMatrix,
		ecef2nedMatrix: ecef2nedMatrix,
	}
}

// ======================================

// ======================================

func NewLocalCoordinatesFromECEF(initECEF []float64) *LocalCoordinates {
	Geodetic := ECEFToGeodetic([][]float64{initECEF}, true)[0]
	return NewLocalCoordinates(Geodetic)
}

// ======================================

// ======================================

func (lc *LocalCoordinates) ECEFToNED(ecef []float64) []float64 {
	deltaECEF := Subtract(ecef, lc.initECEF)
	return Dot(lc.ecef2nedMatrix, deltaECEF)
}

// ======================================

// ======================================

func (lc *LocalCoordinates) NEDToECEF(ned []float64) []float64 {
	ecef := Dot(lc.ned2ecefMatrix, ned)
	return Add(ecef, lc.initECEF)
}

// ======================================

// ======================================

func (lc *LocalCoordinates) GeodeticToNED(geodetic []float64) []float64 {
	ecef := GeodeticToECEF([][]float64{geodetic}, true)[0]
	return lc.ECEFToNED(ecef)
}

// ======================================

// ======================================

func (lc *LocalCoordinates) NEDToGeodetic(ned []float64) []float64 {
	ecef := lc.NEDToECEF(ned)
	return ECEFToGeodetic([][]float64{ecef}, true)[0]
}

// ======================================

// ======================================

func Transpose(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])
	transposed := make([][]float64, cols)
	for i := range transposed {
		transposed[i] = make([]float64, rows)
	}
	for i := range matrix {
		for j := range matrix[i] {
			transposed[j][i] = matrix[i][j]
		}
	}
	return transposed
}

// ======================================

// ======================================

func Dot(matrix [][]float64, vector []float64) []float64 {
	result := make([]float64, len(matrix))
	for i := range matrix {
		for j := range vector {
			result[i] += matrix[i][j] * vector[j]
		}
	}
	return result
}

// ======================================

// ======================================

func Add(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

// ======================================

// ======================================

func Subtract(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

// ==================================== BENCHED FUNCTIONS ====================================

func generateRandomECEF(n int) [][]float64 {
	ecef := make([][]float64, n)
	for i := 0; i < n; i++ {
		ecef[i] = []float64{rand.Float64(), rand.Float64(), rand.Float64()}
	}
	return ecef
}

func BenchECEFToGeodetic(n int) {
	starttime := time.Now()
	for i := 0; i < n; i++ {
		ecef := generateRandomECEF(n)
		fmt.Println(ECEFToGeodetic(ecef, false))

	}
	fmt.Println("ECEFToGeodetic took: ", time.Since(starttime))
}
