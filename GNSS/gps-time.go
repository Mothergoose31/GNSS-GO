package gnss

import (
	"errors"
	"time"
)

// towToDateTime converts a GPS Week and Time Of Week to a Go time.Time object.
// Does *not* convert from GPST to UTC. Fractional seconds are supported.

const SecondsInWeek = 604800

type GPSTime struct {
	Week int
	Tow  float64
}

func NewGPSTime(week int, tow float64) GPSTime {
	return GPSTime{week, tow}
}

func GPSTimeFromDateTime(t time.Time) GPSTime {
	wkRef := time.Date(2014, 2, 16, 0, 0, 0, 0, time.UTC)
	refWk := 1780
	wk := int(t.Sub(wkRef).Hours()/24/7) + refWk
	tow := float64(t.Sub(wkRef).Seconds()) - float64((wk-refWk)*SecondsInWeek)
	return GPSTime{Week: wk, Tow: tow}

}

func TowToDateTime(tow float64, week int) time.Time {

	gpsEpoch := time.Date(1980, 1, 6, 0, 0, 0, 0, time.UTC)

	towDuration := time.Duration(tow * float64(time.Second))
	weekDuration := time.Duration(week) * 7 * 24 * time.Hour

	return gpsEpoch.Add(towDuration + weekDuration)
}

func GetLeapSeconds(t time.Time) (int, error) {
	// Define the date boundaries
	date2006 := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
	date2009 := time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
	date2012 := time.Date(2012, 7, 1, 0, 0, 0, 0, time.UTC)
	date2015 := time.Date(2015, 7, 1, 0, 0, 0, 0, time.UTC)
	date2017 := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)

	switch {
	case t.Before(date2006):
		return 0, errors.New("don't know how many leap seconds to use before 2006")
	case t.Before(date2009):
		return 14, nil
	case t.Before(date2012):
		return 15, nil
	case t.Before(date2015):
		return 16, nil
	case t.Before(date2017):
		return 17, nil
	default:
		return 18, nil
	}
}

func GpstToUtc(tGpst time.Time) (time.Time, error) {
	leapSeconds, err := GetLeapSeconds(tGpst)
	if err != nil {
		return time.Time{}, err
	}
	tUtc := tGpst.Add(-time.Duration(leapSeconds) * time.Second)
	if UtcToGpst(tUtc).Sub(tGpst) != 0 {
		return tUtc.Add(time.Second), nil
	}
	return tUtc, nil
}

// Convert UTC to GPST
func UtcToGpst(tUtc time.Time) time.Time {
	leapSeconds, err := GetLeapSeconds(tUtc)
	if err != nil {
		return time.Time{}
	}
	return tUtc.Add(time.Duration(leapSeconds) * time.Second)
}
