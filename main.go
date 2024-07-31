package main

import (
	"fmt"
	"log"

	// "os"

	// "math"
	// "math/rand"
	"capnproto.org/go/capnp/v3"
	gnss "github.com/mothergoose31/GNNS-GO/GNSS"
	// "github.com/mothergoose31/GNNS-GO/GNSS/helpers"
)

// func printEphemeris(e gnss.Ephemeris) {
// 	fmt.Println("Ephemeris Details:")
// 	fmt.Printf("SatelliteId: %d\n", e.SatelliteId())
// 	fmt.Printf("ToeWeek: %d\n", e.ToeWeek())
// 	fmt.Printf("Toe: %f\n", e.Toe())
// 	fmt.Printf("TimeOfClockWeek: %d\n", e.TimeOfClockWeek())
// 	fmt.Printf("TimeOfClock: %f\n", e.TimeOfClock())
// 	fmt.Printf("Af0: %f\n", e.Af0())
// 	fmt.Printf("Af1: %f\n", e.Af1())
// 	fmt.Printf("Af2: %f\n", e.Af2())
// 	fmt.Printf("Iode: %f\n", e.Iode())
// 	fmt.Printf("Crs: %f\n", e.Crs())
// 	fmt.Printf("DeltaN: %f\n", e.DeltaN())
// 	fmt.Printf("M0: %f\n", e.M0())
// 	fmt.Printf("Cuc: %f\n", e.Cuc())
// 	fmt.Printf("Eccentricity: %f\n", e.Eccentricity())
// 	fmt.Printf("Cus: %f\n", e.Cus())
// 	fmt.Printf("SemiMajorAxis: %f\n", e.SemiMajorAxis())
// 	fmt.Printf("Cic: %f\n", e.Cic())
// 	fmt.Printf("Omega0: %f\n", e.Omega0())
// 	fmt.Printf("Cis: %f\n", e.Cis())
// 	fmt.Printf("Inclination: %f\n", e.Inclination())
// 	fmt.Printf("Crc: %f\n", e.Crc())
// 	fmt.Printf("PerigeeArgument: %f\n", e.PerigeeArgument())
// 	fmt.Printf("RateOfRightAscension: %f\n", e.RateOfRightAscension())
// 	fmt.Printf("RateOfInclination: %f\n", e.RateOfInclination())
// 	fmt.Printf("SatelliteHealth: %f\n", e.SatelliteHealth())
// }

