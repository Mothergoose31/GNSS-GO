package gnss

/*
https://gssc.esa.int/navipedia/index.php/Tropospheric_Delay
http://ftp.aiub.unibe.ch/BERN42/DOCU/DOCU42_5.pdf
https://www.researchgate.net/publication/228746230_Tropospheric_Delay_Estimation_for_Pseudolite_Positioning
https://www.swsc-journal.org/articles/swsc/full_html/2018/01/swsc170068/swsc170068.html */

import (
	"math"

	"github.com/mothergoose31/GNNS-GO/GNSS/helpers"
)

func TropoDelay(lat, lon, height float64, elevation float64) float64 {

	ecef := helpers.GeodeticToECEF([][]float64{{lat, lon, height}}, false)[0]

	geodetic := helpers.ECEFToGeodetic([][]float64{ecef}, false)[0]
	lat = geodetic[0] * math.Pi / 180.0 // Convert to radians
	height = geodetic[2]

	if elevation < 0.0175 {
		elevation = 0.0175
	}

	// TODO make these dynamic

	T0 := 288.15

	beta := 6.5e-3

	P0 := 1013.25

	// Water vapor partial
	e0 := 11.691

	g := EARTH_GM / (EARTH_RADIUS * EARTH_RADIUS)

	// calc temperature at height
	T := T0 - beta*height

	// calc barometric pressure
	P := P0 * math.Pow(1.0-beta*height/T0, g/(EARTH_ROTATION_RATE*beta))

	//  water vapor partial pressure at height
	e := e0 * math.Pow(1.0-beta*height/T0, 4*g/(EARTH_ROTATION_RATE*beta))

	sinEl := math.Sin(elevation)
	tanEl := math.Tan(elevation)

	// zenith hydrostatic delay
	zhd := 0.0022768 * P / (1.0 - 0.00266*math.Cos(2.0*lat) - 0.00028*height/1000.0)

	// Calculate zenith wet delay
	zwd := 0.002277 * (1255.0/T + 0.05) * e

	// Calculate mapping functions
	mh := 1.0 / (sinEl + 0.00143/(tanEl+0.0445))
	mw := 1.0 / (sinEl + 0.00035/(tanEl+0.017))

	// Return total tropospheric delay
	return zhd*mh + zwd*mw
}

// elevationRate: rate of change of elevation angle in radians/second
// Returns delay rate in meters/second
func TropoDelayRate(lat, lon, height float64, elevation, elevationRate float64) float64 {

	delay1 := TropoDelay(lat, lon, height, elevation)

	dt := 0.001

	delay2 := TropoDelay(lat, lon, height, elevation+elevationRate*dt)

	return (delay2 - delay1) / dt
}

// elevation:  angle in radians
// pseudorange: measurement in meters
// carrierPhase:  measurement in cycles
// freq: carrier frequency in Hz
// Returns corrected pseudorange (meters) and carrier phase (cycles)
func TropoCorrection(lat, lon, height float64, elevation float64,
	pseudorange float64, carrierPhase float64, freq float64) (float64, float64) {

	delay := TropoDelay(lat, lon, height, elevation)

	correctedPR := pseudorange - delay
	delayCycles := delay * freq / SPEED_OF_LIGHT

	correctedCP := carrierPhase - delayCycles

	return correctedPR, correctedCP
}
