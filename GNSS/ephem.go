package gnss

import (
	"errors"
	"math"
)

type EphemerisType int

const (
	NAV EphemerisType = iota
	FINAL_ORBIT
	RAPID_ORBIT
	ULTRA_RAPID_ORBIT
	QCOM_POLY
)

type BaseEphemeris struct {
	PRN         string
	Epoch       GPSTime
	EphType     EphemerisType
	Healthy     bool
	MaxTimeDiff float64
	FileEpoch   *GPSTime
	FileName    string
	FileSource  string
}

type GPSEphemeris struct {
	BaseEphemeris
	Ephemeris
	TOE   GPSTime
	TOC   GPSTime
	SqrtA float64
}

type GLONASSEphemeris struct {
	BaseEphemeris
	GlonassEphemeris
	Channel int
}

//=============================================================================

// =========================================================================

func (e *BaseEphemeris) Valid(time GPSTime) bool {
	return math.Abs(time.Sub(e.Epoch)) <= e.MaxTimeDiff
}

func NewGPSEphemeris(eph Ephemeris, fileName string) *GPSEphemeris {
	toe := NewGPSTime(int(eph.ToeWeek()), float64(eph.Toe()))
	toc := NewGPSTime(int(eph.TimeOfClockWeek()), float64(eph.TimeOfClock()))
	epoch := toc

	return &GPSEphemeris{
		BaseEphemeris: BaseEphemeris{
			PRN:         "G" + string(eph.SatelliteId()),
			Epoch:       epoch,
			EphType:     NAV,
			Healthy:     eph.SatelliteHealth() == 0,
			MaxTimeDiff: 2 * SECS_IN_HR,
			FileName:    fileName,
		},
		Ephemeris: eph,
		TOE:       toe,
		TOC:       toc,
		SqrtA:     math.Sqrt(eph.SemiMajorAxis()),
	}
}

// # http://gauss.gge.unb.ca/GLONASS.ICD.pdf
func (e *GPSEphemeris) GetSatInfo(time GPSTime) ([]float64, []float64, float64, float64, error) {
	if !e.Healthy {
		return nil, nil, 0, 0, errors.New("unhealthy ephemeris")
	}

	tdiff := time.Sub(e.TOC)
	clockErr := e.Af0() + tdiff*(e.Af1()+tdiff*e.Af2())
	clockRateErr := e.Af1() + 2*tdiff*e.Af2()

	tdiff = time.Sub(e.TOE)

	a := e.SqrtA * e.SqrtA
	maDot := math.Sqrt(EARTH_GM/(a*a*a)) + e.DeltaN()
	ma := e.M0() + maDot*tdiff

	ea := ma
	eaOld := 2222.0
	for math.Abs(ea-eaOld) > 1.0e-14 {
		eaOld = ea
		tempd1 := 1.0 - e.Eccentricity()*math.Cos(eaOld)
		ea = ea + (ma-eaOld+e.Eccentricity()*math.Sin(eaOld))/tempd1
	}
	eaDot := maDot / (1.0 - e.Eccentricity()*math.Cos(ea))

	einstein := -4.442807633e-10 * e.Eccentricity() * e.SqrtA * math.Sin(ea)

	tempd2 := math.Sqrt(1.0 - e.Eccentricity()*e.Eccentricity())
	al := math.Atan2(tempd2*math.Sin(ea), math.Cos(ea)-e.Eccentricity()) + e.PerigeeArgument()
	alDot := tempd2 * eaDot / (1.0 - e.Eccentricity()*math.Cos(ea))

	cal := al + e.Cus()*math.Sin(2.0*al) + e.Cuc()*math.Cos(2.0*al)
	calDot := alDot * (1.0 + 2.0*(e.Cus()*math.Cos(2.0*al)-e.Cuc()*math.Sin(2.0*al)))

	r := a*(1.0-e.Eccentricity()*math.Cos(ea)) + e.Crc()*math.Cos(2.0*al) + e.Crs()*math.Sin(2.0*al)
	rDot := a*e.Eccentricity()*eaDot*math.Sin(ea) +
		2.0*alDot*(e.Crs()*math.Cos(2.0*al)-e.Crc()*math.Sin(2.0*al))

	inc := e.Inclination() + e.RateOfInclination()*tdiff +
		e.Cic()*math.Cos(2.0*al) + e.Cis()*math.Sin(2.0*al)
	incDot := e.RateOfInclination() +
		2.0*alDot*(e.Cis()*math.Cos(2.0*al)-e.Cic()*math.Sin(2.0*al))

	x := r * math.Cos(cal)
	y := r * math.Sin(cal)
	xDot := rDot*math.Cos(cal) - y*calDot
	yDot := rDot*math.Sin(cal) + x*calDot

	omDot := e.RateOfRightAscension() - EARTH_ROTATION_RATE
	om := e.Omega0() + tdiff*omDot - EARTH_ROTATION_RATE*float64(e.TOE.Tow)

	pos := make([]float64, 3)
	pos[0] = x*math.Cos(om) - y*math.Cos(inc)*math.Sin(om)
	pos[1] = x*math.Sin(om) + y*math.Cos(inc)*math.Cos(om)
	pos[2] = y * math.Sin(inc)

	tempd3 := yDot*math.Cos(inc) - y*math.Sin(inc)*incDot

	// Compute the satellite's velocity in Earth-Centered Earth-Fixed coordinates
	vel := make([]float64, 3)
	vel[0] = -omDot*pos[1] + xDot*math.Cos(om) - tempd3*math.Sin(om)
	vel[1] = omDot*pos[0] + xDot*math.Sin(om) + tempd3*math.Cos(om)
	vel[2] = y*math.Cos(inc)*incDot + yDot*math.Sin(inc)

	clockErr += einstein

	return pos, vel, clockErr, clockRateErr, nil
}
