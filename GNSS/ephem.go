package gnss

import (
	"bufio"
	math "math"

	"fmt"

	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"capnproto.org/go/capnp/v3"
)

// type EphemerisType int

const (
	NAV EphemerisType = iota
	FINAL_ORBIT
	RAPID_ORBIT
	ULTRA_RAPID_ORBIT
	QCOM_POLY
)

// =========================================================================

// =========================================================================

func (g *GPSEphemeris) DebugString() string {
	var result string
	result += fmt.Sprintf("GPSEphemeris:\n")
	BaseEphemeris, err := g.BaseEphemeris()
	if err == nil {
		result += fmt.Sprintf("BaseEphemeris: %v\n", BaseEphemeris)
		prn, err := BaseEphemeris.PseudoRandomNumber()
		result += fmt.Sprintf("   PRN:%s (Error: %v)\n", prn, err)

		epoch, err := BaseEphemeris.Epoch()
		result += fmt.Sprintf("   Epoch:%v (Error: %v)\n", epoch, err)

		ephType := BaseEphemeris.EphemerisType()
		result += fmt.Sprintf("   EphemerisType:%v\n", ephType)

		healthy := BaseEphemeris.IsHealthy()
		result += fmt.Sprintf("   Healthy:%v\n", healthy)

		maxTimeDiff := BaseEphemeris.MaximumTimeDifference()
		result += fmt.Sprintf("   MaximumTimeDifference:%v\n", maxTimeDiff)

		fileName, err := BaseEphemeris.FileName()
		result += fmt.Sprintf("   FileName:%s (Error: %v)\n", fileName, err)
	} else {
		result += fmt.Sprintf("BaseEphemeris: Error: %v\n", err)
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

func (e GPSEphemeris) GetSatInfo(time GPSTime) ([]float64, []float64, float64, float64, error) {
	fmt.Println("GetSatInfo==========")

	ephData, err := e.EphemerisData()
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to get ephemeris data: %v", err)
	}

	baseEph, err := e.BaseEphemeris()
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to get base ephemeris: %v", err)
	}

	if !baseEph.IsHealthy() {
		return nil, nil, 0, 0, errors.New("unhealthy ephemeris")
	}

	toc, err := e.Toc()
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to get time of clock: %v", err)
	}

	toe, err := e.Toe()
	if err != nil {
		return nil, nil, 0, 0, fmt.Errorf("failed to get time of ephemeris: %v", err)
	}

	tdiff := time.Sub(toc)
	clockErr := ephData.Af0() + tdiff*(ephData.Af1()+tdiff*ephData.Af2())
	clockRateErr := ephData.Af1() + 2*tdiff*ephData.Af2()

	tdiff = time.Sub(toe)

	sqrtA := e.SquareRootOfSemiMajorAxis()
	a := sqrtA * sqrtA
	maDot := math.Sqrt(EARTH_GM/(a*a*a)) + ephData.DeltaN()
	ma := ephData.M0() + maDot*tdiff

	ea := ma
	eaOld := 2222.0
	for math.Abs(ea-eaOld) > 1.0e-14 {
		eaOld = ea
		ea = ea + (ma-eaOld+ephData.Ecc()*math.Sin(eaOld))/(1.0-ephData.Ecc()*math.Cos(eaOld))
	}

	eaDot := maDot / (1.0 - ephData.Ecc()*math.Cos(ea))

	// Relativistic correction term
	einstein := -4.442807633e-10 * ephData.Ecc() * sqrtA * math.Sin(ea)

	// Begin calc for True Anomaly and Argument of Latitude
	tempd2 := math.Sqrt(1.0 - ephData.Ecc()*ephData.Ecc())
	al := math.Atan2(tempd2*math.Sin(ea), math.Cos(ea)-ephData.Ecc()) + ephData.Omega()
	alDot := tempd2 * eaDot / (1.0 - ephData.Ecc()*math.Cos(ea))

	// Calculate corrected argument of latitude based on position
	cal := al + ephData.Cus()*math.Sin(2.0*al) + ephData.Cuc()*math.Cos(2.0*al)
	calDot := alDot * (1.0 + 2.0*(ephData.Cus()*math.Cos(2.0*al)-ephData.Cuc()*math.Sin(2.0*al)))

	// Calculate corrected radius based on argument of latitude
	r := a*(1.0-ephData.Ecc()*math.Cos(ea)) + ephData.Crc()*math.Cos(2.0*al) + ephData.Crs()*math.Sin(2.0*al)
	rDot := a*ephData.Ecc()*math.Sin(ea)*eaDot + 2.0*alDot*(ephData.Crs()*math.Cos(2.0*al)-ephData.Crc()*math.Sin(2.0*al))

	// Calculate inclination based on argument of latitude
	inc := ephData.I0() + ephData.IDot()*tdiff + ephData.Cic()*math.Cos(2.0*al) + ephData.Cis()*math.Sin(2.0*al)
	incDot := ephData.IDot() + 2.0*alDot*(ephData.Cis()*math.Cos(2.0*al)-ephData.Cic()*math.Sin(2.0*al))

	// Calculate position and velocity in orbital plane
	x := r * math.Cos(cal)
	y := r * math.Sin(cal)
	xDot := rDot*math.Cos(cal) - y*calDot
	yDot := rDot*math.Sin(cal) + x*calDot

	// Corrected longitude of ascending node
	omDot := ephData.OmegaDot() - EARTH_ROTATION_RATE
	om := ephData.Omega0() + tdiff*omDot - EARTH_ROTATION_RATE*toe.TimeOfWeek()

	// Compute the satellite's position in Earth-Centered Earth-Fixed coordinates
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

