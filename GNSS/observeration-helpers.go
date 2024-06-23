package gnss

type ObservationKind int

const (
	UNKNOWN ObservationKind = iota
	NO_OBSERVATION
	GPS_NED
	ODOMETRIC_SPEED
	PHONE_GYRO
	GPS_VEL
	PSEUDORANGE_GPS
	PSEUDORANGE_RATE_GPS
	SPEED
	NO_ROT
	PHONE_ACCEL
	ORB_POINT
	ECEF_POS
	ORB_ODO_TRANSLATION
	ORB_ODO_ROTATION
	ORB_FEATURES
	MSCKF_TEST
	FEATURE_TRACK_TEST
	LANE_PT
	IMU_FRAME
	PSEUDORANGE_GLONASS
	PSEUDORANGE_RATE_GLONASS
	PSEUDORANGE
	PSEUDORANGE_RATE
)

var observationKindNames = []string{
	"Unknown",
	"No observation",
	"GPS NED",
	"Odometric speed",
	"Phone gyro",
	"GPS velocity",
	"GPS pseudorange",
	"GPS pseudorange rate",
	"Speed",
	"No rotation",
	"Phone acceleration",
	"ORB point",
	"ECEF pos",
	"ORB odometry translation",
	"ORB odometry rotation",
	"ORB features",
	"MSCKF test",
	"Feature track test",
	"Lane ecef point",
	"IMU frame eulers",
	"GLONASS pseudorange",
	"GLONASS pseudorange rate",
	"Pseudorange",
	"Pseudorange rate",
}

func (k ObservationKind) String() string {
	if int(k) < len(observationKindNames) {
		return observationKindNames[k]
	}
	return "Unknown"
}
