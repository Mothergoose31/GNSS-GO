package gnss

const (
	// Physical constants
	SPEED_OF_LIGHT = 2.99792458e8 // m/s

	// Physical parameters of the Earth
	EARTH_GM            = 3.986005e14     // m^3/s^2 (gravitational constant * mass of earth)
	EARTH_RADIUS        = 6.3781e6        // m
	EARTH_ROTATION_RATE = 7.2921151467e-5 // rad/s (WGS84 earth rotation rate)

	// GPS system parameters
	GPS_L1 = 1.57542e9 // Hz
	GPS_L2 = 1.22760e9 // Hz
	GPS_L5 = 1.17645e9 // Hz Also E5

	// GLONASS system parameters
	// L1, L2, L3 BANDS
	//    represents the frequency channel,
	//    1602 MHz for the GLONASS L1 band,
	//    562.5 kHz frequency separation between GLONASS carriers in the L1 band,
	//    1246 MHz for the GLONASS L2 band,
	//    437.5 kHz frequency separation between GLONASS carriers in the L2 band,
	//    1201 MHz for the GLONASS L3 band, and
	//    437.5 kHz frequency separation between GLONASS carriers in the L3 band.
	// https://gssc.esa.int/navipedia/index.php/GLONASS_Signal_Plan
	GLONASS_L1       = 1.602e9
	GLONASS_L1_DELTA = 0.5625e6
	GLONASS_L2       = 1.246e9
	GLONASS_L2_DELTA = 0.4375e6
	GLONASS_L3       = 1.201e9
	GLONASS_L3_DELTA = 0.4375e6

	// Galileo system parameters:  Has additional frequencies on E6
	// Source RINEX 2.11 document
	GALILEO_E5B  = 1.207140e9 // Hz
	GALILEO_E5AB = 1.191795e9 // Hz
	GALILEO_E6   = 1.27875e9  // Hz

	// Time constants
	SECS_IN_MIN  = 60
	SECS_IN_HR   = 60 * SECS_IN_MIN
	SECS_IN_DAY  = 24 * SECS_IN_HR
	SECS_IN_WEEK = 7 * SECS_IN_DAY
	SECS_IN_YEAR = 365 * SECS_IN_DAY
)