// func SortEphemerisBySatelliteID(ephs []*RINEXEphemeris) []GroupedEphemerides {
// 	ephemerisMap := make(map[int][]*RINEXEphemeris)

// 	for _, eph := range ephs {
// 		singleEph := RINEXEphemeris{
// 			SatelliteID:           eph.SatelliteID,
// 			Epoch:                 eph.Epoch,
// 			ClockBias:             eph.ClockBias,
// 			RelativeFreqBias:      eph.RelativeFreqBias,
// 			MessageFrameTime:      eph.MessageFrameTime,
// 			PositionX:             eph.PositionX,
// 			VelocityX:             eph.VelocityX,
// 			AccelerationX:         eph.AccelerationX,
// 			PositionY:             eph.PositionY,
// 			VelocityY:             eph.VelocityY,
// 			AccelerationY:         eph.AccelerationY,
// 			PositionZ:             eph.PositionZ,
// 			VelocityZ:             eph.VelocityZ,
// 			AccelerationZ:         eph.AccelerationZ,
// 			Health:                eph.Health,
// 			FrequencyChannelOfSet: eph.FrequencyChannelOfSet,
// 			InformationAge:        eph.InformationAge,
// 		}

// 		ephemerisMap[eph.SatelliteID] = append(ephemerisMap[eph.SatelliteID], &singleEph)

// 	}
// 	var groupedEphemerides []GroupedEphemerides

// 	var satelliteIDs []int
// 	for id := range ephemerisMap {
// 		satelliteIDs = append(satelliteIDs, id)
// 	}
// 	sort.Ints(satelliteIDs)

// 	for _, id := range satelliteIDs {
// 		group := ephemerisMap[id]

// 		sort.Slice(group, func(i, j int) bool {
// 			return group[i].Epoch.Before(group[j].Epoch)
// 		})

// 		groupedEphemerides = append(groupedEphemerides, GroupedEphemerides{
// 			SatelliteID:       id,
// 			SortedEphemerides: group,
// 		})
// 	}

// 	return groupedEphemerides
// }

// =========================================================================

