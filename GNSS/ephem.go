package gnss

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// type EphemerisType int

// const (
// 	NAV EphemerisType = iota
// 	FINAL_ORBIT
// 	RAPID_ORBIT
// 	ULTRA_RAPID_ORBIT
// 	QCOM_POLY
// )

// type BaseEphemeris struct {
// 	PRN         string
// 	Epoch       GPSTime
// 	EphType     EphemerisType
// 	Healthy     bool
// 	MaxTimeDiff float64
// 	FileEpoch   *GPSTime
// 	FileName    string
// 	FileSource  string
// }

// type GPSEphemeris struct {
// 	BaseEphemeris
// 	Ephemeris
// 	TOE      GPSTime
// 	TOC      GPSTime
// 	SqrtA    float64
// 	capnpMsg *capnp.Message
// }

// type RINEXHeader struct {
// 	Version     float64
// 	Type        string
// 	SatSystem   string
// 	ProgramName string
// 	Agency      string
// 	Date        time.Time
// 	Comments    []string
// 	LeapSeconds int
// }

// type RINEXEphemeris struct {
// 	RINEXHeader           *RINEXHeader
// 	SatelliteID           int       `json:"satellite_id"`
// 	Epoch                 time.Time `json:"epoch"`
// 	ClockBias             float64   `json:"clock_bias"`
// 	RelativeFreqBias      float64   `json:"relativeFreqBias"`
// 	MessageFrameTime      float64   `json:"messageFrameTime"`
// 	PositionX             float64   `json:"positionX"`
// 	VelocityX             float64   `json:"velocityX"`
// 	AccelerationX         float64   `json:"AccelerationX"`
// 	PositionY             float64   `json:"PositionY"`
// 	VelocityY             float64   `json:"VelocityY"`
// 	AccelerationY         float64   `json:"AccelerationY"`
// 	PositionZ             float64   `json:"PositionZ"`
// 	VelocityZ             float64   `json:"VelocityZ"`
// 	AccelerationZ         float64   `json:"AccelerationZ"`
// 	Health                float64   `json:"Health"`
// 	FrequencyChannelOfSet int       `json:"FrequencyChannelOfSet"`
// 	InformationAge        float64   `json:"InformationAge"`
// }

// type GroupedEphemerides struct {
// 	SatelliteID       int               `json:"satelliteID"`
// 	SortedEphemerides []*RINEXEphemeris `json:"sortedEphemerides"`
// }

// type SP3FormatEphemeris struct {
// 	SP3Header *SP3Header
// 	SP3Epochs []*SP3Epoch
// }

// type SP3Header struct {
// 	Version           string
// 	Start             time.Time
// 	NumberOfEpochs    int
// 	DataUsed          string
// 	CoordinateSystem  string
// 	OrbitType         string
// 	Agency            string
// 	GPSWeek           int
// 	SecondsOfWeek     float64
// 	EpochInterval     float64
// 	ModifiedJulianDay int
// 	FractionalDay     float64
// }

// type SP3Epoch struct {
// 	Time    time.Time  `json:"time"`
// 	Entries []SP3Entry `json:"entries"`
// }

// type SP3Entry struct {
// 	SatelliteVehicleNumber string  `json:"satelliteVehicleNumber"`
// 	X                      float64 `json:"x"`
// 	Y                      float64 `json:"y"`
// 	Z                      float64 `json:"z"`
// 	Clock                  float64 `json:"clock"`
// }

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

// func (e *BaseEphemeris) Valid(time GPSTime) bool {
// 	return math.Abs(time.Sub(e.Epoch)) <= e.MaxTimeDiff
// }

// =========================================================================

// =========================================================================

// func NewGPSEphemeris(eph Ephemeris, fileName string) *GPSEphemeris {
// 	toe := NewGPSTime(int(eph.ToeWeek()), float64(eph.Toe()))
// 	toc := NewGPSTime(int(eph.TimeOfClockWeek()), float64(eph.TimeOfClock()))
// 	epoch := toc

