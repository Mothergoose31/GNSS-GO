package gnss

import (
	"errors"
	"fmt"
	"time"
)

const SecondsInWeek = 604800

type GPSTime struct {
	Week int
	Tow  float64
}

type TimeSync struct {
	RefMonoTime time.Duration
	RefGPSTime  GPSTime
}

// =======================================

// ========================================

func GetLeapSeconds(t time.Time) (int, error) {

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

// =======================================

// ========================================

func UtcToGpst(tUtc time.Time) time.Time {
	leapSeconds, err := GetLeapSeconds(tUtc)
	if err != nil {
		return time.Time{}
	}
	return tUtc.Add(time.Duration(leapSeconds) * time.Second)
}

// =======================================

// ========================================

func newTimeSync(monoTime time.Duration, gpsTime GPSTime) *TimeSync {
	return &TimeSync{
		RefMonoTime: monoTime,
		RefGPSTime:  gpsTime,
	}
}

// =======================================

// ========================================

func NewGPSTime(week int, tow float64) GPSTime {
	return GPSTime{week, tow}
}

// =======================================

// ========================================

func (ts *TimeSync) FromDateTime(t time.Time) {
	ts.RefGPSTime = GPSTimeFromDateTime(t)
	ts.RefMonoTime = time.Duration(0)
}

// =======================================

// ========================================
func (ts *TimeSync) FromLogs(rawQcomMeasurementReport, clocks interface{}) error {
	// TODO: LOOK AT QCOM MEASUREMENT REPORT AND CLOCKS
	return errors.New("FromLogs method not implemented")
}

// =======================================

// ========================================

func (ts *TimeSync) Mono2GPS(monoTime time.Duration) GPSTime {
	diffDuration := monoTime - ts.RefMonoTime
	diffSeconds := diffDuration.Seconds()
	return ts.RefGPSTime.Add(diffSeconds)
}

// =======================================

// ========================================

func (ts *TimeSync) GPS2Mono(gpsTime GPSTime) time.Duration {
	diffSeconds := gpsTime.Sub(ts.RefGPSTime)
	return ts.RefMonoTime + time.Duration(diffSeconds*float64(time.Second))
}

// =======================================

// ========================================

func (ts *TimeSync) String() string {
	return fmt.Sprintf("Reference mono time: %v\nReference GPS time: %v", ts.RefMonoTime, ts.RefGPSTime)
}

// =======================================

// ========================================

func GPSTimeFromDateTime(t time.Time) GPSTime {
	wkRef := time.Date(2014, 2, 16, 0, 0, 0, 0, time.UTC)
	refWk := 1780
	wk := int(t.Sub(wkRef).Hours()/24/7) + refWk
	tow := float64(t.Sub(wkRef).Seconds()) - float64((wk-refWk)*SecondsInWeek)
	return GPSTime{Week: wk, Tow: tow}

}

// =======================================

// ========================================

func (t GPSTime) ToDateTime() time.Time {
	gpsEpoch := time.Date(1980, 1, 6, 0, 0, 0, 0, time.UTC)
	duration := time.Duration(t.Week)*7*24*time.Hour + time.Duration(t.Tow*float64(time.Second))
	return gpsEpoch.Add(duration)
}

// =======================================

// ========================================

func (t GPSTime) Sub(other GPSTime) float64 {
	return float64(t.Week-other.Week)*SecondsInWeek + t.Tow - other.Tow
}

// =======================================

// ========================================

func (t GPSTime) Add(seconds float64) GPSTime {
	newWeek := t.Week
	newTow := t.Tow + seconds
	for newTow >= SecondsInWeek {
		newTow -= SecondsInWeek
		newWeek++
	}
	return GPSTime{Week: newWeek, Tow: newTow}
}

// =======================================

// ========================================

func (t GPSTime) ToUTC() time.Time {
	gpst := t.ToDateTime()
	leapSeconds, err := GetLeapSeconds(gpst)
	if err != nil {
		return time.Time{}

	}
	utc := gpst.Add(-time.Duration(leapSeconds) * time.Second)
	if UtcToGpst(utc).Sub(gpst) != 0 {
		return utc.Add(time.Second)
	}
	return utc
}

// =======================================

// ========================================

func UTCToGPST(utc time.Time) GPSTime {
	leapSeconds, err := GetLeapSeconds(utc)
	if err != nil {
		return GPSTime{}
	}
	gpst := utc.Add(time.Duration(leapSeconds) * time.Second)
	return GPSTimeFromDateTime(gpst)
}

// =======================================

// ========================================

func GPSTimeFromGLONASS(cycle, days int, tow float64) GPSTime {
	t := time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC)
	t = t.Add(time.Duration(cycle*(365*4+1)+(days-1)) * 24 * time.Hour)
	t = t.Add(-3 * time.Hour)
	t = t.Add(time.Duration(tow * float64(time.Second)))
	return UTCToGPST(t)
}

// =======================================

// ========================================

func (t GPSTime) String() string {
	return t.ToDateTime().Format("2006-01-02 15:04:05.999999999")
}

// =======================================

// ========================================

func TowToDateTime(tow float64, week int) time.Time {

	gpsEpoch := time.Date(1980, 1, 6, 0, 0, 0, 0, time.UTC)

	towDuration := time.Duration(tow * float64(time.Second))
	weekDuration := time.Duration(week) * 7 * 24 * time.Hour

	return gpsEpoch.Add(towDuration + weekDuration)
}
