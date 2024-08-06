package main

import (
	"fmt"

	gnss "github.com/mothergoose31/GNNS-GO/GNSS"
)

func main() {

	// fmt.Println(gnss.ParseSP3File("WUM0MGXULT_20241042300_01D_05M_ORB.SP3"))
	// _, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	// if err != nil {
	// 	log.Fatalf("Failed to create new message: %v", err)
	// }

	// ephemeris, err := gnss.NewEphemeris(seg)
	// if err != nil {
	// 	log.Fatalf("Failed to create new Ephemeris: %v", err)
	// }

	filename := "./brdc2050.24g"
	// fmt.Printf("Debug: Parsing file %s\n", filename)

	header, ephemerides, err := gnss.ParseRINEXFileV201(filename)
	if err != nil {
		fmt.Printf("Error parsing RINEX file: %v\n", err)
		return
	}
	fmt.Println(ephemerides)
	gnss.PrintRINEXHeader(*header)

	// print the first 5 ephemerides

	for i, eph := range ephemerides {
		if i >= 5 {
			break
		}
		fmt.Println("====================================")
		gnss.PrintRINEXEphemeris(eph)
	}

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