// =========================================================================
func parseRINEXHeader(scanner *bufio.Scanner) (RINEXHeader, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return RINEXHeader{}, fmt.Errorf("failed to create new message: %v", err)
	}
	fmt.Println(msg)
	header, err := NewRINEXHeader(seg)
	if err != nil {
		return RINEXHeader{}, fmt.Errorf("failed to create new RINEXHeader: %v", err)
	}

	var comments []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 60 {
			continue
		}
		label := strings.TrimRight(line[60:], " ")
		switch label {
		case "RINEX VERSION / TYPE":
			version, _ := strconv.ParseFloat(strings.TrimSpace(line[:20]), 64)
			header.SetVersion(version)
			header.SetType(strings.TrimSpace(line[20:40]))
			header.SetSatelliteSystem(strings.TrimSpace(line[40:60]))
		case "PGM / RUN BY / DATE":
			header.SetProgramName(strings.TrimSpace(line[:20]))
			header.SetAgency(strings.TrimSpace(line[20:40]))
			parsedDate, _ := time.Parse("02-Jan-06 15:04", strings.TrimSpace(line[40:60]))

			newDate, err := NewTime(seg)
			if err != nil {
				return RINEXHeader{}, fmt.Errorf("failed to create new Time struct: %v", err)
			}

			unixTime := parsedDate.Unix()
			nanos := int64(parsedDate.Nanosecond())

			newDate.SetSeconds(unixTime)
			newDate.SetNanoseconds(int32(nanos))

			if err := header.SetDate(newDate); err != nil {
				return RINEXHeader{}, fmt.Errorf("failed to set Date: %v", err)
			}
		case "COMMENT":
			comments = append(comments, strings.TrimSpace(line[:60]))
		case "END OF HEADER":
			commentList, err := header.NewComments(int32(len(comments)))
			if err != nil {
				return RINEXHeader{}, fmt.Errorf("failed to create new Comments: %v", err)
			}
			for i, comment := range comments {
				commentList.Set(i, comment)
			}
			header.SetComments(commentList)
			return header, nil
		}
	}
	return RINEXHeader{}, fmt.Errorf("end of header not found")
}

// =========================================================================

// =========================================================================

func ParseRINEXFileV201(filename string) (*RINEXHeader, []RINEXEphemeris, error) {
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

	return &header, ephemerides, nil
}

// =========================================================================

