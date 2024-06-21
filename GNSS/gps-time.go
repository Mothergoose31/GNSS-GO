package gnss

import (
	"errors"
	"time"
)

// towToDateTime converts a GPS Week and Time Of Week to a Go time.Time object.
// Does *not* convert from GPST to UTC. Fractional seconds are supported.
func towToDateTime(tow float64, week int) time.Time {

	gpsEpoch := time.Date(1980, 1, 6, 0, 0, 0, 0, time.UTC)

	towDuration := time.Duration(tow * float64(time.Second))
	weekDuration := time.Duration(week) * 7 * 24 * time.Hour

	return gpsEpoch.Add(towDuration + weekDuration)
}

func getLeapSeconds(t time.Time) (int, error) {
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

func gpstToUtc(tGpst time.Time) (time.Time, error) {
	leapSeconds, err := getLeapSeconds(tGpst)
	if err != nil {
		return time.Time{}, err
	}
	tUtc := tGpst.Add(-time.Duration(leapSeconds) * time.Second)
	if utcToGpst(tUtc).Sub(tGpst) != 0 {
		return tUtc.Add(time.Second), nil
	}
	return tUtc, nil
}

// Convert UTC to GPST
func utcToGpst(tUtc time.Time) time.Time {
	leapSeconds, err := getLeapSeconds(tUtc)
	if err != nil {
		return time.Time{}
	}
	return tUtc.Add(time.Duration(leapSeconds) * time.Second)
}
