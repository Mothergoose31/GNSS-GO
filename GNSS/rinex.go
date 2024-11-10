package gnss

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ParseRINEXFile parses a RINEX format file and returns the header and ephemerides
func ParseRINEXFile(filename string) (*RINEXHeader, []RINEXEphemeris, error) {
	// Determine file version from extension
	ext := filepath.Ext(filename)
	version := ""

	if strings.HasSuffix(filename, ".g") || strings.HasSuffix(filename, ".nav") {
		version = "2.01"
	} else {
		return nil, nil, fmt.Errorf("unsupported RINEX file format: %s", ext)
	}

	// Parse based on version
	switch version {
	case "2.01":
		return ParseRINEXFileV201(filename)
	default:
		return nil, nil, fmt.Errorf("unsupported RINEX version: %s", version)
	}
}

// ParseSP3File parses an SP3 format file and returns the ephemerides
func ParseSP3File(filename string) (*SP3FormatEphemeris, error) {
	// TODO: Implement SP3 parsing
	return nil, fmt.Errorf("SP3 parsing not yet implemented")
}

// GetSatellitePosition calculates satellite position at a given time
func GetSatellitePosition(eph RINEXEphemeris, time GPSTime) ([]float64, []float64, float64, float64, error) {
	gpsEph := GPSEphemeris{}
	// Convert RINEXEphemeris to GPSEphemeris
	// TODO: Implement conversion

	return gpsEph.GetSatInfo(time)
}