// =========================================================================
func parseRINEXEphemeris(scanner *bufio.Scanner) ([]RINEXEphemeris, error) {
	var ephemerides []RINEXEphemeris
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
					return nil, fmt.Errorf("error processing ephemeris lines: %v", err)
				}
				ephemerides = append(ephemerides, eph)
				fmt.Printf("Debug: Processed ephemeris for satellite %d\n", eph.SatelliteId())
			}
			currentLines = []string{line}
		} else {
			currentLines = append(currentLines, line)
		}

		if len(currentLines) == 4 {
			eph, err := processEphemerisLines(currentLines)
			if err != nil {
				return nil, fmt.Errorf("error processing ephemeris lines: %v", err)
			}
			ephemerides = append(ephemerides, eph)
			currentLines = nil
		}
	}

	if len(currentLines) > 0 {
		eph, err := processEphemerisLines(currentLines)
		if err != nil {
			return nil, fmt.Errorf("error processing ephemeris lines: %v", err)
		}
		ephemerides = append(ephemerides, eph)
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
func processEphemerisLines(lines []string) (RINEXEphemeris, error) {
	if len(lines) != 4 {
		return RINEXEphemeris{}, fmt.Errorf("invalid number of lines for ephemeris record: %d", len(lines))
	}

	// Create a new RINEXEphemeris
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return RINEXEphemeris{}, fmt.Errorf("failed to create new message: %v", err)
	}
	fmt.Println(msg)
	eph, err := NewRootRINEXEphemeris(seg)
	if err != nil {
		return RINEXEphemeris{}, fmt.Errorf("failed to create new RINEXEphemeris: %v", err)
	}

	if len(lines[0]) < 79 {
		return RINEXEphemeris{}, fmt.Errorf("line 1 is too short: %d characters", len(lines[0]))
	}

	eph.SetSatelliteId(int32(parseInt(lines[0][0:2])))
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

	epochTime := time.Date(year, time.Month(month), day, hour, min, int(sec), int((sec-float64(int(sec)))*1e9), time.UTC)
	epochCapnp, err := NewTime(seg)
	if err != nil {
		return RINEXEphemeris{}, fmt.Errorf("failed to create new Time: %v", err)
	}
	epochCapnp.SetSeconds(epochTime.Unix())
	epochCapnp.SetNanoseconds(int32(epochTime.Nanosecond()))
	eph.SetEpoch(epochCapnp)

	eph.SetClockBias(parseFloat(lines[0][23:41]))
	eph.SetRelativeFrequencyBias(parseFloat(lines[0][42:61]))
	eph.SetMessageFrameTime(parseFloat(lines[0][61:79]))

	for i := 1; i < 4; i++ {
		if len(lines[i]) < 79 {
			fmt.Printf("Debug: Line %d is short: %d characters\n", i+1, len(lines[i]))
			continue
		}
		switch i {
		case 1:
			eph.SetPositionX(parseFloat(lines[i][3:22]))
			eph.SetVelocityX(parseFloat(lines[i][22:41]))
			eph.SetAccelerationX(parseFloat(lines[i][41:59]))
			eph.SetHealth(parseFloat(lines[i][61:79]))
		case 2:
			eph.SetPositionY(parseFloat(lines[i][3:22]))
			eph.SetVelocityY(parseFloat(lines[i][22:41]))
			eph.SetAccelerationY(parseFloat(lines[i][41:59]))
			freqNum := parseFloat(lines[i][60:79])
			eph.SetFrequencyChannelOffset(int32(freqNum))
		case 3:
			eph.SetPositionZ(parseFloat(lines[i][3:22]))
			eph.SetVelocityZ(parseFloat(lines[i][22:41]))
			eph.SetAccelerationZ(parseFloat(lines[i][41:59]))
			eph.SetInformationAge(parseFloat(lines[i][61:79]))
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

// func ParseSP3FormatEphemeris(filename string) (*SP3FormatEphemeris, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("error opening file: %v", err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	SP3FormatEphemeris := &SP3FormatEphemeris{
// 		Header: SP3Header{},
// 		Epochs: []SP3Epoch{},
// 	}
// 	// var currentEpoch *SP3Epoch
// 	lineCount := 0

// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		lineCount++

// 		if lineCount <= 27 {
// 			if err := parseHeaderLine(line, &SP3FormatEphemeris.Header, lineCount); err != nil {
// 				return nil, fmt.Errorf("error parsing header line: %v", err)
// 			}
// 		}
// 	}

// 	return SP3FormatEphemeris, nil
// }

// func parseHeaderLine(line string, header *SP3Header, lineNumber int) error {
// 	switch {
// 	case lineNumber == 1:
// 		return parseFirstHeaderLine(line, header)
// 	case lineNumber == 2:
// 		return parseSecondLine(line, header)
// 	case strings.HasPrefix(line, "++"):

// 	}
// 	return nil
// }

// func parseFirstHeaderLine(line string, header *SP3Header) error {
// 	fmt.Println(line)
// 	// starTimeStr := strings.Join(fields[2:8], " ")
// 	// fmt.Println(starTimeStr)

// 	return nil
// }

// func parseSecondLine(line string, header *SP3Header) error {
// 	fmt.Println(line)
// 	return nil
// }

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

// func PrintSortedEphemerides(sorted []GroupedEphemerides) {
// 	fmt.Println("Grouped Ephemerides Summary:")
// 	for _, group := range sorted {
// 		fmt.Printf("Satellite ID: %d\n", group.SatelliteID)
// 		fmt.Printf("  Number of Ephemerides: %d\n", len(group.SortedEphemerides))

// 		if len(group.SortedEphemerides) > 0 {
// 			firstEph := group.SortedEphemerides[0]
// 			lastEph := group.SortedEphemerides[len(group.SortedEphemerides)-1]

// 			fmt.Printf("  Time Range: %s to %s\n",
// 				formatTime(firstEph.Epoch),
// 				formatTime(lastEph.Epoch))

// 			fmt.Printf("  First Ephemeris Data:\n")
// 			printData(firstEph)
// 		}

// 		fmt.Println(strings.Repeat("-", 60))
// 	}
// }

// =========================================================================

// =========================================================================
func PrintRINEXHeader(header RINEXHeader) {
	fmt.Println("RINEX Header:")

	fmt.Printf("  Version: %.2f\n", header.Version())

	typ, err := header.Type()
	fmt.Printf("  Type: %s (Error: %v)\n", typ, err)

	satSystem, err := header.SatelliteSystem()
	fmt.Printf("  Satellite System: %s (Error: %v)\n", satSystem, err)

	programName, err := header.ProgramName()
	fmt.Printf("  Program: %s (Error: %v)\n", programName, err)

	agency, err := header.Agency()
	fmt.Printf("  Agency: %s (Error: %v)\n", agency, err)

	date, err := header.Date()
	if err == nil {

		fmt.Printf("  Date: %s (Error: %v)\n", formatTime(date), err)
	} else {
		fmt.Printf("  Date: Error retrieving date: %v\n", err)
	}

	fmt.Println("  Comments:")
	comments, err := header.Comments()
	if err == nil {
		for i := 0; i < comments.Len(); i++ {
			comment, err := comments.At(i)
			fmt.Printf("    %s (Error: %v)\n", comment, err)
		}
	} else {
		fmt.Printf("    Error retrieving comments: %v\n", err)
	}

	fmt.Printf("  Leap Seconds: %d\n", header.LeapSeconds())
	fmt.Println(strings.Repeat("-", 60))
}

// ========================================

// ========================================

func formatTime(t Time) string {
	return fmt.Sprintf("%d.%09d", t.Seconds(), t.Nanoseconds())
}

// ========================================

// ========================================

func PrintRINEXEphemeris(eph RINEXEphemeris) {
	fmt.Println("RINEX Ephemeris Data:")
	fmt.Printf("Satellite ID: %d\n", eph.SatelliteId())

	epoch, err := eph.Epoch()
	if err == nil {
		fmt.Printf("Epoch: %s\n", formatTime(epoch))
	} else {
		fmt.Printf("Epoch: Error retrieving epoch: %v\n", err)
	}

	fmt.Printf("Clock Bias: %.12e\n", eph.ClockBias())
	fmt.Printf("Relative Frequency Bias: %.12e\n", eph.RelativeFrequencyBias())
	fmt.Printf("Message Frame Time: %.6f\n", eph.MessageFrameTime())

	fmt.Printf("Position (X, Y, Z): %.3f, %.3f, %.3f\n", eph.PositionX(), eph.PositionY(), eph.PositionZ())
	fmt.Printf("Velocity (X, Y, Z): %.6f, %.6f, %.6f\n", eph.VelocityX(), eph.VelocityY(), eph.VelocityZ())
	fmt.Printf("Acceleration (X, Y, Z): %.9f, %.9f, %.9f\n", eph.AccelerationX(), eph.AccelerationY(), eph.AccelerationZ())

	fmt.Printf("Health: %.0f\n", eph.Health())
	fmt.Printf("Frequency Channel Offset: %d\n", eph.FrequencyChannelOffset())
	fmt.Printf("Information Age: %.6f\n", eph.InformationAge())
}
