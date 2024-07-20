package gnss

import (
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
