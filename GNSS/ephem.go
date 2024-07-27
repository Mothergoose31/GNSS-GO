package gnss

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"capnproto.org/go/capnp/v3"
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
	TOE      GPSTime
	TOC      GPSTime
	SqrtA    float64
	capnpMsg *capnp.Message
}

type GLONASSEphemeris struct {
	SatelliteID      int
	Epoch            time.Time
	ClockBias        float64
	RelativeFreqBias float64
	MessageFrameTime float64
	PositionX        float64
	VelocityX        float64
	AccelerationX    float64
	PositionY        float64
	VelocityY        float64
	AccelerationY    float64
	PositionZ        float64
	VelocityZ        float64
	AccelerationZ    float64
	Health           float64
	FrequencyNumber  int
	InformationAge   float64
}

type RINEXHeader struct {
	Version     float64
	Type        string
	SatSystem   string
	ProgramName string
	Agency      string
	Date        time.Time
	Comments    []string
	LeapSeconds int
}

// =========================================================================

// =========================================================================

func (g *GPSEphemeris) DebugString() string {
	fmt.Println("DebugString==========")
	var result string
	result += fmt.Sprintf("GPSEphemeris:\n")
	result += fmt.Sprintf("  BaseEphemeris:\n")
	result += fmt.Sprintf("    PRN: %s\n", g.BaseEphemeris.PRN)
	result += fmt.Sprintf("    Epoch: %v\n", g.BaseEphemeris.Epoch)
	result += fmt.Sprintf("    EphType: %v\n", g.BaseEphemeris.EphType)
	result += fmt.Sprintf("    Healthy: %v\n", g.BaseEphemeris.Healthy)
	result += fmt.Sprintf("    MaxTimeDiff: %f\n", g.BaseEphemeris.MaxTimeDiff)
	result += fmt.Sprintf("    FileName: %s\n", g.BaseEphemeris.FileName)
	result += fmt.Sprintf("  TOE: %v\n", g.TOE)
	result += fmt.Sprintf("  TOC: %v\n", g.TOC)
	result += fmt.Sprintf("  SqrtA: %f\n", g.SqrtA)

	v := reflect.ValueOf(g.Ephemeris)
	t := v.Type()
	for i := 0; i < v.NumMethod(); i++ {
		method := t.Method(i)
		if method.Type.NumIn() == 1 && method.Type.NumOut() == 1 {
			result += fmt.Sprintf("  %s: %v\n", method.Name, v.Method(i).Call(nil)[0].Interface())
		}
	}

	return result
}

//=============================================================================

// =========================================================================

func (e *BaseEphemeris) Valid(time GPSTime) bool {
	return math.Abs(time.Sub(e.Epoch)) <= e.MaxTimeDiff
}

// =========================================================================

// =========================================================================

