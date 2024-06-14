package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/mothergoose31/GNNS-GO/GNSS/helpers"
)

func main() {

	r := rand.New(rand.NewSource(rand.Int63()))
	geodeticCoordinates := [][]float64{
		{47.6062, -122.3321},
		{35.6895, 139.6917},
		{51.5074, -0.1278},
		{34.0522, -118.2437},
		{39.9042, 116.4074},
		{19.4326, -99.1332},
	}

	for i := range geodeticCoordinates {
		altitude := r.Float64() * 5000
		geodeticCoordinates[i] = append(geodeticCoordinates[i], altitude)
	}

	for _, coord := range geodeticCoordinates {
		fmt.Printf("Lat: %f, Lon: %f, Alt: %f\n", coord[0], coord[1], coord[2])
	}

	ecefCoordinates := helpers.GeodeticToECEF(geodeticCoordinates, false)
	for _, coord := range ecefCoordinates {
		fmt.Printf("ECEF Coordinates: %v\n", coord)
	}

	axis := []float64{1, 0, 0}
	angle := math.Pi / 4

	rr := helpers.Rotate(axis, angle)
	fmt.Println(rr)

	// helpers.TestMatrixMultiplication1(10000000)
	// helpers.TestMatrixRotation(10000000)
	helpers.BenchQuaternion2Euler(10000000)
}
