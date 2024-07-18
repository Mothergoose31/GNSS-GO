package gnss

type EphemerisType int

const (
	NAV EphemerisType = iota
	FINAL_ORBIT
	RAPID_ORBIT
	ULTRA_RAPID_ORBIT
	QCOM_POLY
)

// type Ephemeris interface {
// 	Valid(time GPSTime) bool
// }