// 	return &GPSEphemeris{
// 		BaseEphemeris: BaseEphemeris{
// 			PRN:         fmt.Sprintf("G%02d", eph.SatelliteId()),
// 			Epoch:       epoch,
// 			EphType:     NAV,
// 			Healthy:     eph.SatelliteHealth() == 0,
// 			MaxTimeDiff: 2 * SECS_IN_HR,
// 			FileName:    fileName,
// 		},
// 		Ephemeris: eph,
// 		TOE:       toe,
// 		TOC:       toc,
// 		SqrtA:     math.Sqrt(eph.SemiMajorAxis()),
// 	}
// }

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

func SortEphemerisBySatelliteID(ephs []*RINEXEphemeris) []GroupedEphemerides {
	ephemerisMap := make(map[int][]*RINEXEphemeris)

	for _, eph := range ephs {
		singleEph := RINEXEphemeris{
			SatelliteID:           eph.SatelliteID,
			Epoch:                 eph.Epoch,
			ClockBias:             eph.ClockBias,
			RelativeFreqBias:      eph.RelativeFreqBias,
			MessageFrameTime:      eph.MessageFrameTime,
			PositionX:             eph.PositionX,
			VelocityX:             eph.VelocityX,
			AccelerationX:         eph.AccelerationX,
			PositionY:             eph.PositionY,
			VelocityY:             eph.VelocityY,
			AccelerationY:         eph.AccelerationY,
			PositionZ:             eph.PositionZ,
			VelocityZ:             eph.VelocityZ,
			AccelerationZ:         eph.AccelerationZ,
			Health:                eph.Health,
			FrequencyChannelOfSet: eph.FrequencyChannelOfSet,
			InformationAge:        eph.InformationAge,
		}

		ephemerisMap[eph.SatelliteID] = append(ephemerisMap[eph.SatelliteID], &singleEph)

	}
	var groupedEphemerides []GroupedEphemerides

	var satelliteIDs []int
	for id := range ephemerisMap {
		satelliteIDs = append(satelliteIDs, id)
	}
	sort.Ints(satelliteIDs)

	for _, id := range satelliteIDs {
		group := ephemerisMap[id]

		sort.Slice(group, func(i, j int) bool {
			return group[i].Epoch.Before(group[j].Epoch)
		})

		groupedEphemerides = append(groupedEphemerides, GroupedEphemerides{
			SatelliteID:       id,
			SortedEphemerides: group,
		})
	}

	return groupedEphemerides
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

func ParseRINEXFileV201(filename string) (*RINEXHeader, []*RINEXEphemeris, error) {
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

	ephemerides, err := parseRINEXEphemeris(scanner)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing ephemerides: %v", err)
	}

	return header, ephemerides, nil
}

// =========================================================================

// =========================================================================

func parseRINEXEphemeris(scanner *bufio.Scanner) ([]*RINEXEphemeris, error) {
	var ephemerides []*RINEXEphemeris
	var currentLines []string
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if len(line) == 0 {
			fmt.Printf("Debug: Skipping empty line %d\n", lineCount)
			continue
		}

		if len(line) > 0 && line[0] >= '1' && line[0] <= '9' {
			if len(currentLines) > 0 {
				eph, err := processEphemerisLines(currentLines)
				if err != nil {
					err = fmt.Errorf("error processing ephemeris lines: %v", err)
					return nil, err
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
				err = fmt.Errorf("error processing ephemeris lines: %v", err)
				return nil, err
			} else {
				ephemerides = append(ephemerides, eph)
			}
			currentLines = nil
		}
	}

	if len(currentLines) > 0 {
		eph, err := processEphemerisLines(currentLines)
		if err != nil {
			err = fmt.Errorf("error processing ephemeris lines: %v", err)
			return nil, err
		} else {
			ephemerides = append(ephemerides, eph)
		}
	}

	fmt.Printf("Debug: Total ephemerides parsed: %d\n", len(ephemerides))
	return ephemerides, nil
}

// =========================================================================