func main() {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		log.Fatalf("Failed to create new message: %v", err)
	}

	ephemeris, err := gnss.NewEphemeris(seg)
	if err != nil {
		log.Fatalf("Failed to create new Ephemeris: %v", err)
	}
	fmt.Println(msg)

	ephemeris.SetSatelliteId(1)
	ephemeris.SetToeWeek(2000)
	ephemeris.SetToe(14400)
	ephemeris.SetTimeOfClockWeek(2000)
	ephemeris.SetTimeOfClock(14400)
	ephemeris.SetAf0(0.000100)
	ephemeris.SetAf1(0.000001)
	ephemeris.SetAf2(0)
	ephemeris.SetIode(1)
	ephemeris.SetCrs(100)
	ephemeris.SetDeltaN(0.000001)
	ephemeris.SetM0(1.0)
	ephemeris.SetCuc(0.000100)
	ephemeris.SetEccentricity(0.001)
	ephemeris.SetCus(0.000100)
	ephemeris.SetSemiMajorAxis(26559800.0)
	ephemeris.SetCic(0.000100)
	ephemeris.SetOmega0(1.0)
	ephemeris.SetCis(0.000100)
	ephemeris.SetInclination(0.96)
	ephemeris.SetCrc(100)
	ephemeris.SetPerigeeArgument(1.0)
	ephemeris.SetRateOfRightAscension(0.000001)
	ephemeris.SetRateOfInclination(0)
	ephemeris.SetSatelliteHealth(0)

	filename := "./brdc2050.24g"
	fmt.Printf("Debug: Parsing file %s\n", filename)

	header, ephemerides, err := gnss.ParseRINEXFile(filename)
	if err != nil {
		fmt.Printf("Error parsing RINEX file: %v\n", err)
		return
	}

	fmt.Printf("\nRINEX Header:\n")
	fmt.Printf("Version: %.2f\n", header.Version)
	fmt.Printf("Type: %s\n", header.Type)
	fmt.Printf("Satellite System: %s\n", header.SatSystem)
	fmt.Printf("Program: %s\n", header.ProgramName)
	fmt.Printf("Agency: %s\n", header.Agency)
	fmt.Printf("Date: %s\n", header.Date)
	fmt.Printf("Comments:\n")
	for _, comment := range header.Comments {
		fmt.Printf("  %s\n", comment)
	}

	fmt.Printf("\nNumber of GLONASS Ephemerides: %d\n", len(ephemerides))
	for i, eph := range ephemerides {
		if i >= 5 {
			break
		}
		fmt.Printf("\nEphemeris %d:\n", i+1)
		fmt.Printf("Satellite ID: %d\n", eph.SatelliteID)
		fmt.Printf("Epoch: %s\n", eph.Epoch)
		fmt.Printf("Clock Bias: %.12f\n", eph.ClockBias)
		fmt.Printf("Position (X, Y, Z): %.3f, %.3f, %.3f\n", eph.PositionX, eph.PositionY, eph.PositionZ)
		fmt.Printf("Velocity (X, Y, Z): %.3f, %.3f, %.3f\n", eph.VelocityX, eph.VelocityY, eph.VelocityZ)
		fmt.Printf("Frequency Number: %d\n", eph.FrequencyChannelOfSet)
	}

	fmt.Println("\nRINEX file parsing completed successfully.")

	// timeOfInterest := gnss.NewGPSTime(2000, 18000)

	// position, velocity, clockError, clockRateError, err := gpsEph.GetSatInfo(timeOfInterest)
	// if err != nil {
	// 	fmt.Printf("Error calculating satellite info: %v\n", err)
	// 	return
	// }

	// fmt.Printf("Satellite Position (X, Y, Z): %.2f, %.2f, %.2f meters\n", position[0], position[1], position[2])
	// fmt.Printf("Satellite Velocity (X, Y, Z): %.2f, %.2f, %.2f meters/second\n", velocity[0], velocity[1], velocity[2])
	// fmt.Printf("Satellite Clock Error: %.9f seconds\n", clockError)
	// fmt.Printf("Satellite Clock Rate Error: %.9f seconds/second\n", clockRateError)

	// r := rand.New(rand.NewSource(rand.Int63()))
	// geodeticCoordinates := [][]float64{
	// 	{47.6062, -122.3321},
	// 	{35.6895, 139.6917},
	// 	{51.5074, -0.1278},
	// 	{34.0522, -118.2437},
	// 	{39.9042, 116.4074},
	// 	{19.4326, -99.1332},
	// }

	// for i := range geodeticCoordinates {
	// 	altitude := r.Float64() * 5000
	// 	geodeticCoordinates[i] = append(geodeticCoordinates[i], altitude)
	// }

	// for _, coord := range geodeticCoordinates {
	// 	fmt.Printf("Lat: %f, Lon: %f, Alt: %f\n", coord[0], coord[1], coord[2])
	// }

	// ecefCoordinates := helpers.GeodeticToECEF(geodeticCoordinates, false)
	// for _, coord := range ecefCoordinates {
	// 	fmt.Printf("ECEF Coordinates: %v\n", coord)
	// }

	// axis := []float64{1, 0, 0}
	// angle := math.Pi / 4

	// rr := helpers.Rotate(axis, angle)
	// fmt.Println(rr)

	// ecef := [][]float64{
	// 	{4510733.0, 4510733.0, 4510733.0},
	// 	{6378137.0, 0.0, 0.0},
	// 	{0.0, 6378137.0, 0.0},
	// 	{0.0, 0.0, 6378137.0},
	// }
	// fmt.Println(helpers.ECEFToGeodetic(ecef, false))

	// geodetic := []float64{47.6062, -122.3321, 10.0}
	// fmt.Println("====================================")
	// fmt.Println(helpers.GeodeticToECEF([][]float64{geodetic}, false))
	// fmt.Println("====================================")
	// // fmt.Println(helpers.GeodeticToECEF2([][]float64{geodetic}, false))
	// fmt.Println("====================================")
	// fmt.Println("== ECEF to Geodetic ===============")
	// fmt.Println(helpers.ECEFToGeodetic(helpers.GeodeticToECEF([][]float64{geodetic}, false), false))
	// fmt.Println("====================================")
	// fmt.Println(helpers.ECEFToGeodetic(helpers.GeodeticToECEF2([][]float64{geodetic}, false), false))
	// fmt.Println(helpers.NewLocalCoordinates(geodetic))

	// helpers.TestMatrixMultiplication1(10000000)
	// helpers.TestMatrixRotation(10000000)
	// // helpers.BenchQuaternion2Euler(10000000)
	// helpers.BenchEuler2Quarterion(10000000)
	// helpers.BenchQuaternion2Rot(100000000)
	// helpers.BenchQuaternion2Rot2(100000000)
	// helpers.BenchECEFToGeodetic(10000000)
}
