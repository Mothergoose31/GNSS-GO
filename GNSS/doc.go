package gnss

// Package gnss provides functionality for parsing and processing GNSS data files
// including RINEX navigation files and SP3 precise orbit files.
//
// The package supports:
// - RINEX 2.01 navigation file parsing
// - GPS satellite position calculation
// - Time system conversions between GPS and UTC
//
// Example usage:
//
//	header, ephemerides, err := gnss.ParseRINEXFile("brdc2050.24g")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get satellite positions at a specific time
//	time := gnss.GPSTimeFromDateTime(time.Now())
//	pos, vel, clockErr, clockRate, err := gnss.GetSatellitePosition(ephemerides[0], time)