func NewGPSEphemeris(eph Ephemeris, fileName string) *GPSEphemeris {
	toe := NewGPSTime(int(eph.ToeWeek()), float64(eph.Toe()))
	toc := NewGPSTime(int(eph.TimeOfClockWeek()), float64(eph.TimeOfClock()))
	epoch := toc

	return &GPSEphemeris{
		BaseEphemeris: BaseEphemeris{
			PRN:         fmt.Sprintf("G%02d", eph.SatelliteId()),
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

// =========================================================================

// =========================================================================

// # http://gauss.gge.unb.ca/GLONASS.ICD.pdf
func (e *GPSEphemeris) GetSatInfo(time GPSTime) ([]float64, []float64, float64, float64, error) {
	fmt.Println("GetSatInfo==========")
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

	einstein := -4.442797633e-10 * e.Eccentricity() * e.SqrtA * math.Sin(ea)

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

// =========================================================================

// =========================================================================

func parseRINEXHeader(scanner *bufio.Scanner) (*RINEXHeader, error) {
	header := &RINEXHeader{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 60 {
			continue
		}
		label := strings.TrimRight(line[60:], " ")
		switch label {
		case "RINEX VERSION / TYPE":
			header.Version, _ = strconv.ParseFloat(strings.TrimSpace(line[:20]), 64)
			header.Type = strings.TrimSpace(line[20:40])
			header.SatSystem = strings.TrimSpace(line[40:60])
		case "PGM / RUN BY / DATE":
			header.ProgramName = strings.TrimSpace(line[:20])
			header.Agency = strings.TrimSpace(line[20:40])
			header.Date, _ = time.Parse("02-Jan-06 15:04", strings.TrimSpace(line[40:60]))
		case "COMMENT":
			header.Comments = append(header.Comments, strings.TrimSpace(line[:60]))
		case "END OF HEADER":
			return header, nil
		}
	}
	return nil, fmt.Errorf("end of header not found")
}

// =========================================================================

// =========================================================================

func ParseRINEXFile(filename string) (*RINEXHeader, []*GLONASSEphemeris, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header, err := parseRINEXHeader(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing header: %v", err)
	}

	ephemerides, err := parseGLONASSEphemeris(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing ephemerides: %v", err)
	}

	return header, ephemerides, nil
}

// =========================================================================

// =========================================================================

// https://files.igs.org/pub/data/format/rinex2.txt
// 1 24  7 23  0 15  0.0 0.922027975321D-04 0.909494701773D-12 0.000000000000D+00
//     0.147764956055D+05 0.151113414764D+01 0.000000000000D+00 0.000000000000D+00
//    -0.159414018555D+05-0.103984928131D+01 0.186264514923D-08 0.100000000000D+01
//     0.133588305664D+05-0.291043663025D+01-0.186264514923D-08 0.000000000000D+00
// Line 1:
// - Epoch: 2024-07-23 00:15:00.0
// - SV clock bias (TauN): 0.922027975321E-04 seconds
// - SV relative frequency bias (GammaN): 0.909494701773E-12
// - Message frame time (tk): 0.0 seconds
//
// Line 2:
// - X coordinate: 14776.4956055 km
// - X velocity: 0.151113414764 km/s
// - X acceleration: 0.0 km/s²
// - Health (0 = OK): 0
//
// Line 3:
// - Y coordinate: -15941.4018555 km
// - Y velocity: -1.03984928131 km/s
// - Y acceleration: 0.186264514923E-08 km/s²
// - Frequency number: 1
//
// Line 4:
// - Z coordinate: 13358.8305664 km
// - Z velocity: -2.91043663025 km/s
// - Z acceleration: -0.186264514923E-08 km/s²
// - Age of operation: 0 days

// lines are a maximum of 80 characters long , NEED TO BE KEPT AT 79 CHARACTERS
func parseGLONASSEphemeris(scanner *bufio.Scanner) ([]*GLONASSEphemeris, error) {
	var ephemerides []*GLONASSEphemeris
	var currentLines []string
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		fmt.Printf("Debug: Line %d: %s\n", lineCount, line)
		fmt.Println(line)

		if len(line) == 0 {
			fmt.Printf("Debug: Skipping empty line %d\n", lineCount)
			continue
		}

		if len(line) > 0 && line[0] >= '1' && line[0] <= '9' {
			if len(currentLines) > 0 {
				eph, err := processEphemerisLines(currentLines)
				if err != nil {

				} else {
					ephemerides = append(ephemerides, eph)
					fmt.Printf("Debug: Processed ephemeris for satellite %d\n", eph.SatelliteID)
				}
			}
			currentLines = []string{line}
		} else {
			currentLines = append(currentLines, line)
		}

		if len(currentLines) == 4 {
			eph, err := processEphemerisLines(currentLines)
			if err != nil {
				fmt.Printf("Debug: Error processing ephemeris lines: %v\n", err)
			} else {
				ephemerides = append(ephemerides, eph)
				fmt.Printf("Debug: Processed ephemeris for satellite %d\n", eph.SatelliteID)
			}
			currentLines = nil
		}
	}

	if len(currentLines) > 0 {
		eph, err := processEphemerisLines(currentLines)
		if err != nil {
			fmt.Printf("Debug: Error processing final ephemeris lines: %v\n", err)
		} else {
			ephemerides = append(ephemerides, eph)
			fmt.Printf("Debug: Processed final ephemeris for satellite %d\n", eph.SatelliteID)
		}
	}

	fmt.Printf("Debug: Total ephemerides parsed: %d\n", len(ephemerides))
	return ephemerides, nil
}

// =========================================================================

// =========================================================================
// 1 24  7 23  0 15  0.0 0.922027975321D-04 0.909494701773D-12 0.000000000000D+00
//     0.147764956055D+05 0.151113414764D+01 0.000000000000D+00 0.000000000000D+00
//    -0.159414018555D+05-0.103984928131D+01 0.186264514923D-08 0.100000000000D+01
//     0.133588305664D+05-0.291043663025D+01-0.186264514923D-08 0.000000000000D+00
// LINE 1:
// - Epoch: 2024-07-23 00:15:00.0
// - SV clock bias (TauN): 0.922027975321E-04 seconds
// - SV relative frequency bias (GammaN): 0.909494701773E-12
// - Message frame time (tk): 0.0 seconds
//
// LINE 2:
// - X coordinate: 14776.4956055 km
// - X velocity: 0.151113414764 km/s
// - X acceleration: 0.0 km/s²
// - Health (0 = OK): 0
//
// LINE 3:
// - Y coordinate: -15941.4018555 km
// - Y velocity: -1.03984928131 km/s
// - Y acceleration: 0.186264514923E-08 km/s²
// - Frequency number: 1
//
// LINE 4:
// - Z coordinate: 13358.8305664 km
// - Z velocity: -2.91043663025 km/s
// - Z acceleration: -0.186264514923E-08 km/s²
// - Age of operation: 0 days

// RINEX version  2.01,  lines are a maximum of 80 characters long , NEED TO BE KEPT AT 79 CHARACTERS WHEN PARSSING
func processEphemerisLines(lines []string) (*GLONASSEphemeris, error) {
	if len(lines) != 4 {
		return nil, fmt.Errorf("invalid number of lines for ephemeris record: %d", len(lines))
	}

	eph := &GLONASSEphemeris{}

	// Process first line
	if len(lines[0]) < 79 {
		return nil, fmt.Errorf("line 1 is too short: %d characters", len(lines[0]))
	}
	eph.SatelliteID = safeParseInt(lines[0][:2])
	year := safeParseInt(lines[0][3:7])
	month := safeParseInt(lines[0][8:10])
	day := safeParseInt(lines[0][11:13])
	hour := safeParseInt(lines[0][14:16])
	min := safeParseInt(lines[0][17:19])
	sec := safeParseFloat(lines[0][19:24])
	eph.Epoch = time.Date(year, time.Month(month), day, hour, min, int(sec), int((sec-float64(int(sec)))*1e9), time.UTC)
	eph.ClockBias = safeParseFloat(lines[0][23:42])
	eph.RelativeFreqBias = safeParseFloat(lines[0][42:61])
	eph.MessageFrameTime = safeParseFloat(lines[0][61:79])

	// Process remaining lines
	for i := 1; i < 4; i++ {
		if len(lines[i]) < 79 {
			fmt.Printf("Debug: Line %d is short: %d characters\n", i+1, len(lines[i]))
			continue
		}
		switch i {
		case 1:
			eph.PositionX = safeParseFloat(lines[i][4:23])
			eph.VelocityX = safeParseFloat(lines[i][23:42])
			eph.AccelerationX = safeParseFloat(lines[i][42:61])
			eph.Health = safeParseFloat(lines[i][61:79])
		case 2:
			eph.PositionY = safeParseFloat(lines[i][4:23])
			eph.VelocityY = safeParseFloat(lines[i][23:42])
			eph.AccelerationY = safeParseFloat(lines[i][42:61])
			eph.FrequencyNumber = safeParseInt(lines[i][61:79])
		case 3:
			eph.PositionZ = safeParseFloat(lines[i][4:23])
			eph.VelocityZ = safeParseFloat(lines[i][23:42])
			eph.AccelerationZ = safeParseFloat(lines[i][42:61])
			eph.InformationAge = safeParseFloat(lines[i][61:79])
		}
	}

	return eph, nil
}

// =========================================================================

// =========================================================================

func safeParseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "D", "E")
	if s == "" {
		return 0
	}

	if strings.HasSuffix(s, "-") {
		s = "-" + strings.TrimSuffix(s, "-")
	}
	s = strings.ReplaceAll(s, " -", "-")

	s = strings.ReplaceAll(s, " ", "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Debug: Error parsing float: %s, error: %v\n", s, err)
		return 0
	}
	return f
}

// =========================================================================

// =========================================================================

func safeParseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Debug: Error parsing integer: %s, error: %v\n", s, err)
		return 0
	}
	return i
}