// =========================================================================
//  GLONASS EPHEMERIS FORMAT (RINEX 2.01)
// https://files.igs.org/pub/data/format/rinex2.txt
// https://igs.org/formats-and-standards/
// https://cddis.nasa.gov/Data_and_Derived_Products/GNSS/daily_30second_data.html
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
func processEphemerisLines(lines []string) (*RINEXEphemeris, error) {
	if len(lines) != 4 {
		return nil, fmt.Errorf("invalid number of lines for ephemeris record: %d", len(lines))
	}

	eph := &RINEXEphemeris{}

	if len(lines[0]) < 79 {
		return nil, fmt.Errorf("line 1 is too short: %d characters", len(lines[0]))
	}
	eph.SatelliteID = parseInt(lines[0][0:2])
	year := parseInt(lines[0][3:5])

	if year < 100 {
		if year < 80 {
			year += 2000
		} else {
			year += 1900
		}
	}

	month := parseInt(lines[0][6:8])

	day := parseInt(lines[0][9:11])
	hour := parseInt(lines[0][12:14])
	min := parseInt(lines[0][16:17])
	sec := parseFloat(lines[0][19:22])
	eph.Epoch = time.Date(year, time.Month(month), day, hour, min, int(sec), int((sec-float64(int(sec)))*1e9), time.UTC)
	eph.ClockBias = parseFloat(lines[0][23:41])
	eph.RelativeFreqBias = parseFloat(lines[0][42:61])
	eph.MessageFrameTime = parseFloat(lines[0][61:79])

	for i := 1; i < 4; i++ {
		if len(lines[i]) < 79 {
			fmt.Printf("Debug: Line %d is short: %d characters\n", i+1, len(lines[i]))
			continue
		}
		switch i {
		case 1:
			eph.PositionX = parseFloat(lines[i][3:22])
			eph.VelocityX = parseFloat(lines[i][22:41])
			eph.AccelerationX = parseFloat(lines[i][41:59])
			eph.Health = parseFloat(lines[i][61:79])
		case 2:
			eph.PositionY = parseFloat(lines[i][3:22])
			eph.VelocityY = parseFloat(lines[i][22:41])
			eph.AccelerationY = parseFloat(lines[i][41:59])
			freqNum := parseFloat(lines[i][60:79])
			eph.FrequencyChannelOfSet = int(freqNum)
		case 3:
			eph.PositionZ = parseFloat(lines[i][3:22])
			eph.VelocityZ = parseFloat(lines[i][22:41])
			eph.AccelerationZ = parseFloat(lines[i][41:59])
			eph.InformationAge = parseFloat(lines[i][61:79])
		}
	}

	return eph, nil
}

// =========================================================================

// =========================================================================

// TODO CREATE PARSER  FOR RINEX SP3 FILE
// https://files.igs.org/pub/data/format/sp3c.txt#:~:text=The%20basic%20format%20of%20an,correction%20rate%2Dof%2Dchange.
// https://files.igs.org/pub/data/format/sp3c.txt
// https://files.igs.org/pub/data/format/sp3d.pdf

//  AFTER THE HEADER THERE ARE SATELITE ID  LINES  PG01, PJ02 , PC01, PE01,
//  THOSE CORESPOND WITH
// G: GPS (US)
// R: GLONASS (Russia)
// E: Galileo (EU)
// C: BeiDou (China)
// J: QZSS (Japan)

func ParseSP3FormatEphemeris(filename string) (*SP3FormatEphemeris, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	SP3FormatEphemeris := &SP3FormatEphemeris{
		Header: SP3Header{},
		Epochs: []SP3Epoch{},
	}
	// var currentEpoch *SP3Epoch
	lineCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if lineCount <= 27 {
			if err := parseHeaderLine(line, &SP3FormatEphemeris.Header, lineCount); err != nil {
				return nil, fmt.Errorf("error parsing header line: %v", err)
			}
		}
	}

	return SP3FormatEphemeris, nil
}

func parseHeaderLine(line string, header *SP3Header, lineNumber int) error {
	switch {
	case lineNumber == 1:
		return parseFirstHeaderLine(line, header)
	case lineNumber == 2:
		return parseSecondLine(line, header)
	case strings.HasPrefix(line, "++"):

	}
	return nil
}

