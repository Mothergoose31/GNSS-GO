package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

//  todo Azimuth-Elevation-Range Coordinates

type NedConverter struct {
	Ecef2NedMatrix *mat.Dense
}

// =====================================================

// =====================================================

func matMult(a, b [3][3]float64) [3][3]float64 {
	var result [3][3]float64

	result[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0] + a[0][2]*b[2][0]
	result[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1] + a[0][2]*b[2][1]
	result[0][2] = a[0][0]*b[0][2] + a[0][1]*b[1][2] + a[0][2]*b[2][2]

	result[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0] + a[1][2]*b[2][0]
	result[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1] + a[1][2]*b[2][1]
	result[1][2] = a[1][0]*b[0][2] + a[1][1]*b[1][2] + a[1][2]*b[2][2]

	result[2][0] = a[2][0]*b[0][0] + a[2][1]*b[1][0] + a[2][2]*b[2][0]
	result[2][1] = a[2][0]*b[0][1] + a[2][1]*b[1][1] + a[2][2]*b[2][1]
	result[2][2] = a[2][0]*b[0][2] + a[2][1]*b[1][2] + a[2][2]*b[2][2]

	return result
}

// =====================================================

// =====================================================

func rotMatrix(roll, pitch, yaw float64) [3][3]float64 {

	cr, sr := math.Cos(roll), math.Sin(roll)
	cp, sp := math.Cos(pitch), math.Sin(pitch)
	cy, sy := math.Cos(yaw), math.Sin(yaw)

	rr := [3][3]float64{
		{1, 0, 0},
		{0, cr, -sr},
		{0, sr, cr},
	}

	rp := [3][3]float64{
		{cp, 0, sp},
		{0, 1, 0},
		{-sp, 0, cp},
	}

	ry := [3][3]float64{
		{cy, -sy, 0},
		{sy, cy, 0},
		{0, 0, 1},
	}

	return matMult(ry, matMult(rp, rr))
}

// =====================================================

// =====================================================

func Rotate(axis []float64, angle float64) *mat.Dense {
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)
	oneMinusCos := 1 - cosAngle

	ret1 := mat.NewDense(3, 3, []float64{
		oneMinusCos * axis[0] * axis[0], oneMinusCos * axis[0] * axis[1], oneMinusCos * axis[0] * axis[2],
		oneMinusCos * axis[1] * axis[0], oneMinusCos * axis[1] * axis[1], oneMinusCos * axis[1] * axis[2],
		oneMinusCos * axis[2] * axis[0], oneMinusCos * axis[2] * axis[1], oneMinusCos * axis[2] * axis[2],
	})

	ret2 := mat.NewDense(3, 3, []float64{
		cosAngle, 0, 0,
		0, cosAngle, 0,
		0, 0, cosAngle,
	})

	ret3 := mat.NewDense(3, 3, []float64{
		0, -sinAngle * axis[2], sinAngle * axis[1],
		sinAngle * axis[2], 0, -sinAngle * axis[0],
		-sinAngle * axis[1], sinAngle * axis[0], 0,
	})

	var result mat.Dense
	result.Add(ret1, ret2)
	result.Add(&result, ret3)
	fmt.Println(result.Dims())

	return &result
}

// ===================== Benchmarking =====================
func generateRandomMatrix() [3][3]float64 {
	var matrix [3][3]float64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			matrix[i][j] = r.Float64()
		}
	}

	return matrix
}

func TestMatrixMultiplication1(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		a := generateRandomMatrix()
		b := generateRandomMatrix()
		matMult(a, b)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)

}

func randomAngle() float64 {
	return (rand.Float64() * 2 * math.Pi) - math.Pi
}

func generateRandomAngles() (float64, float64, float64) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	roll := randomAngle()
	pitch := randomAngle()
	yaw := randomAngle()
	return roll, pitch, yaw
}

func TestMatrixRotation(iterations int) {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		roll, pitch, yaw := generateRandomAngles()
		rotMatrix(roll, pitch, yaw)
	}

	duration := time.Since(start)
	fmt.Printf("Time taken for %d iterations: %v\n", iterations, duration)

}

// ===================== Benchmarking =====================