func parseFirstHeaderLine(line string, header *SP3Header) error {
	fmt.Println(line)
	// starTimeStr := strings.Join(fields[2:8], " ")
	// fmt.Println(starTimeStr)

	return nil
}

func parseSecondLine(line string, header *SP3Header) error {
	fmt.Println(line)
	return nil
}

// =========================================================================
//  HELPERS
// =========================================================================

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "D", "E")
	if s == "" {
		return 0
	}

	if strings.HasPrefix(s, "-") {
		s = "-" + strings.Trim(s, "-")
	}

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

func parseInt(s string) int {
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

// =========================================================================

// =========================================================================

func PrintSortedEphemerides(sorted []GroupedEphemerides) {
	fmt.Println("Grouped Ephemerides Summary:")
	for _, group := range sorted {
		fmt.Printf("Satellite ID: %d\n", group.SatelliteID)
		fmt.Printf("  Number of Ephemerides: %d\n", len(group.SortedEphemerides))

		if len(group.SortedEphemerides) > 0 {
			firstEph := group.SortedEphemerides[0]
			lastEph := group.SortedEphemerides[len(group.SortedEphemerides)-1]

			fmt.Printf("  Time Range: %s to %s\n",
				formatTime(firstEph.Epoch),
				formatTime(lastEph.Epoch))

			fmt.Printf("  First Ephemeris Data:\n")
			printData(firstEph)
		}

		fmt.Println(strings.Repeat("-", 60))
	}
}

// =========================================================================

// =========================================================================

func PrintRINEXHeader(header *RINEXHeader) {
	fmt.Println("RINEX Header:")
	fmt.Printf("  Version: %.2f\n", header.Version)
	fmt.Printf("  Type: %s\n", header.Type)
	fmt.Printf("  Satellite System: %s\n", header.SatSystem)
	fmt.Printf("  Program: %s\n", header.ProgramName)
	fmt.Printf("  Agency: %s\n", header.Agency)
	fmt.Printf("  Date: %s\n", formatTime(header.Date))
	fmt.Println("  Comments:")
	for _, comment := range header.Comments {
		fmt.Printf("    %s\n", comment)
	}
	fmt.Printf("  Leap Seconds: %d\n", header.LeapSeconds)
	fmt.Println(strings.Repeat("-", 60))
}

// =========================================================================

// =========================================================================

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// =========================================================================

// =========================================================================

func printData(eph *RINEXEphemeris) {
	fmt.Printf("    Satellite ID: %d\n", eph.SatelliteID)
	fmt.Printf("    Epoch: %s\n", formatTime(eph.Epoch))
	fmt.Printf("    Clock Bias: %.12e\n", eph.ClockBias)
	fmt.Printf("    Relative Frequency Bias: %.12e\n", eph.RelativeFreqBias)
	fmt.Printf("    Message Frame Time: %.6f\n", eph.MessageFrameTime)
	fmt.Printf("    Position X: %.3f\n", eph.PositionX)
	fmt.Printf("    Velocity X: %.6f\n", eph.VelocityX)
	fmt.Printf("    Acceleration X: %.9f\n", eph.AccelerationX)
	fmt.Printf("    Position Y: %.3f\n", eph.PositionY)
	fmt.Printf("    Velocity Y: %.6f\n", eph.VelocityY)
	fmt.Printf("    Acceleration Y: %.9f\n", eph.AccelerationY)
	fmt.Printf("    Position Z: %.3f\n", eph.PositionZ)
	fmt.Printf("    Velocity Z: %.6f\n", eph.VelocityZ)
	fmt.Printf("    Acceleration Z: %.9f\n", eph.AccelerationZ)
	fmt.Printf("    Health: %.0f\n", eph.Health)
	fmt.Printf("    Frequency Channel Offset: %d\n", eph.FrequencyChannelOfSet)
	fmt.Printf("    Information Age: %.6f\n", eph.InformationAge)
}
